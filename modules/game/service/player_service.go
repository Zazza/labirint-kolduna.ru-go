package service

import (
	"context"
	"errors"
	"gamebook-backend/database/entities"
	playerDto "gamebook-backend/modules/game/dto"
	player2 "gamebook-backend/modules/game/player"
	"gamebook-backend/modules/game/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlayerService interface {
	GetByUserId(ctx context.Context, userId string) (entities.Player, error)
}

type playerService struct {
	playerRepository        repository.PlayerRepository
	playerSectionRepository repository.PlayerSectionRepository
	sectionRepository       repository.SectionRepository
	db                      *gorm.DB
}

func NewPlayerService(
	playerRepo repository.PlayerRepository,
	playerSectionRepo repository.PlayerSectionRepository,
	sectionRepo repository.SectionRepository,
	db *gorm.DB,
) PlayerService {
	return &playerService{
		playerRepository:        playerRepo,
		playerSectionRepository: playerSectionRepo,
		sectionRepository:       sectionRepo,
		db:                      db,
	}
}

func (s *playerService) GetByUserId(ctx context.Context, userId string) (entities.Player, error) {
	var player entities.Player
	player, err := s.playerRepository.GetByUserId(ctx, s.db, userId)
	if errors.Is(err, playerDto.ErrPlayerNotFound) {
		section, err := s.sectionRepository.GetBySectionNumber(ctx, s.db, 0)
		if err != nil {
			return entities.Player{}, err
		}

		player, err = s.playerRepository.Create(ctx, s.db, uuid.MustParse(userId), section)
		if err != nil {
			return entities.Player{}, err
		}

		playerRef, err := player2.ResetPlayer(ctx, s.db, player)
		if err != nil {
			return entities.Player{}, err

		}

		playerRef, err = s.playerRepository.Update(ctx, s.db, *playerRef)
		if err != nil {
			return entities.Player{}, err
		}
		player = *playerRef
	} else if err != nil {
		return entities.Player{}, err
	}

	return player, nil
}
