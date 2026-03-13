package bribe

import (
	"gamebook-backend/database/entities"
	gameDTO "gamebook-backend/modules/game/dto"

	"github.com/google/uuid"
)

func IsPossible(section entities.Section) bool {
	if section.Bribe != nil {
		return true
	}

	return false
}

func GetBribeTransition(sectionID uuid.UUID) gameDTO.TransitionDTO {
	return gameDTO.TransitionDTO{
		Text:         "Попытаться дать взятку",
		TransitionID: sectionID,
		Bribe:        true,
	}
}
