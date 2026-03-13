package custom

import (
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"

	"gorm.io/gorm"
)

func IsCustom(number uint) bool {
	registry := NewCustomSectionRegistry()
	return registry.IsCustom(number)
}

func GetSection(db *gorm.DB, player entities.Player, number uint) (dto.CustomSection, error) {
	registry := NewCustomSectionRegistry()
	return registry.GetSection(db, player, number)
}
