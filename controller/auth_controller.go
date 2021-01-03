package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"timelyship.com/accounts/dto/request"
	"timelyship.com/accounts/dto/response"
	"timelyship.com/accounts/service"
	"timelyship.com/accounts/utility"
)

func Login(c *gin.Context) {
	var loginRequest request.LoginRequest
	if jsonBindingError := c.ShouldBindJSON(&loginRequest); jsonBindingError != nil {
		restErr := utility.NewBadRequestError("Invalid JSON body", &jsonBindingError)
		c.JSON(restErr.Status, restErr)
		return
	}
	loginResponse, err := service.HandleLogin(loginRequest)
	if err != nil {
		c.JSON(err.Status, err)
	} else {
		c.JSON(http.StatusOK, loginResponse)
	}
}

func InitiateLogin(c *gin.Context) {
	resp, err := service.InitiateLogin()
	if err != nil {
		c.JSON(err.Status, err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func RefreshToken(c *gin.Context) {
	var refreshTokenRequest response.LoginResponse
	if jsonBindingError := c.ShouldBindJSON(&refreshTokenRequest); jsonBindingError != nil {
		restErr := utility.NewBadRequestError("Invalid JSON body", &jsonBindingError)
		c.JSON(restErr.Status, restErr)
		return
	}
	loginResponse, err := service.RefreshToken(refreshTokenRequest.AccessToken, refreshTokenRequest.RefreshToken)
	if err != nil {
		c.JSON(err.Status, err)
	} else {
		c.JSON(http.StatusOK, loginResponse)
	}
}

func GenerateCode(c *gin.Context) {
	token, ok := c.MustGet("token").(*jwt.Token)
	aud := c.Query("aud")
	state := c.Query("state")
	if !ok {
		c.JSON(http.StatusUnauthorized, nil)
	} else {
		err := service.GenerateCode(token, aud, state)
		if err != nil {
			c.JSON(err.Status, err)
		} else {
			c.JSON(http.StatusOK, nil)
		}
	}

}

func ExchangeCode(c *gin.Context) {
	state := c.Query("state")
	data, err := service.ExchangeCode(state)
	if err != nil {
		c.JSON(err.Status, err)
	} else {
		c.JSON(http.StatusOK, map[string]string{
			"code": data.Code,
		})
	}
}

func Profile(c *gin.Context) {
	token, ok := c.MustGet("token").(*jwt.Token)
	if !ok {
		c.JSON(http.StatusUnauthorized, "token not ok")
		return
	}
	claims, err := utility.GetProfileClaims(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	c.JSON(http.StatusOK, claims)
}
