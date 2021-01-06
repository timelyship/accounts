package dto

import (
	"github.com/dgrijalva/jwt-go"
)

type Role struct {
	Name string `json:"name"`
}

type Principal struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Picture   string `json:"picture"`
	Roles     []Role `json:"roles"`
	UserID    string `json:"sub"`
	jwt.StandardClaims
}
