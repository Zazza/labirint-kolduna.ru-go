package service

import (
	"context"
	"errors"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dice"
	"gamebook-backend/modules/game/dto"
	sectionDTO "gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/log"
	"gamebook-backend/modules/game/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiceService interface {
	RollTheDice(ctx context.Context, player entities.Player) (dto.RollTheDiceDto, error)
	GetLastDice(ctx context.Context, db *gorm.DB, playerID uuid.UUID, reason dto.ReasonType) (*entities.Dice, error)
	RemoveDice(ctx context.Context, db *gorm.DB, playerID uuid.UUID) error
}

type diceService struct {
	diceRepository repository.DiceRepository
	db             *gorm.DB
	logService     log.PlayerLogService
}

func NewDiceService(
	diceRepo repository.DiceRepository,
	db *gorm.DB,
) DiceService {
	return &diceService{
		diceRepository: diceRepo,
		db:             db,
	}
}

func NewDiceServiceWithLogging(
	diceRepo repository.DiceRepository,
	db *gorm.DB,
	logService log.PlayerLogService,
) DiceService {
	return &diceService{
		diceRepository: diceRepo,
		db:             db,
		logService:     logService,
	}
}

func (s *diceService) RollTheDice(ctx context.Context, player entities.Player) (dto.RollTheDiceDto, error) {
	battleDicesDTO := s.diceRepository.FindBattleDicesByPlayerId(ctx, s.db, player.ID)
	if battleDicesDTO.Error != nil {
		return dto.RollTheDiceDto{}, battleDicesDTO.Error
	}

	if len(player.Section.SectionEnemies) > 0 && battleDicesDTO.Exists && battleDicesDTO.Dices.DiceFirst != battleDicesDTO.Dices.DiceSecond {
		return dto.RollTheDiceDto{}, dto.MessageBattleDicesAlreadyExist
	}

	rollTheDices := dice.NewRollTheDices(s.db, &player)
	diceFirst, diceSecond, err := rollTheDices.RollTheDices(ctx, player)
	if err != nil {
		return dto.RollTheDiceDto{}, err
	}

	reason := sectionDTO.ReasonChoice
	if len(player.Section.SectionEnemies) > 0 {
		reason = sectionDTO.ReasonBattle
	}

	err = rollTheDices.StoreDices(ctx, player, *diceFirst, *diceSecond, reason)
	if err != nil {
		return dto.RollTheDiceDto{}, err
	}

	if s.logService != nil {
		s.logService.LogDiceRoll(player.ID, *diceFirst, *diceSecond, string(reason))
	}

	return dto.RollTheDiceDto{
		DiceFirst:  *diceFirst,
		DiceSecond: *diceSecond,
		Result:     sectionDTO.ResultTrue,
	}, nil
}

func (s *diceService) GetLastDice(ctx context.Context, db *gorm.DB, playerID uuid.UUID, reason dto.ReasonType) (*entities.Dice, error) {
	diceEntity, err := s.diceRepository.GetLastByPlayerId(ctx, db, playerID, reason)
	if err != nil {
		if errors.Is(err, sectionDTO.MessageDicesNotDefined) {
			return nil, sectionDTO.MessageDicesNotDefined
		}
		return nil, err
	}
	return &diceEntity, nil
}

func (s *diceService) RemoveDice(ctx context.Context, db *gorm.DB, playerID uuid.UUID) error {
	return s.diceRepository.Remove(ctx, db, playerID)
}
