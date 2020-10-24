package service

import (
	"timelyship.com/accounts/dto/request"
	"timelyship.com/accounts/utility"
)

func InitiateSignUp(signUpRequest request.SignUpRequest) utility.RestError {
	// check if an user exists with the email
	//create a new user with email unverified
	//send verification email using sqs
	//send 200
}
