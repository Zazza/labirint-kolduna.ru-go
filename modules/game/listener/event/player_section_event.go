package event

import (
	"github.com/google/uuid"
)

type PlayerSectionEvent struct {
	PlayerID        uuid.UUID
	TargetSectionID *uuid.UUID
	Description     *string
}

func (e PlayerSectionEvent) GetName() string {
	return "player_section"
}
