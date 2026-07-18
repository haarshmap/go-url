package server

import "github.com/golang-jwt/jwt/v5"

type RegisterRequest struct {
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
