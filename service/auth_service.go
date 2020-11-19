package service

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"strings"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/dto/request"
	"timelyship.com/accounts/dto/response"
	"timelyship.com/accounts/repository"
	"timelyship.com/accounts/utility"
)

func HandleLogin(loginRequest request.LoginRequest) (*response.LoginResponse, *utility.RestError) {
	user, err := repository.GetUserByEmailOrPhone(loginRequest.EmailOrPhone)
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

func GenerateCode(token *jwt.Token, newAud, state string) (string, *utility.RestError) {
	encKey, encKErr := repository.GetEncKeyByState(state)
	if encKErr != nil || encKey == "" {
		return "", utility.NewUnAuthorizedError("Invalid state", &encKErr.Error)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", utility.NewUnAuthorizedError("Can not get claims", nil)
	}
	curAud := claims["aud"].(string)
	if curAud != "*" {
		return "", utility.NewUnAuthorizedError("Insufficient privilege", nil)
	}
	userId, hErr := primitive.ObjectIDFromHex(claims["sub"].(string))
	if hErr != nil {
		return "", utility.NewUnAuthorizedError("Internal error", &hErr)
	}
	user, uError := repository.GetUserById(userId)
	if uError != nil {
		return "", utility.NewUnAuthorizedError("User could be fetched", &uError.Error)
	}
	tokenDetails, tErr := utility.CreateToken(user, newAud)
	if tErr != nil {
		return "", utility.NewUnAuthorizedError("Could not generate token", &tErr)
	}
	saveErr := repository.SaveToken(tokenDetails)
	if saveErr != nil {
		return "", utility.NewUnAuthorizedError("Token persistence failed", &saveErr.Error)
	}
	loginResponse := response.LoginResponse{
		RefreshToken: tokenDetails.RefreshToken,
		AccessToken:  tokenDetails.AccessToken,
	}
	bytes, jErr := json.Marshal(loginResponse)
	if jErr != nil {
		return "", utility.NewUnAuthorizedError("Could not marshal login response", &jErr)
	}
	code, err := utility.AESEncrypt(bytes, []byte(encKey))
	if err != nil {
		return "", utility.NewUnAuthorizedError("Encryption failed", &err.Error)
	}
	fmt.Println(code)
	fmt.Println(encKey)
	return code, nil
}

func InitiateLogin() (*map[string]string, *utility.RestError) {
	state := strings.Replace(uuid.New().String(), "-", "", -1)
	key := strings.Replace(uuid.New().String(), "-", "", -1)
	loginState := &domain.LoginState{
		State: state,
		Key:   key,
	}
	err := repository.SaveLoginState(loginState)
	if err != nil {
		return nil, err
	}
	return &map[string]string{
		"state": state,
		"key":   key,
	}, nil
}
