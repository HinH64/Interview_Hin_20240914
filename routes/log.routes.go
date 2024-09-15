package routes

import (
	"Interview_Hin_20240914/controllers"
	"Interview_Hin_20240914/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func LogRoute(router *gin.RouterGroup) {
	logs := router.Group("/logs")
	{
		logs.GET("", controllers.GetLogs)
		logs.POST("", validators.CreateLogValidator(), controllers.CreateLog)
	}
}