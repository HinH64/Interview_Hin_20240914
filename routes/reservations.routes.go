package routes

import (
	"Interview_Hin_20240914/controllers"
	"Interview_Hin_20240914/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func ReservationRoute(router *gin.RouterGroup) {
	reservations := router.Group("/reservations")
	{
		reservations.GET("", controllers.ListReservations)
		reservations.POST("", validators.CreateReservationValidator(), controllers.CreateReservation)
	}
}