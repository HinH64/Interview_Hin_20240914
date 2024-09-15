package services

import (
	"Interview_Hin_20240914/enums"
	"Interview_Hin_20240914/models"
	db "Interview_Hin_20240914/models/db"
	"Interview_Hin_20240914/services"
	"testing"
	"time"

	"github.com/kamva/mgm/v3"
	"github.com/stretchr/testify/assert"
)

func TestCreateLog(t *testing.T) {
	services.LoadConfig("../../.env")
	services.InitMongoDBTest()

	// Create test player
	testPlayer := &db.Player{Name: "Test Player"}
	err := mgm.Coll(testPlayer).Create(testPlayer)
	assert.NoError(t, err)

	// Test case
	logRequest := models.LogRequest{
		PlayerID: testPlayer.ID.Hex(),
		Action:   enums.LogActionJoinChallenge,
		Details:  "Joined a challenge",
	}

	log, err := services.CreateLog(logRequest)
	assert.NoError(t, err)
	assert.NotNil(t, log)
	assert.Equal(t, testPlayer.ID, log.PlayerID)
	assert.Equal(t, logRequest.Action, log.Action)
	assert.Equal(t, logRequest.Details, log.Details)

	// Clean up
	err = mgm.Coll(&db.Player{}).Delete(testPlayer)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Log{}).Delete(log)
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
	now := time.Now()
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
		Details:  "Won challenge 1",
	}
	err = mgm.Coll(log2).Create(log2)
	assert.NoError(t, err)

	// Test case 1: Get all logs for the player
	logs, err := services.GetLogs(testPlayer.ID.Hex(), "", now.Add(-2*time.Hour), now.Add(time.Hour), 0, 10)
	assert.NoError(t, err)
	assert.Len(t, logs, 2)

	// Test case 2: Get logs with specific action
	logs2, err := services.GetLogs(testPlayer.ID.Hex(), enums.LogActionJoinChallenge, now.Add(-2*time.Hour), now.Add(time.Hour), 0, 10)
	assert.NoError(t, err)
	assert.Len(t, logs2, 1)
	assert.Equal(t, enums.LogActionJoinChallenge, logs2[0].Action)


	// Clean up
	err = mgm.Coll(&db.Player{}).Delete(testPlayer)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Log{}).Delete(log1)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Log{}).Delete(log2)
	assert.NoError(t, err)
}