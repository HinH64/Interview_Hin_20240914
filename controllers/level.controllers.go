package controllers

import (
	"Interview_Hin_20240914/models"
	"Interview_Hin_20240914/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// ListLevels godoc
// @Summary List all levels
// @Description Get a list of all levels with pagination
// @Tags levels
// @Accept json
// @Produce json
// @Param        page  query    string  false  "Switch page by 'page'"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /levels [get]
func ListLevels(c *gin.Context) {
	pageQuery := c.DefaultQuery("page", "0")
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		models.SendErrorResponse(c, http.StatusBadRequest, "Invalid page number")
		return
	}

	limit := 5 // You might want to make this configurable

	levels, err := services.GetLevels(page, limit) // Get one extra to check for next page
	if err != nil {
		models.SendErrorResponse(c, http.StatusInternalServerError, "Error fetching levels")
		return
	}

	hasNext := len(levels) > limit
	if hasNext {
		levels = levels[:limit] // Remove the extra item
	}

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data: gin.H{
			"levels":   levels,
			"page":     page,
			"limit":    limit,
			"hasNext":  hasNext,
			"hasPrev":  page > 0,
		},
	}
	response.SendResponse(c)
}

// CreateLevel godoc
// @Summary Create a new level
// @Description Create a new level with the input payload
// @Tags levels
// @Accept json
// @Produce json
// @Param level body models.LevelRequest true "Create level"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /levels [post]
func CreateLevel(c *gin.Context) {
	var requestBody models.LevelRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}
	//check level is already in use
	err := services.CheckLevelNameInuse(requestBody.Name)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	level, err := services.CreateLevel(requestBody.Name)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{"level": level}
	response.SendResponse(c)
}