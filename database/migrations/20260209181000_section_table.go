package migrations

import (
	"gamebook-backend/database"
	"gamebook-backend/database/entities"

	"gorm.io/gorm"
)

func init() {
	database.RegisterMigration("20260209181000_section_table", UpSectionTable, DownSectionTable)
}

func UpSectionTable(db *gorm.DB) error {
	return db.AutoMigrate(&entities.Section{})
}

func DownSectionTable(db *gorm.DB) error {

	return db.Migrator().DropTable(&entities.Section{})
}
