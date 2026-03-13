package custom

import (
	"context"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/repository"

	"gorm.io/gorm"
)

type section156 struct {
	db     *gorm.DB
	player entities.Player
}

func NewSection156(db *gorm.DB, player entities.Player) dto.CustomSection {
	return &section156{
		db:     db,
		player: player,
	}
}

func (s *section156) Handle(ctx context.Context) (dto.CustomSectionDTO, error) {
	sectionNumbers := helper.GetVisitedSections(s.player.PlayerSection)

	sectionRepository := repository.NewSectionRepository(s.db)
	sections, err := sectionRepository.GetListBySectionNumbers(ctx, s.db, sectionNumbers)
	if err != nil {
		return dto.CustomSectionDTO{}, err
	}

	var gotoItems []dto.TransitionDTO
	for _, section := range sections {
		gotoItems = append(gotoItems, dto.TransitionDTO{
			Text:         fmt.Sprintf("Секция %d", section.Number),
			TransitionID: section.ID,
		})
	}

	return dto.CustomSectionDTO{
		SectionText: s.player.Section.Text,
		GotoItems:   gotoItems,
	}, nil
}
