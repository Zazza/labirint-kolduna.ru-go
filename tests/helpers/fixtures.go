package helpers

import (
	"gamebook-backend/database/entities"
	gameDto "gamebook-backend/modules/game/dto"

	"github.com/google/uuid"
)

const (
	TestEnemyAliasGoblin   = "goblin"
	TestEnemyAliasOrc      = "orc"
	TestEnemyAliasDragon   = "dragon"
	TestEnemyAliasSkeleton = "skeleton"
	TestEnemyAliasZombie   = "zombie"

	TestWeaponSword = "Sword"
	TestWeaponAxe   = "Axe"
	TestWeaponBow   = "Bow"
	TestWeaponMagic = "Magic"
	TestWeaponFist  = "Fist"

	TestMedsPotion = "Potion"
	TestMedsHerb   = "Herb"

	TestBagItemTorch = "Torch"
	TestBagItemRope  = "Rope"
	TestBagItemKey   = "Key"
)

type TestSection struct {
	Number       uint
	Type         entities.SectionType
	Text         string
	Dices        bool
	BattleStart  *string
	BattleSteps  []*string
	Enemies      []string
	TransitionTo uint
}

type CreateSectionRequest struct {
	Number      uint
	Type        entities.SectionType
	Text        string
	Dices       bool
	BattleStart *string
	BattleSteps []*string
	Transitions []TransitionConfig
}

type TransitionConfig struct {
	TargetSection uint
	Text          string
	IsBattleWin   *bool
	PlayerChange  *entities.PlayerChange
}

var DefaultTestSections = []CreateSectionRequest{
	{
		Number: 0,
		Type:   entities.SectionType("normal"),
		Text:   "You are standing at the entrance of a dark dungeon.",
		Dices:  false,
		Transitions: []TransitionConfig{
			{
				TargetSection: 1,
				Text:          "Enter the dungeon",
			},
		},
	},
	{
		Number:      1,
		Type:        gameDto.SectionTypeBattle,
		Text:        "A goblin appears!",
		Dices:       false,
		BattleStart: stringPtr("dices"),
		Transitions: []TransitionConfig{
			{
				TargetSection: 2,
				Text:          "Continue",
				IsBattleWin:   boolPtr(true),
			},
		},
	},
	{
		Number: 2,
		Type:   entities.SectionType("normal"),
		Text:   "You continue through the dungeon.",
		Dices:  false,
		Transitions: []TransitionConfig{
			{
				TargetSection: 3,
				Text:          "Go left",
			},
			{
				TargetSection: 4,
				Text:          "Go right",
			},
		},
	},
	{
		Number: 3,
		Type:   entities.SectionType("normal"),
		Text:   "You found a treasure chest!",
		Dices:  true,
		Transitions: []TransitionConfig{
			{
				TargetSection: 0,
				Text:          "Return to entrance",
				PlayerChange: &entities.PlayerChange{
					Gold: stringPtr("+10"),
				},
			},
		},
	},
	{
		Number: 4,
		Type:   entities.SectionType("normal"),
		Text:   "A dead end. You must turn back.",
		Dices:  false,
		Transitions: []TransitionConfig{
			{
				TargetSection: 2,
				Text:          "Go back",
			},
		},
	},
}

type CreateEnemyRequest struct {
	Alias       string
	Name        string
	Damage      uint
	MinDiceHits uint
	Health      uint
	Defence     uint
	PlayerArmor bool
}

var DefaultTestEnemies = []CreateEnemyRequest{
	{
		Alias:       TestEnemyAliasGoblin,
		Name:        "Гоблин",
		Damage:      5,
		MinDiceHits: 6,
		Health:      10,
		Defence:     2,
		PlayerArmor: true,
	},
	{
		Alias:       TestEnemyAliasOrc,
		Name:        "Орк",
		Damage:      8,
		MinDiceHits: 5,
		Health:      15,
		Defence:     3,
		PlayerArmor: true,
	},
	{
		Alias:       TestEnemyAliasSkeleton,
		Name:        "Скелет",
		Damage:      4,
		MinDiceHits: 5,
		Health:      8,
		Defence:     1,
		PlayerArmor: true,
	},
}

type CreatePlayerRequest struct {
	UserID    uuid.UUID
	SectionID uuid.UUID
	Health    uint
	Gold      uint
	Weapons   []entities.Weapons
}

func DefaultCreatePlayerRequest(userID, sectionID uuid.UUID) CreatePlayerRequest {
	return CreatePlayerRequest{
		UserID:    userID,
		SectionID: sectionID,
		Health:    20,
		Gold:      50,
		Weapons: []entities.Weapons{
			{Name: TestWeaponSword, Damage: 10, MinCubeHit: 4, Item: "sword"},
		},
	}
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func uintPtr(u uint) *uint {
	return &u
}

func intPtr(i int) *int {
	return &i
}
