package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword Hashes user password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("unable to hash password: %v", err)
	}
	return string(bytes), nil
}

// CheckPassword checks the given password against hashed password
func CheckPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}
