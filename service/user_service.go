package service

import (
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/utility"
)

func GetUser() {

}

func CreateUser(user domain.User) (*domain.User, *utility.RestError) {
	if user.Id == 100 {
		return nil, utility.NewInternalServerError("Failed to save user.")
	}
	return &user, nil
}

func FindUser() {

}
