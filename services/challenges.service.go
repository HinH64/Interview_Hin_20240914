package services

import (
	"Interview_Hin_20240914/models"
	db "Interview_Hin_20240914/models/db"
	"errors"
	"math/rand"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func JoinChallenge(playerID string) (*db.Challenge, error) {
	objectID, err := primitive.ObjectIDFromHex(playerID)
	if err != nil {
		return nil, errors.New("invalid player ID")
	}

	// Check if the player exists
	player := &db.Player{}
	err = mgm.Coll(player).FindByID(objectID, player)
	if err != nil {
		return nil, errors.New("player not found")
	}

	// Check if the player has participated in a challenge within the last minute
	lastMinute := time.Now().Add(-1 * time.Minute)
	var recentChallenge db.Challenge
	err = mgm.Coll(&db.Challenge{}).FindOne(mgm.Ctx(), bson.M{
		"playerId": objectID,
		"created_at": bson.M{
			"$gte": lastMinute,
		},
	}).Decode(&recentChallenge)
	if err == nil {
		return nil, errors.New("player can only participate once per minute")
	}

	// Get the last ChallengeGame by created_at
	var challengeGame db.ChallengeGame
	err = mgm.Coll(&db.ChallengeGame{}).FindOne(
		mgm.Ctx(),
		bson.M{},
		options.FindOne().SetSort(bson.M{"created_at": -1}),
	).Decode(&challengeGame)
	if err != nil {
		return nil, errors.New("failed to get the last challenge game")
	}

	// Check if player has enough balance
	if player.Balance < challengeGame.ChallengeCost {
		return nil, errors.New("insufficient balance")
	}

	// Create a new challenge
	challenge := &db.Challenge{
		PlayerID:        objectID,
		ChallengeGameID: challengeGame.ID,
		WinPrize:        false,
		ChallengeCost:   challengeGame.ChallengeCost,
	}

	err = mgm.Coll(challenge).Create(challenge)
	if err != nil {
		return nil, errors.New("cannot create new challenge")
	}

	// Update player's balance and challenge count
	_, err = mgm.Coll(player).UpdateOne(
		mgm.Ctx(),
		bson.M{"_id": objectID},
		bson.M{
			"$inc": bson.M{
				"balance":       -challengeGame.ChallengeCost,
				"challengeCount": 1,
			},
		},
	)
	if err != nil {
		return nil, errors.New("failed to update player's balance and challenge count")
	}

	// Update ChallengeGame's total prize money
	_, err = mgm.Coll(&challengeGame).UpdateOne(
		mgm.Ctx(),
		bson.M{"_id": challengeGame.ID},
		bson.M{"$inc": bson.M{"totalPrizeMoney": challengeGame.ChallengeCost}},
	)
	if err != nil {
		return nil, errors.New("failed to update challenge game's total prize money")
	}

	// Process the challenge result immediately
	processChallengeResult(challenge, &challengeGame, player)

	return challenge, nil
}

func processChallengeResult(challenge *db.Challenge, challengeGame *db.ChallengeGame, player *db.Player) {
	winProbability := challengeGame.WinProbability + (float64(player.ChallengeCount) * challengeGame.AddWinProbability)
	randNum := rand.Float64() * 100

	if randNum <= winProbability {
		challenge.WinPrize = true
	}

	if challenge.WinPrize {
		// Update player's balance and reset challenge count
		_, err := mgm.Coll(player).UpdateOne(
			mgm.Ctx(),
			bson.M{"_id": player.ID},
			bson.M{
				"$inc": bson.M{"balance": challengeGame.TotalPrizeMoney},
				"$set": bson.M{"challengeCount": 0},
			},
		)
		if err != nil {
			// Handle error (e.g., log it)
		}

		// Reset ChallengeGame's total prize money
		_, err = mgm.Coll(challengeGame).UpdateOne(
			mgm.Ctx(),
			bson.M{"_id": challengeGame.ID},
			bson.M{"$set": bson.M{"totalPrizeMoney": 0}},
		)
		if err != nil {
			// Handle error (e.g., log it)
		}
	}

	// Update the challenge result
	err := mgm.Coll(challenge).Update(challenge)
	if err != nil {
		// Handle error (e.g., log it)
	}
}
func GetLastChallenge() (*db.Challenge, error) {
	var challenge db.Challenge
	err := mgm.Coll(&db.Challenge{}).FindOne(
		mgm.Ctx(),
		bson.M{},
		options.FindOne().SetSort(bson.M{"created_at": -1}),
	).Decode(&challenge)
	if err != nil {
		return nil, errors.New("failed to fetch challenge result")
	}
	return &challenge, nil
}

func CreateChallengeGame(request models.CreateChallengeGameRequest) (*db.ChallengeGame, error) {
	game := &db.ChallengeGame{
		ChallengeTimeDurationSec: request.ChallengeTimeDurationSec,
		ChallengeTimeCostSec: request.ChallengeTimeCostSec,
		ChallengeCost:        request.ChallengeCost,
		WinProbability:       request.WinProbability,
		AddWinProbability:    request.AddWinProbability,
		TotalPrizeMoney:      0, // Initialize with 0
	}

	err := mgm.Coll(game).Create(game)
	if err != nil {
		return nil, errors.New("cannot create new challenge game")
	}
	return game, nil
}

func GetChallengeGames() ([]db.ChallengeGame, error) {
	var games []db.ChallengeGame
	err := mgm.Coll(&db.ChallengeGame{}).SimpleFind(&games, bson.M{})
	if err != nil {
		return nil, errors.New("failed to fetch challenge games")
	}
	return games, nil
}

func GetChallengeGameByID(id string) (*db.ChallengeGame, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid challenge game ID")
	}

	game := &db.ChallengeGame{}
	err = mgm.Coll(game).FindByID(objectID, game)
	if err != nil {
		return nil, errors.New("challenge game not found")
	}
	return game, nil
}

func UpdateChallengeGame(id string, request models.UpdateChallengeGameRequest) (*db.ChallengeGame, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid challenge game ID")
	}

	game := &db.ChallengeGame{}
	err = mgm.Coll(game).FindByID(objectID, game)
	if err != nil {
		return nil, errors.New("challenge game not found")
	}

	// Update fields
	if request.ChallengeTimeDurationSec != nil {
		game.ChallengeTimeDurationSec = *request.ChallengeTimeDurationSec
	}
	if request.ChallengeTimeCostSec != nil {
		game.ChallengeTimeCostSec = *request.ChallengeTimeCostSec
	}
	if request.ChallengeCost != nil {
		game.ChallengeCost = *request.ChallengeCost
	}
	if request.WinProbability != nil {
		game.WinProbability = *request.WinProbability
	}
	if request.AddWinProbability != nil {
		game.AddWinProbability = *request.AddWinProbability
	}

	err = mgm.Coll(game).Update(game)
	if err != nil {
		return nil, errors.New("failed to update challenge game")
	}
	return game, nil
}

func DeleteChallengeGame(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid challenge game ID")
	}

	game := &db.ChallengeGame{}
	err = mgm.Coll(game).FindByID(objectID, game)
	if err != nil {
		return errors.New("challenge game not found")
	}

	err = mgm.Coll(game).Delete(game)
	if err != nil {
		return errors.New("failed to delete challenge game")
	}
	return nil
}