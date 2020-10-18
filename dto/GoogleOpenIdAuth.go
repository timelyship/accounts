package dto

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/url"
	"time"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/repository"
)

type GoogleOpenIdAuth struct {
	responseType string
	clientId     string
	scopes       string
	redirectUri  string
	state        string //=security_token%3D138r5719ru3e1%26url%3Dhttps%3A%2F%2Foauth2-login-demo.example.com%2FmyHome&
	nonce        string //=0394852-3190485-2490358&
	loginHint    string
	hd           string //= example.com
}

func NewGoogleOpenIdAuth(responseType, clientId string, scopes string, redirectUri, state, nonce, hd, loginHint string) *GoogleOpenIdAuth {
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
	scopeEncoded := url.QueryEscape(g.scopes)
	fmt.Printf("b4=%v,af=%v\n", g.scopes, scopeEncoded)
	fmt.Println("scopeEncoded = ", scopeEncoded)
	stateEncoded := url.QueryEscape(g.state)
	fmt.Println("stateEncoded = ", stateEncoded)

	uri := "https://accounts.google.com/o/oauth2/v2/auth?" +
		"response_type=code&client_id=%s&" +
		"scope=%s&redirect_uri=%s&state=%s&nonce=%s&prompt=select_account"

	googleState := domain.GoogleState{
		BaseEntity: domain.BaseEntity{Id: primitive.NewObjectID(), InsertedAt: time.Now().UTC(), LastUpdate: time.Now().UTC()},
		State:      g.state,
	}
	repository.SaveGoogleState(&googleState)

	return fmt.Sprintf(uri, g.clientId, scopeEncoded, g.redirectUri, stateEncoded, g.nonce)
}
