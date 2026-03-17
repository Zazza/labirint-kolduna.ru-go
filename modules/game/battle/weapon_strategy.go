package battle

import (
	"gamebook-backend/database/entities"
	battleDTO "gamebook-backend/modules/game/dto"
)

type WeaponStrategy interface {
	Apply(weapon *entities.Weapons, index int, weapons *[]entities.Weapons) error
}

type WeaponStrategyFactory interface {
	GetStrategy(weaponType string) WeaponStrategy
}

type weaponStrategyFactory struct {
	strategies map[string]WeaponStrategy
}

func NewWeaponStrategyFactory() WeaponStrategyFactory {
	factory := &weaponStrategyFactory{
		strategies: make(map[string]WeaponStrategy),
	}

	factory.registerDefaults()
	return factory
}

func (f *weaponStrategyFactory) registerDefaults() {
	f.strategies[battleDTO.Lightning] = &LightningStrategy{}
	f.strategies[battleDTO.BallLightning] = &BallLightningStrategy{}
}

func (f *weaponStrategyFactory) GetStrategy(weaponType string) WeaponStrategy {
	if strategy, exists := f.strategies[weaponType]; exists {
		return strategy
	}
	return &DefaultWeaponStrategy{}
}

type LightningStrategy struct{}

func (s *LightningStrategy) Apply(weapon *entities.Weapons, index int, weapons *[]entities.Weapons) error {
	if weapon.Count > 0 {
		(*weapons)[index].Count -= 1
	}
	return nil
}

type BallLightningStrategy struct{}

func (s *BallLightningStrategy) Apply(weapon *entities.Weapons, index int, weapons *[]entities.Weapons) error {
	if weapon.Count > 0 {
		(*weapons)[index].Count -= 1
	}
	return nil
}

type DefaultWeaponStrategy struct{}

func (s *DefaultWeaponStrategy) Apply(weapon *entities.Weapons, index int, weapons *[]entities.Weapons) error {
	return nil
}
