package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"timelyship.com/accounts/application"
	"timelyship.com/accounts/appwiring"
	"timelyship.com/accounts/dto"
	"timelyship.com/accounts/dto/request"
	"timelyship.com/accounts/dto/response"
	"timelyship.com/accounts/utility"
)

func Login(c *gin.Context) {
	logger := application.NewTraceableLogger(c.Get("logger"))
	authService := appwiring.InitAuthService(*logger)
	var loginRequest request.LoginRequest
	if jsonBindingError := c.ShouldBindJSON(&loginRequest); jsonBindingError != nil {
		logger.Error("Json bind error, login")
		restErr := utility.NewBadRequestError("Invalid JSON body", &jsonBindingError)
		c.JSON(restErr.Status, restErr)
		return
	}
	loginResponse, err := authService.HandleLogin(loginRequest)
	if err != nil {
		c.JSON(err.Status, err)
	} else {
		c.JSON(http.StatusOK, loginResponse)
	}
}

func InitiateLogin(c *gin.Context) {
	logger := application.NewTraceableLogger(c.Get("logger"))
	authService := appwiring.InitAuthService(*logger)
	resp, err := authService.InitiateLogin()
	if err != nil {
		c.JSON(err.Status, err)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func RefreshToken(c *gin.Context) {
	logger := application.NewTraceableLogger(c.Get("logger"))
	authService := appwiring.InitAuthService(*logger)
	var refreshTokenRequest response.LoginResponse
	if jsonBindingError := c.ShouldBindJSON(&refreshTokenRequest); jsonBindingError != nil {
		restErr := utility.NewBadRequestError("Invalid JSON body", &jsonBindingError)
		c.JSON(restErr.Status, restErr)
		return
	}
	loginResponse, err := authService.RefreshToken(refreshTokenRequest.AccessToken, refreshTokenRequest.RefreshToken)
	if err != nil {
		c.JSON(err.Status, err)
	} else {
		c.JSON(http.StatusOK, loginResponse)
	}
}

func GenerateCode(c *gin.Context) {
	logger := application.NewTraceableLogger(c.Get("logger"))
	authService := appwiring.InitAuthService(*logger)
	token, ok := c.MustGet("token").(*jwt.Token)
	aud := c.Query("aud")
	state := c.Query("state")
	if !ok {
		c.JSON(http.StatusUnauthorized, nil)
	} else {
		err := authService.GenerateCode(token, aud, state)
		if err != nil {
			c.JSON(err.Status, err)
		} else {
			c.JSON(http.StatusOK, nil)
		}
	}

}

func ExchangeCode(c *gin.Context) {
	logger := application.NewTraceableLogger(c.Get("logger"))
	authService := appwiring.InitAuthService(*logger)
	state := c.Query("state")
	data, err := authService.ExchangeCode(state)
	if err != nil {
		c.JSON(err.Status, err)
	} else {
		c.JSON(http.StatusOK, map[string]string{
			"code": data.Code,
		})
	}
}

func Profile(c *gin.Context) {
	principal, ok := c.MustGet("principal").(*dto.Principal)
	if !ok {
		c.JSON(http.StatusUnauthorized, "principal not ok")
		return
	}
	c.JSON(http.StatusOK, principal)
}
