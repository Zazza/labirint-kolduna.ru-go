package seeds

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"gamebook-backend/config"
	"gamebook-backend/database/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransitionSeeder_Success(t *testing.T) {
	db := config.SetUpInMemoryDatabase()

	err := db.Exec("PRAGMA foreign_keys = ON").Error
	require.NoError(t, err)

	section1 := entities.Section{
		Type:   "normal",
		Number: 1,
		Text:   "Section 1",
		Dices:  false,
	}
	err = db.Create(&section1).Error
	require.NoError(t, err)

	section2 := entities.Section{
		Type:   "normal",
		Number: 2,
		Text:   "Section 2",
		Dices:  false,
	}
	err = db.Create(&section2).Error
	require.NoError(t, err)

	jsonData := []TransitionJSON{
		{
			TextOrder:           1,
			SectionNumber:       1,
			TargetSectionNumber: 2,
			AvailableOnce:       false,
			Text:                "Go deeper",
			PlayerInput:         false,
			PlayerDebuff:        []*entities.PlayerDebuff{},
			PlayerBuff:          []*entities.PlayerBuff{},
		},
	}

	jsonBytes, err := json.Marshal(jsonData)
	require.NoError(t, err)

	testDir := filepath.Join(os.TempDir(), "test_transition_seeds")
	seedersDir := filepath.Join(testDir, "database", "seeders")
	jsonDir := filepath.Join(seedersDir, "json")
	err = os.MkdirAll(jsonDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	testFilePath := filepath.Join(jsonDir, "transitions.json")
	err = os.WriteFile(testFilePath, jsonBytes, 0644)
	require.NoError(t, err)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(testDir)
	require.NoError(t, err)
	defer os.Chdir(origWd)

	err = TransitionSeeder(db)
	assert.NoError(t, err)

	var count int64
	db.Model(&entities.Transition{}).Count(&count)
	assert.Equal(t, int64(1), count)

	var transition entities.Transition
	err = db.Where("section_id = ?", section1.ID).First(&transition).Error
	require.NoError(t, err)
	assert.Equal(t, "Go deeper", transition.Text)
	assert.Equal(t, uint(1), transition.TextOrder)
}

func TestTransitionSeeder_Upsert(t *testing.T) {
	db := config.SetUpInMemoryDatabase()

	err := db.Exec("PRAGMA foreign_keys = ON").Error
	require.NoError(t, err)

	section1 := entities.Section{
		Type:   "normal",
		Number: 1,
		Text:   "Section 1",
		Dices:  false,
	}
	err = db.Create(&section1).Error
	require.NoError(t, err)

	section2 := entities.Section{
		Type:   "normal",
		Number: 2,
		Text:   "Section 2",
		Dices:  false,
	}
	err = db.Create(&section2).Error
	require.NoError(t, err)

	jsonData := []TransitionJSON{
		{
			TextOrder:           1,
			SectionNumber:       1,
			TargetSectionNumber: 2,
			AvailableOnce:       false,
			Text:                "Original text",
			PlayerInput:         false,
			PlayerDebuff:        []*entities.PlayerDebuff{},
			PlayerBuff:          []*entities.PlayerBuff{},
		},
	}

	jsonBytes, err := json.Marshal(jsonData)
	require.NoError(t, err)

	testDir := filepath.Join(os.TempDir(), "test_transition_seeds_upsert")
	seedersDir := filepath.Join(testDir, "database", "seeders")
	jsonDir := filepath.Join(seedersDir, "json")
	err = os.MkdirAll(jsonDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	testFilePath := filepath.Join(jsonDir, "transitions.json")
	err = os.WriteFile(testFilePath, jsonBytes, 0644)
	require.NoError(t, err)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(testDir)
	require.NoError(t, err)
	defer os.Chdir(origWd)

	err = TransitionSeeder(db)
	require.NoError(t, err)

	var countBefore int64
	db.Model(&entities.Transition{}).Count(&countBefore)
	assert.Equal(t, int64(1), countBefore)

	updatedJsonData := []TransitionJSON{
		{
			TextOrder:           1,
			SectionNumber:       1,
			TargetSectionNumber: 2,
			AvailableOnce:       true,
			Text:                "Updated text",
			PlayerInput:         true,
			PlayerDebuff:        []*entities.PlayerDebuff{},
			PlayerBuff:          []*entities.PlayerBuff{},
		},
	}

	updatedBytes, err := json.Marshal(updatedJsonData)
	require.NoError(t, err)
	err = os.WriteFile(testFilePath, updatedBytes, 0644)
	require.NoError(t, err)

	err = TransitionSeeder(db)
	assert.NoError(t, err)

	var countAfter int64
	db.Model(&entities.Transition{}).Count(&countAfter)
	assert.Equal(t, countBefore, countAfter)

	var transition entities.Transition
	err = db.Where("section_id = ? AND target_section_id = ? AND text_order = ?", section1.ID, section2.ID, 1).First(&transition).Error
	require.NoError(t, err)
	assert.Equal(t, "Updated text", transition.Text)
	assert.True(t, transition.AvailableOnce)
	assert.True(t, transition.PlayerInput)
}

func TestTransitionSeeder_InvalidSectionNumber(t *testing.T) {
	db := config.SetUpInMemoryDatabase()

	err := db.Exec("PRAGMA foreign_keys = ON").Error
	require.NoError(t, err)

	section1 := entities.Section{
		Type:   "normal",
		Number: 1,
		Text:   "Section 1",
		Dices:  false,
	}
	err = db.Create(&section1).Error
	require.NoError(t, err)

	jsonData := []TransitionJSON{
		{
			TextOrder:           1,
			SectionNumber:       1,
			TargetSectionNumber: 999,
			AvailableOnce:       false,
			Text:                "Go to nowhere",
			PlayerInput:         false,
			PlayerDebuff:        []*entities.PlayerDebuff{},
			PlayerBuff:          []*entities.PlayerBuff{},
		},
	}

	jsonBytes, err := json.Marshal(jsonData)
	require.NoError(t, err)

	testDir := filepath.Join(os.TempDir(), "test_transition_seeds_invalid_section")
	seedersDir := filepath.Join(testDir, "database", "seeders")
	jsonDir := filepath.Join(seedersDir, "json")
	err = os.MkdirAll(jsonDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	testFilePath := filepath.Join(jsonDir, "transitions.json")
	err = os.WriteFile(testFilePath, jsonBytes, 0644)
	require.NoError(t, err)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(testDir)
	require.NoError(t, err)
	defer os.Chdir(origWd)

	err = TransitionSeeder(db)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find target section 999 for transition")
}

func TestTransitionSeeder_InvalidJSON(t *testing.T) {
	db := config.SetUpInMemoryDatabase()

	testDir := filepath.Join(os.TempDir(), "test_transition_seeds_invalid_json")
	seedersDir := filepath.Join(testDir, "database", "seeders")
	jsonDir := filepath.Join(seedersDir, "json")
	err := os.MkdirAll(jsonDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	testFilePath := filepath.Join(jsonDir, "transitions.json")
	err = os.WriteFile(testFilePath, []byte("[]\n{invalid json}"), 0644)
	require.NoError(t, err)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(testDir)
	require.NoError(t, err)
	defer os.Chdir(origWd)

	err = TransitionSeeder(db)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal transitions.json")
}

func TestTransitionSeeder_MultipleTransitions(t *testing.T) {
	db := config.SetUpInMemoryDatabase()

	err := db.Exec("PRAGMA foreign_keys = ON").Error
	require.NoError(t, err)

	sections := []entities.Section{
		{Type: "normal", Number: 1, Text: "Section 1", Dices: false},
		{Type: "normal", Number: 2, Text: "Section 2", Dices: false},
		{Type: "normal", Number: 3, Text: "Section 3", Dices: false},
	}

	for _, section := range sections {
		err := db.Create(&section).Error
		require.NoError(t, err)
	}

	jsonData := []TransitionJSON{
		{
			TextOrder:           1,
			SectionNumber:       1,
			TargetSectionNumber: 2,
			AvailableOnce:       false,
			Text:                "Transition 1",
			PlayerInput:         false,
			PlayerDebuff:        []*entities.PlayerDebuff{},
			PlayerBuff:          []*entities.PlayerBuff{},
		},
		{
			TextOrder:           2,
			SectionNumber:       1,
			TargetSectionNumber: 3,
			AvailableOnce:       true,
			Text:                "Transition 2",
			PlayerInput:         false,
			PlayerDebuff:        []*entities.PlayerDebuff{},
			PlayerBuff:          []*entities.PlayerBuff{},
		},
		{
			TextOrder:           1,
			SectionNumber:       2,
			TargetSectionNumber: 3,
			AvailableOnce:       false,
			Text:                "Transition 3",
			PlayerInput:         true,
			PlayerDebuff:        []*entities.PlayerDebuff{},
			PlayerBuff:          []*entities.PlayerBuff{},
		},
	}

	jsonBytes, err := json.Marshal(jsonData)
	require.NoError(t, err)

	testDir := filepath.Join(os.TempDir(), "test_transition_seeds_multiple")
	seedersDir := filepath.Join(testDir, "database", "seeders")
	jsonDir := filepath.Join(seedersDir, "json")
	err = os.MkdirAll(jsonDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	testFilePath := filepath.Join(jsonDir, "transitions.json")
	err = os.WriteFile(testFilePath, jsonBytes, 0644)
	require.NoError(t, err)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(testDir)
	require.NoError(t, err)
	defer os.Chdir(origWd)

	err = TransitionSeeder(db)
	assert.NoError(t, err)

	var count int64
	db.Model(&entities.Transition{}).Count(&count)
	assert.Equal(t, int64(3), count)

	var transitions []entities.Transition
	err = db.Find(&transitions).Error
	require.NoError(t, err)
	assert.Len(t, transitions, 3)
}

func TestTransitionSeeder_WithComplexFields(t *testing.T) {
	db := config.SetUpInMemoryDatabase()

	err := db.Exec("PRAGMA foreign_keys = ON").Error
	require.NoError(t, err)

	section1 := entities.Section{
		Type:   "normal",
		Number: 1,
		Text:   "Section 1",
		Dices:  false,
	}
	err = db.Create(&section1).Error
	require.NoError(t, err)

	section2 := entities.Section{
		Type:   "normal",
		Number: 2,
		Text:   "Section 2",
		Dices:  false,
	}
	err = db.Create(&section2).Error
	require.NoError(t, err)

	isBattleWin := true
	bribeResult := false
	dice := []string{"1d6", "1d8"}
	condition := "gold >= 10"

	jsonData := []TransitionJSON{
		{
			TextOrder:           1,
			SectionNumber:       1,
			TargetSectionNumber: 2,
			AvailableOnce:       false,
			Text:                "Complex transition",
			IsBattleWin:         &isBattleWin,
			BribeResult:         &bribeResult,
			PlayerInput:         true,
			Dice:                &dice,
			Dices:               nil,
			Condition:           &condition,
			PlayerDebuff:        []*entities.PlayerDebuff{},
			PlayerBuff:          []*entities.PlayerBuff{},
		},
	}

	jsonBytes, err := json.Marshal(jsonData)
	require.NoError(t, err)

	testDir := filepath.Join(os.TempDir(), "test_transition_seeds_complex")
	seedersDir := filepath.Join(testDir, "database", "seeders")
	jsonDir := filepath.Join(seedersDir, "json")
	err = os.MkdirAll(jsonDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	testFilePath := filepath.Join(jsonDir, "transitions.json")
	err = os.WriteFile(testFilePath, jsonBytes, 0644)
	require.NoError(t, err)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(testDir)
	require.NoError(t, err)
	defer os.Chdir(origWd)

	err = TransitionSeeder(db)
	assert.NoError(t, err)

	var transition entities.Transition
	err = db.Where("section_id = ?", section1.ID).First(&transition).Error
	require.NoError(t, err)
	assert.Equal(t, "Complex transition", transition.Text)
	assert.NotNil(t, transition.IsBattleWin)
	assert.True(t, *transition.IsBattleWin)
	assert.NotNil(t, transition.BribeResult)
	assert.False(t, *transition.BribeResult)
	assert.True(t, transition.PlayerInput)
	assert.NotNil(t, transition.Dice)
	assert.Len(t, *transition.Dice, 2)
	assert.Equal(t, "gold >= 10", *transition.Condition)
}
