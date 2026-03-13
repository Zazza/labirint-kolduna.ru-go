package migrations

import (
	"gamebook-backend/database"
	"gamebook-backend/database/entities"

	"gorm.io/gorm"
)

func init() {
	database.RegisterMigration("20260223132101_player_section_enemy_table", UpPlayerSectionEnemyTable, DownPlayerSectionEnemyTable)
}

func UpPlayerSectionEnemyTable(db *gorm.DB) error {
	return db.AutoMigrate(&entities.PlayerSectionEnemy{})
}

func DownPlayerSectionEnemyTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&entities.PlayerSectionEnemy{})
}
