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
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type JWTCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Link struct {
	Short_id string    `json:"shortid"`
	Orig_url string    `json:"origurl"`
	Expiry   time.Time `json:"expiry"`
}
