package controllers

import (
	"Interview_Hin_20240914/enums"
	"Interview_Hin_20240914/models"
	"Interview_Hin_20240914/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// GetLogs godoc
// @Summary List logs
// @Description Get a list of logs with optional filters
// @Tags logs
// @Accept json
// @Produce json
// @Param playerId query string false "Player ID"
// @Param action query string false "Action type"
// @Param start_time query string false "Start time (RFC3339 format)"
// @Param end_time query string false "End time (RFC3339 format)"
// @Param limit query int false "Limit"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /logs [get]
func GetLogs(c *gin.Context) {
	playerID := c.Query("playerId")
	action := c.Query("action")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")
	limitStr := c.DefaultQuery("limit", "5")
	pageStr := c.DefaultQuery("page", "0")


	var startTime, endTime time.Time
	var err error
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, "Invalid start_time format")
			return
		}
	}
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, "Invalid end_time format")
			return
		}
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		models.SendErrorResponse(c, http.StatusBadRequest, "Invalid page number")
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		models.SendErrorResponse(c, http.StatusBadRequest, "Invalid limit number")
		return
	}

	logs, err := services.GetLogs(playerID, enums.LogAction(action), startTime, endTime, page, limit)
	if err != nil {
		models.SendErrorResponse(c, http.StatusInternalServerError, "Error fetching logs")
		return
	}

	hasNext := len(logs) > limit
	if hasNext {
		logs = logs[:limit] // Remove the extra item
	}

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data: gin.H{
			"logs":    logs,
			"page":     page,
			"limit":    limit,
			"hasNext":  hasNext,
			"hasPrev":  page > 0,
		},
	}
	response.SendResponse(c)
}

// CreateLog godoc
// @Summary Create a new log
// @Description Create a new log entry
// @Tags logs
// @Accept json
// @Produce json
// @Param log body models.LogRequest true "Create log"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /logs [post]
func CreateLog(c *gin.Context) {
	var requestBody models.LogRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	log, err := services.CreateLog(requestBody)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{"log": log}
	response.SendResponse(c)
}