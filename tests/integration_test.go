package tests_test

import (
	"fmt"
	"os"
	"testing"

	"gamebook-backend/database/entities"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// testDatabase - глобальное подключение к тестовой БД
var testDatabase *gorm.DB

// setupTestDB подготавливает тестовую БД перед каждым тестом
func setupTestDB(t *testing.T) *gorm.DB {
	if testDatabase == nil {
		// Используем настройки из .env.test или переменных окружения
		dbUser := getEnvOrDefault("DB_USER", "postgres")
		dbPass := getEnvOrDefault("DB_PASS", "password")
		dbHost := getEnvOrDefault("DB_HOST", "localhost")
		dbName := getEnvOrDefault("DB_NAME", "test_db")
		dbPort := getEnvOrDefault("DB_PORT", "5433")

		dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
			dbHost, dbUser, dbPass, dbName, dbPort)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Fatalf("Failed to connect to test database: %v", err)
		}

		// Включаем расширение uuid-ossp
		db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

		testDatabase = db
	}

	// Очищаем таблицы перед каждым тестом
	cleanupTables(t, testDatabase)

	return testDatabase
}

// cleanupTables очищает таблицы перед каждым тестом
func cleanupTables(t *testing.T, db *gorm.DB) {
	tables := []string{
		"player_section_enemies",
		"player_sections",
		"players",
		"sections",
		"enemies",
		"users",
	}

	for _, table := range tables {
		if err := db.Exec("DELETE FROM " + table).Error; err != nil {
			t.Logf("Warning: Failed to delete from %s: %v", table, err)
		}
	}
}

// getEnvOrDefault возвращает значение переменной окружения или значение по умолчанию
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// TestMain - точка входа для всех тестов
func TestMain(m *testing.M) {
	// Можно добавить код для запуска docker-compose перед тестами
	// и остановки после тестов

	// Запускаем тесты
	os.Exit(m.Run())
}

// Пример интеграционного теста для Player

func TestPlayerRepository_CreateAndGet(t *testing.T) {
	db := setupTestDB(t)

	// Создаем тестового пользователя
	user := entities.User{
		ID:       uuid.New(),
		Name:     "Test User",
		Password: "password123",
	}
	err := db.Create(&user).Error
	assert.NoError(t, err)

	// Создаем тестовую секцию
	section := entities.Section{
		ID:     uuid.New(),
		Number: 1,
		Text:   "Test section text",
		Type:   "story",
	}
	err = db.Create(&section).Error
	assert.NoError(t, err)

	// Создаем игрока
	player := entities.Player{
		ID:        uuid.New(),
		UserID:    user.ID,
		SectionID: section.ID,
		Health:    20,
		Gold:      50,
		Bag:       []string{"Sword", "Shield"},
	}
	err = db.Create(&player).Error
	assert.NoError(t, err)

	// Проверяем, что игрок создался
	var foundPlayer entities.Player
	err = db.First(&foundPlayer, player.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, player.Health, foundPlayer.Health)
	assert.Equal(t, 2, len(foundPlayer.Bag))
}

func TestPlayerRepository_UpdateHealth(t *testing.T) {
	db := setupTestDB(t)

	// Создаем тестового пользователя и секцию
	user := entities.User{ID: uuid.New(), Name: "Test User", Password: "password123"}
	section := entities.Section{ID: uuid.New(), Number: 2, Text: "Test", Type: "story"}

	db.Create(&user)
	db.Create(&section)

	// Создаем игрока
	player := entities.Player{
		ID:        uuid.New(),
		UserID:    user.ID,
		SectionID: section.ID,
		Health:    20,
	}
	db.Create(&player)

	// Обновляем здоровье
	newHealth := uint(15)
	err := db.Model(&player).Update("health", newHealth).Error
	assert.NoError(t, err)

	// Проверяем обновление
	var updatedPlayer entities.Player
	db.First(&updatedPlayer, player.ID)
	assert.Equal(t, newHealth, updatedPlayer.Health)
}

func TestPlayerRepository_WithRelations(t *testing.T) {
	db := setupTestDB(t)

	// Создаем связанные данные
	user := entities.User{ID: uuid.New(), Name: "Test User", Password: "password123"}
	section := entities.Section{ID: uuid.New(), Number: 3, Text: "Test", Type: "story"}
	enemy := entities.Enemy{
		ID:     uuid.New(),
		Alias:  "goblin",
		Name:   "Гоблин",
		Damage: 5,
		Health: 10,
	}

	db.Create(&user)
	db.Create(&section)
	db.Create(&enemy)

	// Создаем игрока с врагом
	player := entities.Player{
		ID:        uuid.New(),
		UserID:    user.ID,
		SectionID: section.ID,
		Health:    20,
	}
	db.Create(&player)

	// Создаем связь игрок-секция-враг
	playerSection := entities.PlayerSection{
		ID:        uuid.New(),
		PlayerID:  player.ID,
		SectionID: section.ID,
	}
	db.Create(&playerSection)

	playerSectionEnemy := entities.PlayerSectionEnemy{
		ID:        uuid.New(),
		PlayerID:  player.ID,
		SectionID: section.ID,
		EnemyID:   enemy.ID,
		Health:    10,
	}
	db.Create(&playerSectionEnemy)

	// Проверяем связи
	var foundPlayer entities.Player
	err := db.Preload("Section").Preload("Section.SectionEnemies").First(&foundPlayer, player.ID).Error
	assert.NoError(t, err)
	assert.NotNil(t, foundPlayer.Section)
}
