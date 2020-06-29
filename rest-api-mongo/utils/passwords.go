package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword encrypts the given password
func HashPassword(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), 14)
	return string(bytes), err
}

// CheckPasswordHash is used to compare the encrypted and normal password
//and the one received in request.
func CheckPasswordHash(str, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
	return err == nil
}
