package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"timelyship.com/accounts/application"
	"timelyship.com/accounts/dto/request"
	"timelyship.com/accounts/service"
	"timelyship.com/accounts/utility"
)

func SignUp(c *gin.Context) {
	var signUpRequest request.SignUpRequest
	if jsonBindingError := c.ShouldBindJSON(&signUpRequest); jsonBindingError != nil {
		restErr := utility.NewBadRequestError("Invalid JSON body", &jsonBindingError)
		c.JSON(restErr.Status, restErr)
		return
	}
	application.Logger.Debug("signUpRequest debug log", zap.Any("signUpRequest = ", fmt.Sprintf("%v", signUpRequest)))
	err := service.InitiateSignUp(signUpRequest)
	if err != nil {
		c.JSON(err.Status, err)
	} else {
		c.JSON(http.StatusCreated, nil)
	}
}

func VerifyEmail(c *gin.Context) {
	verificationToken := c.Param("verificationToken")
	err := service.VerifyEmail(verificationToken)
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
