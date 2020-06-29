package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTClaims struct {
	UserID      int64  `json:"id"`
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
	jwt.StandardClaims
}

// NewJWTClaims => Payload passed along with signature and encryption method
// expiration is time when token expires
func NewJWTClaims(userId int64, username string, accessToken string, expirationTime time.Time) *JWTClaims {
	return &JWTClaims{
		UserID:      userId,
		Username:    username,
		AccessToken: accessToken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
}

func (jwtClaim *JWTClaims) GetSignedToken(jwtKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

// ValidateAndRefreshToken validates the token in the request
// If token is about it also refreshes the token
// We ensure that a new token is not issued until enough time has elapsed
// In this case, a new token will only be issued if the old token is within
// 30 seconds of expiry. Otherwise, return a bad request status
// Here bool in return type tell whether token is new or false
// expirationTime is time when refreshToken will expire
func ValidateAndRefreshToken(tokenRequest string, jwtKey []byte, expireThreshold int64, newExpirationTime time.Time) (
	*JWTClaims, string, bool) {
	jwtClaim := &JWTClaims{}
	tkn, err := jwt.ParseWithClaims(tokenRequest, jwtClaim, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, tokenRequest, false
	}

	if !tkn.Valid {
		return jwtClaim, tokenRequest, false
	}

	if time.Unix(jwtClaim.ExpiresAt, 0).Sub(time.Now()) > time.Duration(expireThreshold) {
		return jwtClaim, tokenRequest, true
	}

	jwtClaim.ExpiresAt = newExpirationTime.Unix()
	newToken, _ := jwtClaim.GetSignedToken(jwtKey)
	return jwtClaim, newToken, true
}
