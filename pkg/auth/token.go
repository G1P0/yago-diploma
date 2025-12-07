package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("secret")

type Claims struct {
	PasswordHash string `json:"ph"`
	jwt.RegisteredClaims
}

func hashPassword(pass string) string {
	sum := sha256.Sum256([]byte(pass))
	return hex.EncodeToString(sum[:])
}

func GenerateToken(pass string) (string, error) {
	h := hashPassword(pass)

	claims := Claims{
		PasswordHash: h,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtKey)
}

func ValidateToken(tokenString string, currentPassword string) (bool, error) {
	if tokenString == "" {
		return false, errors.New("empty token")
	}

	t, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return false, err
	}

	claims, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return false, errors.New("invalid token")
	}

	if claims.PasswordHash != hashPassword(currentPassword) {
		return false, errors.New("password hash mismatch")
	}

	return true, nil
}
