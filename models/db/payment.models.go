package models

import (
	"Interview_Hin_20240914/enums"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	mgm.DefaultModel `bson:",inline"`
	PlayerID         primitive.ObjectID  `json:"playerId" bson:"player_id"`
	PaymentMethod    enums.PaymentMethod `json:"paymentMethod" bson:"payment_method"`
	Amount           float64             `json:"amount" bson:"amount"`
	Details          string              `json:"details" bson:"details"`	
	Status           string              `json:"status" bson:"status"`
	TransactionId    string              `json:"transactionId" bson:"transaction_id"`
}