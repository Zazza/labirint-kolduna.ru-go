package channel

import (
	"context"
	"gamebook-backend/modules/game/listener/event"
	"sync"
)

type InMemoryEventChannel struct {
	playerUpdateCh        chan event.PlayerUpdateEvent
	enemyUpdateCh         chan event.EnemyUpdateEvent
	playerSectionUpdateCh chan event.PlayerSectionEvent
	playerLogCh           chan event.PlayerLogEvent
	listener              EventChannel
	wg                    sync.WaitGroup
	ctx                   context.Context
	cancel                context.CancelFunc
}

func NewInMemoryEventChannel(listener EventChannel) EventChannel {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	ch := &InMemoryEventChannel{
		playerUpdateCh:        make(chan event.PlayerUpdateEvent, 100),
		enemyUpdateCh:         make(chan event.EnemyUpdateEvent, 100),
		playerSectionUpdateCh: make(chan event.PlayerSectionEvent, 100),
		playerLogCh:           make(chan event.PlayerLogEvent, 100),
		listener:              listener,
		wg:                    wg,
		ctx:                   ctx,
		cancel:                cancel,
	}

	ch.startConsumers()
	return ch
}

func (c *InMemoryEventChannel) startConsumers() {
	c.wg.Add(4)

	go func() {
		defer c.wg.Done()
		for event := range c.playerUpdateCh {
			if c.listener != nil {
				c.listener.SendPlayerUpdate(c.ctx, event)
			}
		}
	}()

	go func() {
		defer c.wg.Done()
		for event := range c.enemyUpdateCh {
			if c.listener != nil {
				c.listener.SendEnemyUpdate(c.ctx, event)
			}
		}
	}()

	go func() {
		defer c.wg.Done()
		for event := range c.playerSectionUpdateCh {
			if c.listener != nil {
				c.listener.SendPlayerSectionUpdate(c.ctx, event)
			}
		}
	}()

	go func() {
		defer c.wg.Done()
		for event := range c.playerLogCh {
			if c.listener != nil {
				c.listener.SendPlayerLog(c.ctx, event)
			}
		}
	}()
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

func (c *InMemoryEventChannel) SendPlayerLog(_ context.Context, event event.PlayerLogEvent) error {
	c.playerLogCh <- event
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

func (c *InMemoryEventChannel) SubscribePlayerLog(_ context.Context) <-chan event.PlayerLogEvent {
	return c.playerLogCh
}

func (c *InMemoryEventChannel) Close() error {
	c.cancel()
	c.wg.Wait()
	close(c.playerUpdateCh)
	close(c.enemyUpdateCh)
	close(c.playerSectionUpdateCh)
	close(c.playerLogCh)
	return nil
}
