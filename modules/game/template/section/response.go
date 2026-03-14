package section

import (
	"context"
	"errors"
	"gamebook-backend/database/entities"
	gameDto "gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/repository"
	"gamebook-backend/modules/game/sleep"

	"gorm.io/gorm"
)

type SleepySectionResponse interface {
	GetResponse(ctx context.Context) (gameDto.CurrentResponse, error)
	GetSectionText(ctx context.Context) (string, error)
}

type sleepySectionResponse struct {
	db     *gorm.DB
	player entities.Player
}

func NewSectionResponse(db *gorm.DB, player entities.Player) SleepySectionResponse {
	return &sleepySectionResponse{
		db:     db,
		player: player,
	}
}

func (r *sleepySectionResponse) GetResponse(ctx context.Context) (gameDto.CurrentResponse, error) {
	var gotoItems []gameDto.TransitionDTO

	exit := sleep.NewExit(r.db, r.player)
	isExit, err := exit.IsExit(ctx)
	if err != nil {
		return gameDto.CurrentResponse{}, err
	}
	if isExit {
		playerSectionRepository := repository.NewPlayerSectionRepository(r.db)
		playerSection, err := playerSectionRepository.GetLastPlayerSection(ctx, r.db, r.player.ID)
		if err != nil {
			return gameDto.CurrentResponse{}, err
		}

		gotoItems = []gameDto.TransitionDTO{
			{
				Text:             "Выйти из сонного царства",
				SleepyTransition: true,
				TransitionID:     playerSection.TargetSectionID,
			},
		}
	} else {
		gotoItems = []gameDto.TransitionDTO{
			{
				Text:             "Далее",
				SleepyTransition: true,
			},
		}
	}

	sectionText, err := r.GetSectionText(ctx)
	if err != nil {
		return gameDto.CurrentResponse{}, err
	}

	return gameDto.CurrentResponse{
		Section:      r.player.Section.Number,
		Text:         sectionText,
		Type:         gameDto.SectionTypeSleepy,
		Transitions:  gotoItems,
		RollTheDices: false,
		Player: gameDto.PlayerInfo{
			Health: r.player.Health,
			Meds:   uint(r.player.Meds.Count),
			Gold:   r.player.Gold,
		},
		MapAvailable: helper.HasBagItem(r.player.Bag, "mapIngredients"),
	}, nil
}

func (r *sleepySectionResponse) GetSectionText(ctx context.Context) (string, error) {
	playerSectionRepository := repository.NewPlayerSectionRepository(r.db)

	sectionText := r.player.Section.Text

	descriptions, err := playerSectionRepository.GetLastSectionDescriptions(ctx, r.db, r.player.ID)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return sectionText, nil
	} else if err != nil {
		return "", err
	}

	for _, description := range *descriptions {
		sectionText += "<div class='section-log'>" + description.Description + "</div>"
	}

	return sectionText, nil
}
