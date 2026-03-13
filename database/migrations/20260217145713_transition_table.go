package migrations

import (
	"gamebook-backend/database"
	"gamebook-backend/database/entities"

	"gorm.io/gorm"
)

func init() {
	database.RegisterMigration("20260217145713_transition_table", UpTransitionTable, DownTransitionTable)
}

func UpTransitionTable(db *gorm.DB) error {
	return db.AutoMigrate(&entities.Transition{})
}

func DownTransitionTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&entities.Transition{})
}
