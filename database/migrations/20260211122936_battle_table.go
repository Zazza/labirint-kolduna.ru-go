package migrations

import (
	"gamebook-backend/database"
	"gamebook-backend/database/entities"

	"gorm.io/gorm"
)

func init() {
	database.RegisterMigration("20260211122936_battle_table", UpBattleTable, DownBattleTable)
}

func UpBattleTable(db *gorm.DB) error {
	return db.AutoMigrate(&entities.Battle{})
}

func DownBattleTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&entities.Battle{})
}
