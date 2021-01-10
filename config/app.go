package config

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"timelyship.com/accounts/application"
	"timelyship.com/accounts/repository"
	"timelyship.com/accounts/utility"
)

var (
	router = gin.New()
)

func LogInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Entering Log interceptor")
		defer log.Println("Exiting Log interceptor")
		t := time.Now()
		var traceID, spanID string
		if traceID = c.GetHeader("ts-trace-id"); traceID == "" {
			traceID = utility.GetUUIDWithoutDash()
		}
		if spanID = c.GetHeader("ts-trace-id"); spanID == "" {
			spanID = utility.GetUUIDWithoutDash()
		}
		c.Set("logger", application.NewLogger(traceID, spanID))
		c.Writer.Header().Set("traceId", traceID)
		c.Writer.Header().Set("commit-id", os.Getenv("COMMIT_ID"))
		c.Writer.Header().Set("live-since", os.Getenv("LIVE_SINCE"))
		// before request
		c.Next()

		// after request
		latency := time.Since(t).Milliseconds()
		// access the status we are sending
		status := c.Writer.Status()
		logger := application.NewTraceableLogger(c.Get("logger"))
		logger.Info("Endpoint analytics", zap.String("path", c.Request.RequestURI),
			zap.Any("latency", latency), zap.Int("status", status))
	}
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Entering CORS Middleware")
		defer log.Println("Exiting CORS Middleware")

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
	whiteListedUrls := []string{"/ping", "/account/login", "/account/sign-up",
		"/initiate-login", "/decode-code", "/exchange-code", "/logout", "/verify-email"}
	for _, a := range whiteListedUrls {
		if strings.HasPrefix(uri, a) {
			return true
		}
	}
	return false
}

func AuthenticationMiddleWare() gin.HandlerFunc {
	log.Println("Entering Authentication Middleware")
	defer log.Println("Exiting Authentication Middleware")

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
			c.Set("logger", logger.With(zap.String("userID", principal.UserID)))
			c.Next()

		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func Start() {
	repository.InitClient()
	defer repository.DisconnectMongoClient()
	router.Use(LogInterceptor()) // creates a logger , generates traceID,spanId
	router.Use(CORSMiddleware())
	router.Use(AuthenticationMiddleWare()) // decodes user, from token, should be the first one, populates userID
	mapUrls()
	err := router.Run(":8080")
	if err != nil {
		log.Printf("Failed to start app.Error = %v\n", zap.Error(err))
	}
}
