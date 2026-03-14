package sleep

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/listener"
	"gamebook-backend/modules/game/listener/event"
	"gamebook-backend/modules/game/repository"

	"gorm.io/gorm"
)

type Exit interface {
	IsExit(ctx context.Context) (bool, error)
	Return(ctx context.Context) error
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

func (e *exit) Return(ctx context.Context) error {
	sectionRepository := repository.NewSectionRepository(e.db)
	returnToSection, err := sectionRepository.GetBySectionNumber(ctx, e.db, e.player.ReturnToSection)
	if err != nil {
		return err
	}

	playerUpdateListener, err := listener.HandleEvent(e.db, "player_update")
	if err != nil {
		return err
	}

	err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
		PlayerID:  e.player.ID,
		SectionID: &returnToSection.ID,
	})
	if err != nil {
		return err
	}

	return nil
}
