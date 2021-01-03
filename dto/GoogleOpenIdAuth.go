package dto

import (
	"fmt"
	"net/url"
)

type GoogleOpenIDAuth struct {
	responseType string
	clientID     string
	scopes       string
	redirectURI  string
	state        string //=security_token%3D138r5719ru3e1%26url%3Dhttps%3A%2F%2Foauth2-login-demo.example.com%2FmyHome&
	nonce        string //=0394852-3190485-2490358&
	loginHint    string
	hd           string //= example.com
}

func NewGoogleOpenIDAuth(responseType, clientId string, scopes string, redirectUri, state, nonce, hd, loginHint string) *GoogleOpenIDAuth {
	return &GoogleOpenIDAuth{
		responseType: responseType,
		clientID:     clientId,
		scopes:       scopes,
		redirectURI:  redirectUri,
		state:        state,
		nonce:        nonce,
		hd:           hd,
		loginHint:    loginHint,
	}
}

func (g *GoogleOpenIDAuth) BuildURI() string {
	scopeEncoded := url.QueryEscape(g.scopes)
	fmt.Printf("b4=%v,af=%v\n", g.scopes, scopeEncoded)
	fmt.Println("scopeEncoded = ", scopeEncoded)
	stateEncoded := url.QueryEscape(g.state)
	fmt.Println("stateEncoded = ", stateEncoded)

	uri := "https://accounts.google.com/o/oauth2/v2/auth?" +
		"response_type=code&client_id=%s&" +
		"scope=%s&redirect_uri=%s&state=%s&nonce=%s&prompt=select_account"

	return fmt.Sprintf(uri, g.clientID, scopeEncoded, g.redirectURI, stateEncoded, g.nonce)
}

func (g *GoogleOpenIDAuth) GetState() string {
	return g.state
}
