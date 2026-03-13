package helpers

import (
	"gamebook-backend/database/entities"
)

func HasBuff(buffList []entities.Buff, bonus entities.BuffOrDebuffAlias) bool {
	for _, item := range buffList {
		if item.Alias == bonus {
			return true
		}
	}

	return false
}

func HasDebuff(debuffList []entities.Debuff, bonus entities.BuffOrDebuffAlias) bool {
	for _, item := range debuffList {
		if item.Alias == bonus {
			return true
		}
	}

	return false
}

func RemoveDebuff(debuffList []entities.Debuff, debuffAlias entities.BuffOrDebuffAlias) []entities.Debuff {
	var result []entities.Debuff
	for _, item := range debuffList {
		if item.Alias != debuffAlias {
			result = append(result, item)
		}
	}

	return result
}

func RemoveBuff(buffList []entities.Buff, buffAlias entities.BuffOrDebuffAlias) []entities.Buff {
	var result []entities.Buff
	for _, item := range buffList {
		if item.Alias != buffAlias {
			result = append(result, item)
		}
	}

	return result
}

func RemoveBonus(bonusList []entities.PlayerBonus, bonusAlias string) []entities.PlayerBonus {
	var result []entities.PlayerBonus
	for _, item := range bonusList {
		if *item.Alias != bonusAlias {
			result = append(result, item)
		}
	}

	return result
}
