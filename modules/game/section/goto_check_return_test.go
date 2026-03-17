package section

import (
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"github.com/google/uuid"
	"testing"
)

func TestCheck_WithReturnToSection_Pass(t *testing.T) {
	player := createTestPlayer()
	section := createTestSectionWithReturnToSection()
	player.Section = section

	result, err := goto_check.CheckDices(player, section)
	if err != nil {
		t.Fatalf("CheckDices failed: %v", err)
	}

	if !result {
		t.Fatal("expected true for return to section")
	}
}

func TestCheck_WithReturnToSection_Fail(t *testing.T) {
	player := createTestPlayer()
	section := createTestSectionWithoutDices()
	player.Section = section

	result, err := goto_check.CheckDices(player, section)
	if err != nil {
		t.Fatalf("CheckDices failed: %v", err)
	}

	if result {
		t.Fatal("expected false for return to section")
	}
}
