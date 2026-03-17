package dto

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	MessageBonusNotDefined                = errors.New("message bonus not defined")
	MessageSleepyKingdomSectionNotDefined = errors.New("sleepy kingdom section not defined")
	MessageAlreadySleepyKingdom           = errors.New("already sleepy kingdom")
	MessageTheDeadNeverSleep              = errors.New("the dead never sleep")
	MessageCannotUseInSleepyKingdom       = errors.New("cannot use in sleepy kingdom")
	MessageCannotUseMedsIfDead            = errors.New("cannot use meds if dead")
	MessageNoMedsAvailable                = errors.New("no meds available")
	MessageNotSleepySectionType           = errors.New("not sleepy section type")
)

type (
	Bonus interface {
		Execute(ctx context.Context, req BonusRequest) error
	}

	Sleep interface {
		Execute(
			ctx context.Context,
			dice1 uint,
			dice2 uint,
		) (SleepyKingdomDTO, error)
	}

	MedsDTO struct {
		Result ActionResult `json:"result"`
	}

	BonusDTO struct {
		Result  ActionResult `json:"result"`
		Success bool         `json:"success"`
		Message string       `json:"message,omitempty"`
	}

	SleepDTO struct {
		Result         ActionResult `json:"result"`
		SleepyKingdom  bool         `json:"sleepy_kingdom"`
		HealthRecovery *uint        `json:"health_recovery"`
	}

	BribeDTO struct {
		Result  ActionResult `json:"result"`
		Success bool         `json:"success"`
		Message string       `json:"message,omitempty"`
	}

	DiceDTO struct {
		DiceFirst  uint         `json:"dice_first"`
		DiceSecond uint         `json:"dice_second"`
		Result     ActionResult `json:"result"`
	}

	BonusRequest struct {
		Bonus  string `json:"bonus"`
		Option string `json:"option"`
	}

	BribeRequest struct {
		Amount uint   `json:"amount,omitempty"`
		Target string `json:"target,omitempty"`
	}

	SleepyChoiceRequest struct {
		Transition uuid.UUID `json:"transitionID"`
	}

	MedsResultDTO struct {
		Result bool `json:"result"`
	}

	BonusResultDTO struct {
		Success bool `json:"success"`
	}

	SleepResultDTO struct {
		Result bool `json:"result"`
	}

	BribeResultDTO struct {
		Result bool `json:"result"`
	}

	DiceResultDTO struct {
		DiceFirst  uint `json:"dice_first"`
		DiceSecond uint `json:"dice_second"`
		Result     bool `json:"result"`
	}
)
