package application

import "timelyship.com/accounts/controller"

func mapUrls() {
	router.GET("/ping", controller.Ping)
	router.POST("/user", controller.RegisterUser)
	router.GET("/login-google", controller.LoginByGoogle)
	router.GET("/account/google-login/redirect", controller.HandleRedirectFromGoogle)
}
