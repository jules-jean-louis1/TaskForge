package handler

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	ID        uuid.UUID `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	jwt.RegisteredClaims
}

func GenerateAccesToken(id uuid.UUID, firstname string, lastname string, email string, role string) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	var claims Claims

	claims = Claims{
		ID:        id,
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
