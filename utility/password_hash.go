package utility

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		fmt.Printf("Error creating Hash = %v", err)
		panic(err)
	}
	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		return false
	}
	return true
}