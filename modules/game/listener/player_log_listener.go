package listener

import (
	"context"
	"encoding/json"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/channel"
	"gamebook-backend/modules/game/repository"
	"gorm.io/gorm"
)

type PlayerLogListener interface {
	Update(ctx context.Context)
}

type playerLogListener struct {
	db                  *gorm.DB
	playerLogRepository repository.PlayerLogRepository
	eventChannel        channel.EventChannel
}

func NewPlayerLogListener(
	eventChannel channel.EventChannel,
	db *gorm.DB,
	playerLogRepo repository.PlayerLogRepository,
) PlayerLogListener {
	return &playerLogListener{
		db:                  db,
		playerLogRepository: playerLogRepo,
		eventChannel:        eventChannel,
	}
}

func (l *playerLogListener) Update(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case eventPlayerLog := <-l.eventChannel.SubscribePlayerLog(ctx):
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
}
