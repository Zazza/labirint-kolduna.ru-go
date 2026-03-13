package sleep

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

func TestSleep2_Execute_Dice1LessThan3(t *testing.T) {
	db := setupTestDB(t)
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	sleep2 := NewSleep2(db, player)
	result, err := sleep2.Execute(context.Background(), 2, 0)

	assert.NoError(t, err)
	assert.True(t, result.Exit)
	assert.False(t, result.Death)
	assert.False(t, result.NextTry)
}

func TestSleep2_Execute_Dice1GreaterThan3(t *testing.T) {
	db := setupTestDB(t)
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	sleep2 := NewSleep2(db, player)
	result, err := sleep2.Execute(context.Background(), 5, 0)

	assert.NoError(t, err)
	assert.False(t, result.Exit)
	assert.False(t, result.Death)
	assert.True(t, result.NextTry)
}

func TestSleep2_Execute_Dice1Equal3(t *testing.T) {
	db := setupTestDB(t)
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	sleep2 := NewSleep2(db, player)
	result, err := sleep2.Execute(context.Background(), 3, 0)

	assert.NoError(t, err)
	assert.True(t, result.Exit)
	assert.False(t, result.Death)
	assert.False(t, result.NextTry)
}

func TestSleep4_Execute_Dice1GreaterOrEqual5(t *testing.T) {
	db := setupTestDB(t)
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	sleep4 := NewSleep4(db, player)
	result, err := sleep4.Execute(context.Background(), 5, 0)

	assert.NoError(t, err)
	assert.True(t, result.Exit)
	assert.False(t, result.Death)
}

func TestSleep4_Execute_Dice1LessThan5(t *testing.T) {
	db := setupTestDB(t)
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	sleep4 := NewSleep4(db, player)
	result, err := sleep4.Execute(context.Background(), 4, 0)

	assert.NoError(t, err)
	assert.False(t, result.Exit)
	assert.True(t, result.Death)
}

func TestSleep4_Execute_Dice1Equal4(t *testing.T) {
	db := setupTestDB(t)
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	sleep4 := NewSleep4(db, player)
	result, err := sleep4.Execute(context.Background(), 4, 0)

	assert.NoError(t, err)
	assert.False(t, result.Exit)
	assert.True(t, result.Death)
}

func TestGetSection_Sleep2(t *testing.T) {
	db := setupTestDB(t)
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	section, err := GetSection(db, player, 2)

	assert.NoError(t, err)
	assert.NotNil(t, section)
}

func TestGetSection_Sleep3(t *testing.T) {
	db := setupTestDB(t)
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	section, err := GetSection(db, player, 3)

	assert.NoError(t, err)
	assert.NotNil(t, section)
}

func TestGetSection_Sleep4(t *testing.T) {
	db := setupTestDB(t)
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	section, err := GetSection(db, player, 4)

	assert.NoError(t, err)
	assert.NotNil(t, section)
}

func TestGetSection_Sleep5(t *testing.T) {
	db := setupTestDB(t)
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	section, err := GetSection(db, player, 5)

	assert.NoError(t, err)
	assert.NotNil(t, section)
}

func TestGetSection_Sleep6(t *testing.T) {
	db := setupTestDB(t)
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	section, err := GetSection(db, player, 6)

	assert.NoError(t, err)
	assert.NotNil(t, section)
}

func TestGetSection_Sleep7(t *testing.T) {
	db := setupTestDB(t)
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	section, err := GetSection(db, player, 7)

	assert.NoError(t, err)
	assert.NotNil(t, section)
}

func TestGetSection_Sleep8(t *testing.T) {
	db := setupTestDB(t)
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	section, err := GetSection(db, player, 8)

	assert.NoError(t, err)
	assert.NotNil(t, section)
}

func TestGetSection_Sleep9(t *testing.T) {
	db := setupTestDB(t)
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	section, err := GetSection(db, player, 9)

	assert.NoError(t, err)
	assert.NotNil(t, section)
}

func TestGetSection_Sleep10(t *testing.T) {
	db := setupTestDB(t)
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	section, err := GetSection(db, player, 10)

	assert.NoError(t, err)
	assert.NotNil(t, section)
}

func TestGetSection_Sleep11(t *testing.T) {
	db := setupTestDB(t)
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	section, err := GetSection(db, player, 11)

	assert.NoError(t, err)
	assert.NotNil(t, section)
}

func TestGetSection_Sleep12(t *testing.T) {
	db := setupTestDB(t)
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	section, err := GetSection(db, player, 12)

	assert.NoError(t, err)
	assert.NotNil(t, section)
}

func TestGetSection_InvalidSection(t *testing.T) {
	db := setupTestDB(t)
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	section, err := GetSection(db, player, 99)

	assert.Error(t, err)
	assert.Nil(t, section)
	assert.Equal(t, dto.MessageSleepyKingdomSectionNotDefined, err)
}
