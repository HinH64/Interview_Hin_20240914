package routes

import (
	"Interview_Hin_20240914/controllers"
	"Interview_Hin_20240914/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func PaymentRoute(router *gin.RouterGroup) {
	payments := router.Group("/payments")
	{
		payments.POST("", validators.CreatePaymentValidator(), controllers.CreatePayment)
		payments.GET("/:id", validators.PathIdValidator(), controllers.GetPayment)
	}
}