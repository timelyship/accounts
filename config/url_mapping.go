package config

import (
	"timelyship.com/accounts/controller"
)

func mapUrls() {
	router.GET("/ping", controller.Ping)
	//router.POST("/user", controller.RegisterUser)
	router.GET("/login-google", controller.LoginByGoogle)
	router.GET("/login-fb", controller.LoginByFB)
	router.GET("/account/google-login/redirect", controller.HandleRedirectFromGoogle)
	router.GET("/account/fb-login/redirect", controller.HandleRedirectFromDB)
	router.POST("/account/sign-up", controller.SignUp)
	router.POST("/verify-email/:verificationToken", controller.VerifyEmail)
}
