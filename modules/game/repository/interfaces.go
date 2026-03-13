package repository

import (
	"context"
	"gamebook-backend/database/entities"
	diceDTO "gamebook-backend/modules/game/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BonusRepository interface {
	GetByPlayerID(ctx context.Context, playerID uuid.UUID) ([]entities.PlayerBonus, error)
	Create(ctx context.Context, db *gorm.DB, bonus *entities.PlayerBonus) error
	Update(ctx context.Context, db *gorm.DB, bonus *entities.PlayerBonus) error
	Delete(ctx context.Context, db *gorm.DB, playerID uuid.UUID, alias string) error
}

type ChannelRepository interface {
}

type BattleRepository interface {
	GetByPlayerIdAndSectionNumber(tx *gorm.DB, playerId uuid.UUID, section uint) ([]entities.Battle, error)
	FindLastByPlayerIdAndSectionNumber(tx *gorm.DB, playerId uuid.UUID, section uint) (*entities.Battle, error)
	AddRecord(tx *gorm.DB, battle entities.Battle) (entities.Battle, error)
	RemoveSleepyByPlayerIDAndSectionNumber(ctx context.Context, tx *gorm.DB, playerId uuid.UUID, sectionId uint) error
}

type DiceRepository interface {
	GetLastByPlayerId(ctx context.Context, tx *gorm.DB, playerID uuid.UUID, reason diceDTO.ReasonType) (entities.Dice, error)
	FindBattleDicesByPlayerId(ctx context.Context, tx *gorm.DB, playerID uuid.UUID) diceDTO.BattleDicesDTO
	Create(ctx context.Context, tx *gorm.DB, playerID uuid.UUID, diceFirst uint, diceSecond uint, reason diceDTO.ReasonType) (entities.Dice, error)
	Remove(ctx context.Context, tx *gorm.DB, playerID uuid.UUID) error
}

type PlayerRepository interface {
	GetByUserId(ctx context.Context, tx *gorm.DB, userId string) (entities.Player, error)
	GetByPlayerId(ctx context.Context, tx *gorm.DB, playerID uuid.UUID) (entities.Player, error)
	Create(ctx context.Context, tx *gorm.DB, userId uuid.UUID, section entities.Section) (entities.Player, error)
	Update(ctx context.Context, tx *gorm.DB, player entities.Player) (*entities.Player, error)
}

type SectionRepository interface {
	GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entities.Section, error)
	GetBySectionNumber(ctx context.Context, tx *gorm.DB, section uint) (entities.Section, error)
	GetListBySectionNumbers(ctx context.Context, tx *gorm.DB, sections []uint) ([]entities.Section, error)
	GetAllWithTransitions(ctx context.Context, tx *gorm.DB, sectionNumbers []uint) ([]entities.Section, error)
	IsDicesRequired(section entities.Section) bool
}

type PlayerSectionRepository interface {
	Create(ctx context.Context, tx *gorm.DB, playerID uuid.UUID, sectionID uuid.UUID) error
	UpdateLastTargetSection(ctx context.Context, tx *gorm.DB, playerID uuid.UUID, targetSectionID uuid.UUID) error
	AddDescriptionLog(ctx context.Context, tx *gorm.DB, playerID uuid.UUID, description string) error
	GetLasSectionDescriptions(ctx context.Context, tx *gorm.DB, playerID uuid.UUID) (*[]entities.DescriptionLog, error)
	GetLastPlayerSection(ctx context.Context, tx *gorm.DB, playerID uuid.UUID) (*entities.PlayerSection, error)
	RemoveLastPlayerSection(ctx context.Context, tx *gorm.DB, playerID uuid.UUID) error
	GetPreviousSectionIdBySectionId(ctx context.Context, tx *gorm.DB, playerID uuid.UUID, sectionID uuid.UUID) (*entities.PlayerSection, error)
}

type PlayerSectionEnemyRepository interface {
	Create(ctx context.Context, tx *gorm.DB, playerID uuid.UUID, sectionID uuid.UUID, sectionEnemies []*entities.Enemy) ([]entities.PlayerSectionEnemy, error)
	GetEnemiesByPlayerIDAndSectionID(ctx context.Context, tx *gorm.DB, playerID uuid.UUID, sectionID uuid.UUID) (*[]entities.PlayerSectionEnemy, error)
	GetEnemiesByPlayerIDAndAlias(ctx context.Context, tx *gorm.DB, playerID uuid.UUID, enemyAlias string) (*[]entities.PlayerSectionEnemy, error)
	UpdateAll(ctx context.Context, tx *gorm.DB, enemies []entities.PlayerSectionEnemy) error
	Update(ctx context.Context, tx *gorm.DB, enemy entities.PlayerSectionEnemy) error
	RemoveSleepyByPlayerIDAndSectionNumber(ctx context.Context, tx *gorm.DB, playerId uuid.UUID, sectionId uuid.UUID) error
}

type TransitionRepository interface {
	GetByTransitionID(ctx context.Context, tx *gorm.DB, transitionID uuid.UUID) (entities.Transition, error)
	IsDicesRequired(transition entities.Transition) bool
}

type PlayerLogRepository interface {
	Create(ctx context.Context, db *gorm.DB, log *entities.PlayerLog) error
	GetByPlayerID(ctx context.Context, db *gorm.DB, playerID uuid.UUID) ([]entities.PlayerLog, error)
}
