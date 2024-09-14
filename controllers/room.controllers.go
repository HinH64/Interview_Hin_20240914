package controllers

import (
	"Interview_Hin_20240914/models"
	"Interview_Hin_20240914/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ListRooms godoc
// @Summary List all rooms
// @Description Get a list of all rooms with pagination
// @Tags rooms
// @Accept json
// @Produce json
// @Param page query string false "Switch page by 'page'"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /rooms [get]
func ListRooms(c *gin.Context) {
	pageQuery := c.DefaultQuery("page", "0")
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		models.SendErrorResponse(c, http.StatusBadRequest, "Invalid page number")
		return
	}

	limit := 5 // You might want to make this configurable

	rooms, err := services.GetRooms(page, limit+1) // Get one extra to check for next page
	if err != nil {
		models.SendErrorResponse(c, http.StatusInternalServerError, "Error fetching rooms")
		return
	}

	hasNext := len(rooms) > limit
	if hasNext {
		rooms = rooms[:limit] // Remove the extra item
	}

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data: gin.H{
			"rooms":    rooms,
			"page":     page,
			"limit":    limit,
			"hasNext":  hasNext,
			"hasPrev":  page > 0,
		},
	}
	response.SendResponse(c)
}

// CreateRoom godoc
// @Summary Create a new room
// @Description Create a new room with the input payload
// @Tags rooms
// @Accept json
// @Produce json
// @Param room body models.RoomRequest true "Create room"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /rooms [post]
func CreateRoom(c *gin.Context) {
	var requestBody models.RoomRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	room, err := services.CreateRoom(requestBody)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{"room": room}
	response.SendResponse(c)
}

// GetRoom godoc
// @Summary Get a room
// @Description Get a room by ID
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path string true "Room ID"
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /rooms/{id} [get]
func GetRoom(c *gin.Context) {
	id := c.Param("id")
	roomID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		models.SendErrorResponse(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	room, err := services.GetRoom(roomID)
	if err != nil {
		models.SendErrorResponse(c, http.StatusNotFound, "Room not found")
		return
	}

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       gin.H{"room": room},
	}
	response.SendResponse(c)
}

// UpdateRoom godoc
// @Summary Update a room
// @Description updates a room by id
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path string true "Room ID"
// @Param req body models.RoomRequest true "Room Request"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /rooms/{id} [put]
func UpdateRoom(c *gin.Context) {
	id := c.Param("id")
	roomID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		models.SendErrorResponse(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	var requestBody models.RoomRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	err = services.UpdateRoom(roomID, &requestBody)
	if err != nil {
		models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "Room updated successfully",
	}
	response.SendResponse(c)
}

// DeleteRoom godoc
// @Summary Delete a room
// @Description Delete a room by ID
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path string true "Room ID"
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /rooms/{id} [delete]
func DeleteRoom(c *gin.Context) {
	id := c.Param("id")
	roomID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		models.SendErrorResponse(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	err = services.DeleteRoom(roomID)
	if err != nil {
		models.SendErrorResponse(c, http.StatusNotFound, "Room not found")
		return
	}

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "Room deleted successfully",
	}
	response.SendResponse(c)
}