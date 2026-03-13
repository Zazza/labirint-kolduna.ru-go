package migrations

import (
	"gamebook-backend/database"
	"gamebook-backend/database/entities"

	"gorm.io/gorm"
)

func init() {
	database.RegisterMigration("20260212200125_dice_table", UpDiceTable, DownDiceTable)
}

func UpDiceTable(db *gorm.DB) error {
	return db.AutoMigrate(&entities.Dice{})
}

func DownDiceTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&entities.Dice{})
}
