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
}

func ProvideFbLoginService(fbLoginRepository repository.FbLoginRepository, logger zap.Logger) FbLoginService {
	return FbLoginService{
		fbLoginRepository: fbLoginRepository,
		logger:            logger,
	}
}

func (s *FbLoginService) GetFBRedirectURI(uiState string) (string, error) {
	state, eUUID := uuid.NewRandom()
	if eUUID == nil {
		fmt.Println(uiState)
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
	fmt.Println("UUID failed")
	return "", nil
}

func (s *FbLoginService) HandleFbRedirect(values url.Values) string {
	receivedState := values["state"][0]
	code := values["code"][0]
	if expected, err := repository.GetByFBState(receivedState); err == nil {
		fmt.Printf("%v %v", receivedState, expected.State)
		if receivedState == expected.State {
			// get user info from google
			userMap := s.exchangeTokenWithFb(code)
			fmt.Sprintf("\nfb data = %v\n", userMap)
			/*
				fbId := userMap["sub"]
				fmt.Println(fbId)
				existingUser, _ := repository.GetUserByGoogleId(fmt.Sprintf("%v", fbId))
				if existingUser == nil {
					//create new user
					existingUser = &domain.User{
						BaseEntity:             domain.BaseEntity{Id: primitive.NewObjectID(), InsertedAt: time.Now().UTC(), LastUpdate: time.Now().UTC()},
						FirstName:              userMap["given_name"].(string),
						LastName:               userMap["family_name"].(string),
						PrimaryEmail:           userMap["email"].(string),
						IsPrimaryEmailVerified: userMap["email_verified"].(bool),
						PrimaryPicture:         userMap["picture"].(string),
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
	accessTokenURL := fmt.Sprintf("https://graph.facebook.com/v8.0/oauth/access_token?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s",
		os.Getenv("FB_OAUTH_CLIENT_ID"), os.Getenv("FB_OAUTH_REDIRECT_URI"), os.Getenv("FB_OAUTH_CLIENT_SECRET"), code)
	resp, err := http.Get(accessTokenURL)

	if err != nil {
		panic(err)
	}
	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	accessToken := response["access_token"]
	fmt.Printf("facebook data = %v", response)

	//userInfoUrl := fmt.Sprintf("https://graph.facebook.com/debug_token?input_token=%s&access_token=%s|%s",
	//	accessToken, os.Getenv("FB_OAUTH_CLIENT_ID"), os.Getenv("FB_OAUTH_CLIENT_SECRET"))
	userInfoUrl := fmt.Sprintf("https://graph.facebook.com/me?fields=id,first_name,last_name,picture,email&access_token=%s",
		accessToken)

	resp2, uError := http.Get(userInfoUrl)
	fmt.Printf("Error = %v", uError)
	var userData map[string]interface{}
	json.NewDecoder(resp2.Body).Decode(&userData)
	fmt.Printf("user data = %v", userData)
	return userData

	//var res map[string]interface{}
	//json.NewDecoder(resp.Body).Decode(&res)
	//idToken := res["id_token"]
	//user := extractToken(idToken.(string))
	//fmt.Printf("\nUSER :\n %v", user)
	//return user
}

//func extractToken(token string) map[string]interface{} {
//	splits := strings.Split(token, ".")
//	userJsonStrBytes, _ := base64.StdEncoding.DecodeString(splits[1])
//	var result map[string]interface{}
//	json.Unmarshal(userJsonStrBytes, &result)
//	return result
//}
