package routes

import (
	"Interview_Hin_20240914/controllers"
	"Interview_Hin_20240914/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func LevelRoute(router *gin.RouterGroup) {
	levels := router.Group("/levels")
	{
		levels.GET("", controllers.ListLevels)
		levels.POST("", validators.CreateLevelValidator(), controllers.CreateLevel)
	}
}