package bonus

import (
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/repository"

	"gorm.io/gorm"
)

type BonusRegistry interface {
	Register(alias string, factory BonusFactoryFunc) error
	Get(alias string) (BonusFactoryFunc, error)
	GetBonus(db *gorm.DB, player entities.Player, alias string) (dto.Bonus, error)
}

type BonusFactoryFunc func(db *gorm.DB, player entities.Player) dto.Bonus

type bonusRegistry struct {
	factories map[string]BonusFactoryFunc
	bonusRepo repository.BonusRepository
}

func NewBonusRegistry(
	bonusRepo repository.BonusRepository,
) BonusRegistry {
	registry := &bonusRegistry{
		factories: make(map[string]BonusFactoryFunc),
		bonusRepo: bonusRepo,
	}

	registry.RegisterDefaults()

	return registry
}

func (r *bonusRegistry) Register(alias string, factory BonusFactoryFunc) error {
	r.factories[alias] = factory
	return nil
}

func (r *bonusRegistry) Get(alias string) (BonusFactoryFunc, error) {
	factory, exists := r.factories[alias]
	if !exists {
		return nil, dto.MessageBonusNotDefined
	}
	return factory, nil
}

func (r *bonusRegistry) GetBonus(db *gorm.DB, player entities.Player, alias string) (dto.Bonus, error) {
	factory, err := r.Get(alias)
	if err != nil {
		return nil, err
	}
	return factory(db, player), nil
}

func (r *bonusRegistry) RegisterDefaults() {
	r.factories[DeathSpellAlias] = NewDeathSpell
	r.factories[AntiPoisonSpellAlias] = NewAntiPoisonSpell
	r.factories[InstantMovementAlias] = NewInstantMovement
	r.factories[InstantHypnosisSpellAlias] = NewInstantHypnosisSpell
	r.factories[InstantRecoveryAlias] = NewInstantRecoverySpell
	r.factories[MagicDuckAlias] = NewMagicDuck
	r.factories[WandAlias] = NewWand
	r.factories[MagicRingAlias] = NewMagicRing
	r.factories[LuckyStoneAlias] = NewLuckyStone
	r.factories[DeathTeleportAlias] = NewDeathTeleport
}
