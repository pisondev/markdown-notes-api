package web

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserID int `json:"userId"`
	jwt.RegisteredClaims
}
