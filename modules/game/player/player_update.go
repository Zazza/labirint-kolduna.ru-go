package player

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/listener/event"
	"gamebook-backend/modules/game/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Update interface {
	Update(ctx context.Context, eventPlayerUpdate event.PlayerUpdateEvent) (*entities.Player, error)
}

type playerUpdate struct {
	db                *gorm.DB
	playerID          uuid.UUID
	playerRepository  repository.PlayerRepository
	sectionRepository repository.SectionRepository
}

func NewPlayerUpdate(
	db *gorm.DB,
	playerID uuid.UUID,
) Update {
	playerRepository := repository.NewPlayerRepository(db)
	sectionRepository := repository.NewSectionRepository(db)
	return &playerUpdate{
		db:                db,
		playerID:          playerID,
		playerRepository:  playerRepository,
		sectionRepository: sectionRepository,
	}
}

func NewPlayerUpdateWithRepositories(
	playerRepository repository.PlayerRepository,
	sectionRepository repository.SectionRepository,
	playerID uuid.UUID,
) Update {
	return &playerUpdate{
		playerID:          playerID,
		playerRepository:  playerRepository,
		sectionRepository: sectionRepository,
	}
}

func (l *playerUpdate) Update(ctx context.Context, eventPlayerUpdate event.PlayerUpdateEvent) (*entities.Player, error) {
	player, err := l.playerRepository.GetByPlayerId(ctx, l.db, l.playerID)
	if err != nil {
		return nil, err
	}

	playerChanged, err := l.playerUpdateActions(ctx, &player, eventPlayerUpdate)
	if err != nil {
		return nil, err
	}

	_, err = l.playerRepository.Update(ctx, l.db, *playerChanged)
	if err != nil {
		return nil, err
	}

	return playerChanged, nil
}

func (l *playerUpdate) playerUpdateActions(
	ctx context.Context,
	player *entities.Player,
	eventPlayerUpdate event.PlayerUpdateEvent,
) (*entities.Player, error) {
	if eventPlayerUpdate.Health != nil {
		player.Health = *eventPlayerUpdate.Health
	}

	if eventPlayerUpdate.HealthMax != nil {
		player.HealthMax = *eventPlayerUpdate.HealthMax
	}

	if eventPlayerUpdate.SectionID != nil {
		player.SectionID = *eventPlayerUpdate.SectionID
	}

	if eventPlayerUpdate.ReturnToSection != nil {
		section, err := l.sectionRepository.GetByID(ctx, l.db, *eventPlayerUpdate.ReturnToSection)
		if err != nil {
			return nil, err
		}

		player.ReturnToSection = section.Number
	}

	if eventPlayerUpdate.Meds != nil {
		player.Meds = *eventPlayerUpdate.Meds
	}

	if eventPlayerUpdate.Bag != nil {
		player.Bag = *eventPlayerUpdate.Bag
	}
	if eventPlayerUpdate.Weapons != nil {
		player.Weapons = *eventPlayerUpdate.Weapons
	}
	if eventPlayerUpdate.Debuff != nil {
		player.Debuff = *eventPlayerUpdate.Debuff
	}
	if eventPlayerUpdate.Buff != nil {
		player.Buff = *eventPlayerUpdate.Buff
	}

	if eventPlayerUpdate.BonusList != nil {
		player.Bonus = *eventPlayerUpdate.BonusList
	}
	if eventPlayerUpdate.Bonus != nil {
		for _, bonus := range *eventPlayerUpdate.Bonus {
			player.Bonus = append(player.Bonus, bonus)
		}
	}
	if eventPlayerUpdate.Gold != nil {
		player.Gold = *eventPlayerUpdate.Gold
	}

	return player, nil
}
