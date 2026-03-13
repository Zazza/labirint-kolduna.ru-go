package channel

import (
	"gamebook-backend/modules/game/listener/event"
)

var ChPlayerSection = make(chan event.PlayerSectionEvent, 100)

func init() {
	go func() {
		for event := range ChPlayerSection {
			if globalChannel != nil {
				globalChannel.SendPlayerSectionUpdate(nil, event)
			}
		}
	}()
}
