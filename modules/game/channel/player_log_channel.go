package channel

import (
	"gamebook-backend/modules/game/listener/event"
)

var ChPlayerLog = make(chan event.PlayerLogEvent, 100)

func init() {
	go func() {
		for event := range ChPlayerLog {
			if globalChannel != nil {
				globalChannel.SendPlayerLog(nil, event)
			}
		}
	}()
}
