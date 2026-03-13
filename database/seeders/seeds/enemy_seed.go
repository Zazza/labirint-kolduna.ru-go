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

type EnemyJSON struct {
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

func EnemySeeder(db *gorm.DB) error {
	data, err := os.ReadFile(filepath.Clean("./database/seeders/json/enemies.json"))
	if err != nil {
		return fmt.Errorf("failed to read enemies.json: %w", err)
	}

	var enemiesJSON []EnemyJSON
	if err := json.Unmarshal(data, &enemiesJSON); err != nil {
		return fmt.Errorf("failed to unmarshal enemies.json: %w", err)
	}

	if err := db.AutoMigrate(&entities.Enemy{}); err != nil {
		return fmt.Errorf("failed to auto migrate enemies: %w", err)
	}

	for _, enemyJSON := range enemiesJSON {
		enemy := entities.Enemy{
			Alias:        enemyJSON.Alias,
			Name:         enemyJSON.Name,
			Damage:       enemyJSON.Damage,
			DamageType:   enemyJSON.DamageType,
			MinDiceHits:  enemyJSON.MinDiceHits,
			Health:       enemyJSON.Health,
			Defence:      enemyJSON.Defence,
			PlayerArmor:  enemyJSON.PlayerArmor,
			OnlyDiceHits: enemyJSON.OnlyDiceHits,
			MagicHit:     enemyJSON.MagicHit,
			Weapons:      enemyJSON.Weapons,
		}

		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "alias"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "damage", "damage_type", "min_dice_hits", "health", "defence", "player_armor", "only_dice_hits", "magic_hit", "weapons"}),
		}).Create(&enemy).Error; err != nil {
			return fmt.Errorf("failed to upsert enemy %s: %w", enemyJSON.Alias, err)
		}
	}

	return nil
}
