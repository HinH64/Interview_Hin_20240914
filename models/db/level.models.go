package models

import (
	"github.com/kamva/mgm/v3"
)

// Level represents a game level
// @Description Level model
type Level struct {
	// DefaultModel contains fields for the MongoDB document
	mgm.DefaultModel `bson:",inline"`
	// Name of the level
	Name string `json:"name" bson:"name" example:"Beginner"`
}

func NewLevel(name string) *Level {
	return &Level{
		Name: name,
	}
}