package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"timelyship.com/accounts/service"
)

func LoginByFB(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	uiState := queryParams["r"][0]
	redirectUri, _ := service.GetFBRedirectUri(uiState)
	c.JSON(http.StatusOK, redirectUri)
	c.Abort()
}

func HandleRedirectFromDB(c *gin.Context) {
	fmt.Println("Login redirect log...")
	queryParams := c.Request.URL.Query()
	redirectUri := service.HandleFbRedirect(queryParams)
	c.Redirect(http.StatusTemporaryRedirect, redirectUri)
	//c.Abort()
}
