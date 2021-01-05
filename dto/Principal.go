package dto

import (
	"github.com/dgrijalva/jwt-go"
)

type Role struct {
	Name string `json:"name"`
}

type Principal struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Picture   string `json:"picture"`
	Roles     []Role `json:"roles"`
	UserID    string `json:"sub"`
	jwt.StandardClaims
}
