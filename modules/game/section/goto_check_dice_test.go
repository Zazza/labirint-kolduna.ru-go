package section

import (
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"github.com/google/uuid"
	"testing"
)

func TestCheck_TwoDice_Pass(t *testing.T) {
	player := createTestPlayer()
	section := createTestSectionWithTwoDice()
	player.Section = section

	result, err := goto_check.CheckDices(player, section)
	if err != nil {
		t.Fatalf("CheckDices failed: %v", err)
	}

	if !result {
		t.Fatal("expected true for two dice")
	}
}

func TestCheck_TwoDice_Fail(t *testing.T) {
	player := createTestPlayer()
	section := createTestSectionWithoutDices()
	player.Section = section

	result, err := goto_check.CheckDices(player, section)
	if err != nil {
		t.Fatalf("CheckDices failed: %v", err)
	}

	if result {
		t.Fatal("expected false for two dice")
	}
}

func TestCheck_MultipleDiceConditions_Pass(t *testing.T) {
	player := createTestPlayer()
	section := createTestSectionWithMultipleDiceConditions()
	player.Section = section

	result, err := goto_check.CheckDices(player, section)
	if err != nil {
		t.Fatalf("CheckDices failed: %v", err)
	}

	if !result {
		t.Fatal("expected true for multiple dice conditions")
	}
}

func TestCheck_ExactDiceMatch_Pass(t *testing.T) {
	player := createTestPlayer()
	section := createTestSectionWithExactDiceMatch()
	player.Section = section

	result, err := goto_check.CheckDices(player, section)
	if err != nil {
		t.Fatalf("CheckDices failed: %v", err)
	}

	if !result {
		t.Fatal("expected true for exact dice match")
	}
}

func TestCheck_ExactDiceMatch_Fail(t *testing.T) {
	player := createTestPlayer()
	section := createTestSectionWithoutDices()
	player.Section = section

	result, err := goto_check.CheckDices(player, section)
	if err != nil {
		t.Fatalf("CheckDices failed: %v", err)
	}

	if result {
		t.Fatal("expected false for exact dice match")
	}
}

func TestCheck_LessThanCondition_Pass(t *testing.T) {
	player := createTestPlayer()
	section := createTestSectionWithLessThanCondition()
	player.Section = section

	result, err := goto_check.CheckDices(player, section)
	if err != nil {
		t.Fatalf("CheckDices failed: %v", err)
	}

	if !result {
		t.Fatal("expected true for less than condition")
	}
}
