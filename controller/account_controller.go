package controller

import (
	"github.com/gin-gonic/gin"
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
	validationError := signUpRequest.Validate()
	if validationError != nil {
		c.JSON(validationError.Status, validationError)
		return
	}
	err := service.InitiateSignUp(signUpRequest)
}

func VerifyEmail() {

}

func ChangePassword() {

}

// sends a password reset email to user
func ForgotPassword() {

}

// resets password
func ResetPassword() {

}
