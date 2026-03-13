package section

import (
	"context"
	"gamebook-backend/database/entities"
	gameDto "gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/repository"
	"gamebook-backend/modules/game/template/section"

	"gorm.io/gorm"
)

type SleepySectionService interface {
	GetSection(ctx context.Context, player entities.Player) (gameDto.CurrentResponse, error)
}

type sleepySectionService struct {
	gameRepository          repository.SectionRepository
	playerRepository        repository.PlayerRepository
	diceRepository          repository.DiceRepository
	playerSectionRepository repository.PlayerSectionRepository
	db                      *gorm.DB
}

func NewSleepySectionService(
	gameRepo repository.SectionRepository,
	playerRepo repository.PlayerRepository,
	diceRepo repository.DiceRepository,
	playerSectionRepo repository.PlayerSectionRepository,
	db *gorm.DB,
) SleepySectionService {
	return &sleepySectionService{
		gameRepository:          gameRepo,
		playerRepository:        playerRepo,
		diceRepository:          diceRepo,
		playerSectionRepository: playerSectionRepo,
		db:                      db,
	}
}

func (s *sleepySectionService) GetSection(ctx context.Context, player entities.Player) (gameDto.CurrentResponse, error) {
	sectionResponse := section.NewSectionResponse(s.db, player)
	response, err := sectionResponse.GetResponse(ctx)
	if err != nil {
		return gameDto.CurrentResponse{}, err
	}

	return response, nil
}
