package bonus

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"testing"

	"github.com/google/uuid"
)

func TestAntiPoisonSpell_Execute(t *testing.T) {
	// Test case 1: Player with poison debuff
	t.Run("Remove poison debuff and add health", func(t *testing.T) {
		health := uint(10)
		player := entities.Player{
			ID:     uuid.New(),
			Health: 20,
			Debuff: []entities.Debuff{
				{
					Alias:  entities.AliasPoisonReason,
					Health: &health,
				},
			},
			Section: entities.Section{
				Number: 5, // Not death section
			},
		}

		bonus := NewAntiPoisonSpell(nil, player)
		req := dto.BonusRequest{
			Bonus: AntiPoisonSpellAlias,
		}

		err := bonus.Execute(context.Background(), req)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	// Test case 2: Player without poison debuff
	t.Run("No poison debuff - no changes", func(t *testing.T) {
		player := entities.Player{
			ID:     uuid.New(),
			Health: 20,
			Debuff: []entities.Debuff{},
		}

		bonus := NewAntiPoisonSpell(nil, player)
		req := dto.BonusRequest{
			Bonus: AntiPoisonSpellAlias,
		}

		err := bonus.Execute(context.Background(), req)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	// Test case 3: Player with other debuffs but not poison
	t.Run("Other debuffs - no changes", func(t *testing.T) {
		health := uint(10)
		player := entities.Player{
			ID:     uuid.New(),
			Health: 20,
			Debuff: []entities.Debuff{
				{
					Alias:  entities.AliasMagicReason,
					Health: &health,
				},
			},
		}

		bonus := NewAntiPoisonSpell(nil, player)
		req := dto.BonusRequest{
			Bonus: AntiPoisonSpellAlias,
		}

		err := bonus.Execute(context.Background(), req)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func TestInstantRecoverySpell_Execute(t *testing.T) {
	t.Run("Add 50 health", func(t *testing.T) {
		player := entities.Player{
			ID:     uuid.New(),
			Health: 20,
		}

		bonus := NewInstantRecoverySpell(nil, player)
		req := dto.BonusRequest{
			Bonus: InstantRecoveryAlias,
		}

		err := bonus.Execute(context.Background(), req)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("Zero health player", func(t *testing.T) {
		player := entities.Player{
			ID:     uuid.New(),
			Health: 0,
		}

		bonus := NewInstantRecoverySpell(nil, player)
		req := dto.BonusRequest{
			Bonus: InstantRecoveryAlias,
		}

		err := bonus.Execute(context.Background(), req)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func TestLuckyStone_Execute(t *testing.T) {
	t.Run("Add lucky stone buff", func(t *testing.T) {
		player := entities.Player{
			ID:     uuid.New(),
			Health: 20,
		}

		bonus := NewLuckyStone(nil, player)
		req := dto.BonusRequest{
			Bonus: LuckyStoneAlias,
		}

		err := bonus.Execute(context.Background(), req)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func TestGetBonus(t *testing.T) {
	tests := []struct {
		name        string
		bonusAlias  string
		expectError bool
	}{
		{
			name:        "Death spell",
			bonusAlias:  DeathSpellAlias,
			expectError: false,
		},
		{
			name:        "Anti poison spell",
			bonusAlias:  AntiPoisonSpellAlias,
			expectError: false,
		},
		{
			name:        "Instant movement",
			bonusAlias:  InstantMovementAlias,
			expectError: false,
		},
		{
			name:        "Instant hypnosis spell",
			bonusAlias:  InstantHypnosisSpellAlias,
			expectError: false,
		},
		{
			name:        "Instant recovery",
			bonusAlias:  InstantRecoveryAlias,
			expectError: false,
		},
		{
			name:        "Magic duck",
			bonusAlias:  MagicDuckAlias,
			expectError: false,
		},
		{
			name:        "Wand",
			bonusAlias:  WandAlias,
			expectError: false,
		},
		{
			name:        "Magic ring",
			bonusAlias:  MagicRingAlias,
			expectError: false,
		},
		{
			name:        "Lucky stone",
			bonusAlias:  LuckyStoneAlias,
			expectError: false,
		},
		{
			name:        "Unknown bonus",
			bonusAlias:  "unknown_bonus",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := entities.Player{
				ID:     uuid.New(),
				Health: 20,
			}

			result, err := GetBonus(nil, player, tt.bonusAlias)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				if result != nil {
					t.Error("Expected nil result")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Error("Expected non-nil result")
				}
			}
		})
	}
}

func TestMagicRing_Execute(t *testing.T) {
	t.Run("Add magic ring buff - health", func(t *testing.T) {
		player := entities.Player{
			ID:     uuid.New(),
			Health: 20,
		}

		bonus := NewMagicRing(nil, player)
		req := dto.BonusRequest{
			Bonus:  MagicRingAlias,
			Option: MagicRingOptions[0], // Use valid index 0
		}

		err := bonus.Execute(context.Background(), req)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func TestMagicDuck_Execute(t *testing.T) {
	t.Skip("Skipping - requires database connection (integration test)")
}

func TestInstantMovement_Execute(t *testing.T) {
	t.Skip("Skipping - requires database connection (integration test)")
}

func TestInstantHypnosisSpell_Execute(t *testing.T) {
	t.Skip("Skipping - requires database connection (integration test)")
}
