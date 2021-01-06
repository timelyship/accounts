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
	c.JSON(http.StatusOK, principal)
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
