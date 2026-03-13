package custom

import (
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"

	"gorm.io/gorm"
)

type CustomSectionRegistry interface {
	Register(number uint, factory CustomSectionFactoryFunc) error
	Get(number uint) (CustomSectionFactoryFunc, error)
	GetSection(db *gorm.DB, player entities.Player, number uint) (dto.CustomSection, error)
	IsCustom(number uint) bool
}

type CustomSectionFactoryFunc func(db *gorm.DB, player entities.Player) dto.CustomSection

type customSectionRegistry struct {
	factories map[uint]CustomSectionFactoryFunc
}

func NewCustomSectionRegistry() CustomSectionRegistry {
	registry := &customSectionRegistry{
		factories: make(map[uint]CustomSectionFactoryFunc),
	}

	registry.RegisterDefaults()

	return registry
}

func (r *customSectionRegistry) Register(number uint, factory CustomSectionFactoryFunc) error {
	r.factories[number] = factory
	return nil
}

func (r *customSectionRegistry) Get(number uint) (CustomSectionFactoryFunc, error) {
	factory, exists := r.factories[number]
	if !exists {
		return nil, dto.ErrCustomSectionNotFound
	}
	return factory, nil
}

func (r *customSectionRegistry) GetSection(db *gorm.DB, player entities.Player, number uint) (dto.CustomSection, error) {
	factory, err := r.Get(number)
	if err != nil {
		return nil, err
	}
	return factory(db, player), nil
}

func (r *customSectionRegistry) IsCustom(number uint) bool {
	_, exists := r.factories[number]
	return exists
}

func (r *customSectionRegistry) RegisterDefaults() {
	r.factories[156] = NewSection156
	r.factories[157] = NewSection157
	r.factories[158] = NewSection158
}
