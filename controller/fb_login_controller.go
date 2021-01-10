package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"timelyship.com/accounts/application"
	"timelyship.com/accounts/appwiring"
	"timelyship.com/accounts/service"
)

func LoginByFB(c *gin.Context) {
	logger := application.NewTraceableLogger(c.Get("logger"))
	httpClient := &service.HTTPClientImpl{}
	fbAuthService := appwiring.InitFbLoginService(*logger, httpClient)
	queryParams := c.Request.URL.Query()
	uiState := queryParams["r"][0]
	redirectURI, _ := fbAuthService.GetFBRedirectURI(uiState)
	c.JSON(http.StatusOK, redirectURI)
	c.Abort()
}

func HandleRedirectFromFB(c *gin.Context) {
	logger := application.NewTraceableLogger(c.Get("logger"))
	httpClient := &service.HTTPClientImpl{}
	fbAuthService := appwiring.InitFbLoginService(*logger, httpClient)
	logger.Info("Login redirect log...")
	queryParams := c.Request.URL.Query()
	redirectURI := fbAuthService.HandleFbRedirect(queryParams)
	c.Redirect(http.StatusTemporaryRedirect, redirectURI)
	// c.Abort()
}
