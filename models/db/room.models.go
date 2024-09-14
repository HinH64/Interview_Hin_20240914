package models

import (
	"Interview_Hin_20240914/enums"

	"github.com/kamva/mgm/v3"
)

type Room struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string            `json:"name" bson:"name"`
	Description      string            `json:"description" bson:"description"`
	Status           enums.RoomStatus  `json:"status" bson:"status"`
}