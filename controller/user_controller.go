package controller

//
//import (
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"timelyship.com/accounts/domain"
//	"timelyship.com/accounts/service"
//	"timelyship.com/accounts/utility"
//)
//
//func RegisterUser(c *gin.Context) {
//	var user domain.User
//	if err := c.ShouldBindJSON(&user); err != nil {
//		restErr := utility.NewBadRequestError("Invalid JSON Body")
//		c.JSON(restErr.Status, restErr)
//		return
//	}
//	validationError := domain.ValidateUser(&user)
//	if validationError != nil {
//		c.JSON(validationError.Status, validationError)
//		return
//	}
//	result, restErr := service.CreateUser(user)
//	if restErr != nil {
//		c.JSON(restErr.Status, restErr)
//		return
//	}
//	fmt.Print(user)
//	c.JSON(http.StatusCreated, result)
//}
