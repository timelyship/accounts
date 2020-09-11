package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterUser(c *gin.Context) {
	c.String(http.StatusOK, "Pong")
}
