package listener

import (
	"context"
	"encoding/json"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/listener/event"
	"gamebook-backend/modules/game/repository"

	"gorm.io/gorm"
)

type PlayerLogListener interface {
	Update(ctx context.Context)
}

type playerLogListener struct {
	chPlayerLog         chan event.PlayerLogEvent
	db                  *gorm.DB
	playerLogRepository repository.PlayerLogRepository
}

func NewPlayerLogListener(
	chPlayerLog chan event.PlayerLogEvent,
	db *gorm.DB,
	playerLogRepo repository.PlayerLogRepository,
) PlayerLogListener {
	return &playerLogListener{
		chPlayerLog:         chPlayerLog,
		db:                  db,
		playerLogRepository: playerLogRepo,
	}
}

func (l *playerLogListener) Update(ctx context.Context) {
	for {
		eventPlayerLog := <-l.chPlayerLog

		detailsJSON, err := json.Marshal(eventPlayerLog.Details)
		if err != nil {
			continue
		}

		logEntity := entities.PlayerLog{
			PlayerID:    eventPlayerLog.PlayerID,
			ActionType:  eventPlayerLog.ActionType,
			Description: eventPlayerLog.Description,
			Details:     string(detailsJSON),
		}

		_ = l.playerLogRepository.Create(ctx, l.db, &logEntity)
	}
}
