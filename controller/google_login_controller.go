package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"timelyship.com/accounts/application"
	"timelyship.com/accounts/appwiring"
)

func LoginByGoogle(c *gin.Context) {
	logger := application.NewTraceableLogger(c.Get("logger"))
	googleAuthService := appwiring.InitGoogleLoginService(*logger)
	queryParams := c.Request.URL.Query()
	uiState := queryParams["r"][0]
	redirectURI, _ := googleAuthService.GetGoogleRedirectURI(uiState)
	c.JSON(http.StatusOK, redirectURI)
	c.Abort()
}

func HandleRedirectFromGoogle(c *gin.Context) {
	logger := application.NewTraceableLogger(c.Get("logger"))
	googleAuthService := appwiring.InitGoogleLoginService(*logger)
	logger.Info("Login redirect log...")
	queryParams := c.Request.URL.Query()
	redirectURI := googleAuthService.HandleGoogleRedirect(queryParams)
	c.Redirect(http.StatusTemporaryRedirect, redirectURI)
	// c.Abort()
}
