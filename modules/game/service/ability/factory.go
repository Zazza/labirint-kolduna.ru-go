package ability

import (
	"gamebook-backend/modules/game/listener"
	"gamebook-backend/modules/game/repository"
	"gorm.io/gorm"
)

type AbilityFactory struct {
	db                           *gorm.DB
	diceRepository               repository.DiceRepository
	playerRepository             repository.PlayerRepository
	playerUpdateListener         listener.PlayerUpdateListener
	bonusRepository              repository.BonusRepository
	sectionRepository            repository.SectionRepository
	playerSectionRepository      repository.PlayerSectionRepository
	battleRepository             repository.BattleRepository
	playerSectionEnemyRepository repository.PlayerSectionEnemyRepository
}

func NewAbilityFactory(
	db *gorm.DB,
	diceRepo repository.DiceRepository,
	playerRepo repository.PlayerRepository,
	playerUpdateListener listener.PlayerUpdateListener,
	bonusRepo repository.BonusRepository,
	sectionRepo repository.SectionRepository,
	playerSectionRepo repository.PlayerSectionRepository,
	battleRepo repository.BattleRepository,
	playerSectionEnemyRepo repository.PlayerSectionEnemyRepository,
) *AbilityFactory {
	return &AbilityFactory{
		db:                           db,
		diceRepository:               diceRepo,
		playerRepository:             playerRepo,
		playerUpdateListener:         playerUpdateListener,
		bonusRepository:              bonusRepo,
		sectionRepository:            sectionRepo,
		playerSectionRepository:      playerSectionRepo,
		battleRepository:             battleRepo,
		playerSectionEnemyRepository: playerSectionEnemyRepo,
	}
}

func (f *AbilityFactory) CreateMedsAbility() MedsAbility {
	return NewMedsLogic(f.diceRepository, f.playerRepository, f.playerUpdateListener)
}

func (f *AbilityFactory) CreateBonusAbility() BonusAbility {
	return NewBonusLogic(f.diceRepository, f.playerRepository, f.playerUpdateListener, f.bonusRepository)
}

func (f *AbilityFactory) CreateSleepAbility() SleepAbility {
	return NewSleepLogic(f.db, f.sectionRepository, f.playerUpdateListener, f.playerSectionRepository, f.battleRepository, f.playerSectionEnemyRepository)
}

func (f *AbilityFactory) CreateBribeAbility() BribeAbility {
	return NewBribeLogic(f.db)
}

func (f *AbilityFactory) CreateDiceAbility() DiceAbility {
	return NewDiceLogic(f.db, f.diceRepository)
}
