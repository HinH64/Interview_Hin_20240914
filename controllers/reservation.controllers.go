package controllers

import (
	"Interview_Hin_20240914/models"
	"Interview_Hin_20240914/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// ListReservations godoc
// @Summary List reservations
// @Description Get a list of reservations with optional filters
// @Tags reservations
// @Accept json
// @Produce json
// @Param room_id query string false "Room ID"
// @Param date query string false "Date (YYYY-MM-DD)"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /reservations [get]
func ListReservations(c *gin.Context) {
	roomID := c.Query("room_id")
	dateStr := c.Query("date")
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "0")

	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)

	var date *time.Time
	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, "Invalid date format")
			return
		}
		date = &parsedDate
	}

	reservations, err := services.GetReservations(roomID, date, page, limit)
	if err != nil {
		models.SendErrorResponse(c, http.StatusInternalServerError, "Error fetching reservations")
		return
	}

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data: gin.H{
			"reservations": reservations,
			"page":         page,
			"limit":        limit,
		},
	}
	response.SendResponse(c)
}

// CreateReservation godoc
// @Summary Create a new reservation
// @Description Create a new reservation with the input payload
// @Tags reservations
// @Accept json
// @Produce json
// @Param reservation body models.ReservationRequest true "Create reservation"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /reservations [post]
func CreateReservation(c *gin.Context) {
	var requestBody models.ReservationRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	reservation, err := services.CreateReservation(requestBody)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{"reservation": reservation}
	response.SendResponse(c)
}