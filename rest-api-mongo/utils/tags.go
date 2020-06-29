package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	ResetPassword = 0
	VerifyEmail   = 1
	UserInvite    = 2
)

type TagClaim struct {
	Identity string
	Type     int
	ID       string
	jwt.StandardClaims
}

func generateNewTag(userIdentity string, id string, tagType int, expirationTime time.Time) *TagClaim {
	tagClaim := &TagClaim{
		Identity: userIdentity,
		Type:     tagType,
		ID:       id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	return tagClaim
}

func (tc *TagClaim) GetSignedToken(jwtKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, tc)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

// DecodeTag decodes the tag passed in url.
// Returns true/false depending on whether tag is expired or not.
func (tc *TagClaim) DecodeTag(encodedTag string, tagType int, tagKey []byte) (bool, error) {
	tkn, err := jwt.ParseWithClaims(encodedTag, tc, func(token *jwt.Token) (interface{}, error) {
		return tagKey, nil
	})

	if err != nil {
		return false, err
	}

	return tc.Type == tagType && tkn.Valid, nil
}

func GenerateResetPasswordTag(identity string, id string) *TagClaim {
	expirationTime := time.Now().Add(24 * time.Hour)
	return generateNewTag(identity, id, ResetPassword, expirationTime)
}

func GenerateVerifyEmailTag(email, id string) *TagClaim {
	expirationTime := time.Now().Add(24 * time.Hour)
	return generateNewTag(email, id, VerifyEmail, expirationTime)
}
