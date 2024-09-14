package validators

import (
	"Interview_Hin_20240914/models"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func PathIdValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")
		err := validation.Validate(id, is.MongoID)
		if err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, "invalid id: "+id)
			return
		}

		c.Next()
	}
}
