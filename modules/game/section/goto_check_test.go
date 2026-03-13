package section

import (
	"gamebook-backend/database/entities"
	"testing"

	"github.com/google/uuid"
)

func TestCheck_NoDices(t *testing.T) {
	transition := entities.Transition{
		ID: uuid.New(),
	}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	result, err := Check(transition, nil, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Error("Expected true when no dices")
	}
}

func TestCheck_SingleDice_Pass(t *testing.T) {
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:   uuid.New(),
		Dice: &diceCheck,
	}
	dices := &entities.Dice{DiceFirst: 5}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Error("Expected true for dice >= 5 with roll 5")
	}
}

func TestCheck_SingleDice_Fail(t *testing.T) {
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:   uuid.New(),
		Dice: &diceCheck,
	}
	dices := &entities.Dice{DiceFirst: 4}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Error("Expected false for dice >= 5 with roll 4")
	}
}

func TestCheck_TwoDice_Pass(t *testing.T) {
	diceCheck := []string{">=8"}
	transition := entities.Transition{
		ID:    uuid.New(),
		Dices: &diceCheck,
	}
	dices := &entities.Dice{DiceFirst: 5, DiceSecond: 4}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Error("Expected true for dices >= 8 with sum 9")
	}
}

func TestCheck_TwoDice_Fail(t *testing.T) {
	diceCheck := []string{">=10"}
	transition := entities.Transition{
		ID:    uuid.New(),
		Dices: &diceCheck,
	}
	dices := &entities.Dice{DiceFirst: 4, DiceSecond: 5}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Error("Expected false for dices >= 10 with sum 9")
	}
}

func TestCheck_MultipleDiceConditions(t *testing.T) {
	diceCheck := []string{">=5", "<=10"}
	transition := entities.Transition{
		ID:    uuid.New(),
		Dices: &diceCheck,
	}
	dices := &entities.Dice{DiceFirst: 5, DiceSecond: 5}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Error("Expected true for dices >=5 && <=10 with sum 10")
	}
}

func TestCheck_WithBagCondition_Pass(t *testing.T) {
	bagCondition := "Bag.Sword"
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:        uuid.New(),
		Dice:      &diceCheck,
		Condition: &bagCondition,
	}
	dices := &entities.Dice{DiceFirst: 5}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
		Bag:    []entities.Bag{{Name: "Sword"}, {Name: "Shield"}},
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Error("Expected true with sword in bag and dice pass")
	}
}

func TestCheck_WithBagCondition_FailBag(t *testing.T) {
	bagCondition := "Bag.Sword"
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:        uuid.New(),
		Dice:      &diceCheck,
		Condition: &bagCondition,
	}
	dices := &entities.Dice{DiceFirst: 5}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
		Bag:    []entities.Bag{{Name: "Shield"}, {Name: "Potion"}},
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Error("Expected false without sword in bag")
	}
}

func TestCheck_WithBagCondition_FailDice(t *testing.T) {
	bagCondition := "Bag.Sword"
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:        uuid.New(),
		Dice:      &diceCheck,
		Condition: &bagCondition,
	}
	dices := &entities.Dice{DiceFirst: 3}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
		Bag:    []entities.Bag{{Name: "Sword"}, {Name: "Shield"}},
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Error("Expected false with sword in bag but dice fail")
	}
}

func TestCheck_WithHistoryCondition_Pass(t *testing.T) {
	historyCondition := "History.10"
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:        uuid.New(),
		Dice:      &diceCheck,
		Condition: &historyCondition,
	}
	dices := &entities.Dice{DiceFirst: 5}
	sectionID := uuid.New()
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
		PlayerSection: []entities.PlayerSection{
			{SectionID: sectionID, Section: entities.Section{Number: 10}},
		},
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Error("Expected true with section 10 in history and dice pass")
	}
}

func TestCheck_WithHistoryCondition_Fail(t *testing.T) {
	historyCondition := "History.10"
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:        uuid.New(),
		Dice:      &diceCheck,
		Condition: &historyCondition,
	}
	dices := &entities.Dice{DiceFirst: 5}
	sectionID := uuid.New()
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
		PlayerSection: []entities.PlayerSection{
			{SectionID: sectionID, Section: entities.Section{Number: 5}},
		},
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Error("Expected false without section 10 in history")
	}
}

func TestCheck_WithReturnToSection_Pass(t *testing.T) {
	returnToSection := uint(5)
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:           uuid.New(),
		Dice:         &diceCheck,
		PlayerChange: &entities.PlayerChange{ReturnToSection: &returnToSection},
	}
	dices := &entities.Dice{DiceFirst: 5}
	player := entities.Player{
		ID:              uuid.New(),
		Health:          20,
		ReturnToSection: 5,
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Error("Expected true when ReturnToSection matches")
	}
}

func TestCheck_WithReturnToSection_Fail(t *testing.T) {
	returnToSection := uint(5)
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:           uuid.New(),
		Dice:         &diceCheck,
		PlayerChange: &entities.PlayerChange{ReturnToSection: &returnToSection},
	}
	dices := &entities.Dice{DiceFirst: 5}
	player := entities.Player{
		ID:              uuid.New(),
		Health:          20,
		ReturnToSection: 0,
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Error("Expected false when ReturnToSection doesn't match")
	}
}

func TestCheck_WithMultipleConditions_Pass(t *testing.T) {
	bagCondition := "Bag.Sword"
	returnToSection := uint(5)
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:           uuid.New(),
		Dice:         &diceCheck,
		Condition:    &bagCondition,
		PlayerChange: &entities.PlayerChange{ReturnToSection: &returnToSection},
	}
	dices := &entities.Dice{DiceFirst: 5}
	player := entities.Player{
		ID:              uuid.New(),
		Health:          20,
		Bag:             []entities.Bag{{Name: "Sword"}},
		ReturnToSection: 5,
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Error("Expected true when all conditions pass")
	}
}

func TestCheck_WithMultipleConditions_FailOne(t *testing.T) {
	bagCondition := "Bag.Sword"
	returnToSection := uint(5)
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:           uuid.New(),
		Dice:         &diceCheck,
		Condition:    &bagCondition,
		PlayerChange: &entities.PlayerChange{ReturnToSection: &returnToSection},
	}
	dices := &entities.Dice{DiceFirst: 3}
	player := entities.Player{
		ID:              uuid.New(),
		Health:          20,
		Bag:             []entities.Bag{{Name: "Sword"}},
		ReturnToSection: 5,
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Error("Expected false when dice fails")
	}
}

func TestCheck_InvalidExpression(t *testing.T) {
	diceCheck := []string{">==5"} // Invalid syntax
	transition := entities.Transition{
		ID:   uuid.New(),
		Dice: &diceCheck,
	}
	dices := &entities.Dice{DiceFirst: 5}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	_, err := Check(transition, dices, player)

	if err == nil {
		t.Error("Expected error for invalid expression")
	}
}

func TestCheck_ExactDiceMatch(t *testing.T) {
	diceCheck := []string{"==5"}
	transition := entities.Transition{
		ID:   uuid.New(),
		Dice: &diceCheck,
	}
	dices := &entities.Dice{DiceFirst: 5}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Error("Expected true for exact dice match")
	}
}

func TestCheck_LessThanCondition(t *testing.T) {
	diceCheck := []string{"<5"}
	transition := entities.Transition{
		ID:   uuid.New(),
		Dice: &diceCheck,
	}
	dices := &entities.Dice{DiceFirst: 4}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Error("Expected true for dice < 5 with roll 4")
	}
}

func TestCheck_NegatedBagCondition(t *testing.T) {
	bagCondition := "!Bag.Sword"
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:        uuid.New(),
		Dice:      &diceCheck,
		Condition: &bagCondition,
	}
	dices := &entities.Dice{DiceFirst: 5}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
		Bag:    []entities.Bag{{Name: "Shield"}, {Name: "Potion"}},
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Error("Expected true without sword in bag (negated condition)")
	}
}

func TestCheck_NegatedBagCondition_Fail(t *testing.T) {
	bagCondition := "!Bag.Sword"
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:        uuid.New(),
		Dice:      &diceCheck,
		Condition: &bagCondition,
	}
	dices := &entities.Dice{DiceFirst: 5}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
		Bag:    []entities.Bag{{Name: "Sword"}, {Name: "Shield"}},
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Error("Expected false with sword in bag (negated condition)")
	}
}

func TestCheck_NegatedHistoryCondition(t *testing.T) {
	historyCondition := "!History.10"
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:        uuid.New(),
		Dice:      &diceCheck,
		Condition: &historyCondition,
	}
	dices := &entities.Dice{DiceFirst: 5}
	sectionID := uuid.New()
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
		PlayerSection: []entities.PlayerSection{
			{SectionID: sectionID, Section: entities.Section{Number: 5}},
		},
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Error("Expected true without section 10 in history (negated condition)")
	}
}

func TestCheck_OrConditions(t *testing.T) {
	bagCondition := "Bag.Sword || Bag.Axe"
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:        uuid.New(),
		Dice:      &diceCheck,
		Condition: &bagCondition,
	}
	dices := &entities.Dice{DiceFirst: 5}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
		Bag:    []entities.Bag{{Name: "Shield"}, {Name: "Axe"}, {Name: "Potion"}},
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Error("Expected true with axe in bag (OR condition)")
	}
}

func TestCheck_OrConditions_Fail(t *testing.T) {
	bagCondition := "Bag.Sword || Bag.Axe"
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:        uuid.New(),
		Dice:      &diceCheck,
		Condition: &bagCondition,
	}
	dices := &entities.Dice{DiceFirst: 5}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
		Bag:    []entities.Bag{{Name: "Shield"}, {Name: "Potion"}},
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Error("Expected false without sword or axe in bag")
	}
}

func TestCheck_AndConditions(t *testing.T) {
	bagCondition := "Bag.Sword && Bag.Shield"
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:        uuid.New(),
		Dice:      &diceCheck,
		Condition: &bagCondition,
	}
	dices := &entities.Dice{DiceFirst: 5}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
		Bag:    []entities.Bag{{Name: "Sword"}, {Name: "Shield"}, {Name: "Potion"}},
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Error("Expected true with both sword and shield in bag (AND condition)")
	}
}

func TestCheck_AndConditions_Fail(t *testing.T) {
	bagCondition := "Bag.Sword && Bag.Shield"
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:        uuid.New(),
		Dice:      &diceCheck,
		Condition: &bagCondition,
	}
	dices := &entities.Dice{DiceFirst: 5}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
		Bag:    []entities.Bag{{Name: "Sword"}, {Name: "Potion"}},
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Error("Expected false without shield in bag (AND condition)")
	}
}

func TestCheck_ComplexConditions(t *testing.T) {
	bagCondition := "Bag.Sword && Bag.Shield || Bag.MagicWand"
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:        uuid.New(),
		Dice:      &diceCheck,
		Condition: &bagCondition,
	}
	dices := &entities.Dice{DiceFirst: 5}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
		Bag:    []entities.Bag{{Name: "MagicWand"}, {Name: "Potion"}},
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Error("Expected true with magic wand in bag (complex OR/AND condition)")
	}
}

func TestCheck_ComplexConditions_Fail(t *testing.T) {
	bagCondition := "Bag.Sword && Bag.Shield || Bag.MagicWand"
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:        uuid.New(),
		Dice:      &diceCheck,
		Condition: &bagCondition,
	}
	dices := &entities.Dice{DiceFirst: 5}
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
		Bag:    []entities.Bag{{Name: "Sword"}, {Name: "Potion"}},
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Error("Expected false without matching condition")
	}
}

func TestCheck_MixedBagAndHistoryConditions(t *testing.T) {
	condition := "Bag.Sword && History.10"
	diceCheck := []string{">=5"}
	transition := entities.Transition{
		ID:        uuid.New(),
		Dice:      &diceCheck,
		Condition: &condition,
	}
	dices := &entities.Dice{DiceFirst: 5}
	sectionID := uuid.New()
	player := entities.Player{
		ID:     uuid.New(),
		Health: 20,
		Bag:    []entities.Bag{{Name: "Sword"}, {Name: "Shield"}},
		PlayerSection: []entities.PlayerSection{
			{SectionID: sectionID, Section: entities.Section{Number: 10}},
		},
	}

	result, err := Check(transition, dices, player)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Error("Expected true with sword and history 10")
	}
}
