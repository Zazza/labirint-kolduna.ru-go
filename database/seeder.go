package database

import (
	"gamebook-backend/database/seeders/seeds"

	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	if err := seeds.EnemySeeder(db); err != nil {
		return err
	}

	if err := seeds.SectionSeeder(db); err != nil {
		return err
	}

	if err := seeds.TransitionSeeder(db); err != nil {
		return err
	}

	return nil
}
