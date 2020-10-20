package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/dto"
	"timelyship.com/accounts/repository"
)

func GetGoogleRedirectUri(uiState string) (string, error) {
	state, eUuid := uuid.NewRandom()
	nonce, eNonce := uuid.NewRandom()
	if eUuid == nil && eNonce == nil {
		scopes := os.Getenv("GOOGLE_OAUTH_SCOPES")
		fmt.Println(uiState)
		gAuth := dto.NewGoogleOpenIdAuth("code",
			os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
			scopes,
			os.Getenv("GOOGLE_OAUTH_REDIRECT_URI"),
			fmt.Sprintf("security_token=%v&ui_state=%v", state, uiState),
			nonce.String(),
			"",
			"",
		)
		googleState := domain.GoogleState{
			BaseEntity: domain.BaseEntity{Id: primitive.NewObjectID(), InsertedAt: time.Now().UTC(), LastUpdate: time.Now().UTC()},
			State:      gAuth.GetState(),
		}
		repository.SaveGoogleState(&googleState)
		return gAuth.BuildUri(), nil
	} else {
		fmt.Println("UUID failed")
	}
	return "", nil
}

func HandleGoogleRedirect(values url.Values) string {
	receivedState := values["state"][0]
	code := values["code"][0]
	if expected, err := repository.GetByGoogleState(receivedState); err == nil {
		fmt.Printf("%v %v", receivedState, expected.State)
		if receivedState == expected.State {
			// get user info from google
			userMap := exchangeCode(code)
			googleId := userMap["sub"]
			exstingGoogleUser := repository.GetUserByGoogleId(googleId)
			// create or update user
			// give exchange token to the user
			splits := strings.Split(receivedState, "&")
			return splits[1]
		} else {
			fmt.Println(err)
		}
	}

	return ""

}

func exchangeCode(code string) map[string]interface{} {
	data := url.Values{
		"code":          {code},
		"client_id":     {os.Getenv("GOOGLE_OAUTH_CLIENT_ID")},
		"client_secret": {os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")},
		"redirect_uri":  {os.Getenv("GOOGLE_OAUTH_REDIRECT_URI")},
		"grant_type":    {"authorization_code"},
	}
	fmt.Printf("\ndebug %v\n", data)

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)

	if err != nil {
		panic(err)
	}

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)
	idToken := res["id_token"]
	user := extractToken(idToken.(string))
	fmt.Printf("\nUSER :\n %v", user)
	return user
}

func extractToken(token string) map[string]interface{} {
	splits := strings.Split(token, ".")
	userJsonStrBytes, _ := base64.StdEncoding.DecodeString(splits[1])
	var result map[string]interface{}
	json.Unmarshal(userJsonStrBytes, result)
	return result
}
