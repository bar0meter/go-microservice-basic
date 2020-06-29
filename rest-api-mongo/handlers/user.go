package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/frost060/go-microservice-basic/rest-api-mongo/configs"
	"github.com/frost060/go-microservice-basic/rest-api-mongo/models"
	"github.com/frost060/go-microservice-basic/rest-api-mongo/services/email"
	"github.com/frost060/go-microservice-basic/rest-api-mongo/utils"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userRepo     models.UserRepository
	serverConfig *configs.Config
}

func NewUserHandler(
	userRepo models.UserRepository, serverConfig *configs.Config) *UserHandler {
	return &UserHandler{userRepo, serverConfig}
}

func (u *UserHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	jwtClaim, valid := u.extractJWTClaimFromRequest(r)
	if !valid {
		http.Error(rw, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	user, err := u.userRepo.FindByID(jwtClaim.UserID)
	if err != nil {
		http.Error(rw, "Failed to read database", http.StatusInternalServerError)
		return
	}

	err = user.Serialize(rw)
	if err != nil {
		http.Error(rw, "Failed to marshal data", http.StatusInternalServerError)
	}
}

// extractJWTClaimFromRequest extract payload from request context which was decoded in the jwt middleware
func (u *UserHandler) extractJWTClaimFromRequest(r *http.Request) (*utils.JWTClaims, bool) {
	claim := r.Context().Value("claim")
	if claim == nil {
		return nil, false
	}

	jwtClaim, ok := claim.(*utils.JWTClaims)
	if !ok {
		return nil, false
	}
	return jwtClaim, true
}

func (u *UserHandler) PerformLogin(rw http.ResponseWriter, r *http.Request) {
	var userRequest models.User
	err := userRequest.ParseJSON(r.Body)
	if err != nil {
		http.Error(rw, "Error occurred while unmarshalling data", http.StatusBadRequest)
		return
	}

	err = userRequest.Validate()
	if err != nil {
		http.Error(rw, fmt.Sprintf("Error validating user: %s", err), http.StatusBadRequest)
		return
	}

	userDB, err := u.userRepo.FindByUsername(userRequest.Username)
	if err != nil {
		http.Error(rw, "Invalid username/password", http.StatusForbidden)
		return
	}

	if !utils.CheckPasswordHash(userRequest.Password, userDB.Password) {
		http.Error(rw, "Invalid username/password", http.StatusForbidden)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	expirationTime := time.Now().Add(time.Duration(u.serverConfig.JWT.ExpirationTime) * time.Second)
	jwtClaim := utils.NewJWTClaims(userDB.ID, userDB.Username, "", expirationTime)
	signedToken, err := jwtClaim.GetSignedToken(u.serverConfig.JWT.SecretKey)

	if err != nil {
		http.Error(rw, "Error occurred while creating jwt token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:    "token",
		Value:   signedToken,
		Expires: expirationTime,
		Path:    "/",
	})

	rw.WriteHeader(http.StatusOK)
}

func (u *UserHandler) SaveUser(rw http.ResponseWriter, r *http.Request) {
	var user models.User
	err := user.ParseJSON(r.Body)
	if err != nil {
		http.Error(rw, "Error occurred while unmarshalling data", http.StatusBadRequest)
		return
	}

	// Encrypts password for saving
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(
			rw,
			fmt.Sprintf("Error occurred while hashing password for user: %s", user.Username),
			http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	err = u.userRepo.Save(&user)
	if err != nil {
		http.Error(rw, "Error occurred while saving user to DB", http.StatusInternalServerError)
	}

	rw.WriteHeader(http.StatusOK)
}

// TODO implement this
func (u *UserHandler) NewUserSignUp(rw http.ResponseWriter, r *http.Request) {

}

func (u *UserHandler) RequestResetPassword(rw http.ResponseWriter, r *http.Request) {
	var user models.User
	err := user.ParseJSON(r.Body)
	if err != nil {
		http.Error(rw, "Error occurred while unmarshalling data", http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(user.Username)
	if username == "" {
		http.Error(rw, "Error occurred while unmarshalling data", http.StatusBadRequest)
		return
	}

	tagClaim, name, err := u.userRepo.ResetPassword(username)
	if err != nil || tagClaim == nil {
		rw.WriteHeader(http.StatusOK)
		return
	}

	resetTag, err := tagClaim.GetSignedToken(u.serverConfig.JWT.SecretKey)
	if err != nil {
		http.Error(rw, "Error occurred while generating reset password token", http.StatusInternalServerError)
		return
	}

	templatePath := u.serverConfig.RootPath + "/web/templates/forgot_password.html"
	html := utils.GenerateResetPasswordMail(username, name, resetTag, "http://"+r.Host, templatePath)
	body := email.GetHtmlBody(username, "Reset Password from TEST", html)
	ok, err := email.SendMail(body, u.serverConfig.SendGrid.ApiKey)

	if err != nil || !ok {
		http.Error(rw, "Error occurred while "+
			"sending reset password mail", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (u *UserHandler) ResetPassword(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encodedTag := vars["slug"]

	var tagClaim utils.TagClaim
	ok, err := tagClaim.DecodeTag(encodedTag, utils.ResetPassword, u.serverConfig.JWT.SecretKey)

	if err != nil || !ok {
		http.Error(rw, "Invalid token", http.StatusBadRequest)
		return
	}

	ok, err = u.userRepo.ValidateResetPassword(&tagClaim)
	if err != nil || !ok {
		http.Error(rw, "Invalid token", http.StatusBadRequest)
		return
	}

	// TODO => handle password change in front end.
	_, _ = fmt.Fprintln(rw, "Reset successful for user: "+tagClaim.Identity)
}

func (u *UserHandler) RequestVerifyEmail(rw http.ResponseWriter, r *http.Request) {
	jwtClaim, valid := u.extractJWTClaimFromRequest(r)
	if !valid {
		http.Error(rw, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	tagClaim, name, err := u.userRepo.VerifyUser(jwtClaim.UserID, "")
	if err != nil {
		http.Error(rw, "Failed to read database", http.StatusInternalServerError)
		return
	}

	resetTag, err := tagClaim.GetSignedToken(u.serverConfig.JWT.SecretKey)
	if err != nil {
		http.Error(rw, "Error occurred while generating verify email token", http.StatusInternalServerError)
		return
	}

	templatePath := u.serverConfig.RootPath + "/web/templates/verify_email.html"
	html := utils.GenerateVerifyEmailMail(tagClaim.Identity, name, resetTag, "http://"+r.Host, templatePath)
	body := email.GetHtmlBody(tagClaim.Identity, "Verify Email Address from TEST", html)
	ok, err := email.SendMail(body, u.serverConfig.SendGrid.ApiKey)

	if err != nil || !ok {
		http.Error(rw, "Error occurred while "+
			"sending reset password mail", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (u *UserHandler) VerifyEmail(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encodedTag := vars["slug"]

	var tagClaim utils.TagClaim
	ok, err := tagClaim.DecodeTag(encodedTag, utils.VerifyEmail, u.serverConfig.JWT.SecretKey)

	if err != nil || !ok {
		http.Error(rw, "Invalid token", http.StatusBadRequest)
		return
	}

	ok, err = u.userRepo.ValidateVerifyUser(&tagClaim)
	if err != nil || !ok {
		http.Error(rw, "Invalid token", http.StatusBadRequest)
		return
	}

	_, _ = fmt.Fprintln(rw, "Verify Email successful for user: "+tagClaim.Identity)
}
