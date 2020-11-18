package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
		c.JSON(200, loginResponse)
	}
}

func InitiateLogin(c *gin.Context) {
	resp, err := service.InitiateLogin()
	if err != nil {
		c.JSON(err.Status, err)
	} else {
		c.JSON(200, resp)
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
		c.JSON(200, loginResponse)
	}

}

func GenerateCode(c *gin.Context) {
	token, ok := c.MustGet("token").(*jwt.Token)
	aud := c.Query("aud")
	state := c.Query("state")
	if !ok {
		c.JSON(401, nil)
	} else {
		code, err := service.GenerateCode(token, aud, state)
		if err != nil {
			c.JSON(err.Status, err)
		} else {
			c.JSON(200, map[string]string{
				"code": code,
			})
		}
	}

}
