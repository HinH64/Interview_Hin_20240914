package routes

import (
	"Interview_Hin_20240914/controllers"
	"Interview_Hin_20240914/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func RoomRoute(router *gin.RouterGroup) {
	rooms := router.Group("/rooms")
	{
		rooms.GET("", controllers.ListRooms)
		rooms.POST("", validators.CreateRoomValidator(), controllers.CreateRoom)
		rooms.GET("/:id", validators.PathIdValidator(), controllers.GetRoom)
		rooms.PUT("/:id", validators.PathIdValidator(), validators.UpdateRoomValidator(), controllers.UpdateRoom)
		rooms.DELETE("/:id", validators.PathIdValidator(), controllers.DeleteRoom)
	}
}