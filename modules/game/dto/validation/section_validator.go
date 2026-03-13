package validation

import (
	"errors"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"

	"github.com/google/uuid"
)

type SectionValidator interface {
	ValidateSection(player entities.Player) error
	ValidateTransition(transition *entities.Transition) error
	ValidateDiceRequest(req dto.ActionRequest) error
}

type sectionValidator struct{}

func NewSectionValidator() SectionValidator {
	return &sectionValidator{}
}

func (v *sectionValidator) ValidateSection(player entities.Player) error {
	if player.Section.ID == (uuid.UUID{}) {
		return dto.ErrSectionNotFound
	}
	return nil
}

func (v *sectionValidator) ValidateTransition(transition *entities.Transition) error {
	if transition == nil {
		return errors.New("transition is required")
	}
	if transition.ID == (uuid.UUID{}) {
		return errors.New("transition ID is required")
	}
	return nil
}

func (v *sectionValidator) ValidateDiceRequest(req dto.ActionRequest) error {
	return nil
}
