package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type LevelRequest struct {
	Name   string `json:"name"`
}

func (a LevelRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required),
	)
}

type PlayerRequest struct {
	Name    string `json:"name"`
	LevelID string `json:"levelId"`
}

func (p PlayerRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.LevelID, validation.Required, is.MongoID),
	)
}