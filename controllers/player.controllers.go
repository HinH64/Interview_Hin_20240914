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

// ListPlayers godoc
// @Summary List all players
// @Description Get a list of all players with pagination
// @Tags players
// @Accept json
// @Produce json
// @Param page query string false "Switch page by 'page'"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /players [get]
func ListPlayers(c *gin.Context) {
	pageQuery := c.DefaultQuery("page", "0")
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		models.SendErrorResponse(c, http.StatusBadRequest, "Invalid page number")
		return
	}

	limit := 5 // You might want to make this configurable

	players, err := services.GetPlayers(page, limit+1) // Get one extra to check for next page
	if err != nil {
		models.SendErrorResponse(c, http.StatusInternalServerError, "Error fetching players")
		return
	}

	hasNext := len(players) > limit
	if hasNext {
		players = players[:limit] // Remove the extra item
	}

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data: gin.H{
			"players":  players,
			"page":     page,
			"limit":    limit,
			"hasNext":  hasNext,
			"hasPrev":  page > 0,
		},
	}
	response.SendResponse(c)
}

// CreatePlayer godoc
// @Summary Create a new player
// @Description Create a new player with the input payload
// @Tags players
// @Accept json
// @Produce json
// @Param player body models.PlayerRequest true "Create player"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /players [post]
func CreatePlayer(c *gin.Context) {
	
	var requestBody models.PlayerRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	player, err := services.CreatePlayer(requestBody)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{"player": player}
	response.SendResponse(c)
}

// GetPlayer godoc
// @Summary Get a player
// @Description Get a player by ID
// @Tags players
// @Accept json
// @Produce json
// @Param id path string true "Player ID"
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /players/{id} [get]
func GetPlayer(c *gin.Context) {
	idHex := c.Param("id")
	playerId, _ := primitive.ObjectIDFromHex(idHex)
	player, err := services.GetPlayerByID(playerId)
	if err != nil {
		models.SendErrorResponse(c, http.StatusNotFound, "Player not found")
		return
	}
	models.SendResponseData(c, gin.H{"player": player})
}

// UpdatePlayer godoc
// @Summary      Update a player
// @Description  updates a player by id
// @Tags         players
// @Accept       json
// @Produce      json
// @Param        id     path    string  true  "Player ID"
// @Param        req    body    models.PlayerRequest true "Player Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /players/{id} [put]
func UpdatePlayer(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	idHex := c.Param("id")
	playerId, _ := primitive.ObjectIDFromHex(idHex)

	var playerRequest models.PlayerRequest
	_ = c.ShouldBindBodyWith(&playerRequest, binding.JSON)

	err := services.UpdatePlayer(playerId, &playerRequest)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.SendResponse(c)
}

// DeletePlayer godoc
// @Summary Delete a player
// @Description Delete a player by ID
// @Tags players
// @Accept json
// @Produce json
// @Param id path string true "Player ID"
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /players/{id} [delete]
func DeletePlayer(c *gin.Context) {
	idHex := c.Param("id")
	playerId, _ := primitive.ObjectIDFromHex(idHex)
	err := services.DeletePlayer(playerId)
	if err != nil {
		models.SendErrorResponse(c, http.StatusNotFound, "Player not found")
		return
	}

	models.SendResponseData(c, gin.H{"message": "Player deleted successfully"})
}