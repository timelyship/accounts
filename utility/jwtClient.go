package utility

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"strconv"
	"time"
	"timelyship.com/accounts/domain"
)

// todo - refactor this one.
func CreateToken(user *domain.User, aud string) (*domain.TokenDetails, error) {
	td := &domain.TokenDetails{
		BaseEntity: domain.BaseEntity{Id: primitive.NewObjectID(), InsertedAt: time.Now().UTC(), LastUpdate: time.Now().UTC()},
	}
	accessTokenExpInMinute, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP_MINUTE"))
	refreshTokenExpInMinute, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP_MINUTE"))
	td.AtExpires = time.Now().Add(time.Minute * time.Duration(accessTokenExpInMinute)).Unix()
	td.AccessUuid = uuid.New().String()

	td.RtExpires = time.Now().Add(time.Hour * time.Duration(refreshTokenExpInMinute)).Unix()
	td.RefreshUuid = uuid.New().String()

	var err error
	//Creating Access Token
	//iss, sub string, aud string, exp, nbf, iat int64, jti string, typ string
	atClaims := mapClaims(os.Getenv("TOKEN_ISSUER"), user.Id.Hex(), aud,
		td.AtExpires, time.Now().Unix(), time.Now().Unix(), td.AccessUuid, "jwt")
	addProfileClaims(&atClaims, user)
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

func CreateAccessToken(user *domain.User, aud string) (*domain.TokenDetails, error) {
	td := &domain.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.New().String()

	var err error
	//Creating Access Token
	//iss, sub string, aud string, exp, nbf, iat int64, jti string, typ string
	atClaims := mapClaims(os.Getenv("TOKEN_ISSUER"), user.Id.Hex(), aud,
		td.AtExpires, time.Now().Unix(), time.Now().Unix(), td.AccessUuid, "jwt")
	addProfileClaims(&atClaims, user)
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}
func addProfileClaims(claims *jwt.MapClaims, user *domain.User) {
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

func DecodeToken(jwtTokenRaw string, secret string) (*jwt.MapClaims, *RestError) {
	token, err := jwt.Parse(jwtTokenRaw, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		return nil, NewUnAuthorizedError("Unauthorized", &err)
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, NewUnAuthorizedError("Unauthorized", nil)
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return &claims, nil
	}
	return nil, NewUnAuthorizedError("Invalid token", &err)
}

func ValidateToken(encodedToken, secret string) (*jwt.Token, error) {
	tokenValidator := func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("Invalid token %v", token.Header["alg"])
		}
		return []byte(secret), nil
	}
	return jwt.Parse(encodedToken, tokenValidator)

}

func GetProfileClaims(token *jwt.Token) (map[string]interface{}, *RestError) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, NewUnAuthorizedError("Token could not be decoded", nil)
	}
	return map[string]interface{}{
		"firstName": claims["first_name"],
		"lastName":  claims["last_name"],
		"email":      claims["email"],
		"picture":    claims["picture"],
		"roles":      claims["first_name"],
	}, nil
}
