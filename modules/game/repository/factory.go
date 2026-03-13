package repository

import (
	"gorm.io/gorm"
)

type RepositoryFactory interface {
	NewBonusRepository(db *gorm.DB) BonusRepository
	NewBattleRepository(db *gorm.DB) BattleRepository
	NewDiceRepository(db *gorm.DB) DiceRepository
	NewPlayerRepository(db *gorm.DB) PlayerRepository
	NewSectionRepository(db *gorm.DB) SectionRepository
	NewPlayerSectionRepository(db *gorm.DB) PlayerSectionRepository
	NewPlayerSectionEnemyRepository(db *gorm.DB) PlayerSectionEnemyRepository
	NewTransitionRepository(db *gorm.DB) TransitionRepository
}

type repositoryFactory struct{}

func NewRepositoryFactory() RepositoryFactory {
	return &repositoryFactory{}
}

func (f *repositoryFactory) NewBonusRepository(db *gorm.DB) BonusRepository {
	return NewBonusRepository(db)
}

func (f *repositoryFactory) NewBattleRepository(db *gorm.DB) BattleRepository {
	return NewBattleRepository(db)
}

func (f *repositoryFactory) NewDiceRepository(db *gorm.DB) DiceRepository {
	return NewDiceRepository(db)
}

func (f *repositoryFactory) NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return NewPlayerRepository(db)
}

func (f *repositoryFactory) NewSectionRepository(db *gorm.DB) SectionRepository {
	return NewSectionRepository(db)
}

func (f *repositoryFactory) NewPlayerSectionRepository(db *gorm.DB) PlayerSectionRepository {
	return NewPlayerSectionRepository(db)
}

func (f *repositoryFactory) NewPlayerSectionEnemyRepository(db *gorm.DB) PlayerSectionEnemyRepository {
	return NewPlayerSectionEnemyRepository(db)
}

func (f *repositoryFactory) NewTransitionRepository(db *gorm.DB) TransitionRepository {
	return NewTransitionRepository(db)
}
