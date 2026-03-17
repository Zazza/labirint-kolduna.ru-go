package section

import (
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"github.com/google/uuid"
	"testing"
)

func TestCheck_WithBagCondition_Pass(t *testing.T) {
	player := createTestPlayerWithBag()
	section := createTestSectionWithBagCondition()
	player.Section = section

	result, err := goto_check.CheckDices(player, section)
	if err != nil {
		t.Fatalf("CheckDices failed: %v", err)
	}

	if !result {
		t.Fatal("expected true for bag condition")
	}
}

func TestCheck_WithBagCondition_FailBag(t *testing.T) {
	player := createTestPlayerWithoutBag()
	section := createTestSectionWithBagCondition()
	player.Section = section

	result, err := goto_check.CheckDices(player, section)
	if err != nil {
		t.Fatalf("CheckDices failed: %v", err)
	}

	if result {
		t.Fatal("expected false for bag condition")
	}
}

func TestCheck_WithBagCondition_FailDice(t *testing.T) {
	player := createTestPlayerWithoutBag()
	section := createTestSectionWithBagCondition()
	player.Section = section

	result, err := goto_check.CheckDices(player, section)
	if err != nil {
		t.Fatalf("CheckDices failed: %v", err)
	}

	if result {
		t.Fatal("expected false for bag condition")
	}
}

func TestCheck_NegatedBagCondition_Pass(t *testing.T) {
	player := createTestPlayerWithoutBag()
	section := createTestSectionWithNegatedBagCondition()
	player.Section = section

	result, err := goto_check.CheckDices(player, section)
	if err != nil {
		t.Fatalf("CheckDices failed: %v", err)
	}

	if !result {
		t.Fatal("expected true for negated bag condition")
	}
}

func TestCheck_NegatedBagCondition_Fail(t *testing.T) {
	player := createTestPlayer()
	section := createTestSectionWithNegatedBagCondition()
	player.Section = section

	result, err := goto_check.CheckDices(player, section)
	if err != nil {
		t.Fatalf("CheckDices failed: %v", err)
	}

	if result {
		t.Fatal("expected false for negated bag condition")
	}
}
