package handler

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func isItTheSamePw(hashedPassword []byte, providedPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, providedPassword)
	if err != nil {
		log.Printf("error comparing passwords: %s \n", err.Error())
		return false
	}
	return true
}
