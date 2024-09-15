package services

import (
	"Interview_Hin_20240914/models"
	db "Interview_Hin_20240914/models/db"
	"errors"
	"fmt"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreatePlayer(request models.PlayerRequest) (*db.Player, error) {
	levelID, err := primitive.ObjectIDFromHex(request.LevelID)
	if err != nil {
		return nil, errors.New("invalid levelId")
	}

	player := &db.Player{
		Name:    request.Name,
		LevelID: levelID,
		Balance: request.Balance,
		ChallengeCount: 0,
	}
	err = CheckLevelExist(player.LevelID)
	if err != nil {
		return nil, errors.New("cannot find level")
	}

	err = mgm.Coll(player).Create(player)
	if err != nil {
		return nil, errors.New("cannot create new player")
	}

	return player, nil
}

func GetPlayers(page int, limit int) ([]db.GetPlayer, error) {
	var players []db.GetPlayer

	if page < 0 {
		page = 0
	}

	skip := page * limit

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit + 1))

	err := mgm.Coll(&db.Player{}).SimpleFind(&players, bson.M{}, findOptions)

	if err != nil {
		return nil, fmt.Errorf("cannot find players: %w", err)
	}

	for i := range players {
		level := &db.Level{}
		err := mgm.Coll(level).FindByID(players[i].LevelID, level)
		if err == nil {
			players[i].LevelName = level.Name
		}
	}

	return players, nil
}

func GetPlayerByID(id primitive.ObjectID) (*db.GetPlayer, error) {
    var player db.Player
    err := mgm.Coll(&player).FindByID(id, &player)
    if err != nil {
        return nil, err
    }
	var getPlayer db.GetPlayer
	getPlayer.ID = player.ID
	getPlayer.Name = player.Name
	getPlayer.LevelID = player.LevelID
	getPlayer.Balance = player.Balance
	getPlayer.ChallengeCount = player.ChallengeCount

	level := &db.Level{}
	err = mgm.Coll(level).FindByID(player.LevelID, level)
	if err == nil {
		getPlayer.LevelName = level.Name
	}

    return &getPlayer, nil
}

func UpdatePlayer(playerId primitive.ObjectID, request *models.PlayerRequest) error {
	player := &db.Player{}
	err := mgm.Coll(player).FindByID(playerId, player)
	if err != nil {
		return errors.New("cannot find player")
	}

	player.Name = request.Name
	player.LevelID, err = primitive.ObjectIDFromHex(request.LevelID)
	if err != nil {
		return errors.New("invalid levelId")
	}
	if request.Balance != 0 {
		player.Balance = request.Balance
	}
	if request.ChallengeCount != 0 {
		player.ChallengeCount = request.ChallengeCount
	}

	err = CheckLevelExist(player.LevelID)
	if err != nil {
		return errors.New("cannot find level")
	}

	err = mgm.Coll(player).Update(player)

	if err != nil {
		return errors.New("cannot update")
	}

	return nil
}

func DeletePlayer(playerId primitive.ObjectID) error {
	player, err := GetPlayerByID(playerId)
	if err != nil {
		return err
	}

	err = mgm.Coll(player).Delete(player)
	if err != nil {
		return err
	}

	return nil
}