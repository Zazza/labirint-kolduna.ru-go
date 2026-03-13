package channel

import (
	"context"
	"gamebook-backend/modules/game/listener/event"
)

type InMemoryEventChannel struct {
	playerUpdateCh        chan event.PlayerUpdateEvent
	enemyUpdateCh         chan event.EnemyUpdateEvent
	playerSectionUpdateCh chan event.PlayerSectionEvent
	playerLogCh           chan event.PlayerLogEvent
}

func NewInMemoryEventChannel() EventChannel {
	return &InMemoryEventChannel{
		playerUpdateCh:        make(chan event.PlayerUpdateEvent, 100),
		enemyUpdateCh:         make(chan event.EnemyUpdateEvent, 100),
		playerSectionUpdateCh: make(chan event.PlayerSectionEvent, 100),
		playerLogCh:           make(chan event.PlayerLogEvent, 100),
	}
}

func (c *InMemoryEventChannel) SendPlayerUpdate(_ context.Context, event event.PlayerUpdateEvent) error {
	c.playerUpdateCh <- event
	return nil
}

func (c *InMemoryEventChannel) SendEnemyUpdate(_ context.Context, event event.EnemyUpdateEvent) error {
	c.enemyUpdateCh <- event
	return nil
}

func (c *InMemoryEventChannel) SendPlayerSectionUpdate(_ context.Context, event event.PlayerSectionEvent) error {
	c.playerSectionUpdateCh <- event
	return nil
}

func (c *InMemoryEventChannel) SubscribePlayerUpdate(_ context.Context) <-chan event.PlayerUpdateEvent {
	return c.playerUpdateCh
}

func (c *InMemoryEventChannel) SubscribeEnemyUpdate(_ context.Context) <-chan event.EnemyUpdateEvent {
	return c.enemyUpdateCh
}

func (c *InMemoryEventChannel) SubscribePlayerSectionUpdate(_ context.Context) <-chan event.PlayerSectionEvent {
	return c.playerSectionUpdateCh
}

func (c *InMemoryEventChannel) SendPlayerLog(_ context.Context, event event.PlayerLogEvent) error {
	c.playerLogCh <- event
	return nil
}

func (c *InMemoryEventChannel) SubscribePlayerLog(_ context.Context) <-chan event.PlayerLogEvent {
	return c.playerLogCh
}
