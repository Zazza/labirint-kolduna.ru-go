package migrations

import (
	"gamebook-backend/database"
	"gamebook-backend/database/entities"

	"gorm.io/gorm"
)

func init() {
	database.RegisterMigration("20260226165412_description_log_table", UpDescriptionLogTable, DownDescriptionLogTable)
}

func UpDescriptionLogTable(db *gorm.DB) error {
	return db.AutoMigrate(&entities.DescriptionLog{})
}

func DownDescriptionLogTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&entities.DescriptionLog{})
}
