package service

import (
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/bribe"
	"gamebook-backend/modules/game/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BribeService interface {
	IsPossible(section entities.Section) bool
	IsBribe(db *gorm.DB, playerID uuid.UUID) (bool, error)
	GetBribeResultTransition(db *gorm.DB, player entities.Player) (dto.TransitionDTO, error)
	GetBribeTransition(sectionID uuid.UUID) dto.TransitionDTO
	BribeAction(db *gorm.DB, player entities.Player) error
}

type bribeService struct{}

func NewBribeService() BribeService {
	return &bribeService{}
}

func (s *bribeService) IsPossible(section entities.Section) bool {
	return bribe.IsPossible(section)
}

func (s *bribeService) IsBribe(db *gorm.DB, playerID uuid.UUID) (bool, error) {
	return bribe.IsBribe(db, playerID)
}

func (s *bribeService) GetBribeResultTransition(db *gorm.DB, player entities.Player) (dto.TransitionDTO, error) {
	transition, err := bribe.GetBribeResultTransition(db, player)
	if err != nil {
		return dto.TransitionDTO{}, err
	}
	return transition, nil
}

func (s *bribeService) GetBribeTransition(sectionID uuid.UUID) dto.TransitionDTO {
	return bribe.GetBribeTransition(sectionID)
}

func (s *bribeService) BribeAction(db *gorm.DB, player entities.Player) error {
	return bribe.BribeAction(db, player)
}
