package jwt

import "github.com/golang-jwt/jwt/v5"

type claims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
}

type registerUser struct {
	Username string `json:"user"`
	Password string `json:"password"`
}

type loginUser registerUser
