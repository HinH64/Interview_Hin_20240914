package models

import (
	"time"

	"Interview_Hin_20240914/enums"

	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

var statusRule = []validation.Rule{
	validation.Required,
	validation.In(enums.RoomStatusAvailable, enums.RoomStatusUnavailable, enums.RoomStatusMaintenance).Error("invalid status"),
}
type LevelRequest struct {
	Name   string `json:"name"`
}

func (a LevelRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required),
	)
}

type PlayerRequest struct {
	Name    string  `json:"name" binding:"required"`
	LevelID string  `json:"levelId" binding:"required"`
	Balance float64 `json:"balance"`
	ChallengeCount int `json:"challengeCount"`
}

func (p PlayerRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.LevelID, validation.Required, is.MongoID),
	)
}

type RoomRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      enums.RoomStatus  `json:"status"`
}

func (r RoomRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Status, statusRule...),
	)
}

type ReservationRequest struct {
	RoomID           string    `json:"roomId"`
	PlayerIDs        []string  `json:"playerIds"`
	BookingDatetime  string    `json:"bookingDatetime"`
}

func (r ReservationRequest) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.RoomID, validation.Required, is.MongoID),
		validation.Field(&r.PlayerIDs, validation.Required, validation.Each(is.MongoID)),
		validation.Field(&r.BookingDatetime, validation.Required, validation.Date("2006/01/02 15:04")),
	)
	if err != nil {
		return err
	}

	// Parse the date string
	_, err = time.Parse("2006/01/02 15:04", r.BookingDatetime)
	if err != nil {
		return errors.New("invalid date format for bookingDatetime")
	}

	return nil
}

type ChallengeRequest struct {
	PlayerID string `json:"playerId"`
}

func (c ChallengeRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.PlayerID, validation.Required, is.MongoID),
	)
}

type CreateChallengeGameRequest struct {
	ChallengeTimeDurationSec int     `json:"challengeTimeDurationSec" binding:"required"`
	ChallengeTimeCostSec int     `json:"challengeTimeCostSec" binding:"required"`
	ChallengeCost        float64 `json:"challengeCost" binding:"required"`
	WinProbability       float64 `json:"winProbability" binding:"required"`
	AddWinProbability    float64 `json:"addWinProbability" binding:"required"`
}

type UpdateChallengeGameRequest struct {
	ChallengeTimeDurationSec *int     `json:"challengeTimeDurationSec,omitempty"`
	ChallengeTimeCostSec *int     `json:"challengeTimeCostSec,omitempty"`
	ChallengeCost        *float64 `json:"challengeCost,omitempty"`
	WinProbability       *float64 `json:"winProbability,omitempty"`
	AddWinProbability    *float64 `json:"addWinProbability,omitempty"`
}

type LogRequest struct {
	PlayerID string      `json:"playerId" binding:"required"`
	Action   enums.LogAction `json:"action" binding:"required"`
	Details  string      `json:"details"`
}

func (l LogRequest) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.PlayerID, validation.Required, is.MongoID),
		validation.Field(&l.Action, validation.Required, validation.By(func(value interface{}) error {
			if !value.(enums.LogAction).IsValid() {
				return errors.New("invalid action")
			}
			return nil
		})),
		validation.Field(&l.Details, validation.Required),
	)
}

type PaymentRequest struct {
	PlayerID      string              `json:"playerId" binding:"required"`
	PaymentMethod enums.PaymentMethod `json:"paymentMethod" binding:"required"`
	Amount        float64             `json:"amount" binding:"required"`
}

func (p PaymentRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.PlayerID, validation.Required, is.MongoID),
		validation.Field(&p.PaymentMethod, validation.Required, validation.By(func(value interface{}) error {
			if !value.(enums.PaymentMethod).IsValid() {
				return errors.New("invalid payment method")
			}
			return nil
		})),
		validation.Field(&p.Amount, validation.Required),
	)
}