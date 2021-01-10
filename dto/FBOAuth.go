package dto

import (
	"fmt"
	"log"
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
	log.Printf("b4=%v,af=%v\n", g.scopes, scopeEncoded)
	log.Printf("scopeEncoded = %s", scopeEncoded)
	stateEncoded := url.QueryEscape(g.state)
	log.Printf("stateEncoded = %s", stateEncoded)

	uri := "https://www.facebook.com/v8.0/dialog/oauth?" +
		"client_id=%s&redirect_uri=%s&state=%s&response_type=code&scope=%s&display=popup"

	return fmt.Sprintf(uri, g.clientID, g.redirectURI, stateEncoded, scopeEncoded)
}

func (g *FBOAuth) GetState() string {
	return g.state
}
