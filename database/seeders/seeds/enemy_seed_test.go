package seeds

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"testing"

	"gamebook-backend/config"
	"gamebook-backend/database/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	err := os.Setenv("APP_ENV", "test")
	if err != nil {
		return nil
	}

	config.LoadEnv()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("❌ Configuration error: %v", err)
	}

	db := config.SetUpTestDatabaseConnection(cfg)

	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Exec("DROP TABLE IF EXISTS enemies CASCADE")
		db.Exec("DROP TABLE IF EXISTS sections CASCADE")
		db.Exec("DROP TABLE IF EXISTS transitions CASCADE")
		db.Exec("DROP TABLE IF EXISTS section_enemies CASCADE")
		config.CloseDatabaseConnection(db)
	})

	return db
}

func TestEnemySeeder_Success(t *testing.T) {
	db := setupTestDB(t)

	jsonData := []EnemyJSON{
		{
			Alias:       "goblin",
			Name:        "Goblin",
			Damage:      2,
			DamageType:  "normal",
			MinDiceHits: 6,
			Health:      5,
			Defence:     0,
			PlayerArmor: true,
			Weapons:     []entities.EnemyWeapon{},
		},
	}

	jsonBytes, err := json.Marshal(jsonData)
	require.NoError(t, err)

	testDir := filepath.Join(os.TempDir(), "test_seeds")
	seedersDir := filepath.Join(testDir, "database", "seeders")
	jsonDir := filepath.Join(seedersDir, "json")
	err = os.MkdirAll(jsonDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	testFilePath := filepath.Join(jsonDir, "enemies.json")
	err = os.WriteFile(testFilePath, jsonBytes, 0644)
	require.NoError(t, err)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(testDir)
	require.NoError(t, err)
	defer os.Chdir(origWd)

	err = EnemySeeder(db)
	assert.NoError(t, err)

	var count int64
	db.Model(&entities.Enemy{}).Count(&count)
	assert.Equal(t, int64(1), count)

	var enemy entities.Enemy
	err = db.Where("alias = ?", "goblin").First(&enemy).Error
	require.NoError(t, err)
	assert.Equal(t, "Goblin", enemy.Name)
	assert.Equal(t, uint(2), enemy.Damage)
	assert.Equal(t, uint(5), enemy.Health)
}

func TestEnemySeeder_Upsert(t *testing.T) {
	db := setupTestDB(t)

	jsonData := []EnemyJSON{
		{
			Alias:       "goblin",
			Name:        "Goblin",
			Damage:      2,
			DamageType:  "normal",
			MinDiceHits: 6,
			Health:      5,
			Defence:     0,
			PlayerArmor: true,
			Weapons:     []entities.EnemyWeapon{},
		},
	}

	jsonBytes, err := json.Marshal(jsonData)
	require.NoError(t, err)

	testDir := filepath.Join(os.TempDir(), "test_seeds_upsert")
	seedersDir := filepath.Join(testDir, "database", "seeders")
	jsonDir := filepath.Join(seedersDir, "json")
	err = os.MkdirAll(jsonDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	testFilePath := filepath.Join(jsonDir, "enemies.json")
	err = os.WriteFile(testFilePath, jsonBytes, 0644)
	require.NoError(t, err)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(testDir)
	require.NoError(t, err)
	defer os.Chdir(origWd)

	err = EnemySeeder(db)
	require.NoError(t, err)

	var countBefore int64
	db.Model(&entities.Enemy{}).Count(&countBefore)
	assert.Equal(t, int64(1), countBefore)

	updatedJsonData := []EnemyJSON{
		{
			Alias:       "goblin",
			Name:        "Goblin Warrior",
			Damage:      3,
			DamageType:  "normal",
			MinDiceHits: 5,
			Health:      7,
			Defence:     1,
			PlayerArmor: false,
			Weapons:     []entities.EnemyWeapon{},
		},
	}

	updatedBytes, err := json.Marshal(updatedJsonData)
	require.NoError(t, err)
	err = os.WriteFile(testFilePath, updatedBytes, 0644)
	require.NoError(t, err)

	err = EnemySeeder(db)
	assert.NoError(t, err)

	var countAfter int64
	db.Model(&entities.Enemy{}).Count(&countAfter)
	assert.Equal(t, countBefore, countAfter)

	var enemy entities.Enemy
	err = db.Where("alias = ?", "goblin").First(&enemy).Error
	require.NoError(t, err)
	assert.Equal(t, "Goblin Warrior", enemy.Name)
	assert.Equal(t, uint(3), enemy.Damage)
	assert.Equal(t, uint(7), enemy.Health)
}

func TestEnemySeeder_InvalidJSON(t *testing.T) {
	db := setupTestDB(t)

	testDir := filepath.Join(os.TempDir(), "test_seeds_invalid_json")
	seedersDir := filepath.Join(testDir, "database", "seeders")
	jsonDir := filepath.Join(seedersDir, "json")
	err := os.MkdirAll(jsonDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	testFilePath := filepath.Join(jsonDir, "enemies.json")
	err = os.WriteFile(testFilePath, []byte("[]\n{invalid json}"), 0644)
	require.NoError(t, err)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(testDir)
	require.NoError(t, err)
	defer os.Chdir(origWd)

	err = EnemySeeder(db)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal enemies.json")
}

func TestEnemySeeder_FileNotFound(t *testing.T) {
	db := setupTestDB(t)

	testDir := filepath.Join(os.TempDir(), "test_seeds_file_not_found")
	err := os.MkdirAll(testDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(testDir)
	require.NoError(t, err)
	defer os.Chdir(origWd)

	err = EnemySeeder(db)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read enemies.json")
}

func TestEnemySeeder_MultipleEnemies(t *testing.T) {
	db := setupTestDB(t)

	jsonData := []EnemyJSON{
		{
			Alias:       "goblin",
			Name:        "Goblin",
			Damage:      2,
			DamageType:  "normal",
			MinDiceHits: 6,
			Health:      5,
			Defence:     0,
			PlayerArmor: true,
			Weapons:     []entities.EnemyWeapon{},
		},
		{
			Alias:       "orc",
			Name:        "Orc",
			Damage:      4,
			DamageType:  "normal",
			MinDiceHits: 5,
			Health:      10,
			Defence:     2,
			PlayerArmor: false,
			Weapons:     []entities.EnemyWeapon{},
		},
		{
			Alias:       "dragon",
			Name:        "Dragon",
			Damage:      8,
			DamageType:  "fire",
			MinDiceHits: 4,
			Health:      20,
			Defence:     5,
			PlayerArmor: true,
			Weapons:     []entities.EnemyWeapon{},
		},
	}

	jsonBytes, err := json.Marshal(jsonData)
	require.NoError(t, err)

	testDir := filepath.Join(os.TempDir(), "test_seeds_multiple")
	seedersDir := filepath.Join(testDir, "database", "seeders")
	jsonDir := filepath.Join(seedersDir, "json")
	err = os.MkdirAll(jsonDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	testFilePath := filepath.Join(jsonDir, "enemies.json")
	err = os.WriteFile(testFilePath, jsonBytes, 0644)
	require.NoError(t, err)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(testDir)
	require.NoError(t, err)
	defer os.Chdir(origWd)

	err = EnemySeeder(db)
	assert.NoError(t, err)

	var count int64
	db.Model(&entities.Enemy{}).Count(&count)
	assert.Equal(t, int64(3), count)

	var enemies []entities.Enemy
	err = db.Find(&enemies).Error
	require.NoError(t, err)
	assert.Len(t, enemies, 3)
}
