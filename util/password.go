package util

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword return the bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hashedString, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedString), err
}

// CheckPassword check if the provided password is correct or not
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
