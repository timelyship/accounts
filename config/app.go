package config

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"timelyship.com/accounts/utility"
)

var (
	router = gin.New()
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// before request
		c.Next()
		// after request
		latency := time.Since(t)
		// access the status we are sending
		status := c.Writer.Status()
		log.Println(latency, status)
	}
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func isWhiteListed(uri string)bool  {
	whiteListedUrls := []string{"/account/login","/account/sign-up","/initiate-login","/decode-code","/exchange-code"}
	for _, a := range whiteListedUrls {
		if strings.HasPrefix(uri,a) {
			return true
		}
	}
	return false
}

func AuthenticationMiddleWare() gin.HandlerFunc {

	return func(c *gin.Context) {
		if isWhiteListed(c.Request.RequestURI){
			c.Next()
			return
		}
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, BEARER_SCHEMA) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString := authHeader[len(BEARER_SCHEMA):]
		token, err := utility.ValidateToken(tokenString, os.Getenv("ACCESS_SECRET"))
		if err == nil && token.Valid {
			c.Set("token",token)
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func Start() {
	router.Use(Logger())
	router.Use(CORSMiddleware())
	router.Use(AuthenticationMiddleWare())
	mapUrls()
	router.Run(":8080")
}
