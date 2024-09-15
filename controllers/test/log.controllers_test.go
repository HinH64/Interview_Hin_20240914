package controllers

import (
	"Interview_Hin_20240914/controllers"
	"Interview_Hin_20240914/enums"
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

func setupTestLogRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/logs", validators.CreateLogValidator(), controllers.CreateLog)
	r.GET("/logs", controllers.GetLogs)
	return r
}

func TestCreateLog(t *testing.T) {
	services.LoadConfig("../../.env")
	services.InitMongoDBTest()

	// Create test player
	testPlayer := &db.Player{Name: "Test Player"}
	err := mgm.Coll(testPlayer).Create(testPlayer)
	assert.NoError(t, err)

	// Setup router
	router := setupTestLogRouter()

	// Test case
	logRequest := models.LogRequest{
		PlayerID: testPlayer.ID.Hex(),
		Action:   enums.LogActionJoinChallenge,
		Details:  "Joined a challenge",
	}
	jsonValue, _ := json.Marshal(logRequest)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logs", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Verify log creation
	logData, ok := response.Data["log"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, testPlayer.ID.Hex(), logData["playerId"])
	assert.Equal(t, string(enums.LogActionJoinChallenge), logData["action"])
	assert.Equal(t, "Joined a challenge", logData["details"])

	// Clean up
	err = mgm.Coll(&db.Player{}).Delete(testPlayer)
	assert.NoError(t, err)
	log := &db.Log{}
	err = mgm.Coll(log).FindByID(logData["id"], log)
	assert.NoError(t, err)
	err = mgm.Coll(log).Delete(log)
	assert.NoError(t, err)
}

func TestGetLogs(t *testing.T) {
	services.LoadConfig("../../.env")
	services.InitMongoDBTest()

	// Create test player
	testPlayer := &db.Player{Name: "Test Player"}
	err := mgm.Coll(testPlayer).Create(testPlayer)
	assert.NoError(t, err)

	// Create test logs
	log1 := &db.Log{
		PlayerID: testPlayer.ID,
		Action:   enums.LogActionJoinChallenge,
		Details:  "Joined challenge 1",
	}
	err = mgm.Coll(log1).Create(log1)
	assert.NoError(t, err)

	log2 := &db.Log{
		PlayerID: testPlayer.ID,
		Action:   enums.LogActionLogin,
		Details:  "Logged in",
	}
	err = mgm.Coll(log2).Create(log2)
	assert.NoError(t, err)

	// Setup router
	router := setupTestLogRouter()

	// Test case
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/logs?playerID="+testPlayer.ID.Hex(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Verify logs retrieval
	logs, ok := response.Data["logs"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, logs, 2)

	// Clean up
	err = mgm.Coll(&db.Player{}).Delete(testPlayer)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Log{}).Delete(log1)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Log{}).Delete(log2)
	assert.NoError(t, err)
}