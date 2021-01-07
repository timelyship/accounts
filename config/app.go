package config

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"timelyship.com/accounts/application"
	"timelyship.com/accounts/utility"
)

var (
	router = gin.New()
)

func LogInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		var traceId, spanId string
		if traceId = c.GetHeader("ts-trace-id"); traceId == "" {
			traceId = strings.ReplaceAll(uuid.New().String(), "-", "")
		}
		if spanId = c.GetHeader("ts-trace-id"); spanId == "" {
			spanId = strings.ReplaceAll(uuid.New().String(), "-", "")
		}

		c.Set("logger", application.NewLogger(traceId, spanId))
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
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, "+
			"X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)

			return
		}
		c.Next()
	}
}

func isWhiteListed(uri string) bool {
	whiteListedUrls := []string{"/account/login", "/account/sign-up",
		"/initiate-login", "/decode-code", "/exchange-code", "/logout", "/verify-email"}
	for _, a := range whiteListedUrls {
		if strings.HasPrefix(uri, a) {
			return true
		}
	}
	return false
}

func AuthenticationMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isWhiteListed(c.Request.RequestURI) {
			c.Next()
			return
		}
		const BearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, BearerSchema) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString := authHeader[len(BearerSchema):]
		principal, err := utility.ExtractPrincipalFromToken(tokenString, os.Getenv("ACCESS_SECRET"))
		if err == nil {
			logger := application.NewTraceableLogger(c.Get("logger"))
			c.Set("principal", principal)
			c.Set("logger", logger.With(zap.String("user-id", principal.UserID)))
			// todo : set user id in the logger here.
			c.Next()

		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func Start() {
	router.Use(LogInterceptor()) // creates a logger , generates traceID,spanId
	router.Use(CORSMiddleware())
	router.Use(AuthenticationMiddleWare()) // decodes user, from token, should be the first one, populates userID
	mapUrls()
	err := router.Run(":8080")
	if err != nil {
		zap.L().Error("Error starting app", zap.Error(err))
	}
}
