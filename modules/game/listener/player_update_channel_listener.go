package listener

import (
	"context"
	"gamebook-backend/modules/game/listener/event"
	"gamebook-backend/modules/game/repository"
	"log"

	"gorm.io/gorm"
)

type PlayerSectionChannelListener interface {
	Update(ctx context.Context)
}

type playerSectionChannelListener struct {
	chPlayerSection         chan event.PlayerSectionEvent
	db                      *gorm.DB
	playerSectionRepository repository.PlayerSectionRepository
}

func NewPlayerSectionChannelListener(
	chPlayerSection chan event.PlayerSectionEvent,
	db *gorm.DB,
	playerSectionRepo repository.PlayerSectionRepository,
) PlayerSectionChannelListener {
	return &playerSectionChannelListener{
		chPlayerSection:         chPlayerSection,
		db:                      db,
		playerSectionRepository: playerSectionRepo,
	}
}

func (l *playerSectionChannelListener) Update(ctx context.Context) {
	for {
		eventPlayerSection := <-l.chPlayerSection

		if eventPlayerSection.TargetSectionID != nil {
			err := l.playerSectionRepository.UpdateLastTargetSection(
				ctx,
				l.db,
				eventPlayerSection.PlayerID,
				*eventPlayerSection.TargetSectionID,
			)
			if err != nil {
				log.Println("[PlayerSectionListener] UpdateTargetSection err:", err)
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
				log.Println("[PlayerSectionListener] AddDescriptionLog err:", err)
			}
		}
	}
}
