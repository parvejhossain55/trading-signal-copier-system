package types

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
	jwt.RegisteredClaims
}

type ContextKey string

const UserContextKey ContextKey = "user"
