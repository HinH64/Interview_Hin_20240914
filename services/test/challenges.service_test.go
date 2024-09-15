package services

import (
	"Interview_Hin_20240914/models"
	db "Interview_Hin_20240914/models/db"
	"Interview_Hin_20240914/services"
	"testing"

	"github.com/kamva/mgm/v3"
	"github.com/stretchr/testify/assert"
)

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

	// Test case
	challenge, err := services.JoinChallenge(testPlayer.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, challenge)
	assert.Equal(t, testPlayer.ID, challenge.PlayerID)
	assert.Equal(t, testChallengeGame.ID, challenge.ChallengeGameID)
	assert.Equal(t, testChallengeGame.ChallengeCost, challenge.ChallengeCost)

	// Verify player's balance and challenge count
	updatedPlayer := &db.Player{}
	err = mgm.Coll(updatedPlayer).FindByID(testPlayer.ID, updatedPlayer)
	assert.NoError(t, err)
	assert.Equal(t, testPlayer.Balance-testChallengeGame.ChallengeCost, updatedPlayer.Balance)
	assert.Equal(t, testPlayer.ChallengeCount+1, updatedPlayer.ChallengeCount)

	assert.Equal(t, testChallengeGame.ChallengeCost, challenge.ChallengeCost)

	// Clean up
	err = mgm.Coll(&db.Level{}).Delete(testLevel)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Player{}).Delete(testPlayer)
	assert.NoError(t, err)
	err = mgm.Coll(&db.ChallengeGame{}).Delete(testChallengeGame)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Challenge{}).Delete(challenge)
	assert.NoError(t, err)
}

func TestGetLastChallenge(t *testing.T) {
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

	// Test case
	lastChallenge, err := services.GetLastChallenge()
	assert.NoError(t, err)
	assert.NotNil(t, lastChallenge)
	assert.Equal(t, testChallenge.ID, lastChallenge.ID)
	assert.Equal(t, testPlayer.ID, lastChallenge.PlayerID)
	assert.Equal(t, testChallengeGame.ID, lastChallenge.ChallengeGameID)
	assert.Equal(t, testChallenge.WinPrize, lastChallenge.WinPrize)
	assert.Equal(t, testChallenge.ChallengeCost, lastChallenge.ChallengeCost)

	// Clean up
	err = mgm.Coll(&db.Player{}).Delete(testPlayer)
	assert.NoError(t, err)
	err = mgm.Coll(&db.ChallengeGame{}).Delete(testChallengeGame)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Challenge{}).Delete(testChallenge)
	assert.NoError(t, err)
}

func TestCreateChallengeGame(t *testing.T) {
	services.LoadConfig("../../.env")
	services.InitMongoDBTest()

	// Test case
	request := models.CreateChallengeGameRequest{
		ChallengeTimeDurationSec: 60,
		ChallengeTimeCostSec:     10,
		ChallengeCost:            100,
		WinProbability:           0.5,
		AddWinProbability:        0.1,
	}

	challengeGame, err := services.CreateChallengeGame(request)
	assert.NoError(t, err)
	assert.NotNil(t, challengeGame)
	assert.Equal(t, request.ChallengeTimeDurationSec, challengeGame.ChallengeTimeDurationSec)
	assert.Equal(t, request.ChallengeTimeCostSec, challengeGame.ChallengeTimeCostSec)
	assert.Equal(t, request.ChallengeCost, challengeGame.ChallengeCost)
	assert.Equal(t, request.WinProbability, challengeGame.WinProbability)
	assert.Equal(t, request.AddWinProbability, challengeGame.AddWinProbability)

	// Clean up
	err = mgm.Coll(&db.ChallengeGame{}).Delete(challengeGame)
	assert.NoError(t, err)
}

func TestGetChallengeGames(t *testing.T) {
	services.LoadConfig("../../.env")
	services.InitMongoDBTest()

	// Create test challenge games
	game1 := &db.ChallengeGame{ChallengeTimeDurationSec: 60, ChallengeCost: 100}
	game2 := &db.ChallengeGame{ChallengeTimeDurationSec: 120, ChallengeCost: 200}
	err := mgm.Coll(game1).Create(game1)
	assert.NoError(t, err)
	err = mgm.Coll(game2).Create(game2)
	assert.NoError(t, err)

	// Test case
	games, err := services.GetChallengeGames()
	assert.NoError(t, err)
	assert.Len(t, games, 2)

	// Clean up
	err = mgm.Coll(&db.ChallengeGame{}).Delete(game1)
	assert.NoError(t, err)
	err = mgm.Coll(&db.ChallengeGame{}).Delete(game2)
	assert.NoError(t, err)
}

func TestGetChallengeGameByID(t *testing.T) {
	services.LoadConfig("../../.env")
	services.InitMongoDBTest()

	// Create test challenge game
	testGame := &db.ChallengeGame{ChallengeTimeDurationSec: 60, ChallengeCost: 100}
	err := mgm.Coll(testGame).Create(testGame)
	assert.NoError(t, err)

	// Test case
	game, err := services.GetChallengeGameByID(testGame.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, game)
	assert.Equal(t, testGame.ID, game.ID)
	assert.Equal(t, testGame.ChallengeTimeDurationSec, game.ChallengeTimeDurationSec)
	assert.Equal(t, testGame.ChallengeCost, game.ChallengeCost)

	// Clean up
	err = mgm.Coll(&db.ChallengeGame{}).Delete(testGame)
	assert.NoError(t, err)
}

func TestUpdateChallengeGame(t *testing.T) {
	services.LoadConfig("../../.env")
	services.InitMongoDBTest()

	// Create test challenge game
	testGame := &db.ChallengeGame{
		ChallengeTimeDurationSec: 60,
		ChallengeTimeCostSec:     10,
		ChallengeCost:            100,
		WinProbability:           0.5,
		AddWinProbability:        0.1,
	}
	err := mgm.Coll(testGame).Create(testGame)
	assert.NoError(t, err)

	// Test case
	newDuration := 120
	newCost := 200.0
	updateRequest := models.UpdateChallengeGameRequest{
		ChallengeTimeDurationSec: &newDuration,
		ChallengeCost:            &newCost,
	}

	updatedGame, err := services.UpdateChallengeGame(testGame.ID.Hex(), updateRequest)
	assert.NoError(t, err)
	assert.NotNil(t, updatedGame)
	assert.Equal(t, newDuration, updatedGame.ChallengeTimeDurationSec)
	assert.Equal(t, newCost, updatedGame.ChallengeCost)

	// Clean up
	err = mgm.Coll(&db.ChallengeGame{}).Delete(testGame)
	assert.NoError(t, err)
}

func TestDeleteChallengeGame(t *testing.T) {
	services.LoadConfig("../../.env")
	services.InitMongoDBTest()

	// Create test challenge game
	testGame := &db.ChallengeGame{ChallengeTimeDurationSec: 60, ChallengeCost: 100}
	err := mgm.Coll(testGame).Create(testGame)
	assert.NoError(t, err)

	// Test case
	err = services.DeleteChallengeGame(testGame.ID.Hex())
	assert.NoError(t, err)

	// Verify deletion
	deletedGame := &db.ChallengeGame{}
	err = mgm.Coll(deletedGame).FindByID(testGame.ID, deletedGame)
	assert.Error(t, err)
}