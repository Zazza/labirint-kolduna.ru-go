package service

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/bonus"
	"gamebook-backend/modules/game/dto"

	"gorm.io/gorm"
)

type PlayerProfileService interface {
	GetProfile(ctx context.Context, db *gorm.DB, player entities.Player) (dto.ProfileResponse, error)
}

type playerProfileService struct {
	bonusService BonusService
}

func NewPlayerProfileService(
	bonusSvc BonusService,
) PlayerProfileService {
	return &playerProfileService{
		bonusService: bonusSvc,
	}
}

func (s *playerProfileService) GetProfile(ctx context.Context, db *gorm.DB, player entities.Player) (dto.ProfileResponse, error) {
	weapons := s.convertWeaponsToDTO(player.Weapons)
	debuffs := s.convertDebuffsToDTO(player.Debuff)
	buffs := s.convertBuffsToDTO(player.Buff)
	bonuses := s.bonusService.ConvertBonusesToDTO(player.Bonus)

	return dto.ProfileResponse{
		Health:    player.Health,
		MaxHealth: player.HealthMax,
		Meds: dto.ProfileMeds{
			Count: player.Meds.Count,
			Name:  player.Meds.Name,
		},
		Weapons: weapons,
		Bag:     player.Bag,
		Debuff:  debuffs,
		Buff:    buffs,
		Gold:    player.Gold,
		Bonus:   bonuses,
	}, nil
}

func (s *playerProfileService) convertWeaponsToDTO(weapons []entities.Weapons) []dto.ProfileWeapons {
	var result []dto.ProfileWeapons
	for _, item := range weapons {
		result = append(result, dto.ProfileWeapons{
			Name:       item.Name,
			Damage:     item.Damage,
			MinCubeHit: item.MinCubeHit,
			Item:       item.Item,
			Count:      item.Count,
		})
	}
	return result
}

func (s *playerProfileService) convertDebuffsToDTO(debuffs []entities.Debuff) []dto.ProfileDebuff {
	var result []dto.ProfileDebuff
	var bonusName string
	for _, item := range debuffs {
		bonusName = bonus.GetBonusNameByAlias(string(item.Alias))
		result = append(result, dto.ProfileDebuff{
			Health:     item.Health,
			MinCubeHit: item.MinCubeHit,
			Duration:   item.Duration,
			Alias:      &item.Alias,
			Name:       &bonusName,
		})
	}
	return result
}

func (s *playerProfileService) convertBuffsToDTO(buffs []entities.Buff) []dto.ProfileBuff {
	var result []dto.ProfileBuff
	var bonusName string
	for _, item := range buffs {
		bonusName = bonus.GetBonusNameByAlias(string(item.Alias))
		result = append(result, dto.ProfileBuff{
			Health:     item.Health,
			MinCubeHit: item.MinCubeHit,
			Duration:   item.Duration,
			Alias:      &item.Alias,
			Name:       &bonusName,
		})
	}
	return result
}
