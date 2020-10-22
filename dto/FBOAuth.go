package dto

import (
	"fmt"
	"net/url"
)

type FBOAuth struct {
	responseType string
	clientId     string
	scopes       string
	redirectUri  string
	state        string
}

func NewFBOAuth(responseType, clientId, scopes, redirectUri, state string) *FBOAuth {
	return &FBOAuth{
		responseType: responseType,
		clientId:     clientId,
		scopes:       scopes,
		redirectUri:  redirectUri,
		state:        state,
	}
}

func (g *FBOAuth) BuildUri() string {
	scopeEncoded := url.QueryEscape(g.scopes)
	fmt.Printf("b4=%v,af=%v\n", g.scopes, scopeEncoded)
	fmt.Println("scopeEncoded = ", scopeEncoded)
	stateEncoded := url.QueryEscape(g.state)
	fmt.Println("stateEncoded = ", stateEncoded)

	uri := "https://www.facebook.com/v8.0/dialog/oauth?client_id=%s&redirect_uri=%s&state=%s&response_type=code&scope=%s&display=popup"

	return fmt.Sprintf(uri, g.clientId, g.redirectUri, stateEncoded, scopeEncoded)
}

func (g *FBOAuth) GetState() string {
	return g.state
}
