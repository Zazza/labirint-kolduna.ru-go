package bonus

import (
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/repository"

	"gorm.io/gorm"
)

type BonusFactory interface {
	GetBonus(db *gorm.DB, player entities.Player, bonusAlias string) (dto.Bonus, error)
}

type bonusFactory struct {
	registry BonusRegistry
}

func NewBonusFactory(
	bonusRepo repository.BonusRepository,
) BonusFactory {
	return &bonusFactory{
		registry: NewBonusRegistry(bonusRepo),
	}
}

func (f *bonusFactory) GetBonus(
	db *gorm.DB,
	player entities.Player,
	bonusAlias string,
) (dto.Bonus, error) {
	return f.registry.GetBonus(db, player, bonusAlias)
}

func GetBonus(
	db *gorm.DB,
	player entities.Player,
	bonusAlias string,
) (dto.Bonus, error) {
	factory := NewBonusFactory(repository.NewBonusRepository(db))
	return factory.GetBonus(db, player, bonusAlias)
}
