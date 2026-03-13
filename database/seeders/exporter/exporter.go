package exporter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gamebook-backend/database/entities"

	"gorm.io/gorm"
)

type Exporter struct {
	db     *gorm.DB
	output string
}

func NewExporter(db *gorm.DB, output string) *Exporter {
	return &Exporter{
		db:     db,
		output: output,
	}
}

func (e *Exporter) ExportAll() error {
	if err := e.EnsureOutputDir(); err != nil {
		return err
	}

	if err := e.ExportEnemies(); err != nil {
		return fmt.Errorf("failed to export enemies: %w", err)
	}

	if err := e.ExportSections(); err != nil {
		return fmt.Errorf("failed to export sections: %w", err)
	}

	if err := e.ExportTransitions(); err != nil {
		return fmt.Errorf("failed to export transitions: %w", err)
	}

	return nil
}

func (e *Exporter) EnsureOutputDir() error {
	if e.output == "" {
		e.output = "./database/seeders/json"
	}

	if err := os.MkdirAll(e.output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return nil
}

type EnemyExport struct {
	ID           *string                    `json:"id"`
	Alias        string                     `json:"alias"`
	Name         string                     `json:"name"`
	Damage       uint                       `json:"damage"`
	DamageType   entities.BuffOrDebuffAlias `json:"damage_type"`
	MinDiceHits  uint                       `json:"min_dice_hits"`
	Health       uint                       `json:"health"`
	Defence      uint                       `json:"defence"`
	PlayerArmor  bool                       `json:"player_armor"`
	OnlyDiceHits *uint                      `json:"only_dice_hits"`
	MagicHit     *entities.MagicHit         `json:"magic_hit"`
	Weapons      []entities.EnemyWeapon     `json:"weapons"`
}

func (e *Exporter) ExportEnemies() error {
	var enemies []entities.Enemy
	if err := e.db.Find(&enemies).Error; err != nil {
		return fmt.Errorf("failed to fetch enemies: %w", err)
	}

	if len(enemies) == 0 {
		log.Println("WARNING: enemies table is empty")
		exports := []EnemyExport{}
		return e.writeJSON("enemies.json", exports)
	}

	exports := make([]EnemyExport, 0, len(enemies))
	for _, enemy := range enemies {
		export := EnemyExport{
			ID:           nil,
			Alias:        enemy.Alias,
			Name:         enemy.Name,
			Damage:       enemy.Damage,
			DamageType:   enemy.DamageType,
			MinDiceHits:  enemy.MinDiceHits,
			Health:       enemy.Health,
			Defence:      enemy.Defence,
			PlayerArmor:  enemy.PlayerArmor,
			OnlyDiceHits: enemy.OnlyDiceHits,
			MagicHit:     enemy.MagicHit,
			Weapons:      enemy.Weapons,
		}
		exports = append(exports, export)
	}

	return e.writeJSON("enemies.json", exports)
}

type SectionExport struct {
	ID           *string              `json:"id"`
	Type         entities.SectionType `json:"type"`
	Number       uint                 `json:"number"`
	Text         string               `json:"text"`
	EnemyAliases []string             `json:"enemy_aliases"`
	Choice       *entities.Choice     `json:"choice"`
	BattleStart  *string              `json:"battle_start"`
	BattleSteps  []*string            `json:"battle_steps"`
	Dices        bool                 `json:"dices"`
	Bribe        *entities.Bribe      `json:"bribe"`
}

func (e *Exporter) ExportSections() error {
	var sections []entities.Section
	if err := e.db.Preload("SectionEnemies").Find(&sections).Error; err != nil {
		return fmt.Errorf("failed to fetch sections: %w", err)
	}

	if len(sections) == 0 {
		return fmt.Errorf("sections table is empty")
	}

	exports := make([]SectionExport, 0, len(sections))
	for _, section := range sections {
		enemyAliases := make([]string, 0, len(section.SectionEnemies))
		for _, enemy := range section.SectionEnemies {
			enemyAliases = append(enemyAliases, enemy.Alias)
		}

		export := SectionExport{
			ID:           nil,
			Type:         section.Type,
			Number:       section.Number,
			Text:         section.Text,
			EnemyAliases: enemyAliases,
			Choice:       section.Choice,
			BattleStart:  section.BattleStart,
			BattleSteps:  section.BattleSteps,
			Dices:        section.Dices,
			Bribe:        section.Bribe,
		}
		exports = append(exports, export)
	}

	return e.writeJSON("sections.json", exports)
}

type TransitionExport struct {
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

func (e *Exporter) ExportTransitions() error {
	var transitions []entities.Transition
	if err := e.db.Preload("Section").Preload("TargetSection").Find(&transitions).Error; err != nil {
		return fmt.Errorf("failed to fetch transitions: %w", err)
	}

	if len(transitions) == 0 {
		return fmt.Errorf("transitions table is empty")
	}

	exports := make([]TransitionExport, 0, len(transitions))
	for _, transition := range transitions {
		export := TransitionExport{
			ID:                  nil,
			TextOrder:           transition.TextOrder,
			SectionNumber:       transition.Section.Number,
			TargetSectionNumber: transition.TargetSection.Number,
			AvailableOnce:       transition.AvailableOnce,
			Text:                transition.Text,
			IsBattleWin:         transition.IsBattleWin,
			BribeResult:         transition.BribeResult,
			PlayerInput:         transition.PlayerInput,
			Dice:                transition.Dice,
			Dices:               transition.Dices,
			Condition:           transition.Condition,
			PlayerChange:        transition.PlayerChange,
			PlayerDebuff:        transition.PlayerDebuff,
			PlayerBuff:          transition.PlayerBuff,
		}
		exports = append(exports, export)
	}

	return e.writeJSON("transitions.json", exports)
}

func (e *Exporter) writeJSON(filename string, data interface{}) error {
	outputPath := filepath.Join(e.output, filename)
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	log.Printf("Exported %d records to %s", countRecords(data), outputPath)
	return nil
}

func countRecords(data interface{}) int {
	switch v := data.(type) {
	case []interface{}:
		return len(v)
	case []EnemyExport:
		return len(v)
	case []SectionExport:
		return len(v)
	case []TransitionExport:
		return len(v)
	default:
		return 0
	}
}
