package validators

import (
	"net/http"

	"Interview_Hin_20240914/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CreateLevelValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var createLevelRequest models.LevelRequest
		_ = c.ShouldBindBodyWith(&createLevelRequest, binding.JSON)

		if err := createLevelRequest.Validate(); err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}