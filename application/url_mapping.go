package application

import "timelyship.com/accounts/controller"

func mapUrls() {
	router.GET("/ping", controller.Ping)
	router.POST("/user", controller.RegisterUser)
}
