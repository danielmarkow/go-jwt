package token

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/golang-jwt/jwt/v5"
	"testing"
	"time"
)

func TestToken(t *testing.T) {
	t.Run("Generate Access Token", func(t *testing.T) {
		token, err := CreateToken("test@test.de", 1)
		if err != nil {
			t.Error("expected CreateToken to not throw errors")
		}
		err = VerifyToken(token)
		if err != nil {
			t.Error("expected token to be valid")
		}
	})

	t.Run("Generate Refresh Token", func(t *testing.T) {
		token, err := CreateRefreshToken(1)
		if err != nil {
			t.Error("expected CreateRefreshToken not to throw errors")
		}
		err = VerifyToken(token)
		if err != nil {
			t.Error("expected token to be valid")
		}
	})

	t.Run("Verify token that is signed with a different secret key", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"exp": time.Now().Add(time.Hour * 24).Unix(),
				"sub": 1,
			})
		tokenString, err := token.SignedString([]byte("different-key"))
		if err != nil {
			t.Fatalf("expected signing to work: %s \n", err.Error())
		}

		err = VerifyToken(tokenString)
		if err == nil {
			t.Error("expected verification to fail")
		}
	})

	t.Run("Verify token that is signed with a different method (ES256 vs HMAC)", func(t *testing.T) {
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatalf("failed to generate ECDSA private key: %s \n", err.Error())
		}

		token := jwt.NewWithClaims(jwt.SigningMethodES256,
			jwt.MapClaims{
				"exp": time.Now().Add(time.Hour * 24).Unix(),
				"sub": 1,
			})
		tokenString, err := token.SignedString(privateKey)
		if err != nil {
			t.Fatalf("expected signing to work: %s \n", err.Error())
		}

		err = VerifyToken(tokenString)
		if err == nil {
			t.Error("expected verification to fail")
		}
	})
}
