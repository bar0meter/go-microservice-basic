package social_logins

import (
	"context"
	"fmt"
	log "github.com/frost060/go-microservice-basic/rest-api-mongo/logging"
	"io/ioutil"
	"net/http"

	"github.com/frost060/go-microservice-basic/rest-api-mongo/configs"
)

// GoogleHandler model
type GoogleHandler struct {
	google *configs.GoogleConfig
}

// NewGoogleHandler creates a new GoogleHandler for handling route '/google'
func NewGoogleHandler(google *configs.GoogleConfig) *GoogleHandler {
	return &GoogleHandler{google}
}

// GoogleHome => /google
func (g *GoogleHandler) GoogleHome(rw http.ResponseWriter, r *http.Request) {
	log.Info("Google Home Page")
	var html = `<html><body><a href="/google/login">Google Login In</a></body></html>`
	_, _ = fmt.Fprint(rw, html)
}

// GoogleLogin => /google/login
func (g *GoogleHandler) GoogleLogin(rw http.ResponseWriter, r *http.Request) {
	log.Info("Google Login Page")
	url := g.google.GoogleOauth.AuthCodeURL(g.google.RandomState)
	http.Redirect(rw, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback => /google/callback
func (g *GoogleHandler) GoogleCallback(rw http.ResponseWriter, r *http.Request) {
	log.Info("Google Callback url hit")
	if r.FormValue("state") != g.google.RandomState {
		log.Error("Google State is not valid")
		http.Redirect(rw, r, "/google", http.StatusTemporaryRedirect)
		return
	}

	token, err := g.google.GoogleOauth.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		log.Error("Google Couldn't get token", "error", err)
		http.Redirect(rw, r, "/google", http.StatusTemporaryRedirect)
		return
	}

	res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Error("Google Couldn't create get r", "error", err)
		http.Redirect(rw, r, "/google", http.StatusTemporaryRedirect)
		return
	}

	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("Google Couldn't parse rw", "error", err)
		http.Redirect(rw, r, "/google", http.StatusTemporaryRedirect)
		return
	}

	log.Info("Google access token", "token", token.AccessToken)
	_, _ = fmt.Fprintf(rw, "Response: %s", content)
}
