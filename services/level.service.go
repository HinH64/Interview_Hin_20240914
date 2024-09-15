package services

import (
	db "Interview_Hin_20240914/models/db"
	"errors"
	"fmt"

	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateLevel(name string) (*db.Level, error) {
	level := db.NewLevel(name)
	err := mgm.Coll(level).Create(level)
	if err != nil {
		return nil, errors.New("cannot create new level")
	}

	return level, nil
}

func CheckLevelNameInuse(name string) error {
	level := &db.Level{}
	levelCollection := mgm.Coll(level)
	err := levelCollection.First(bson.M{"name": name}, level)
	if err == nil {
		return errors.New("level name is already in use")
	}

	return nil
}
func CheckLevelExist(levelId primitive.ObjectID) error {
	level := &db.Level{}
	levelCollection := mgm.Coll(level)
	err := levelCollection.First(bson.M{field.ID: levelId}, level)
	if err != nil {
		return errors.New("cannot find level")
	}

	return nil
}
func GetLevels(page int, limit int) ([]db.Level, error) {
	var levels []db.Level

	// Ensure page is at least 0
	if page < 0 {
		page = 0
	}

	// Calculate skip
	skip := page * limit

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit + 1))

	err := mgm.Coll(&db.Level{}).SimpleFind(&levels, bson.M{}, findOptions)

	if err != nil {
		return nil, fmt.Errorf("cannot find levels: %w", err)
	}

	return levels, nil
}


