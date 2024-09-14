package validators

import (
	"net/http"

	"Interview_Hin_20240914/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func JoinChallengeValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var joinChallengeRequest models.ChallengeRequest
		_ = c.ShouldBindBodyWith(&joinChallengeRequest, binding.JSON)

		if err := joinChallengeRequest.Validate(); err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}