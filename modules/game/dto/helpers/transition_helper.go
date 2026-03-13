package helpers

import (
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"
)

func TransitionToDTO(transition entities.Transition, sectionNumber uint, playerSections []entities.PlayerSection) dto.TransitionDTO {
	visited := helper.IsVisited(sectionNumber, playerSections)
	return dto.TransitionDTO{
		Text:         transition.Text,
		TransitionID: transition.ID,
		Visited:      visited,
	}
}

func TransitionsToDTO(transitions []entities.Transition, sectionNumbers []uint, playerSections []entities.PlayerSection) []dto.TransitionDTO {
	var result []dto.TransitionDTO

	for i, transition := range transitions {
		sectionNumber := sectionNumbers[i]
		result = append(result, TransitionToDTO(transition, sectionNumber, playerSections))
	}
	return result
}
