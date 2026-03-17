package ability

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dice"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/listener"
	"gamebook-backend/modules/game/listener/event"
	"gamebook-backend/modules/game/repository"
	"gorm.io/gorm"
)

type MedsLogic struct {
	diceRepository       repository.DiceRepository
	playerRepository     repository.PlayerRepository
	playerUpdateListener listener.PlayerUpdateListener
	db                   *gorm.DB
}

func NewMedsLogic(
	diceRepo repository.DiceRepository,
	playerRepo repository.PlayerRepository,
	playerUpdateListener listener.PlayerUpdateListener,
) *MedsLogic {
	return &MedsLogic{
		diceRepository:       diceRepo,
		playerRepository:     playerRepo,
		playerUpdateListener: playerUpdateListener,
	}
}

func (m *MedsLogic) Validate(ctx context.Context, player entities.Player) error {
	if player.Section.Type == dto.SectionTypeSleepy {
		return dto.MessageCannotUseInSleepyKingdom
	}

	if player.Health == 0 {
		return dto.MessageCannotUseMedsIfDead
	}

	if player.Meds.Count == 0 {
		return dto.MessageNoMedsAvailable
	}

	return nil
}

func (m *MedsLogic) Execute(ctx context.Context, player entities.Player) (dto.MedsDTO, error) {
	if player.Section.Type == dto.SectionTypeSleepy {
		return dto.MedsDTO{Result: false}, dto.MessageCannotUseInSleepyKingdom
	}

	if player.Health == 0 {
		return dto.MedsDTO{Result: false}, dto.MessageCannotUseMedsIfDead
	}

	if player.Meds.Count == 0 {
		return dto.MedsDTO{
			Result: false,
		}, nil
	}

	rollTheDices := dice.NewRollTheDices(m.db, &player)
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

	player.Health = health
	player.Meds.Count--

	err = m.playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
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

func (m *MedsLogic) Result() dto.MedsResultDTO {
	return dto.MedsResultDTO{
		Result: true,
	}
}
