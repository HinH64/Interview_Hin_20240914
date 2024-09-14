package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Player represents a game player
// @Description Player model
type Player struct {
	// DefaultModel contains fields for the MongoDB document
	mgm.DefaultModel `bson:",inline"`
	// Name of the player
	Name string `json:"name" bson:"name" example:"John Doe"`
	// LevelID references the Level document
	LevelID primitive.ObjectID `json:"levelId" bson:"levelId" example:"5f5e7e9b9b9b9b9b9b9b9b9b"`
}

type GetPlayer struct {
	// DefaultModel contains fields for the MongoDB document
	mgm.DefaultModel `bson:",inline"`
	// Name of the player
	Name string `json:"name" bson:"name" example:"John Doe"`
	// LevelID references the Level document
	LevelID primitive.ObjectID `json:"levelId" bson:"levelId" example:"5f5e7e9b9b9b9b9b9b9b9b9b"`
	LevelName string            `json:"levelName" bson:"-"`
}