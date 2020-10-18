package dto

import (
	"fmt"
	"net/url"
	"strings"
)

type GoogleOpenIdAuth struct {
	responseType string
	clientId     string
	scopes       []string
	redirectUri  string
	state        string //=security_token%3D138r5719ru3e1%26url%3Dhttps%3A%2F%2Foauth2-login-demo.example.com%2FmyHome&
	nonce        string //=0394852-3190485-2490358&
	loginHint    string
	hd           string //= example.com
}

func NewGoogleOpenIdAuth(responseType, clientId string, scopes []string, redirectUri, state, nonce, hd, loginHint string) *GoogleOpenIdAuth {
	return &GoogleOpenIdAuth{
		responseType: responseType,
		clientId:     clientId,
		scopes:       scopes,
		redirectUri:  redirectUri,
		state:        state,
		nonce:        nonce,
		hd:           hd,
		loginHint:    loginHint,
	}
}

func (g *GoogleOpenIdAuth) BuildUri() string {
	scopeEncoded := url.QueryEscape(strings.Join(g.scopes, " "))
	stateEncoded := url.QueryEscape(g.state)
	uri := "https://accounts.google.com/o/oauth2/v2/auth?" +
		"response_type=code&client_id=%s&" +
		"scope=%s&redirect_uri=%s&state=%s&nonce=%s"
	return fmt.Sprintf(uri, g.clientId, scopeEncoded, g.redirectUri, stateEncoded, g.nonce)
}
