package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"timelyship.com/accounts/application"
	"timelyship.com/accounts/appwiring"
	"timelyship.com/accounts/dto/request"
	"timelyship.com/accounts/utility"
)

func SignUp(c *gin.Context) {
	logger := application.NewTraceableLogger(c.Request.Context().Value("logger"))
	accountService := appwiring.InitAccountService(logger)
	var signUpRequest request.SignUpRequest
	if jsonBindingError := c.ShouldBindJSON(&signUpRequest); jsonBindingError != nil {
		restErr := utility.NewBadRequestError("Invalid JSON body", &jsonBindingError)
		c.JSON(restErr.Status, restErr)
		return
	}
	logger.Debug("signUpRequest debug log", zap.Any("signUpRequest = ", fmt.Sprintf("%v", signUpRequest)))
	err := accountService.InitiateSignUp(signUpRequest)
	if err != nil {
		c.JSON(err.Status, err)
	} else {
		c.JSON(http.StatusCreated, nil)
	}
}

func VerifyEmail(c *gin.Context) {
	logger := application.NewTraceableLogger(c.Request.Context().Value("logger"))
	accountService := appwiring.InitAccountService(logger)
	verificationToken := c.Param("verificationToken")
	err := accountService.VerifyEmail(verificationToken)
	if err != nil {
		c.JSON(err.Status, err)
	} else {
		c.JSON(http.StatusOK, nil)
	}
}

func ChangePassword() {

}

// sends a password reset email to user
func ForgotPassword() {

}

// resets password
func ResetPassword() {

}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
