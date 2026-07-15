package server

import "github.com/golang-jwt/jwt/v5"

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,min=3"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type jwtCustomClaims struct {
	jwt.RegisteredClaims
}
