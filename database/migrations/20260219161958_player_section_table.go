package migrations

import (
	"gamebook-backend/database"
	"gamebook-backend/database/entities"

	"gorm.io/gorm"
)

func init() {
	database.RegisterMigration("20260219161958_player_section_table", UpPlayerSectionTable, DownPlayerSectionTable)
}

func UpPlayerSectionTable(db *gorm.DB) error {
	return db.AutoMigrate(&entities.PlayerSection{})
}

func DownPlayerSectionTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&entities.PlayerSection{})
}
