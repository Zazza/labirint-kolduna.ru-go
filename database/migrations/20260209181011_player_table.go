package migrations

import (
	"gamebook-backend/database"
	"gamebook-backend/database/entities"

	"gorm.io/gorm"
)

func init() {
	database.RegisterMigration("20260209181011_player_table", UpPlayerTable, DownPlayerTable)
}

func UpPlayerTable(db *gorm.DB) error {
	return db.AutoMigrate(&entities.Player{})
}

func DownPlayerTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&entities.Player{})
}
