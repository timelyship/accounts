package controller

import (
	"github.com/gin-gonic/gin"
	"timelyship.com/accounts/service"
)

func LoginByGoogle(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	uiState := queryParams["r"][0]
	redirectUri, _ := service.GetGoogleRedirectUri(uiState)
	c.String(200, redirectUri)
	//c.Redirect(http.StatusTemporaryRedirect, "")
	//c.Abort()
}
