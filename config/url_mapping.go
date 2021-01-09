package config

import (
	"timelyship.com/accounts/controller"
)

func mapUrls() {
	router.GET("/ping", controller.Ping)
	// router.POST("/user", controller.RegisterUser)
	router.GET("/login-google", controller.LoginByGoogle)
	router.GET("/login-fb", controller.LoginByFB)
	router.GET("/account/google-login/redirect", controller.HandleRedirectFromGoogle)
	router.GET("/account/fb-login/redirect", controller.HandleRedirectFromFB)
	router.POST("/account/sign-up", controller.SignUp)
	router.POST("/account/login", controller.Login)
	router.POST("/account/refresh-token", controller.RefreshToken)
	router.POST("/verify-email/:verificationToken", controller.VerifyEmail)
	router.GET("/initiate-login", controller.InitiateLogin)
	router.GET("/generate-code", controller.GenerateCode)
	router.GET("/decode-code", controller.Decode)
	router.GET("/exchange-code", controller.ExchangeCode)
	router.GET("/profile", controller.Profile)
	router.PATCH("/profile-patch", controller.PatchProfile)
	router.PATCH("/change-phone", controller.ChangePhone)
	router.PUT("/upload-profile-photo", controller.UploadProfilePhoto)
	router.GET("/logout", controller.Logout)
}
