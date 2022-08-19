package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(p string) string {
	hashBytes, _ := bcrypt.GenerateFromPassword([]byte(p), 10)
	hash := string(hashBytes)
	return hash
}
