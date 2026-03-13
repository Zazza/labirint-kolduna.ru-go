package sleep

import (
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/repository"

	"gorm.io/gorm"
)

type SleepRegistry interface {
	Register(number uint, factory SleepFactoryFunc) error
	Get(number uint) (SleepFactoryFunc, error)
	GetSection(db *gorm.DB, player entities.Player, number uint) (dto.Sleep, error)
}

type SleepFactoryFunc func(db *gorm.DB, player entities.Player) dto.Sleep

type sleepRegistry struct {
	factories  map[uint]SleepFactoryFunc
	repository repository.PlayerSectionRepository
}

func NewSleepRegistry(
	playerSectionRepo repository.PlayerSectionRepository,
) SleepRegistry {
	registry := &sleepRegistry{
		factories:  make(map[uint]SleepFactoryFunc),
		repository: playerSectionRepo,
	}

	registry.RegisterDefaults()

	return registry
}

func (r *sleepRegistry) Register(number uint, factory SleepFactoryFunc) error {
	r.factories[number] = factory
	return nil
}

func (r *sleepRegistry) Get(number uint) (SleepFactoryFunc, error) {
	factory, exists := r.factories[number]
	if !exists {
		return nil, dto.MessageSleepyKingdomSectionNotDefined
	}
	return factory, nil
}

func (r *sleepRegistry) GetSection(db *gorm.DB, player entities.Player, number uint) (dto.Sleep, error) {
	factory, err := r.Get(number)
	if err != nil {
		return nil, err
	}
	return factory(db, player), nil
}

func (r *sleepRegistry) RegisterDefaults() {
	r.factories[2] = NewSleep2
	r.factories[3] = NewSleep3
	r.factories[4] = NewSleep4
	r.factories[5] = NewSleep5
	r.factories[6] = NewSleep6
	r.factories[7] = NewSleep7
	r.factories[8] = NewSleep8
	r.factories[9] = NewSleep9
	r.factories[10] = NewSleep10
	r.factories[11] = NewSleep11
	r.factories[12] = NewSleep12
}
