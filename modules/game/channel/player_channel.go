package channel

import (
	"gamebook-backend/modules/game/listener/event"
	"sync"
)

var (
	globalChannel EventChannel
	once          sync.Once
)

func InitGlobalChannel(channel EventChannel) {
	once.Do(func() {
		globalChannel = channel
	})
}

var ChPlayerUpdate = make(chan event.PlayerUpdateEvent, 100)

func init() {
	go func() {
		for event := range ChPlayerUpdate {
			if globalChannel != nil {
				globalChannel.SendPlayerUpdate(nil, event)
			}
		}
	}()
}
