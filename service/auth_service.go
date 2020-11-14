package service

import (
	"timelyship.com/accounts/dto/request"
	"timelyship.com/accounts/dto/response"
	"timelyship.com/accounts/repository"
	"timelyship.com/accounts/utility"
)

func HandleLogin(loginRequest request.LoginRequest) (*response.LoginResponse, *utility.RestError) {
	user, err := repository.GetUserByEmailOrPhone(loginRequest.Email, loginRequest.Phone)
	if err != nil {
		return nil, utility.NewUnAuthorizedError("User not found", nil)
	}
	if utility.ComparePasswords(user.Password, loginRequest.Password) == false {
		return nil, utility.NewUnAuthorizedError("Password not valid", nil)
	}
	tokenDetails, tokenError := utility.CreateToken(user, "*")
	if tokenError != nil {
		return nil, utility.NewUnAuthorizedError("Token creation failed", &tokenError)
	}
	saveErr := repository.SaveToken(tokenDetails)
	if saveErr != nil {
		return nil, utility.NewUnAuthorizedError("Token persistence failed", &saveErr.Error)
	}
	return &response.LoginResponse{
		RefreshToken: tokenDetails.RefreshToken,
		AccessToken:  tokenDetails.AccessToken,
	}, nil
}
