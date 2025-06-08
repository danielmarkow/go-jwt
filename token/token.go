package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

var secretKey = []byte("secret-key") // TODO change to actual secret key!!!

// CreateToken needs to be changed in case uuids are used
func CreateToken(email string, userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(time.Minute * 15).Unix(),
			"sub":   userId,
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("error signing token: %s \n", err.Error())
		return "", err
	}

	return tokenString, nil
}

// CreateRefreshToken needs to be changed in case uuids are used
func CreateRefreshToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 24).Unix(),
			"sub": userId,
		})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("error signing refresh token: %s \n", err.Error())
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
