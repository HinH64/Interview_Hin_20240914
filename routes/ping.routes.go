package routes

import (
	"Interview_Hin_20240914/controllers"

	"github.com/gin-gonic/gin"
)

func PingRoute(router *gin.RouterGroup) {
	auth := router.Group("/ping")
	{
		auth.GET(
			"",
			controllers.Ping,
		)
	}
}
