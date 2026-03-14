package channel

import (
	"gamebook-backend/modules/game/listener/event"
)

var ChPlayerSection chan event.PlayerSectionEvent

func init() {
	ChPlayerSection = make(chan event.PlayerSectionEvent, 100)
}
