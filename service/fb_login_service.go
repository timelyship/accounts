package service

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"os"
	"time"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/dto"
	"timelyship.com/accounts/repository"
)

type FbLoginService struct {
	fbLoginRepository repository.FbLoginRepository
	logger            zap.Logger
	httpClient        HTTPClient
}

func ProvideFbLoginService(
	fbLoginRepository repository.FbLoginRepository, logger zap.Logger, httpClient HTTPClient) FbLoginService {
	return FbLoginService{
		fbLoginRepository: fbLoginRepository,
		logger:            logger,
	}
}

func (s *FbLoginService) GetFBRedirectURI(uiState string) (string, error) {
	state, eUUIDError := uuid.NewRandom()
	if eUUIDError != nil {
		s.logger.Error("Could not generate new uuid")
		return "", eUUIDError
	}
	s.logger.Info("uiState", zap.String("uiState", uiState))
	fAuth := dto.NewFBOAuth("token",
		os.Getenv("FB_OAUTH_CLIENT_ID"),
		os.Getenv("FB_OAUTH_SCOPES"),
		os.Getenv("FB_OAUTH_REDIRECT_URI"),
		fmt.Sprintf("security_token=%v&ui_state=%v", state, uiState),
	)
	fbState := domain.FBState{
		BaseEntity: domain.BaseEntity{
			ID: primitive.NewObjectID(), InsertedAt: time.Now().UTC(), LastUpdate: time.Now().UTC()},
		State: fAuth.GetState(),
	}
	repository.SaveFBState(&fbState)
	return fAuth.BuildURI(), nil
}

// Todo : Refactor this to proper golang coding pattern

func (s *FbLoginService) HandleFbRedirect(values url.Values) string {
	receivedState := values["state"][0]
	code := values["code"][0]
	if expected, err := repository.GetByFBState(receivedState); err == nil {
		s.logger.Debug(fmt.Sprintf("%v %v", receivedState, expected.State))
		if receivedState == expected.State {
			// get user info from google
			userMap := s.exchangeTokenWithFb(code)
			s.logger.Info(fmt.Sprintf("\nfb data = %v\n", userMap))
			/*
					fbId := userMap["sub"]
					fmt.Println(fbId)
					existingUser, _ := repository.GetUserByGoogleId(fmt.Sprintf("%v", fbId))
					if existingUser == nil {
						//create new user
						existingUser = &domain.User{
							BaseEntity:             domain.BaseEntity{Id: primitive.NewObjectID(),
				InsertedAt: time.Now().UTC(), LastUpdate: time.Now().UTC()},
							FirstName:              userMap["given_name"].(string),
							LastName:               userMap["family_name"].(string),
							Email:           userMap["email"].(string),
							IsPrimaryEmailVerified: userMap["email_verified"].(bool),
							Picture:         userMap["picture"].(string),
							DateCreated:            time.Now(),
							DateUpdated:            time.Now(),
							GoogleAuthInfo: domain.GoogleAuthInfo{
								Id:      userMap["sub"].(string),
								Email:   userMap["email"].(string),
								Picture: userMap["picture"].(string),
							},
						}
						if existingUser.IsPrimaryEmailVerified {
							// save user
							repository.SaveUser(existingUser)
						} else {
							// raise panic
						}
					} else {
						// create exchange code
					}
					// create or update user
					// give exchange code to the user
					splits := strings.Split(receivedState, "&")
					return splits[1]
			*/
		} else {
			fmt.Println(err)
		}
	}

	return ""

}

func (s *FbLoginService) exchangeTokenWithFb(code string) map[string]interface{} {
	accessTokenURL := fmt.Sprintf(
		"https://graph.facebook.com/v8.0/oauth/access_token?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s",
		os.Getenv("FB_OAUTH_CLIENT_ID"), os.Getenv("FB_OAUTH_REDIRECT_URI"), os.Getenv("FB_OAUTH_CLIENT_SECRET"), code)
	var response map[string]interface{}
	accessTokenBytes, atErr := s.httpClient.SendWithContext(http.MethodGet, accessTokenURL, nil)
	if atErr != nil {
		// todo : handle error in go way
		s.logger.Error("could not decode response from facebook", zap.Error(atErr))
	}
	respDecodeErr := json.Unmarshal(accessTokenBytes, &response)
	if respDecodeErr != nil {
		// todo: Refactor it properly instead of just logging
		s.logger.Error("could not decode response from facebook", zap.Error(respDecodeErr))
	}
	accessToken := response["access_token"]
	fmt.Printf("facebook data = %v", response)

	userInfoURL := fmt.Sprintf(
		"https://graph.facebook.com/me?fields=id,first_name,last_name,picture,email&access_token=%s",
		accessToken)

	userDataBytes, uError := s.httpClient.SendWithContext(http.MethodGet, userInfoURL, nil)
	if uError != nil {
		// todo: Refactor it properly instead of just logging
		s.logger.Error("could not decode response from facebook", zap.Error(uError))
	}
	fmt.Printf("Error = %v", uError)
	var userData map[string]interface{}
	userDataUmErr := json.Unmarshal(userDataBytes, &userData)
	if userDataUmErr != nil {
		// todo: Refactor it properly instead of just logging
		s.logger.Error("could not decode response from facebook", zap.Error(userDataUmErr))
	}
	fmt.Printf("user data = %v", userData)
	return userData
}
