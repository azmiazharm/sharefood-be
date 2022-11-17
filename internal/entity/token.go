package entity

import (
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
)

type TokenClaims struct {
	ID uuid.UUID `json:"id_user"`
	jwt.StandardClaims
}

type TokenResponse struct {
	Type      string    `json:"type"`
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}
