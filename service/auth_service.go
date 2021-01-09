package service

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"os"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/dto/request"
	"timelyship.com/accounts/dto/response"
	"timelyship.com/accounts/repository"
	"timelyship.com/accounts/utility"
)

type AuthService struct {
	AuthRepository repository.AuthRepository
	logger         zap.Logger
}

func ProvideAuthService(a repository.AuthRepository, l zap.Logger) AuthService {
	return AuthService{
		AuthRepository: a,
		logger:         l,
	}
}

func (s *AuthService) HandleLogin(loginRequest request.LoginRequest) (*response.LoginResponse, *utility.RestError) {
	s.logger.Info("Going to fetch user by username or email", zap.String("EmailOrPhone", loginRequest.EmailOrPhone))
	user, err := repository.GetUserByEmailOrPhone(loginRequest.EmailOrPhone)
	if err != nil {
		return nil, utility.NewUnAuthorizedError("User not found", nil)
	}
	if !utility.ComparePasswords(user.Password, loginRequest.Password) {
		return nil, utility.NewUnAuthorizedError("Password not valid", nil)
	}
	s.logger.Info("Password ok,going to create token")
	tokenDetails, tokenError := utility.CreateToken(user, "*")
	if tokenError != nil {
		return nil, utility.NewUnAuthorizedError("Token creation failed", &tokenError)
	}
	s.logger.Info("Token creation success")
	saveErr := repository.SaveToken(tokenDetails)
	if saveErr != nil {
		return nil, utility.NewUnAuthorizedError("Token persistence failed", &saveErr.Error)
	}
	s.logger.Info("Token save success")
	return &response.LoginResponse{
		RefreshToken: tokenDetails.RefreshToken,
		AccessToken:  tokenDetails.AccessToken,
	}, nil
}

func (s *AuthService) InitiateLogin() (*map[string]string, *utility.RestError) {
	state := fmt.Sprintf("%s%s%s%s", utility.GetUUIDWithoutDash(),
		utility.GetUUIDWithoutDash(), utility.GetUUIDWithoutDash(), utility.GetUUIDWithoutDash())
	key := utility.GetUUIDWithoutDash()
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

func (s *AuthService) RefreshToken(accessToken, refreshToken string) (*response.LoginResponse, *utility.RestError) {
	s.logger.Info("Getting old token by refresh token")
	token, err := repository.GetTokenByRefreshToken(refreshToken)
	if err != nil {
		return nil, utility.NewUnAuthorizedError("Token persistence failed", &err.Error)
	}
	if token.AccessToken != accessToken {
		s.logger.Error("Access token not matching")
		return nil, utility.NewUnAuthorizedError("Invalid at,rt pair", nil)
	}
	s.logger.Info("Decoding refresh token")
	claims, err := utility.DecodeToken(token.RefreshToken, os.Getenv("REFRESH_SECRET"))
	if err != nil {
		s.logger.Info("error decoding refresh token", zap.Error(err.Error))
	}
	sub := (*claims)["sub"]
	userID, hexErr := primitive.ObjectIDFromHex(sub.(string))
	if hexErr != nil {
		return nil, utility.NewUnAuthorizedError("Invalid subject", &hexErr)
	}
	s.logger.Info("Going to fetch user details by refresh token subject",
		zap.Any("refreshTokenSub", sub))
	user, userErr := repository.GetUserByID(userID)
	if userErr != nil {
		return nil, utility.NewUnAuthorizedError("User not found for subject", &userErr.Error)
	}
	s.logger.Info("Fetched user details by refresh token subject",
		zap.Any("refreshTokenSub", sub))
	newAccessToken, tErr := utility.CreateAccessToken(user, "*")
	if tErr != nil {
		return nil, utility.NewUnAuthorizedError("Failed to create access token", &tErr)
	}
	s.logger.Info("New token creation done")
	token.AccessToken = newAccessToken.AccessToken
	token.AccessUUID = newAccessToken.AccessUUID
	token.AtExpires = newAccessToken.AtExpires
	s.logger.Info("Updating token")
	updErr := repository.UpdateToken(token)
	if updErr != nil {
		return nil, utility.NewUnAuthorizedError("Failed to create access token", &updErr.Error)
	}
	return &response.LoginResponse{
		RefreshToken: token.RefreshToken,
		AccessToken:  token.AccessToken,
	}, nil
}

// Todo : consider some code documentation
func (s AuthService) GenerateCode(token *jwt.Token, newAud, state string) *utility.RestError {
	s.logger.Info("Going to fetch login state by state", zap.String("state", state))
	loginState, encKErr := repository.GetLoginState(state)
	if encKErr != nil || loginState.Key == "" {
		return utility.NewUnAuthorizedError("Invalid state", &encKErr.Error)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return utility.NewUnAuthorizedError("Can not get claims", nil)
	}
	curAud := claims["aud"].(string)
	if curAud != "*" {
		return utility.NewUnAuthorizedError("Insufficient privilege", nil)
	}
	userID, hErr := primitive.ObjectIDFromHex(claims["sub"].(string))
	if hErr != nil {
		return utility.NewUnAuthorizedError("Internal error", &hErr)
	}
	user, uError := repository.GetUserByID(userID)
	if uError != nil {
		return utility.NewUnAuthorizedError("User could be fetched", &uError.Error)
	}
	tokenDetails, tErr := utility.CreateToken(user, newAud)
	if tErr != nil {
		return utility.NewUnAuthorizedError("Could not generate token", &tErr)
	}
	saveErr := repository.SaveToken(tokenDetails)
	if saveErr != nil {
		return utility.NewUnAuthorizedError("Token persistence failed", &saveErr.Error)
	}
	loginResponse := response.LoginResponse{
		RefreshToken: tokenDetails.RefreshToken,
		AccessToken:  tokenDetails.AccessToken,
	}
	bytes, jErr := json.Marshal(loginResponse)
	fmt.Println("jwts", string(bytes))
	if jErr != nil {
		return utility.NewUnAuthorizedError("Could not marshal login response", &jErr)
	}
	code, err := utility.SimpleAESEncrypt([]byte(loginState.Key), string(bytes))
	if err != nil {
		return utility.NewUnAuthorizedError("Encryption failed", &err.Error)
	}
	loginState.Code = code
	updErr := repository.UpdateLoginState(loginState)
	if updErr != nil {
		return utility.NewUnAuthorizedError("Login state upd failed", &updErr.Error)
	}
	return nil
}

// Todo : Consider some documentation
func (s AuthService) ExchangeCode(state string) (*domain.LoginState, *utility.RestError) {
	data, err := repository.GetLoginState(state)
	if err != nil {
		return nil, utility.NewUnAuthorizedError("Invalid state", &err.Error)
	}
	return data, nil
}
