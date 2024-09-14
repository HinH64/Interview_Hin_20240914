package models

import (
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reservation struct {
	mgm.DefaultModel `bson:",inline"`
	RoomID           primitive.ObjectID   `json:"roomId" bson:"roomId"`
	PlayerIDs        []primitive.ObjectID `json:"playerIds" bson:"playerIds"`
	BookingDatetime  time.Time            `json:"bookingDatetime" bson:"bookingDatetime"`
}

type GetReservation struct {
	mgm.DefaultModel `bson:",inline"`
	RoomID           primitive.ObjectID   `json:"roomId" bson:"roomId"`
	RoomName         string               `json:"roomName" bson:"roomName"`
	PlayerIDs        []primitive.ObjectID `json:"playerIds" bson:"playerIds"`
	PlayerNames      []string             `json:"playerNames" bson:"playerNames"`
	BookingDatetime  time.Time            `json:"bookingDatetime" bson:"bookingDatetime"`

}