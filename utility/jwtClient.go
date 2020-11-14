package utility

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"os"
	"time"
	"timelyship.com/accounts/domain"
)

func CreateToken(user *domain.User, aud string) (*domain.TokenDetails, error) {
	td := &domain.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.New().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.New().String()

	var err error
	//Creating Access Token
	//iss, sub string, aud string, exp, nbf, iat int64, jti string, typ string
	atClaims := mapClaims(os.Getenv("TOKEN_ISSUER"), user.Id.Hex(), aud,
		td.AtExpires, time.Now().Unix(), time.Now().Unix(), td.AccessUuid, "jwt")
	addAdditionalClaims(&atClaims, user)
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	rtClaims := mapClaims(os.Getenv("TOKEN_ISSUER"), user.Id.Hex(), aud,
		td.RtExpires, time.Now().Unix(), time.Now().Unix(), td.RefreshUuid, "jwt")
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func addAdditionalClaims(claims *jwt.MapClaims, user *domain.User) {
	(*claims)["first_name"] = user.FirstName
	(*claims)["last_name"] = user.LastName
	(*claims)["email"] = user.PrimaryEmail
	(*claims)["picture"] = FirstNotNullString(user.PrimaryPicture, user.FacebookAuthInfo.Picture, user.GoogleAuthInfo.Picture) //primary_picture
	(*claims)["roles"] = user.Roles

}

//https://en.wikipedia.org/wiki/JSON_Web_Token
func mapClaims(iss, sub, aud string, exp, nbf, iat int64, jti, typ string) jwt.MapClaims {
	atClaims := jwt.MapClaims{}
	atClaims["iss"] = iss
	atClaims["sub"] = sub
	atClaims["aud"] = aud
	atClaims["exp"] = exp
	atClaims["nbf"] = nbf
	atClaims["iat"] = iat
	atClaims["jti"] = jti
	atClaims["typ"] = typ
	return atClaims
}
