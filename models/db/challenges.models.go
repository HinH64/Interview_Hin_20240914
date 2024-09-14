package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChallengeGame struct {
	mgm.DefaultModel `bson:",inline"`
	ChallengeTimeDurationSec int     `json:"challengeTimeDurationSec" bson:"challengeTimeDurationSec"`
	ChallengeTimeCostSec     int     `json:"challengeTimeCostSec" bson:"challengeTimeCostSec"`
	ChallengeCost            float64 `json:"challengeCost" bson:"challengeCost"`
	WinProbability            float64 `json:"winProbability" bson:"winProbability"`
	AddWinProbability    float64 `json:"addWinProbability" bson:"addWinProbability"`
	TotalPrizeMoney      float64 `json:"totalPrizeMoney" bson:"totalPrizeMoney" default:"0"`
}

type Challenge struct {
	mgm.DefaultModel `bson:",inline"`
	PlayerID         primitive.ObjectID `json:"playerId" bson:"playerId"`
	ChallengeGameID  primitive.ObjectID `json:"challengeGameId" bson:"challengeGameId"`
	WinPrize         bool               `json:"winPrize" bson:"winPrize"`
	ChallengeCost    float64            `json:"challengeCost" bson:"challengeCost"`
}