package migrations

import (
	"gamebook-backend/database"
	"gamebook-backend/database/entities"

	"gorm.io/gorm"
)

func init() {
	database.RegisterMigration("20260310000000_player_log_table", UpPlayerLogTable, DownPlayerLogTable)
}

func UpPlayerLogTable(db *gorm.DB) error {
	return db.AutoMigrate(&entities.PlayerLog{})
}

func DownPlayerLogTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&entities.PlayerLog{})
}
