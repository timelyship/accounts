package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"timelyship.com/accounts/service"
)

func LoginByGoogle(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	uiState := queryParams["r"][0]
	redirectURI, _ := service.GetGoogleRedirectUri(uiState)
	c.JSON(http.StatusOK, redirectURI)
	c.Abort()
}

func HandleRedirectFromGoogle(c *gin.Context) {
	fmt.Println("Login redirect log...")
	queryParams := c.Request.URL.Query()
	redirectURI := service.HandleGoogleRedirect(queryParams)
	c.Redirect(http.StatusTemporaryRedirect, redirectURI)
	// c.Abort()
}
