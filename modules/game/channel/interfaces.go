package channel

import (
	"context"
	"gamebook-backend/modules/game/listener/event"
)

type EventChannel interface {
	SendPlayerUpdate(ctx context.Context, event event.PlayerUpdateEvent) error
	SendEnemyUpdate(ctx context.Context, event event.EnemyUpdateEvent) error
	SendPlayerSectionUpdate(ctx context.Context, event event.PlayerSectionEvent) error
	SendPlayerLog(ctx context.Context, event event.PlayerLogEvent) error
	SubscribePlayerUpdate(ctx context.Context) <-chan event.PlayerUpdateEvent
	SubscribeEnemyUpdate(ctx context.Context) <-chan event.EnemyUpdateEvent
	SubscribePlayerSectionUpdate(ctx context.Context) <-chan event.PlayerSectionEvent
	SubscribePlayerLog(ctx context.Context) <-chan event.PlayerLogEvent
}
