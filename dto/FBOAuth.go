package dto

import (
	"fmt"
	"net/url"
)

type FBOAuth struct {
	responseType string
	clientID     string
	scopes       string
	redirectURI  string
	state        string
}

func NewFBOAuth(responseType, clientID, scopes, redirectURI, state string) *FBOAuth {
	return &FBOAuth{
		responseType: responseType,
		clientID:     clientID,
		scopes:       scopes,
		redirectURI:  redirectURI,
		state:        state,
	}
}

func (g *FBOAuth) BuildURI() string {
	scopeEncoded := url.QueryEscape(g.scopes)
	fmt.Printf("b4=%v,af=%v\n", g.scopes, scopeEncoded)
	fmt.Println("scopeEncoded = ", scopeEncoded)
	stateEncoded := url.QueryEscape(g.state)
	fmt.Println("stateEncoded = ", stateEncoded)

	uri := "https://www.facebook.com/v8.0/dialog/oauth?" +
		"client_id=%s&redirect_uri=%s&state=%s&response_type=code&scope=%s&display=popup"

	return fmt.Sprintf(uri, g.clientID, g.redirectURI, stateEncoded, scopeEncoded)
}

func (g *FBOAuth) GetState() string {
	return g.state
}
