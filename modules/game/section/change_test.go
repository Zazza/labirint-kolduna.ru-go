package section

import (
	"context"
	"gamebook-backend/database/entities"
	"testing"

	"github.com/google/uuid"
)

func TestChange_HealthModification(t *testing.T) {
	tests := []struct {
		name         string
		healthChange *string
		playerHealth uint
		diceFirst    uint
		diceSecond   uint
		expectedMin  uint
		expectedMax  uint
		expectError  bool
	}{
		{
			name:         "Add health",
			healthChange: stringPtr("+5"),
			playerHealth: 20,
			diceFirst:    0,
			diceSecond:   0,
			expectedMin:  25,
			expectedMax:  25,
			expectError:  false,
		},
		{
			name:         "Subtract health",
			healthChange: stringPtr("-5"),
			playerHealth: 20,
			diceFirst:    0,
			diceSecond:   0,
			expectedMin:  15,
			expectedMax:  15,
			expectError:  false,
		},
		{
			name:         "health *0",
			healthChange: stringPtr("*0"),
			playerHealth: 10,
			diceFirst:    0,
			diceSecond:   0,
			expectedMin:  0,
			expectedMax:  0,
			expectError:  false,
		},
		{
			name:         "Multiply health",
			healthChange: stringPtr("*2"),
			playerHealth: 10,
			diceFirst:    0,
			diceSecond:   0,
			expectedMin:  20,
			expectedMax:  20,
			expectError:  false,
		},
		{
			name:         "No health change",
			healthChange: nil,
			playerHealth: 20,
			diceFirst:    0,
			diceSecond:   0,
			expectedMin:  20,
			expectedMax:  20,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transition := entities.Transition{
				PlayerChange: &entities.PlayerChange{
					Health: tt.healthChange,
				},
			}
			player := entities.Player{
				ID:     uuid.New(),
				Health: tt.playerHealth,
			}

			result, err := Change(context.Background(), transition, player, nil)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result.Player.Health < tt.expectedMin || result.Player.Health > tt.expectedMax {
					t.Errorf("Expected health between %d and %d, got %d", tt.expectedMin, tt.expectedMax, result.Player.Health)
				}
			}
		})
	}
}

func TestChange_WeaponsModification(t *testing.T) {
	change := "+1"
	item := "sword"

	tests := []struct {
		name          string
		weaponChange  *[]entities.PlayerChangeWeapon
		playerWeapons []entities.Weapons
		expectedCount int
	}{
		{
			name: "Add weapon count",
			weaponChange: &[]entities.PlayerChangeWeapon{
				{Item: &item, Change: &change},
			},
			playerWeapons: []entities.Weapons{
				{Item: "sword", Count: 1},
			},
			expectedCount: 2,
		},
		{
			name: "Subtract weapon count",
			weaponChange: &[]entities.PlayerChangeWeapon{
				{Item: &item, Change: stringPtr("-1")},
			},
			playerWeapons: []entities.Weapons{
				{Item: "sword", Count: 3},
			},
			expectedCount: 2,
		},
		{
			name:         "No weapon change",
			weaponChange: nil,
			playerWeapons: []entities.Weapons{
				{Item: "sword", Count: 1},
			},
			expectedCount: 1,
		},
		{
			name: "Weapon not found - no change",
			weaponChange: &[]entities.PlayerChangeWeapon{
				{Item: stringPtr("axe"), Change: &change},
			},
			playerWeapons: []entities.Weapons{
				{Item: "sword", Count: 1},
			},
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transition := entities.Transition{
				PlayerChange: &entities.PlayerChange{
					Weapons: tt.weaponChange,
				},
			}
			player := entities.Player{
				ID:      uuid.New(),
				Health:  20,
				Weapons: tt.playerWeapons,
			}

			result, err := Change(context.Background(), transition, player, nil, nil)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if len(result.Weapons) > 0 && result.Weapons[0].Count != uint(tt.expectedCount) {
				t.Errorf("Expected weapon count %d, got %d", tt.expectedCount, result.Weapons[0].Count)
			}
		})
	}
}

func TestChange_BagModification(t *testing.T) {
	tests := []struct {
		name        string
		bagItems    []entities.Bag
		playerBag   []entities.Bag
		expectedLen int
	}{
		{
			name:        "Add items to bag",
			bagItems:    []entities.Bag{{Name: "Sword"}, {Name: "Shield"}},
			playerBag:   []entities.Bag{{Name: "Potion"}},
			expectedLen: 3,
		},
		{
			name:        "No bag change",
			bagItems:    nil,
			playerBag:   []entities.Bag{{Name: "Potion"}},
			expectedLen: 1,
		},
		{
			name:        "Empty bag items",
			bagItems:    nil,
			playerBag:   []entities.Bag{{Name: "Potion"}},
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transition := entities.Transition{
				PlayerChange: &entities.PlayerChange{
					Bag: &tt.bagItems,
				},
			}
			player := entities.Player{
				ID:     uuid.New(),
				Health: 20,
				Bag:    tt.playerBag,
			}

			result, err := Change(context.Background(), transition, player, nil, nil)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if len(result.Bag) != tt.expectedLen {
				t.Errorf("Expected bag length %d, got %d", tt.expectedLen, len(result.Bag))
			}
		})
	}
}

func TestChange_GoldModification(t *testing.T) {
	tests := []struct {
		name        string
		goldChange  *string
		playerGold  uint
		expected    uint
		expectError bool
	}{
		{
			name:       "Add gold",
			goldChange: stringPtr("+10"),
			playerGold: 100,
			expected:   110,
		},
		{
			name:       "Subtract gold",
			goldChange: stringPtr("-10"),
			playerGold: 100,
			expected:   90,
		},
		{
			name:       "No gold change",
			goldChange: nil,
			playerGold: 100,
			expected:   100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transition := entities.Transition{
				PlayerChange: &entities.PlayerChange{
					Gold: tt.goldChange,
				},
			}
			player := entities.Player{
				ID:     uuid.New(),
				Health: 20,
				Gold:   tt.playerGold,
			}

			result, err := Change(context.Background(), transition, player, nil, nil)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result.Gold != tt.expected {
					t.Errorf("Expected gold %d, got %d", tt.expected, result.Gold)
				}
			}
		})
	}
}

func TestChange_BonusModification(t *testing.T) {
	alias := "Strength"
	name := "Strength Boost"

	tests := []struct {
		name          string
		bonusChange   *[]entities.PlayerBonus
		playerBonuses []entities.PlayerBonus
		expectedLen   int
	}{
		{
			name: "Add bonus",
			bonusChange: &[]entities.PlayerBonus{
				{Alias: &alias, Name: &name},
			},
			playerBonuses: []entities.PlayerBonus{},
			expectedLen:   1,
		},
		{
			name:          "No bonus change",
			bonusChange:   nil,
			playerBonuses: []entities.PlayerBonus{},
			expectedLen:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transition := entities.Transition{
				PlayerChange: &entities.PlayerChange{
					Bonus: tt.bonusChange,
				},
			}
			player := entities.Player{
				ID:     uuid.New(),
				Health: 20,
				Bonus:  tt.playerBonuses,
			}

			result, err := Change(context.Background(), transition, player, nil, nil)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if len(result.Bonus) != tt.expectedLen {
				t.Errorf("Expected bonus count %d, got %d", tt.expectedLen, len(result.Bonus))
			}
		})
	}
}

func TestChange_ReturnToSection(t *testing.T) {
	tests := []struct {
		name            string
		returnToSection *uint
		playerReturn    uint
		expected        uint
	}{
		{
			name:            "Set return to section",
			returnToSection: uintPtr(5),
			playerReturn:    0,
			expected:        5,
		},
		{
			name:            "No return to section change",
			returnToSection: nil,
			playerReturn:    3,
			expected:        3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transition := entities.Transition{
				PlayerChange: &entities.PlayerChange{
					ReturnToSection: tt.returnToSection,
				},
			}
			player := entities.Player{
				ID:              uuid.New(),
				Health:          20,
				ReturnToSection: tt.playerReturn,
			}

			result, err := Change(context.Background(), transition, player, nil, nil)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result.ReturnToSection != tt.expected {
				t.Errorf("Expected ReturnToSection %d, got %d", tt.expected, result.ReturnToSection)
			}
		})
	}
}

func TestChange_InvalidExpression(t *testing.T) {
	tests := []struct {
		name         string
		healthChange *string
		expectError  bool
	}{
		{
			name:         "Invalid health expression",
			healthChange: stringPtr("+invalid"),
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transition := entities.Transition{
				PlayerChange: &entities.PlayerChange{
					Health: tt.healthChange,
				},
			}
			player := entities.Player{
				ID:     uuid.New(),
				Health: 20,
			}

			_, err := Change(context.Background(), transition, player, nil, nil)

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}

func TestChange_ComplexPlayer(t *testing.T) {
	change := "+1"
	item := "sword"
	alias := "Strength"
	name := "Strength Boost"

	transition := entities.Transition{
		PlayerChange: &entities.PlayerChange{
			Health: stringPtr("+5"),
			Gold:   stringPtr("+10"),
			Bag:    &[]entities.Bag{{Name: "Action Item"}},
			Bonus:  &[]entities.PlayerBonus{{Alias: &alias, Name: &name}},
			Weapons: &[]entities.PlayerChangeWeapon{
				{Item: &item, Change: &change},
			},
		},
	}

	player := entities.Player{
		ID:              uuid.New(),
		Health:          20,
		Gold:            100,
		Bag:             []entities.Bag{{Name: "Old Item"}},
		Bonus:           []entities.PlayerBonus{},
		Weapons:         []entities.Weapons{{Item: "sword", Count: 1}},
		ReturnToSection: 0,
	}

	result, err := Change(context.Background(), transition, player, nil, nil)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result.Health != 25 {
		t.Errorf("Expected health 25, got %d", result.Health)
	}
	if result.Gold != 110 {
		t.Errorf("Expected gold 110, got %d", result.Gold)
	}
	if len(result.Bag) != 2 {
		t.Errorf("Expected 2 bag items, got %d", len(result.Bag))
	}
	if len(result.Bonus) != 1 {
		t.Errorf("Expected 1 bonus, got %d", len(result.Bonus))
	}
	if len(result.Weapons) > 0 && result.Weapons[0].Count != 2 {
		t.Errorf("Expected weapon count 2, got %d", result.Weapons[0].Count)
	}
}

func TestCheckConditions(t *testing.T) {
	tests := []struct {
		name          string
		bagItems      *string
		playerBag     []entities.Bag
		playerSection []entities.PlayerSection
		expected      bool
	}{
		{
			name:          "Nil condition - always true",
			bagItems:      nil,
			playerBag:     []entities.Bag{},
			playerSection: []entities.PlayerSection{},
			expected:      true,
		},
		{
			name:          "Bag item exists - true",
			bagItems:      stringPtr("Bag.Potion"),
			playerBag:     []entities.Bag{{Name: "Potion"}, {Name: "Sword"}},
			playerSection: []entities.PlayerSection{},
			expected:      true,
		},
		{
			name:          "Bag item does not exist - false",
			bagItems:      stringPtr("Bag.Potion"),
			playerBag:     []entities.Bag{{Name: "Sword"}},
			playerSection: []entities.PlayerSection{},
			expected:      false,
		},
		{
			name:          "Negated bag item exists - false",
			bagItems:      stringPtr("!Bag.Potion"),
			playerBag:     []entities.Bag{{Name: "Potion"}, {Name: "Sword"}},
			playerSection: []entities.PlayerSection{},
			expected:      false,
		},
		{
			name:          "Negated bag item does not exist - true",
			bagItems:      stringPtr("!Bag.Potion"),
			playerBag:     []entities.Bag{{Name: "Sword"}},
			playerSection: []entities.PlayerSection{},
			expected:      true,
		},
		{
			name:          "History contains section - true",
			bagItems:      stringPtr("History.5"),
			playerBag:     []entities.Bag{},
			playerSection: []entities.PlayerSection{{Section: entities.Section{Number: 5}}},
			expected:      true,
		},
		{
			name:          "History does not contain section - false",
			bagItems:      stringPtr("History.5"),
			playerBag:     []entities.Bag{},
			playerSection: []entities.PlayerSection{{Section: entities.Section{Number: 3}}},
			expected:      false,
		},
		{
			name:          "Negated history contains section - false",
			bagItems:      stringPtr("!History.5"),
			playerBag:     []entities.Bag{},
			playerSection: []entities.PlayerSection{{Section: entities.Section{Number: 5}}},
			expected:      false,
		},
		{
			name:          "OR condition - one item exists",
			bagItems:      stringPtr("Bag.Potion || Bag.Sword"),
			playerBag:     []entities.Bag{{Name: "Potion"}},
			playerSection: []entities.PlayerSection{},
			expected:      true,
		},
		{
			name:          "OR condition - no items exist",
			bagItems:      stringPtr("Bag.Potion || Bag.Sword"),
			playerBag:     []entities.Bag{},
			playerSection: []entities.PlayerSection{},
			expected:      false,
		},
		{
			name:          "AND condition - all items exist",
			bagItems:      stringPtr("Bag.Potion && Bag.Sword"),
			playerBag:     []entities.Bag{{Name: "Potion"}, {Name: "Sword"}},
			playerSection: []entities.PlayerSection{},
			expected:      true,
		},
		{
			name:          "AND condition - not all items exist",
			bagItems:      stringPtr("Bag.Potion && Bag.Sword"),
			playerBag:     []entities.Bag{{Name: "Potion"}},
			playerSection: []entities.PlayerSection{},
			expected:      false,
		},
		{
			name:          "Multiple OR with AND",
			bagItems:      stringPtr("(Bag.Potion && Bag.Shield) || Bag.Sword"),
			playerBag:     []entities.Bag{{Name: "Sword"}},
			playerSection: []entities.PlayerSection{},
			expected:      true,
		},
		{
			name:          "Multiple !AND with AND",
			bagItems:      stringPtr("(!Bag.Potion && !Bag.Shield) && (!Bag.Sword && !Bag.Shield)"),
			playerBag:     []entities.Bag{{Name: "Sword"}, {Name: "Shield"}},
			playerSection: []entities.PlayerSection{},
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckConditions(tt.bagItems, tt.playerBag, tt.playerSection)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func uintPtr(u uint) *uint {
	return &u
}
