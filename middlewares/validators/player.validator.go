package validators

import (
	"net/http"

	"Interview_Hin_20240914/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CreatePlayerValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var createPlayerRequest models.PlayerRequest
		_ = c.ShouldBindBodyWith(&createPlayerRequest, binding.JSON)

		if err := createPlayerRequest.Validate(); err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}

func UpdatePlayerValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var updatePlayerRequest models.PlayerRequest
		_ = c.ShouldBindBodyWith(&updatePlayerRequest, binding.JSON)

		if err := updatePlayerRequest.Validate(); err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}