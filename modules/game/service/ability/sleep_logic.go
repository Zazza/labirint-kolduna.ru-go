package ability

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dice"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/listener"
	"gamebook-backend/modules/game/listener/event"
	"gamebook-backend/modules/game/repository"
	"gamebook-backend/modules/game/sleep"
	"gorm.io/gorm"
)

type SleepLogic struct {
	db                           *gorm.DB
	sectionRepository            repository.SectionRepository
	playerUpdateListener         listener.PlayerUpdateListener
	playerSectionRepository      repository.PlayerSectionRepository
	battleRepository             repository.BattleRepository
	playerSectionEnemyRepository repository.PlayerSectionEnemyRepository
}

func NewSleepLogic(
	db *gorm.DB,
	sectionRepo repository.SectionRepository,
	playerUpdateListener listener.PlayerUpdateListener,
	playerSectionRepo repository.PlayerSectionRepository,
	battleRepo repository.BattleRepository,
	playerSectionEnemyRepo repository.PlayerSectionEnemyRepository,
) *SleepLogic {
	return &SleepLogic{
		db:                           db,
		sectionRepository:            sectionRepo,
		playerUpdateListener:         playerUpdateListener,
		playerSectionRepository:      playerSectionRepo,
		battleRepository:             battleRepo,
		playerSectionEnemyRepository: playerSectionEnemyRepo,
	}
}

func (s *SleepLogic) Validate(ctx context.Context, player entities.Player) error {
	if player.Section.Type != dto.SectionTypeSleepy {
		return dto.MessageNotSleepySectionType
	}
	return nil
}

func (s *SleepLogic) Execute(ctx context.Context, player entities.Player) (dto.SleepDTO, error) {
	return s.SleepValidate(ctx, player)
}

func (s *SleepLogic) SleepValidate(ctx context.Context, player entities.Player) (dto.SleepDTO, error) {
	if player.Section.Type != dto.SectionTypeSleepy {
		return dto.SleepDTO{}, dto.MessageNotSleepySectionType
	}

	if player.Health == 0 {
		return dto.SleepDTO{}, dto.MessageTheDeadNeverSleep
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

func (s *SleepLogic) SleepChoice(ctx context.Context, player entities.Player) (dto.SleepDTO, error) {
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

	if resultDTO.Exit {
		exit := sleep.NewExit(s.db, player)
		err := exit.Return(ctx)
		if err != nil {
			return dto.SleepDTO{}, err
		}

		helper.DescriptionMessageWithContext(
			ctx,
			player.ID,
			"<p>🏆 Успешно вернулся из сонного царства</p>",
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

		helper.DescriptionMessageWithContext(
			ctx,
			player.ID,
			"<p>💀 Погиб в сонном царстве</p>",
		)

		deathSection, err := s.sectionRepository.GetBySectionNumber(ctx, s.db, 9)
		if err != nil {
			return dto.SleepDTO{}, err
		}

		err = s.playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
			PlayerID:  player.ID,
			SectionID: &deathSection.ID,
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

func (s *SleepLogic) Result() dto.SleepResultDTO {
	return dto.SleepResultDTO{
		Result: true,
	}
}
