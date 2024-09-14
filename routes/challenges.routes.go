package routes

import (
	"Interview_Hin_20240914/controllers"
	"Interview_Hin_20240914/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func ChallengesRoute(router *gin.RouterGroup) {
	challengesGroup := router.Group("/challenges")
	{
		challengesGroup.POST("", validators.JoinChallengeValidator(), controllers.JoinChallenge)
		challengesGroup.GET("/results", controllers.GetChallengeResults)
	}

	challengeGamesGroup := router.Group("/challengegames")
	{
		challengeGamesGroup.POST("", controllers.CreateChallengeGame)
		challengeGamesGroup.GET("", controllers.GetChallengeGames)
		challengeGamesGroup.GET("/:id", controllers.GetChallengeGameByID)
		challengeGamesGroup.PUT("/:id", controllers.UpdateChallengeGame)
		challengeGamesGroup.DELETE("/:id", controllers.DeleteChallengeGame)
	}
}