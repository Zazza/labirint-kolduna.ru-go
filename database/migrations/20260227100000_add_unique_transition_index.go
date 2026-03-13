package migrations

import (
	"gamebook-backend/database"

	"gorm.io/gorm"
)

func AddUniqueTransitionIndex(db *gorm.DB) error {
	return db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_transitions_section_target_order 
		ON transitions (section_id, target_section_id, text_order)
	`).Error
}

func RollbackUniqueTransitionIndex(db *gorm.DB) error {
	return db.Exec(`
		DROP INDEX IF EXISTS idx_transitions_section_target_order
	`).Error
}

func init() {
	database.RegisterMigration("20260227100000", AddUniqueTransitionIndex, RollbackUniqueTransitionIndex)
}
