package helpers

import (
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
)

func BonusToDTO(bonuses []entities.PlayerBonus) []dto.PlayerInfoBonus {
	var result []dto.PlayerInfoBonus
	for _, playerBonus := range bonuses {
		if playerBonus.Option != nil {
			for _, option := range *playerBonus.Option {
				result = append(result, dto.PlayerInfoBonus{
					Alias:  *playerBonus.Alias,
					Name:   *playerBonus.Name,
					Option: &option,
				})
			}
		} else {
			result = append(result, dto.PlayerInfoBonus{
				Alias: *playerBonus.Alias,
				Name:  *playerBonus.Name,
			})
		}
	}
	return result
}
