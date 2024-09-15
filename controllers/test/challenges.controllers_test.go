package controllers

import (
	"Interview_Hin_20240914/controllers"
	"Interview_Hin_20240914/middlewares/validators"
	"Interview_Hin_20240914/models"
	db "Interview_Hin_20240914/models/db"
	"Interview_Hin_20240914/services"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"github.com/stretchr/testify/assert"
)

func setupTestChallengesRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/challenges", validators.JoinChallengeValidator(), controllers.JoinChallenge)
	r.GET("/challenges/results", controllers.GetChallengeResults)
	return r
}

func TestJoinChallenge(t *testing.T) {
	services.LoadConfig("../../.env")
	services.InitMongoDBTest()

	// Create test level
	testLevel := &db.Level{Name: "Test Level"}
	err := mgm.Coll(testLevel).Create(testLevel)
	assert.NoError(t, err)

	// Create test player
	testPlayer := &db.Player{
		Name:           "Test Player",
		Balance:        1000,
		ChallengeCount: 0,
		LevelID:        testLevel.ID,
	}
	err = mgm.Coll(testPlayer).Create(testPlayer)
	assert.NoError(t, err)

	// Create test challenge game
	testChallengeGame := &db.ChallengeGame{
		ChallengeTimeDurationSec: 60,
		ChallengeTimeCostSec:     10,
		ChallengeCost:            100,
		WinProbability:           0.5,
		AddWinProbability:        0.1,
	}
	err = mgm.Coll(testChallengeGame).Create(testChallengeGame)
	assert.NoError(t, err)

	// Setup router
	router := setupTestChallengesRouter()

	// Test case
	challengeRequest := models.ChallengeRequest{
		PlayerID: testPlayer.ID.Hex(),
	}
	jsonValue, _ := json.Marshal(challengeRequest)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/challenges", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Verify challenge creation
	challengeData, ok := response.Data["challenge"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, testPlayer.ID.Hex(), challengeData["playerId"])
	assert.Equal(t, testChallengeGame.ID.Hex(), challengeData["challengeGameId"])

	// Clean up
	err = mgm.Coll(&db.Level{}).Delete(testLevel)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Player{}).Delete(testPlayer)
	assert.NoError(t, err)
	err = mgm.Coll(&db.ChallengeGame{}).Delete(testChallengeGame)
	assert.NoError(t, err)
}

func TestGetChallengeResults(t *testing.T) {
	services.LoadConfig("../../.env")
	services.InitMongoDBTest()

	// Create test player
	testPlayer := &db.Player{Name: "Test Player", Balance: 1000, ChallengeCount: 1}
	err := mgm.Coll(testPlayer).Create(testPlayer)
	assert.NoError(t, err)

	// Create test challenge game
	testChallengeGame := &db.ChallengeGame{
		ChallengeTimeDurationSec: 60,
		ChallengeTimeCostSec:     10,
		ChallengeCost:            100,
		WinProbability:           0.5,
		AddWinProbability:        0.1,
	}
	err = mgm.Coll(testChallengeGame).Create(testChallengeGame)
	assert.NoError(t, err)

	// Create test challenge
	testChallenge := &db.Challenge{
		PlayerID:        testPlayer.ID,
		ChallengeGameID: testChallengeGame.ID,
		WinPrize:        true,
		ChallengeCost:   100,
	}
	err = mgm.Coll(testChallenge).Create(testChallenge)
	assert.NoError(t, err)

	// Setup router
	router := setupTestChallengesRouter()

	// Test case
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/challenges/results", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Verify challenge result
	resultData, ok := response.Data["result"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, testChallenge.ID.Hex(), resultData["id"])
	assert.Equal(t, testPlayer.ID.Hex(), resultData["playerId"])
	assert.Equal(t, testChallengeGame.ID.Hex(), resultData["challengeGameId"])
	assert.Equal(t, testChallenge.WinPrize, resultData["winPrize"])
	assert.Equal(t, testChallenge.ChallengeCost, resultData["challengeCost"])

	// Clean up
	err = mgm.Coll(&db.Player{}).Delete(testPlayer)
	assert.NoError(t, err)
	err = mgm.Coll(&db.ChallengeGame{}).Delete(testChallengeGame)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Challenge{}).Delete(testChallenge)
	assert.NoError(t, err)
}