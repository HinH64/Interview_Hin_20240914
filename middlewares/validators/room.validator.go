package validators

import (
	"net/http"

	"Interview_Hin_20240914/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CreateRoomValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var createRoomRequest models.RoomRequest
		_ = c.ShouldBindBodyWith(&createRoomRequest, binding.JSON)

		if err := createRoomRequest.Validate(); err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}

func UpdateRoomValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var updateRoomRequest models.RoomRequest
		_ = c.ShouldBindBodyWith(&updateRoomRequest, binding.JSON)

		if err := updateRoomRequest.Validate(); err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}