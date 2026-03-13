package helpers

import (
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
)

func PlayerInfoToDTO(player entities.Player) dto.PlayerInfo {
	return dto.PlayerInfo{
		Health: player.Health,
		Meds:   uint(player.Meds.Count),
		Gold:   player.Gold,
		Bonus:  BonusToDTO(player.Bonus),
	}
}
