package dto

import (
	"database/sql"
	"errors"
	"gamebook-backend/database/entities"

	"github.com/google/uuid"
)

const (
	MessageFailedGetPlayer          = "failed get player"
	MessageFailedCurrentRequest     = "failed current request"
	MessageFailedRollTheDiceRequest = "failed roll the dice request"
	MessageFailedActionRequest      = "failed action request"
	MessageFailedGetMap             = "failed get map"

	MessageSuccessCurrentRequest = "success current request"
	MessageSuccessActionRequest  = "success action request"
	MessageSuccessGetMap         = "success get map"

	MessageMapIngredientNotFound = "map ingredients not found"

	SectionTypeChoice entities.SectionType = "choice"
	SectionTypeBattle entities.SectionType = "battle"
	SectionTypeSleepy entities.SectionType = "sleepy"

	SectionDeath = 9

	ResultTrue  ActionResult = true
	ResultFalse ActionResult = false

	CHOICE = "choice"
	BATTLE = "battle"

	StepTypeNormal = "normal"
)

var (
	ErrSectionNotFound             = errors.New("section not found")
	ErrBattleNotFound              = errors.New("battle not found")
	ErrBattleWeaponNotDefined      = errors.New("battle weapon not defined")
	ErrStepNotDefined              = errors.New("step not defined")
	ErrEnemyNotDefined             = errors.New("enemy not defined")
	ErrBattleStartNotDefined       = errors.New("battle start not defined")
	ErrGetActivityBySectionId      = errors.New("get activity by section id error")
	ErrCustomSectionNotFound       = errors.New("custom section not found")
	ErrPlayerChangesDicesFirstChar = errors.New("first character must be '-' or '+'")
)

type (
	CurrentRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	PlayerInfoBonus struct {
		Alias  string                      `json:"alias"`
		Name   string                      `json:"name"`
		Option *entities.PlayerBonusOption `json:"option,omitempty"`
	}

	PlayerInfo struct {
		Health uint              `json:"health"`
		Meds   uint              `json:"meds"`
		Gold   uint              `json:"gold"`
		Bonus  []PlayerInfoBonus `json:"bonus"`
	}

	ProfileWeapons struct {
		Name       string
		Damage     uint
		MinCubeHit uint
		Item       string
		Count      uint `json:"Count,omitempty"`
	}
	ProfileMeds struct {
		Name  string
		Item  string
		Count int
	}

	ProfileDebuff struct {
		Health     *uint                       `json:"Health,omitempty"`
		MinCubeHit *uint                       `json:"MinCubeHit,omitempty"`
		Duration   *uint                       `json:"Duration,omitempty"`
		Alias      *entities.BuffOrDebuffAlias `json:"Alias,omitempty"`
		Name       *string                     `json:"Name,omitempty"`
	}

	ProfileBuff struct {
		Health     *uint                       `json:"Health,omitempty"`
		MinCubeHit *uint                       `json:"MinCubeHit,omitempty"`
		Duration   *uint                       `json:"Duration,omitempty"`
		Alias      *entities.BuffOrDebuffAlias `json:"Alias,omitempty"`
		Name       *string                     `json:"Name,omitempty"`
	}

	CurrentResponse struct {
		Section      uint                 `json:"section"`
		Text         string               `json:"text"`
		Type         entities.SectionType `json:"type"`
		Transitions  []TransitionDTO      `json:"transitions"`
		Dices        []uint               `json:"dices"`
		Choice       entities.Choice      `json:"choice"`
		RollTheDices bool                 `json:"roll_the_dices"`
		Player       PlayerInfo           `json:"player"`
		MapAvailable bool                 `json:"map_available"`
	}

	ProfileResponse struct {
		Health    uint              `json:"health"`
		MaxHealth uint              `json:"max_health"`
		Weapons   []ProfileWeapons  `json:"weapons"`
		Meds      ProfileMeds       `json:"meds"`
		Bag       []entities.Bag    `json:"bag"`
		Debuff    []ProfileDebuff   `json:"debuff"`
		Buff      []ProfileBuff     `json:"buff"`
		Gold      uint              `json:"gold"`
		Bonus     []PlayerInfoBonus `json:"bonus"`
	}

	ActionResult bool

	ActionRequest struct {
		Transition uuid.UUID `json:"transitionID"`
		Data       []string  `json:"data"`
		Weapon     string    `json:"weapon"`
	}

	ActionResponse struct {
		Result  ActionResult    `json:"result"`
		Error   string          `json:"error"`
		Content string          `json:"content,omitempty"`
		Actions []TransitionDTO `json:"actions,omitempty"`
	}

	PlayerActivityResponse struct {
		Section        *uint
		GameCubeFirst  *sql.NullInt16
		GameCubeSecond *sql.NullInt16
	}

	TransitionDTO struct {
		Text             string
		TransitionID     uuid.UUID
		Visited          bool
		Weapon           string
		SleepyTransition bool
		Bribe            bool `json:"Bribe,omitempty"`
		Death            bool `json:"Death,omitempty"`
	}

	ActivityResponse struct {
		Text         string
		Type         entities.SectionType
		Transitions  []TransitionDTO
		Player       PlayerActivityResponse
		RollTheDices bool
	}

	SectionTemplate struct {
		Section      uint
		Text         string
		Type         entities.SectionType
		Transitions  []TransitionDTO
		Choice       entities.Choice
		Player       PlayerInfo
		Dices        []uint
		RollTheDices bool
	}
)
