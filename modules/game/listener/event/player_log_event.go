package event

import (
	"github.com/google/uuid"
)

type PlayerLogEvent struct {
	PlayerID    uuid.UUID
	ActionType  string
	Description string
	Details     map[string]interface{}
}
