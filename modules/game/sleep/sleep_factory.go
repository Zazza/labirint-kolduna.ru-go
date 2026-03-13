package sleep

import (
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/repository"

	"gorm.io/gorm"
)

type SleepFactory interface {
	GetSection(db *gorm.DB, player entities.Player, number uint) (dto.Sleep, error)
	IsCustom(number uint) bool
}

type sleepFactory struct {
	registry SleepRegistry
}

func NewSleepFactory(
	playerSectionRepo repository.PlayerSectionRepository,
) SleepFactory {
	return &sleepFactory{
		registry: NewSleepRegistry(playerSectionRepo),
	}
}

func (f *sleepFactory) GetSection(db *gorm.DB, player entities.Player, number uint) (dto.Sleep, error) {
	return f.registry.GetSection(db, player, number)
}

func (f *sleepFactory) IsCustom(number uint) bool {
	_, err := f.registry.Get(number)
	return err == nil
}

func GetSection(db *gorm.DB, player entities.Player, number uint) (dto.Sleep, error) {
	factory := NewSleepFactory(repository.NewPlayerSectionRepository(db))
	return factory.GetSection(db, player, number)
}
