package listener

import (
	"context"
	"fmt"
	"gamebook-backend/modules/game/listener/event"
	"gamebook-backend/modules/game/player"
	"gamebook-backend/modules/game/repository"

	"gorm.io/gorm"
)

type PlayerUpdateListener interface {
	Handle(ctx context.Context, e event.Event) error
}

type playerListener struct {
	db                *gorm.DB
	playerRepository  repository.PlayerRepository
	sectionRepository repository.SectionRepository
}

func NewPlayerUpdateListener(
	db *gorm.DB,
) PlayerUpdateListener {
	playerRepository := repository.NewPlayerRepository(db)
	sectionRepository := repository.NewSectionRepository(db)

	return &playerListener{
		db:                db,
		playerRepository:  playerRepository,
		sectionRepository: sectionRepository,
	}
}

func (l *playerListener) Handle(ctx context.Context, e event.Event) error {
	eventPlayerUpdate, ok := e.(event.PlayerUpdateEvent)
	if !ok {
		fmt.Println("eventPlayerUpdate event is not of type PlayerUpdateEvent")
		return nil
	}

	playerEntity, err := l.playerRepository.GetByPlayerId(ctx, l.db, eventPlayerUpdate.PlayerID)
	if err != nil {
		return err
	}

	playerUpdate := player.NewPlayerUpdate(l.db, playerEntity.ID)
	_, err = playerUpdate.Update(ctx, eventPlayerUpdate)
	return err
}
