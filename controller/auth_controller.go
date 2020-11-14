package controller

import (
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

func InitiateLogin() {

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
