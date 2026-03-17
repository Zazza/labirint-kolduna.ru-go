package listener

import (
	"context"
	"gamebook-backend/modules/game/channel"
	"gamebook-backend/modules/game/listener/event"
	"gamebook-backend/modules/game/repository"
	"gorm.io/gorm"
)

type PlayerSectionListener interface {
	Handle(ctx context.Context, e event.Event) error
}

type playerSectionListener struct {
	db                      *gorm.DB
	playerSectionRepository repository.PlayerSectionRepository
	eventChannel            channel.EventChannel
}

func NewPlayerSectionListener(
	eventChannel channel.EventChannel,
	db *gorm.DB,
) PlayerSectionListener {
	playerSectionRepository := repository.NewPlayerSectionRepository(db)
	return &playerSectionListener{
		db:                      db,
		playerSectionRepository: playerSectionRepository,
		eventChannel:            eventChannel,
	}
}

func (l *playerSectionListener) Handle(ctx context.Context, e event.Event) error {
	eventPlayerSection, ok := e.(event.PlayerSectionEvent)
	if !ok {
		return nil
	}

	if l.eventChannel != nil {
		if err := l.eventChannel.SendPlayerSectionUpdate(ctx, eventPlayerSection); err != nil {
			return err
		}
	}

	if eventPlayerSection.TargetSectionID != nil {
		err := l.playerSectionRepository.UpdateLastTargetSection(
			ctx,
			l.db,
			eventPlayerSection.PlayerID,
			*eventPlayerSection.TargetSectionID,
		)
		if err != nil {
			return err
		}
	}

	if eventPlayerSection.Description != nil {
		err := l.playerSectionRepository.AddDescriptionLog(
			ctx,
			l.db,
			eventPlayerSection.PlayerID,
			*eventPlayerSection.Description,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
