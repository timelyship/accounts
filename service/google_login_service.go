package service

import (
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/dto"
	"timelyship.com/accounts/repository"
)

type GoogleLoginService struct {
	accountRepository     repository.AccountRepository
	googleLoginRepository repository.GoogleLoginRepository
	logger                zap.Logger
}

func ProvideGoogleLoginService(
	a repository.AccountRepository, g repository.GoogleLoginRepository, l zap.Logger) GoogleLoginService {
	return GoogleLoginService{
		accountRepository:     a,
		googleLoginRepository: g,
		logger:                l,
	}
}

func (s *GoogleLoginService) GetGoogleRedirectURI(uiState string) (string, error) {
	state, eUUID := uuid.NewRandom()
	nonce, eNonce := uuid.NewRandom()
	if eUUID == nil && eNonce == nil {
		scopes := os.Getenv("GOOGLE_OAUTH_SCOPES")
		s.logger.Debug("uiState", zap.String("uiState", uiState))
		s.logger.Info(uiState)
		gAuth := dto.NewGoogleOpenIDAuth("code",
			os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
			scopes,
			os.Getenv("GOOGLE_OAUTH_REDIRECT_URI"),
			fmt.Sprintf("security_token=%v&ui_state=%v", state, uiState),
			nonce.String(),
			"",
			"",
		)
		googleState := domain.GoogleState{
			BaseEntity: domain.BaseEntity{
				ID: primitive.NewObjectID(), InsertedAt: time.Now().UTC(), LastUpdate: time.Now().UTC()},
			State: gAuth.GetState(),
		}
		s.googleLoginRepository.SaveGoogleState(&googleState)
		return gAuth.BuildURI(), nil
	}
	s.logger.Error("UUID failed")
	return "", nil
}

func (s *GoogleLoginService) HandleGoogleRedirect(values url.Values) string {
	receivedState := values["state"][0]
	code := values["code"][0]
	expected, err := s.googleLoginRepository.GetByGoogleState(receivedState)
	if err == nil { //nolint:nestif
		s.logger.Info(fmt.Sprintf("%v %v", receivedState, expected.State))
		if receivedState == expected.State {
			// get user info from google
			userMap := s.exchangeCode(code)
			googleID := userMap["sub"]
			s.logger.Debug("googleID", zap.Any("googleID", googleID))
			existingUser, _ := repository.GetUserByGoogleID(fmt.Sprintf("%v", googleID))
			if existingUser == nil {
				// create new user
				existingUser = &domain.User{
					BaseEntity: domain.BaseEntity{
						ID: primitive.NewObjectID(), InsertedAt: time.Now().UTC(), LastUpdate: time.Now().UTC()},
					FirstName:       userMap["given_name"].(string),
					LastName:        userMap["family_name"].(string),
					Email:           userMap["email"].(string),
					IsEmailVerified: userMap["email_verified"].(bool),
					Picture:         userMap["picture"].(string),
					DateCreated:     time.Now(),
					DateUpdated:     time.Now(),
					GoogleAuthInfo: domain.GoogleAuthInfo{
						ID:      userMap["sub"].(string),
						Email:   userMap["email"].(string),
						Picture: userMap["picture"].(string),
					},
				}
				if existingUser.IsEmailVerified {
					// save user
					s.accountRepository.SaveUser(existingUser)
				} else {
					s.logger.Error("possible error")
				}
			}
			// create or update user
			// give exchange token to the user
			splits := strings.Split(receivedState, "&")
			return splits[1]
		}
		// s.logger.Error("err", zap.Error(err.Error))
	}

	return ""

}

func (s *GoogleLoginService) exchangeCode(code string) map[string]interface{} {
	log.Println(code)
	/*
		data := url.Values{
			"code":          {code},
			"client_id":     {os.Getenv("GOOGLE_OAUTH_CLIENT_ID")},
			"client_secret": {os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")},
			"redirect_uri":  {os.Getenv("GOOGLE_OAUTH_REDIRECT_URI")},
			"grant_type":    {"authorization_code"},
		}
		s.logger.Debug(fmt.Sprintf("\ndebug %v\n", data))

		resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)

		if err != nil {
			panic(err)
		}

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		idToken := res["id_token"]
		user := s.extractToken(idToken.(string))
		s.logger.Info(fmt.Sprintf("\nUSER :\n %v", user))
		return user
	*/
	return nil
}

/*
func (s *GoogleLoginService) extractToken(token string) map[string]interface{} {
	splits := strings.Split(token, ".")
	userJSONStrBytes, _ := base64.StdEncoding.DecodeString(splits[1])
	var result map[string]interface{}
	json.Unmarshal(userJSONStrBytes, &result) //nolint:errcheck
	return result
}
*/
