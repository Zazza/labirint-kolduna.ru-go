package dto

import (
	"errors"
	"gamebook-backend/database/entities"
)

var (
	MessageFailedRollTheDice       = errors.New("roll the dice error")
	MessageDicesNotDefined         = errors.New("dices not defined")
	MessageBattleDicesAlreadyExist = errors.New("battle dices already exist")
	MessageBattleDicesRequired     = errors.New("battle dices required")
)

const (
	ReasonBattle ReasonType = "battle"
	ReasonBribe  ReasonType = "bribe"
	ReasonChoice ReasonType = "choice"
)

type (
	ReasonType string

	BattleDicesDTO struct {
		Exists bool
		Dices  entities.Dice
		Error  error
	}

	RollTheDiceDto struct {
		DiceFirst  uint         `json:"dice_first"`
		DiceSecond uint         `json:"dice_second"`
		Result     ActionResult `json:"result"`
	}
)
