package channel

import (
	"gamebook-backend/modules/game/listener/event"
)

var ChEnemyUpdate = make(chan event.EnemyUpdateEvent, 100)

func init() {
	go func() {
		for event := range ChEnemyUpdate {
			if globalChannel != nil {
				globalChannel.SendEnemyUpdate(nil, event)
			}
		}
	}()
}
