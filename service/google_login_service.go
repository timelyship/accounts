package service

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"
	"timelyship.com/accounts/dto"
)

func GetGoogleRedirectUri(uiState string) (string, error) {

	state, eUuid := uuid.NewRandom()
	nonce, eNonce := uuid.NewRandom()
	if eUuid == nil && eNonce == nil {
		scopeList := os.Getenv("GOOGLE_OAUTH_SCOPES")
		scopes := strings.Split(scopeList, " ")
		fmt.Println(uiState)
		gAuth := dto.NewGoogleOpenIdAuth("code",
			os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
			scopes,
			"",
			fmt.Sprintf("security_token=%s&ui_state=%s", state.String(), uiState),
			nonce.String(),
			"",
			"",
		)
		return gAuth.BuildUri(), nil
	} else {
		fmt.Println("UUID failed")
	}
	return "", nil
}
