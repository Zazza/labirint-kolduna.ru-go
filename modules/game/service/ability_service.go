package service

import (
	"context"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/bonus"
	"gamebook-backend/modules/game/bribe"
	"gamebook-backend/modules/game/dice"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/listener"
	"gamebook-backend/modules/game/listener/event"
	"gamebook-backend/modules/game/repository"
	"gamebook-backend/modules/game/sleep"
	template2 "gamebook-backend/modules/game/template"

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
	db                           *gorm.DB
	playerUpdateListener         listener.PlayerUpdateListener
}

func NewAbilityService(
	playerRepo repository.PlayerRepository,
	diceRepo repository.DiceRepository,
	sectionRepo repository.SectionRepository,
	battleRepo repository.BattleRepository,
	transitionRepo repository.TransitionRepository,
	playerSectionRepo repository.PlayerSectionRepository,
	playerSectionEnemyRepo repository.PlayerSectionEnemyRepository,
	db *gorm.DB,
) AbilityService {
	playerUpdateListener, _ := listener.HandleEvent(db, "player_update")

	return &abilityService{
		playerRepository:             playerRepo,
		diceRepository:               diceRepo,
		sectionRepository:            sectionRepo,
		battleRepository:             battleRepo,
		transitionRepository:         transitionRepo,
		playerSectionRepository:      playerSectionRepo,
		playerSectionEnemyRepository: playerSectionEnemyRepo,
		db:                           db,
		playerUpdateListener:         playerUpdateListener,
	}
}

func (s *abilityService) Meds(ctx context.Context, player entities.Player) (dto.MedsDTO, error) {
	if player.Section.Type == dto.SectionTypeSleepy {
		return dto.MedsDTO{Result: false}, dto.MessageCannotUseInSleepyKingdom
	}

	if player.Meds.Count == 0 {
		return dto.MedsDTO{
			Result: false,
		}, nil
	}

	rollTheDices := dice.NewRollTheDices(s.db, &player)
	diceFirst, diceSecond, err := rollTheDices.RollTheDices(ctx, player)
	if err != nil {
		return dto.MedsDTO{
			Result: false,
		}, err
	}

	health := player.Health + *diceFirst + *diceSecond
	if player.HealthMax < health {
		health = player.HealthMax
	}

	player.Health += health
	player.Meds.Count--

	err = s.playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
		PlayerID: player.ID,
		Health:   &health,
		Meds:     &player.Meds,
	})
	if err != nil {
		return dto.MedsDTO{
			Result: false,
		}, err
	}

	return dto.MedsDTO{
		Result: true,
	}, nil
}

func (s *abilityService) Bonus(ctx context.Context, req dto.BonusRequest, player entities.Player) (dto.BonusDTO, error) {
	if player.Section.Type == dto.SectionTypeSleepy {
		return dto.BonusDTO{Result: false}, dto.MessageCannotUseInSleepyKingdom
	}

	var bonusAlias string
	for _, item := range player.Bonus {
		if *item.Alias == req.Bonus {
			bonusAlias = *item.Alias
		}
	}

	if &bonusAlias == nil {
		return dto.BonusDTO{
			Result: false,
		}, nil
	}

	bonusInstance, err := bonus.GetBonus(s.db, player, bonusAlias)
	if err != nil {
		return dto.BonusDTO{}, err
	}

	err = bonusInstance.Execute(ctx, req)
	if err != nil {
		return dto.BonusDTO{}, err
	}

	return dto.BonusDTO{
		Result: true,
	}, nil
}

func (s *abilityService) Sleep(ctx context.Context, player entities.Player) (dto.SleepDTO, error) {
	if player.Section.Type == dto.SectionTypeSleepy {
		return dto.SleepDTO{}, dto.MessageAlreadySleepyKingdom
	}

	entrance := sleep.NewEntrance(s.db, player)
	err := entrance.Handle(ctx)
	if err != nil {
		return dto.SleepDTO{}, err
	}

	return dto.SleepDTO{
		Result: true,
	}, nil
}

func (s *abilityService) SleepChoice(ctx context.Context, player entities.Player) (dto.SleepDTO, error) {
	playerSection, err := s.playerSectionRepository.GetLastPlayerSection(ctx, s.db, player.ID)
	if err != nil {
		return dto.SleepDTO{}, err
	}

	emptyPlayerTargetSectionID := entities.PlayerSection{}.TargetSectionID
	if playerSection.TargetSectionID != emptyPlayerTargetSectionID {
		err = s.playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
			PlayerID:  player.ID,
			SectionID: &playerSection.TargetSectionID,
		})
		if err != nil {
			return dto.SleepDTO{}, err
		}

		return dto.SleepDTO{
			Result: true,
		}, nil
	}

	sleepSectionInstance, err := sleep.GetSection(s.db, player, player.Section.Number-200)
	if err != nil {
		return dto.SleepDTO{}, err
	}

	rollTheDices := dice.NewRollTheDices(s.db, &player)

	dice1, dice2, err := rollTheDices.RollTheDices(ctx, player)
	if err != nil {
		return dto.SleepDTO{}, err
	}

	resultDTO, err := sleepSectionInstance.Execute(
		ctx,
		*dice1,
		*dice2,
	)
	if err != nil {
		return dto.SleepDTO{}, err
	}

	template, err := template2.GetDicesTemplate(ctx, *dice1, *dice2, true)
	if err != nil {
		return dto.SleepDTO{}, err
	}

	helper.DescriptionMessage(
		player.ID,
		fmt.Sprintf("<p>%s</p>", template),
	)

	if resultDTO.Exit {
		sectionRepository := repository.NewSectionRepository(s.db)
		returnToSection, err := sectionRepository.GetBySectionNumber(ctx, s.db, player.ReturnToSection)
		if err != nil {
			return dto.SleepDTO{}, err
		}

		err = s.playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
			PlayerID:  player.ID,
			SectionID: &returnToSection.ID,
		})
		if err != nil {
			return dto.SleepDTO{}, err
		}

		helper.DescriptionMessage(
			player.ID,
			"<p>Успешно вернулся из сонного царства</p>",
		)
	}

	if resultDTO.Death {
		health := uint(0)
		err = s.playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
			PlayerID: player.ID,
			Health:   &health,
		})
		if err != nil {
			return dto.SleepDTO{}, err
		}
	}

	if resultDTO.NextTry {
		entrance := sleep.NewEntrance(s.db, player)
		err := entrance.Handle(ctx)
		if err != nil {
			return dto.SleepDTO{}, err
		}
	}

	err = s.battleRepository.RemoveSleepyByPlayerIDAndSectionNumber(ctx, s.db, player.ID, player.Section.Number)
	if err != nil {
		return dto.SleepDTO{}, err
	}

	err = s.playerSectionEnemyRepository.RemoveSleepyByPlayerIDAndSectionNumber(ctx, s.db, player.ID, player.SectionID)
	if err != nil {
		return dto.SleepDTO{}, err
	}

	err = s.playerSectionRepository.RemoveLastPlayerSection(ctx, s.db, player.ID)
	if err != nil {
		return dto.SleepDTO{}, err
	}

	return dto.SleepDTO{
		Result: true,
	}, nil
}

func (s *abilityService) Bribe(ctx context.Context, player entities.Player) (dto.BribeDTO, error) {
	rollTheDices := dice.NewRollTheDices(s.db, &player)
	diceFirst, diceSecond, err := rollTheDices.RollTheDices(ctx, player)
	if err != nil {
		return dto.BribeDTO{}, err
	}

	template, err := template2.GetDicesTemplate(ctx, *diceFirst, *diceSecond, true)
	if err != nil {
		return dto.BribeDTO{}, err
	}

	helper.DescriptionMessage(player.ID, fmt.Sprintf("<p>%s</p>", template))

	err = bribe.BribeAction(s.db, player)
	if err != nil {
		return dto.BribeDTO{}, err
	}

	err = rollTheDices.StoreDices(ctx, player, *diceFirst, *diceSecond, dto.ReasonBribe)
	if err != nil {
		return dto.BribeDTO{}, err
	}

	return dto.BribeDTO{
		Result: true,
	}, nil
}

func (s *abilityService) RollTheDices(ctx context.Context, player entities.Player) (dto.RollTheDiceDto, error) {
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

	reason := dto.ReasonChoice
	if len(player.Section.SectionEnemies) > 0 {
		reason = dto.ReasonBattle
	}

	err = rollTheDices.StoreDices(ctx, player, *diceFirst, *diceSecond, reason)
	if err != nil {
		return dto.RollTheDiceDto{}, err
	}

	return dto.RollTheDiceDto{
		DiceFirst:  *diceFirst,
		DiceSecond: *diceSecond,
		Result:     dto.ResultTrue,
	}, nil
}
