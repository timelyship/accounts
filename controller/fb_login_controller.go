package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"timelyship.com/accounts/application"
	"timelyship.com/accounts/appwiring"
)

func LoginByFB(c *gin.Context) {
	logger := application.NewTraceableLogger(c.Request.Context().Value("logger"))
	fbAuthService := appwiring.InitFbLoginService(logger)

	queryParams := c.Request.URL.Query()
	uiState := queryParams["r"][0]
	redirectURI, _ := fbAuthService.GetFBRedirectURI(uiState)
	c.JSON(http.StatusOK, redirectURI)
	c.Abort()
}

func HandleRedirectFromDB(c *gin.Context) {
	logger := application.NewTraceableLogger(c.Request.Context().Value("logger"))
	fbAuthService := appwiring.InitFbLoginService(logger)
	fmt.Println("Login redirect log...")
	queryParams := c.Request.URL.Query()
	redirectURI := fbAuthService.HandleFbRedirect(queryParams)
	c.Redirect(http.StatusTemporaryRedirect, redirectURI)
	// c.Abort()
}
