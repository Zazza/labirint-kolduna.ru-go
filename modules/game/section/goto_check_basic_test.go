package section

import (
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"github.com/google/uuid"
	"testing"
)

func TestCheck_NoDices(t *testing.T) {
	player := createTestPlayer()
	section := createTestSectionWithoutDices()
	player.Section = section

	result, err := goto_check.CheckDices(player, section)
	if err != nil {
		t.Fatalf("CheckDices failed: %v", err)
	}

	if !result {
		t.Fatal("expected false when no dices required")
	}
}

func TestCheck_SingleDice_Pass(t *testing.T) {
	player := createTestPlayer()
	section := createTestSectionWithSingleDice()
	player.Section = section

	result, err := goto_check.CheckDices(player, section)
	if err != nil {
		t.Fatalf("CheckDices failed: %v", err)
	}

	if !result {
		t.Fatal("expected true for single dice")
	}
}

func TestCheck_SingleDice_Fail(t *testing.T) {
	player := createTestPlayer()
	section := createTestSectionWithMultipleDice()
	player.Section = section

	result, err := goto_check.CheckDices(player, section)
	if err != nil {
		t.Fatalf("CheckDices failed: %v", err)
	}

	if result {
		t.Fatal("expected false for multiple dice")
	}
}
