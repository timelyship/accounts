package service

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"timelyship.com/accounts/dto"
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

		return gAuth.BuildUri(), nil
	} else {
		fmt.Println("UUID failed")
	}
	return "", nil
}
