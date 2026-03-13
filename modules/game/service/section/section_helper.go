package section

import (
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"
)

func BuildCurrentResponse(player entities.Player, sectionText string, gotoItems []dto.TransitionDTO, choice entities.Choice, bonuses []dto.PlayerInfoBonus, rollTheDicesNeeded bool, dataType entities.SectionType) dto.CurrentResponse {
	return dto.CurrentResponse{
		Section:      player.Section.Number,
		Text:         sectionText,
		Type:         dataType,
		Transitions:  gotoItems,
		Choice:       choice,
		RollTheDices: rollTheDicesNeeded,
		Player: dto.PlayerInfo{
			Health: player.Health,
			Meds:   uint(player.Meds.Count),
			Gold:   player.Gold,
			Bonus:  bonuses,
		},
		MapAvailable: helper.HasBagItem(player.Bag, "mapIngredients"),
	}
}
