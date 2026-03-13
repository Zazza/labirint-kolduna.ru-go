package validation

import (
	"errors"
)

type BonusValidator interface {
	ValidateBonusAlias(alias string) error
	ValidateBonusRequest(alias string) error
}

type bonusValidator struct{}

func NewBonusValidator() BonusValidator {
	return &bonusValidator{}
}

func (v *bonusValidator) ValidateBonusAlias(alias string) error {
	if alias == "" {
		return errors.New("bonus alias is required")
	}
	return nil
}

func (v *bonusValidator) ValidateBonusRequest(alias string) error {
	return v.ValidateBonusAlias(alias)
}
