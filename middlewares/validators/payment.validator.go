package validators

import (
	"net/http"

	"Interview_Hin_20240914/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CreatePaymentValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var createPaymentRequest models.PaymentRequest
		_ = c.ShouldBindBodyWith(&createPaymentRequest, binding.JSON)

		if err := createPaymentRequest.Validate(); err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}