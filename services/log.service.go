package services

import (
	"Interview_Hin_20240914/enums"
	"Interview_Hin_20240914/models"
	db "Interview_Hin_20240914/models/db"
	"errors"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetLogs(playerID string, action enums.LogAction, startTime, endTime time.Time, page int, limit int) ([]db.Log, error) {
	var logs []db.Log

	if page < 0 {
		page = 0
	}

	skip := page * limit

	filter := bson.M{}
	if playerID != "" {
		playerObjID, err := primitive.ObjectIDFromHex(playerID)
		if err != nil {
			return nil, errors.New("invalid player ID")
		}
		filter["playerId"] = playerObjID
	}
	if action != "" {
		filter["action"] = action
	}
	if !startTime.IsZero() {
		filter["created_at"] = bson.M{"$gte": startTime}
	}
	if !endTime.IsZero() {
		if filter["created_at"] == nil {
			filter["created_at"] = bson.M{}
		}
		filter["created_at"].(bson.M)["$lte"] = endTime
	}

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit + 1)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	err := mgm.Coll(&db.Log{}).SimpleFind(&logs, filter, findOptions)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func CreateLog(request models.LogRequest) (*db.Log, error) {
	playerID, err := primitive.ObjectIDFromHex(request.PlayerID)
	if err != nil {
		return nil, errors.New("invalid player ID")
	}
	log := &db.Log{
		PlayerID:  playerID,
		Action:    request.Action,
		Details:   request.Details,
	}

	err = mgm.Coll(log).Create(log)
	if err != nil {
		return nil, errors.New("cannot create new log")
	}


	return log, nil
}