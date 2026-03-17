package service

import (
	"context"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/listener"
	"gamebook-backend/modules/game/repository"
	"gamebook-backend/modules/game/service/ability"

	"gorm.io/gorm"
)

type AbilityService interface {
	Meds(ctx context.Context, player entities.Player) (dto.MedsDTO, error)
	Bonus(ctx context.Context, req dto.BonusRequest, player entities.Player) (dto.BonusDTO, error)
	Sleep(ctx context.Context, player entities.Player) (dto.SleepDTO, error)
	SleepChoice(ctx context.Context, player entities.Player) (dto.SleepDTO, error)
	Bribe(ctx context.Context, player entities.Player) (dto.BribeDTO, error)
	RollTheDices(ctx context.Context, player entities.Player) (dto.RollTheDiceDto, error)
}

type abilityService struct {
	playerRepository             repository.PlayerRepository
	diceRepository               repository.DiceRepository
	sectionRepository            repository.SectionRepository
	battleRepository             repository.BattleRepository
	transitionRepository         repository.TransitionRepository
	playerSectionRepository      repository.PlayerSectionRepository
	playerSectionEnemyRepository repository.PlayerSectionEnemyRepository
	bonusRepository              repository.BonusRepository
	db                           *gorm.DB
	playerUpdateListener         listener.PlayerUpdateListener
	abilityFactory               *ability.AbilityFactory
}

func NewAbilityService(
	playerRepo repository.PlayerRepository,
	diceRepo repository.DiceRepository,
	sectionRepo repository.SectionRepository,
	battleRepo repository.BattleRepository,
	transitionRepo repository.TransitionRepository,
	playerSectionRepo repository.PlayerSectionRepository,
	playerSectionEnemyRepo repository.PlayerSectionEnemyRepository,
	bonusRepo repository.BonusRepository,
	db *gorm.DB,
) AbilityService {
	playerUpdateListener, _ := listener.HandleEvent(db, "player_update")

	abilityFactory := ability.NewAbilityFactory(
		db,
		diceRepo,
		playerRepo,
		playerUpdateListener,
		bonusRepo,
		sectionRepo,
		playerSectionRepo,
		battleRepo,
		playerSectionEnemyRepo,
	)

	return &abilityService{
		playerRepository:             playerRepo,
		diceRepository:               diceRepo,
		sectionRepository:            sectionRepo,
		battleRepository:             battleRepo,
		transitionRepository:         transitionRepo,
		playerSectionRepository:      playerSectionRepo,
		playerSectionEnemyRepository: playerSectionEnemyRepo,
		bonusRepository:              bonusRepo,
		db:                           db,
		playerUpdateListener:         playerUpdateListener,
		abilityFactory:               abilityFactory,
	}
}

func (s *abilityService) Meds(ctx context.Context, player entities.Player) (dto.MedsDTO, error) {
	medsAbility := s.abilityFactory.CreateMedsAbility()
	return medsAbility.Execute(ctx, player)
}

func (s *abilityService) Bonus(ctx context.Context, req dto.BonusRequest, player entities.Player) (dto.BonusDTO, error) {
	bonusAbility := s.abilityFactory.CreateBonusAbility()
	err := bonusAbility.Execute(ctx, player, req)
	if err != nil {
		return dto.BonusDTO{}, err
	}
	return bonusAbility.Result(), nil
}

func (s *abilityService) Sleep(ctx context.Context, player entities.Player) (dto.SleepDTO, error) {
	sleepAbility := s.abilityFactory.CreateSleepAbility()
	return sleepAbility.Execute(ctx, player)
}

func (s *abilityService) SleepChoice(ctx context.Context, player entities.Player) (dto.SleepDTO, error) {
	return dto.SleepDTO{}, fmt.Errorf("SleepChoice functionality moved to separate method")
}

func (s *abilityService) Bribe(ctx context.Context, player entities.Player) (dto.BribeDTO, error) {
	bribeAbility := s.abilityFactory.CreateBribeAbility()
	err := bribeAbility.Execute(ctx, player, dto.BribeRequest{})
	if err != nil {
		return dto.BribeDTO{}, err
	}
	return bribeAbility.Result(), nil
}

func (s *abilityService) RollTheDices(ctx context.Context, player entities.Player) (dto.RollTheDiceDto, error) {
	diceAbility := s.abilityFactory.CreateDiceAbility()
	diceDTO, err := diceAbility.Execute(ctx, player)
	return dto.RollTheDiceDto{
		DiceFirst:  diceDTO.DiceFirst,
		DiceSecond: diceDTO.DiceSecond,
		Result:     diceDTO.Result,
	}, err
}
