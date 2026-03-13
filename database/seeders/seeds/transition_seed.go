package seeds

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gamebook-backend/database/entities"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransitionJSON struct {
	ID                  *string                  `json:"id"`
	TextOrder           uint                     `json:"text_order"`
	SectionNumber       uint                     `json:"section_number"`
	TargetSectionNumber uint                     `json:"target_section_number"`
	AvailableOnce       bool                     `json:"available_once"`
	Text                string                   `json:"text"`
	IsBattleWin         *bool                    `json:"is_battle_win"`
	BribeResult         *bool                    `json:"bribe_result"`
	PlayerInput         bool                     `json:"player_input"`
	Dice                *[]string                `json:"dice"`
	Dices               *[]string                `json:"dices"`
	Condition           *string                  `json:"condition"`
	PlayerChange        *entities.PlayerChange   `json:"player_change"`
	PlayerDebuff        []*entities.PlayerDebuff `json:"player_debuff"`
	PlayerBuff          []*entities.PlayerBuff   `json:"player_buff"`
}

func TransitionSeeder(db *gorm.DB) error {
	data, err := os.ReadFile(filepath.Clean("./database/seeders/json/transitions.json"))
	if err != nil {
		return fmt.Errorf("failed to read transitions.json: %w", err)
	}

	var transitionsJSON []TransitionJSON
	if err := json.Unmarshal(data, &transitionsJSON); err != nil {
		return fmt.Errorf("failed to unmarshal transitions.json: %w", err)
	}

	if err := db.AutoMigrate(&entities.Transition{}); err != nil {
		return fmt.Errorf("failed to auto migrate transitions: %w", err)
	}

	for _, transitionJSON := range transitionsJSON {
		var section entities.Section
		if err := db.Where("number = ?", transitionJSON.SectionNumber).First(&section).Error; err != nil {
			return fmt.Errorf("failed to find section %d for transition: %w", transitionJSON.SectionNumber, err)
		}

		var targetSection entities.Section
		if err := db.Where("number = ?", transitionJSON.TargetSectionNumber).First(&targetSection).Error; err != nil {
			return fmt.Errorf("failed to find target section %d for transition: %w", transitionJSON.TargetSectionNumber, err)
		}

		transition := entities.Transition{
			TextOrder:       transitionJSON.TextOrder,
			SectionID:       section.ID,
			TargetSectionID: targetSection.ID,
			AvailableOnce:   transitionJSON.AvailableOnce,
			Text:            transitionJSON.Text,
			IsBattleWin:     transitionJSON.IsBattleWin,
			BribeResult:     transitionJSON.BribeResult,
			PlayerInput:     transitionJSON.PlayerInput,
			Dice:            transitionJSON.Dice,
			Dices:           transitionJSON.Dices,
			Condition:       transitionJSON.Condition,
			PlayerChange:    transitionJSON.PlayerChange,
			PlayerDebuff:    transitionJSON.PlayerDebuff,
			PlayerBuff:      transitionJSON.PlayerBuff,
		}

		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "section_id"}, {Name: "target_section_id"}, {Name: "text_order"}},
			DoUpdates: clause.AssignmentColumns([]string{"available_once", "text", "is_battle_win", "bribe_result", "player_input", "dice", "dices", "condition", "player_change", "player_debuff", "player_buff"}),
		}).Create(&transition).Error; err != nil {
			return fmt.Errorf("failed to upsert transition from section %d to %d: %w", transitionJSON.SectionNumber, transitionJSON.TargetSectionNumber, err)
		}
	}

	return nil
}
