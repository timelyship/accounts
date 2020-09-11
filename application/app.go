package application

import "github.com/gin-gonic/gin"

var (
	router = gin.Default()
)

func Start() {
	mapUrls()
	router.Run(":8080")

}
