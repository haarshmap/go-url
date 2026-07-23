package server

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	UserID   int64  `json:"userid"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type JWTCustomClaims struct {
	UserID   int64  `json:"userid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type Links struct {
	ShortID string    `json:"shortid"`
	OrigURL string    `json:"origurl"`
	UserID  int64     `json:"userid"`
	Expiry  time.Time `json:"expiry"`
}
