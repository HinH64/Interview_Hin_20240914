package models

import (
	"Interview_Hin_20240914/enums"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Log struct {
	mgm.DefaultModel `bson:",inline"`
	PlayerID         primitive.ObjectID `json:"playerId" bson:"playerId"`
	Action           enums.LogAction    `json:"action" bson:"action"`
	Details          string    `json:"details" bson:"details"`
}
