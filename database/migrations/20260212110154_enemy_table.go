package migrations

import (
	"gamebook-backend/database"
	"gamebook-backend/database/entities"

	"gorm.io/gorm"
)

func init() {
	database.RegisterMigration("20260212110154_enemy_table", UpEnemyTable, DownEnemyTable)
}

func UpEnemyTable(db *gorm.DB) error {
	return db.AutoMigrate(&entities.Enemy{})
}

func DownEnemyTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&entities.Enemy{})
}
