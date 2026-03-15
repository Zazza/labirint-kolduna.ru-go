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
		Result ActionResult `json:"result"`
	}

	SleepDTO struct {
		Result         ActionResult `json:"result"`
		SleepyKingdom  bool         `json:"sleepy_kingdom"`
		HealthRecovery *uint        `json:"health_recovery"`
	}

	BribeDTO struct {
		Result ActionResult `json:"result"`
	}

	BonusRequest struct {
		Bonus  string `json:"bonus"`
		Option string `json:"option"`
	}
	SleepyChoiceRequest struct {
		Transition uuid.UUID `json:"transitionID"`
	}
)
