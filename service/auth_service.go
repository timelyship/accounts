package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
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

func RefreshToken(accessToken, refreshToken string) (*response.LoginResponse, *utility.RestError) {
	token, err := repository.GetTokenByRefreshToken(refreshToken)
	if err != nil {
		return nil, utility.NewUnAuthorizedError("Token persistence failed", &err.Error)
	}
	if token.AccessToken != accessToken {
		return nil, utility.NewUnAuthorizedError("Invalid at,rt pair", nil)
	}
	claims, err := utility.DecodeToken(token.RefreshToken, os.Getenv("REFRESH_SECRET"))
	sub := (*claims)["sub"]
	userId, hexErr := primitive.ObjectIDFromHex(sub.(string))
	if hexErr != nil {
		return nil, utility.NewUnAuthorizedError("Invalid subject", &hexErr)
	}
	user, userErr := repository.GetUserById(userId)
	if userErr != nil {
		return nil, utility.NewUnAuthorizedError("User not found for subject", &userErr.Error)
	}
	newAccessToken, tErr := utility.CreateAccessToken(user, "*")
	if tErr != nil {
		return nil, utility.NewUnAuthorizedError("Failed to create access token", &tErr)
	}
	token.AccessToken = newAccessToken.AccessToken
	token.AccessUuid = newAccessToken.AccessUuid
	token.AtExpires = newAccessToken.AtExpires
	updErr := repository.UpdateToken(token)
	if updErr != nil {
		return nil, utility.NewUnAuthorizedError("Failed to create access token", &updErr.Error)
	}
	return &response.LoginResponse{
		RefreshToken: token.RefreshToken,
		AccessToken:  token.AccessToken,
	}, nil

}
