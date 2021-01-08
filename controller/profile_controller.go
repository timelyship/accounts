package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"timelyship.com/accounts/application"
	"timelyship.com/accounts/appwiring"
	"timelyship.com/accounts/dto"
	"timelyship.com/accounts/dto/request"
	"timelyship.com/accounts/utility"
)

func Profile(c *gin.Context) {
	principal, ok := c.MustGet("principal").(*dto.Principal)
	if !ok {
		c.JSON(http.StatusUnauthorized, "principal not ok")
		return
	}
	logger := application.NewTraceableLogger(c.Get("logger"))
	profileService := appwiring.InitProfileService(*logger)
	profile, err := profileService.GetProfileById(principal.UserID)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	logger.Debug("Profile debug", zap.Any("profile", profile))
	c.JSON(http.StatusOK, profile)
}

func PatchProfile(c *gin.Context) {
	logger := application.NewTraceableLogger(c.Get("logger"))
	profileService := appwiring.InitProfileService(*logger)
	principal, ok := c.MustGet("principal").(*dto.Principal)
	if !ok {
		c.JSON(http.StatusUnauthorized, "principal not ok")
		return
	}
	var profilePatchRequest []*request.ProfilePatchRequest
	if jsonBindingError := c.ShouldBindJSON(&profilePatchRequest); jsonBindingError != nil {
		restErr := utility.NewBadRequestError("Invalid JSON body", &jsonBindingError)
		logger.Error("Failed to bind json", zap.Error(restErr.Error))
		c.JSON(restErr.Status, restErr)
		return
	}
	patchErr := profileService.Patch(principal.UserID, profilePatchRequest)
	if patchErr != nil {
		c.JSON(patchErr.Status, patchErr)
		return
	}
	c.JSON(http.StatusAccepted, nil)
}

func ChangePhone(c *gin.Context) {
	logger := application.NewTraceableLogger(c.Get("logger"))
	profileService := appwiring.InitProfileService(*logger)
	principal, ok := c.MustGet("principal").(*dto.Principal)
	if !ok {
		c.JSON(http.StatusUnauthorized, "principal not ok")
		return
	}
	var changePhoneRequest *request.ChangePhoneRequest
	if jsonBindingError := c.ShouldBindJSON(&changePhoneRequest); jsonBindingError != nil {
		restErr := utility.NewBadRequestError("Invalid JSON body", &jsonBindingError)
		logger.Error("Failed to bind json", zap.Error(restErr.Error))
		c.JSON(restErr.Status, restErr)
		return
	}
	patchErr := profileService.ChangePhoneNumber(principal.UserID, changePhoneRequest.Phone)
	if patchErr != nil {
		c.JSON(patchErr.Status, patchErr)
		return
	}
	c.JSON(http.StatusAccepted, nil)
}
