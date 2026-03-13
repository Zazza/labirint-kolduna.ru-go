package helpers

import (
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
)

func WeaponsToDTO(weapons []entities.Weapons) []dto.ProfileWeapons {
	var result []dto.ProfileWeapons
	for _, weapon := range weapons {
		result = append(result, dto.ProfileWeapons{
			Name:       weapon.Name,
			Damage:     weapon.Damage,
			MinCubeHit: weapon.MinCubeHit,
			Item:       weapon.Item,
			Count:      weapon.Count,
		})
	}
	return result
}

func DebuffsToDTO(debuffs []entities.Debuff) []dto.ProfileDebuff {
	var result []dto.ProfileDebuff
	for _, debuff := range debuffs {
		result = append(result, dto.ProfileDebuff{
			Health:     debuff.Health,
			MinCubeHit: debuff.MinCubeHit,
			Duration:   debuff.Duration,
			Alias:      &debuff.Alias,
		})
	}
	return result
}

func BuffsToDTO(buffs []entities.Buff) []dto.ProfileBuff {
	var result []dto.ProfileBuff
	for _, buff := range buffs {
		result = append(result, dto.ProfileBuff{
			Health:     buff.Health,
			MinCubeHit: buff.MinCubeHit,
			Duration:   buff.Duration,
			Alias:      &buff.Alias,
		})
	}
	return result
}
