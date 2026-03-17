package section

import (
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"github.com/google/uuid"
	"testing"
)

func TestCheck_WithMultipleConditions_Pass(t *testing.T) {
	player := createTestPlayer()
	section := createTestSectionWithMultipleConditions()
	player.Section = section

	result, err := goto_check.CheckDices(player, section)
	if err != nil {
		t.Fatalf("CheckDices failed: %v", err)
	}

	if !result {
		t.Fatal("expected true for multiple conditions")
	}
}

func TestCheck_InvalidExpression(t *testing.T) {
	player := createTestPlayer()
	section := createTestSectionWithInvalidExpression()
	player.Section = section

	result, err := goto_check.CheckDices(player, section)
	if err != nil {
		t.Fatalf("CheckDices failed: %v", err)
	}

	if result {
		t.Fatal("expected false for invalid expression")
	}
}
