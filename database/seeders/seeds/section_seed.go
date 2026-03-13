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

type SectionJSON struct {
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

func SectionSeeder(db *gorm.DB) error {
	data, err := os.ReadFile(filepath.Clean("./database/seeders/json/sections.json"))
	if err != nil {
		return fmt.Errorf("failed to read sections.json: %w", err)
	}

	var sectionsJSON []SectionJSON
	if err := json.Unmarshal(data, &sectionsJSON); err != nil {
		return fmt.Errorf("failed to unmarshal sections.json: %w", err)
	}

	if err := db.AutoMigrate(&entities.Section{}); err != nil {
		return fmt.Errorf("failed to auto migrate sections: %w", err)
	}

	for _, sectionJSON := range sectionsJSON {
		var enemies []*entities.Enemy
		if len(sectionJSON.EnemyAliases) > 0 {
			if err := db.Where("alias IN ?", sectionJSON.EnemyAliases).Find(&enemies).Error; err != nil {
				return fmt.Errorf("failed to find enemies for section %d: %w", sectionJSON.Number, err)
			}
		}

		section := entities.Section{
			Type:           sectionJSON.Type,
			Number:         sectionJSON.Number,
			Text:           sectionJSON.Text,
			SectionEnemies: enemies,
			Choice:         sectionJSON.Choice,
			BattleStart:    sectionJSON.BattleStart,
			BattleSteps:    sectionJSON.BattleSteps,
			Dices:          sectionJSON.Dices,
			Bribe:          sectionJSON.Bribe,
		}

		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "number"}},
			DoUpdates: clause.AssignmentColumns([]string{"type", "text", "choice", "battle_start", "battle_steps", "dices", "bribe"}),
		}).Create(&section).Error; err != nil {
			return fmt.Errorf("failed to upsert section %d: %w", sectionJSON.Number, err)
		}

		if len(enemies) > 0 {
			if err := db.Model(&section).Association("SectionEnemies").Replace(enemies); err != nil {
				return fmt.Errorf("failed to replace enemies for section %d: %w", sectionJSON.Number, err)
			}
		}
	}

	return nil
}
