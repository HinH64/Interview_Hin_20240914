package controllers

import (
	"Interview_Hin_20240914/models"
	"Interview_Hin_20240914/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// JoinChallenge godoc
// @Summary Join a challenge
// @Description Player joins a challenge with a fixed payment amount
// @Tags challenges
// @Accept json
// @Produce json
// @Param challenge body models.ChallengeRequest true "Join challenge"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /challenges [post]
func JoinChallenge(c *gin.Context) {
	var requestBody models.ChallengeRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	result, err := services.JoinChallenge(requestBody.PlayerID)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"challenge": result}
	response.SendResponse(c)
}

// GetChallengeResults godoc
// @Summary Get last challenge result
// @Description Get the most recent challenge result for a player
// @Tags challenges
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /challenges/results [get]
func GetChallengeResults(c *gin.Context) {

	result, err := services.GetLastChallenge()
	if err != nil {
		models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       gin.H{"result": result},
	}
	response.SendResponse(c)
}