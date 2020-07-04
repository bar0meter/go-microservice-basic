package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/frost060/go-microservice-basic/rest-api-mongo/configs"
	"github.com/frost060/go-microservice-basic/rest-api-mongo/logging"
	"github.com/frost060/go-microservice-basic/rest-api-mongo/utils"
)

type JWTMiddleWare struct {
	jwtConfig *configs.JWTConfig
	log       *logging.LogWrapper
}

func NewJWTMiddleWare(jwtConfig *configs.JWTConfig, l *logging.LogWrapper) *JWTMiddleWare {
	return &JWTMiddleWare{jwtConfig, l}
}

// Middleware for jwt validation and token refresh

func (jwt *JWTMiddleWare) ValidateAndRefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		c, err := request.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				http.Error(response, "Unauthorized request", http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			http.Error(response, "Bad request", http.StatusBadRequest)
			return
		}

		tknStr := c.Value
		newExpirationTime := time.Now().Add(time.Duration(jwt.jwtConfig.ExpirationTime) * time.Second)
		jwtClaim, newToken, ok := utils.ValidateAndRefreshToken(
			tknStr, jwt.jwtConfig.SecretKey,
			jwt.jwtConfig.ExpireInThreshold, newExpirationTime)

		if jwtClaim == nil || !ok {
			http.Error(response, "Unauthorized request", http.StatusUnauthorized)
			return
		}

		jwt.log.Info(fmt.Sprintf("User Payload in request cookie: %v", jwtClaim.Username))
		jwt.log.Info(fmt.Sprintf("Refreshed token: %s", newToken))
		http.SetCookie(response, &http.Cookie{
			Name:    "token",
			Value:   newToken,
			Expires: newExpirationTime,
			Path:    "/",
		})

		// Passing decode claim from middleware to handler. (Dont need to decode it again in the handler)
		ctx := context.WithValue(request.Context(), "claim", jwtClaim)
		next.ServeHTTP(response, request.WithContext(ctx))
	})
}
