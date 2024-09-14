package controllers

import (
	"Interview_Hin_20240914/models"
	"Interview_Hin_20240914/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// CreateChallengeGame godoc
// @Summary Create a new challenge game
// @Description Create a new challenge game with the given details
// @Tags challengegames
// @Accept json
// @Produce json
// @Param challengegame body models.CreateChallengeGameRequest true "Create challenge game"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /challengegames [post]
func CreateChallengeGame(c *gin.Context) {
	var request models.CreateChallengeGameRequest
	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		models.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := services.CreateChallengeGame(request)
	if err != nil {
		models.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       gin.H{"result": result},
	}
	response.SendResponse(c)
}

// GetChallengeGames godoc
// @Summary List all challenge games
// @Description Get a list of all challenge games
// @Tags challengegames
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /challengegames [get]
func GetChallengeGames(c *gin.Context) {
	games, err := services.GetChallengeGames()
	if err != nil {
		models.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       gin.H{"games": games},
	}
	response.SendResponse(c)
}

// GetChallengeGameByID godoc
// @Summary Get a challenge game by ID
// @Description Get details of a specific challenge game
// @Tags challengegames
// @Accept json
// @Produce json
// @Param id path string true "Challenge Game ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /challengegames/{id} [get]
func GetChallengeGameByID(c *gin.Context) {
	id := c.Param("id")
	game, err := services.GetChallengeGameByID(id)
	if err != nil {
		models.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       gin.H{"game": game},
	}
	response.SendResponse(c)
}

// UpdateChallengeGame godoc
// @Summary Update a challenge game
// @Description Update details of a specific challenge game
// @Tags challengegames
// @Accept json
// @Produce json
// @Param id path string true "Challenge Game ID"
// @Param challengegame body models.UpdateChallengeGameRequest true "Update challenge game"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /challengegames/{id} [put]
func UpdateChallengeGame(c *gin.Context) {
	id := c.Param("id")
	var request models.UpdateChallengeGameRequest
	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		models.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	game, err := services.UpdateChallengeGame(id, request)
	if err != nil {
		models.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       gin.H{"game": game},
	}
	response.SendResponse(c)
}

// DeleteChallengeGame godoc
// @Summary Delete a challenge game
// @Description Delete a specific challenge game
// @Tags challengegames
// @Accept json
// @Produce json
// @Param id path string true "Challenge Game ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /challengegames/{id} [delete]
func DeleteChallengeGame(c *gin.Context) {
	id := c.Param("id")
	err := services.DeleteChallengeGame(id)
	if err != nil {
		models.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       gin.H{"game": nil},
	}
	response.SendResponse(c)
}