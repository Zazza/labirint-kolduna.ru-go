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

func setupSectionTestDB(t *testing.T) *gorm.DB {
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
		db.Exec("DROP TABLE IF EXISTS sections CASCADE")
		db.Exec("DROP TABLE IF EXISTS enemies CASCADE")
		db.Exec("DROP TABLE IF EXISTS transitions CASCADE")
		db.Exec("DROP TABLE IF EXISTS section_enemies CASCADE")
		config.CloseDatabaseConnection(db)
	})

	return db
}

func TestSectionSeeder_Success(t *testing.T) {
	db := setupSectionTestDB(t)

	jsonData := []SectionJSON{
		{
			Type:         "normal",
			Number:       1,
			Text:         "You are at the entrance of a mysterious dungeon.",
			EnemyAliases: []string{},
			Dices:        false,
		},
	}

	jsonBytes, err := json.Marshal(jsonData)
	require.NoError(t, err)

	testDir := filepath.Join(os.TempDir(), "test_section_seeds")
	seedersDir := filepath.Join(testDir, "database", "seeders")
	jsonDir := filepath.Join(seedersDir, "json")
	err = os.MkdirAll(jsonDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	testFilePath := filepath.Join(jsonDir, "sections.json")
	err = os.WriteFile(testFilePath, jsonBytes, 0644)
	require.NoError(t, err)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(testDir)
	require.NoError(t, err)
	defer os.Chdir(origWd)

	err = SectionSeeder(db)
	assert.NoError(t, err)

	var count int64
	db.Model(&entities.Section{}).Count(&count)
	assert.Equal(t, int64(1), count)

	var section entities.Section
	err = db.Where("number = ?", 1).First(&section).Error
	require.NoError(t, err)
	assert.Equal(t, "normal", string(section.Type))
	assert.Equal(t, "You are at the entrance of a mysterious dungeon.", section.Text)
	assert.Equal(t, uint(1), section.Number)
}

func TestSectionSeeder_Upsert(t *testing.T) {
	db := setupSectionTestDB(t)

	jsonData := []SectionJSON{
		{
			Type:         "normal",
			Number:       1,
			Text:         "Original text",
			EnemyAliases: []string{},
			Dices:        false,
		},
	}

	jsonBytes, err := json.Marshal(jsonData)
	require.NoError(t, err)

	testDir := filepath.Join(os.TempDir(), "test_section_seeds_upsert")
	seedersDir := filepath.Join(testDir, "database", "seeders")
	jsonDir := filepath.Join(seedersDir, "json")
	err = os.MkdirAll(jsonDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	testFilePath := filepath.Join(jsonDir, "sections.json")
	err = os.WriteFile(testFilePath, jsonBytes, 0644)
	require.NoError(t, err)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(testDir)
	require.NoError(t, err)
	defer os.Chdir(origWd)

	err = SectionSeeder(db)
	require.NoError(t, err)

	var countBefore int64
	db.Model(&entities.Section{}).Count(&countBefore)
	assert.Equal(t, int64(1), countBefore)

	updatedJsonData := []SectionJSON{
		{
			Type:         "battle",
			Number:       1,
			Text:         "Updated text",
			EnemyAliases: []string{},
			Dices:        true,
		},
	}

	updatedBytes, err := json.Marshal(updatedJsonData)
	require.NoError(t, err)
	err = os.WriteFile(testFilePath, updatedBytes, 0644)
	require.NoError(t, err)

	err = SectionSeeder(db)
	assert.NoError(t, err)

	var countAfter int64
	db.Model(&entities.Section{}).Count(&countAfter)
	assert.Equal(t, countBefore, countAfter)

	var section entities.Section
	err = db.Where("number = ?", 1).First(&section).Error
	require.NoError(t, err)
	assert.Equal(t, "battle", string(section.Type))
	assert.Equal(t, "Updated text", section.Text)
	assert.True(t, section.Dices)
}

func TestSectionSeeder_WithEnemies(t *testing.T) {
	db := setupSectionTestDB(t)

	err := db.AutoMigrate(&entities.Enemy{})
	require.NoError(t, err)

	enemies := []entities.Enemy{
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
	}

	for _, enemy := range enemies {
		err = db.Create(&enemy).Error
		require.NoError(t, err)
	}

	jsonData := []SectionJSON{
		{
			Type:         "battle",
			Number:       2,
			Text:         "Enemies attack!",
			EnemyAliases: []string{"goblin", "orc"},
			Dices:        false,
		},
	}

	jsonBytes, err := json.Marshal(jsonData)
	require.NoError(t, err)

	testDir := filepath.Join(os.TempDir(), "test_section_seeds_enemies")
	seedersDir := filepath.Join(testDir, "database", "seeders")
	jsonDir := filepath.Join(seedersDir, "json")
	err = os.MkdirAll(jsonDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	testFilePath := filepath.Join(jsonDir, "sections.json")
	err = os.WriteFile(testFilePath, jsonBytes, 0644)
	require.NoError(t, err)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(testDir)
	require.NoError(t, err)
	defer os.Chdir(origWd)

	err = SectionSeeder(db)
	assert.NoError(t, err)

	var section entities.Section
	err = db.Where("number = ?", 2).Preload("SectionEnemies").First(&section).Error
	require.NoError(t, err)
	assert.Len(t, section.SectionEnemies, 2)

	enemyAliases := []string{}
	for _, enemy := range section.SectionEnemies {
		enemyAliases = append(enemyAliases, enemy.Alias)
	}
	assert.Contains(t, enemyAliases, "goblin")
	assert.Contains(t, enemyAliases, "orc")
}

func TestSectionSeeder_InvalidJSON(t *testing.T) {
	db := setupSectionTestDB(t)

	testDir := filepath.Join(os.TempDir(), "test_section_seeds_invalid_json")
	seedersDir := filepath.Join(testDir, "database", "seeders")
	jsonDir := filepath.Join(seedersDir, "json")
	err := os.MkdirAll(jsonDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	testFilePath := filepath.Join(jsonDir, "sections.json")
	err = os.WriteFile(testFilePath, []byte("[]\n{invalid json}"), 0644)
	require.NoError(t, err)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(testDir)
	require.NoError(t, err)
	defer os.Chdir(origWd)

	err = SectionSeeder(db)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal sections.json")
}

func TestSectionSeeder_EnemyNotFound(t *testing.T) {
	db := setupSectionTestDB(t)

	err := db.AutoMigrate(&entities.Enemy{})
	require.NoError(t, err)

	jsonData := []SectionJSON{
		{
			Type:         "battle",
			Number:       1,
			Text:         "Unknown enemy",
			EnemyAliases: []string{"nonexistent_enemy"},
			Dices:        false,
		},
	}

	jsonBytes, err := json.Marshal(jsonData)
	require.NoError(t, err)

	testDir := filepath.Join(os.TempDir(), "test_section_seeds_enemy_not_found")
	seedersDir := filepath.Join(testDir, "database", "seeders")
	jsonDir := filepath.Join(seedersDir, "json")
	err = os.MkdirAll(jsonDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	testFilePath := filepath.Join(jsonDir, "sections.json")
	err = os.WriteFile(testFilePath, jsonBytes, 0644)
	require.NoError(t, err)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(testDir)
	require.NoError(t, err)
	defer os.Chdir(origWd)

	err = SectionSeeder(db)
	assert.NoError(t, err)

	var section entities.Section
	err = db.Where("number = ?", 1).Preload("SectionEnemies").First(&section).Error
	require.NoError(t, err)
	assert.Len(t, section.SectionEnemies, 0)
}

func TestSectionSeeder_MultipleSections(t *testing.T) {
	db := setupSectionTestDB(t)

	jsonData := []SectionJSON{
		{
			Type:         "normal",
			Number:       1,
			Text:         "Section 1",
			EnemyAliases: []string{},
			Dices:        false,
		},
		{
			Type:         "battle",
			Number:       2,
			Text:         "Section 2",
			EnemyAliases: []string{},
			Dices:        false,
		},
		{
			Type:         "normal",
			Number:       3,
			Text:         "Section 3",
			EnemyAliases: []string{},
			Dices:        true,
		},
	}

	jsonBytes, err := json.Marshal(jsonData)
	require.NoError(t, err)

	testDir := filepath.Join(os.TempDir(), "test_section_seeds_multiple")
	seedersDir := filepath.Join(testDir, "database", "seeders")
	jsonDir := filepath.Join(seedersDir, "json")
	err = os.MkdirAll(jsonDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	testFilePath := filepath.Join(jsonDir, "sections.json")
	err = os.WriteFile(testFilePath, jsonBytes, 0644)
	require.NoError(t, err)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(testDir)
	require.NoError(t, err)
	defer os.Chdir(origWd)

	err = SectionSeeder(db)
	assert.NoError(t, err)

	var count int64
	db.Model(&entities.Section{}).Count(&count)
	assert.Equal(t, int64(3), count)
}
