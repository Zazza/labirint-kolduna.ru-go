package validation

import (
	"errors"
	"gamebook-backend/modules/game/dto"
)

type BattleValidator interface {
	ValidateBattleRequest(req dto.ActionRequest) error
	ValidateWeapon(weapon string) error
}

type battleValidator struct{}

func NewBattleValidator() BattleValidator {
	return &battleValidator{}
}

func (v *battleValidator) ValidateBattleRequest(req dto.ActionRequest) error {
	if req.Weapon == "" {
		return errors.New("weapon is required")
	}
	return v.ValidateWeapon(req.Weapon)
}

func (v *battleValidator) ValidateWeapon(weapon string) error {
	if weapon == "" {
		return errors.New("weapon is required")
	}
	return nil
}
