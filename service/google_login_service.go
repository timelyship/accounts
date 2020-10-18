package service

import (
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if expected, err := repository.GetByGoogleState(receivedState); err == nil {
		fmt.Printf("%v %v", receivedState, expected.State)
		if receivedState == expected.State {
			splits := strings.Split(receivedState, "&")
			return splits[1]
		} else {
			fmt.Println(err)
		}
	}

	return ""

}
