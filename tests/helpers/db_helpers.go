package helpers

import (
	"context"
	"fmt"
	"gamebook-backend/config"
	"gamebook-backend/database/entities"
	"gamebook-backend/pkg/constants"

	"os"

	"testing"

	"github.com/google/uuid"
	"github.com/samber/do"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitTestDatabase(injector *do.Injector) error {
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPass := getEnvOrDefault("DB_PASS", "password")
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbName := getEnvOrDefault("DB_NAME", "test_db")
	dbPort := getEnvOrDefault("DB_PORT", "5432")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v", dbHost, dbUser, dbPass, dbName, dbPort)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to test database: %w", err)
	}

	config.RunExtension(db)

	do.ProvideNamed(injector, constants.DB, func(i *do.Injector) (*gorm.DB, error) {
		return db, nil
	})

	return nil
}

func CleanupAllTables(db *gorm.DB, ctx context.Context) {
	tables := []string{
		"battles",
		"player_section_enemies",
		"player_sections",
		"player_logs",
		"players",
		"section_enemies",
		"transitions",
		"sections",
		"enemies",
		"refresh_tokens",
		"users",
	}

	for _, table := range tables {
		db.Exec("DELETE FROM " + table)
		db.Exec("ALTER SEQUENCE " + table + "_id_seq RESTART WITH 1")
		db.Exec("TRUNCATE TABLE " + table + " CASCADE")
	}
}

func LoadSeedData(db *gorm.DB) error {
	var count int64
	db.Model(&entities.Section{}).Count(&count)
	if count > 0 {
		return nil
	}

	return nil
}

func FindSectionByNumber(t *testing.T, db *gorm.DB, number uint) *entities.Section {
	var section entities.Section
	err := db.Where("number = ?", number).First(&section).Error
	require.NoError(t, err, "Section with number %d not found", number)
	return &section
}

func CreatePlayerSectionEnemy(t *testing.T, db *gorm.DB, playerID, sectionID, enemyID uuid.UUID) *entities.PlayerSectionEnemy {
	playerSectionEnemyID := uuid.New()
	playerSectionEnemy := &entities.PlayerSectionEnemy{
		ID:        playerSectionEnemyID,
		PlayerID:  playerID,
		SectionID: sectionID,
		EnemyID:   enemyID,
		Health:    10,
	}
	err := db.Create(playerSectionEnemy).Error
	require.NoError(t, err)
	return playerSectionEnemy
}

func CreateBattleLog(t *testing.T, db *gorm.DB, playerID uuid.UUID, section uint, step uint, attacking string, damage uint, description string) *entities.Battle {
	battleID := uuid.New()
	battle := &entities.Battle{
		ID:          battleID,
		PlayerID:    playerID,
		Section:     section,
		Type:        "normal",
		Step:        step,
		Attacking:   attacking,
		Dice1:       3,
		Dice2:       4,
		Damage:      damage,
		Description: description,
		Weapon:      "Sword",
	}
	err := db.Create(battle).Error
	require.NoError(t, err)
	return battle
}

func UpdateDiceRoll(t *testing.T, db *gorm.DB, playerID uuid.UUID, diceFirst, diceSecond int) {
	var player entities.Player
	err := db.First(&player, playerID).Error
	require.NoError(t, err)

	diceValue := fmt.Sprintf("%d%d", diceFirst, diceSecond)
	err = db.Model(&player).Updates(map[string]interface{}{
		"dice_first":  diceFirst,
		"dice_second": diceSecond,
		"dice_value":  diceValue,
	}).Error
	require.NoError(t, err)
}

func CreatePlayerSection(t *testing.T, db *gorm.DB, playerSection entities.PlayerSection) {
	err := db.Create(&playerSection).Error
	require.NoError(t, err)
}

func UpdatePlayerSection(t *testing.T, db *gorm.DB, playerID uuid.UUID, sectionID uuid.UUID) {
	var player entities.Player
	err := db.First(&player, playerID).Error
	require.NoError(t, err)

	err = db.Model(&player).Update("section_id", sectionID).Error
	require.NoError(t, err)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
