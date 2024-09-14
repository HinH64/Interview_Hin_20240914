package validators

import (
	"net/http"

	"Interview_Hin_20240914/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CreateReservationValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var createReservationRequest models.ReservationRequest
		_ = c.ShouldBindBodyWith(&createReservationRequest, binding.JSON)

		if err := createReservationRequest.Validate(); err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}