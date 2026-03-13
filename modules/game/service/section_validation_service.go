package service

import (
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"

	"github.com/google/uuid"
)

type SectionValidationService interface {
	ValidateSection(player entities.Player) error
	ValidateTransition(transition *entities.Transition) error
	ValidateDiceRequest(req dto.ActionRequest) error
}

type sectionValidationService struct{}

func NewSectionValidationService() SectionValidationService {
	return &sectionValidationService{}
}

func (s *sectionValidationService) ValidateSection(player entities.Player) error {
	if player.Section.ID == uuid.Nil {
		return dto.ErrSectionNotFound
	}
	return nil
}

func (s *sectionValidationService) ValidateTransition(transition *entities.Transition) error {
	if transition.ID == uuid.Nil {
		return dto.ErrSectionNotFound
	}
	return nil
}

func (s *sectionValidationService) ValidateDiceRequest(req dto.ActionRequest) error {
	return nil
}
