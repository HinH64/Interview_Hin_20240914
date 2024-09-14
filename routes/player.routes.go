package routes

import (
	"Interview_Hin_20240914/controllers"
	"Interview_Hin_20240914/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func PlayerRoute(router *gin.RouterGroup) {
	players := router.Group("/players")
	{
		players.GET("", controllers.ListPlayers)
		players.POST("", validators.CreatePlayerValidator(), controllers.CreatePlayer)
		players.GET("/:id", validators.PathIdValidator(), controllers.GetPlayer)
		players.PUT("/:id", validators.PathIdValidator(), validators.UpdatePlayerValidator(), controllers.UpdatePlayer)
		players.DELETE("/:id", validators.PathIdValidator(), controllers.DeletePlayer)
	}
}