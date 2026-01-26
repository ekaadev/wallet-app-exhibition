package model

import "github.com/golang-jwt/jwt/v5"

type Auth struct {
	UserID   *uint  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
