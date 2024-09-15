package validators

import (
	"net/http"

	"Interview_Hin_20240914/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CreateLogValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var createLogRequest models.LogRequest
		_ = c.ShouldBindBodyWith(&createLogRequest, binding.JSON)

		if err := createLogRequest.Validate(); err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}