package sleep

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/repository"

	"gorm.io/gorm"
)

type Exit interface {
	IsExit(ctx context.Context) (bool, error)
}

type exit struct {
	db     *gorm.DB
	player entities.Player
}

func NewExit(db *gorm.DB, player entities.Player) Exit {
	return &exit{
		db:     db,
		player: player,
	}
}

func (e *exit) IsExit(ctx context.Context) (bool, error) {
	playerSectionRepository := repository.NewPlayerSectionRepository(e.db)

	playerSection, err := playerSectionRepository.GetLastPlayerSection(ctx, e.db, e.player.ID)
	if err != nil {
		return false, err
	}

	emptyPlayerTargetSectionID := entities.PlayerSection{}.TargetSectionID
	//if playerSection.SectionID == e.player.Section.ID && playerSection.TargetSectionID != emptyPlayerTargetSectionID {
	if playerSection.TargetSectionID != emptyPlayerTargetSectionID {
		return true, nil
	}

	return false, nil
}
