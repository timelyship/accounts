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
	redirectUri, _ := service.GetGoogleRedirectUri(uiState)
	//c.String(200, redirectUri)
	//c.JSON(redirectUri,redirectUri)
	//c.Redirect(http.StatusTemporaryRedirect, redirectUri)
	c.JSON(http.StatusOK, redirectUri)
	c.Abort()
}

func HandleRedirectFromGoogle(c *gin.Context) {
	fmt.Println("Login redirect log...")
	queryParams := c.Request.URL.Query()
	for k, v := range queryParams {
		fmt.Printf("%v %v %T %T\n", k, v, k, v)
	}
	redirectUri := service.HandleGoogleRedirect(queryParams)
	c.Redirect(http.StatusTemporaryRedirect, redirectUri)
	//c.Abort()
}
