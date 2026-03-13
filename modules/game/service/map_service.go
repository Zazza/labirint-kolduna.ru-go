package service

import (
	"context"
	"gamebook-backend/database/entities"
	mapDTO "gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/repository"
	"gamebook-backend/modules/game/section"

	"gorm.io/gorm"
)

type MapService interface {
	GetMap(ctx context.Context, db *gorm.DB, player entities.Player) (mapDTO.MapResponse, error)
}

type mapService struct {
	sectionRepository       repository.SectionRepository
	playerSectionRepository repository.PlayerSectionRepository
}

func NewMapService(sectionRepository repository.SectionRepository, playerSectionRepository repository.PlayerSectionRepository) MapService {
	return &mapService{
		sectionRepository:       sectionRepository,
		playerSectionRepository: playerSectionRepository,
	}
}

func (s *mapService) GetMap(ctx context.Context, db *gorm.DB, player entities.Player) (mapDTO.MapResponse, error) {
	visitedSectionMap := make(map[uint]bool)
	for _, ps := range player.PlayerSection {
		visitedSectionMap[ps.Section.Number] = true
	}

	availableSections := s.getAvailableSections(player, visitedSectionMap)

	var allSectionNumbers []uint
	for sectionNumber := range visitedSectionMap {
		allSectionNumbers = append(allSectionNumbers, sectionNumber)
	}
	for _, sectionNumber := range availableSections {
		allSectionNumbers = append(allSectionNumbers, sectionNumber)
	}

	sectionsMap, err := s.getAllSections(ctx, db, allSectionNumbers)
	if err != nil {
		return mapDTO.MapResponse{}, err
	}

	var sections []mapDTO.SectionNode
	for number, sectionItem := range sectionsMap {
		var isCurrent bool
		if player.SectionID == sectionItem.ID {
			isCurrent = true
		} else {
			isCurrent = false
		}

		sections = append(sections, mapDTO.SectionNode{
			ID:        sectionItem.ID,
			Number:    sectionItem.Number,
			Title:     sectionItem.Text,
			IsVisited: visitedSectionMap[number],
			IsCurrent: isCurrent,
		})
	}

	var transitions []mapDTO.TransitionEdge
	for number, sectionItem := range sectionsMap {
		for _, transition := range sectionItem.Transitions {
			targetNumber := transition.TargetSection.Number
			//if _, exists := sectionsMap[targetNumber]; exists {
			isAvailable := s.checkTransitionAvailability(ctx, db, player, transition)
			transitions = append(transitions, mapDTO.TransitionEdge{
				FromSection: number,
				ToSection:   targetNumber,
				Text:        transition.Text,
				IsAvailable: isAvailable,
			})
			//}
		}
	}

	return mapDTO.MapResponse{
		Sections:    sections,
		Transitions: transitions,
	}, nil
}

func (s *mapService) getAvailableSections(player entities.Player, visitedSectionMap map[uint]bool) []uint {
	var availableSections []uint
	for _, transition := range player.Section.Transitions {
		targetNumber := transition.TargetSection.Number
		if !visitedSectionMap[targetNumber] {
			availableSections = append(availableSections, targetNumber)
		}
	}
	return availableSections
}

func (s *mapService) getAllSections(ctx context.Context, db *gorm.DB, sectionNumbers []uint) (map[uint]entities.Section, error) {
	sections, err := s.sectionRepository.GetAllWithTransitions(ctx, db, sectionNumbers)
	if err != nil {
		return nil, err
	}

	result := make(map[uint]entities.Section)
	for _, section := range sections {
		result[section.Number] = section
	}
	return result, nil
}

func (s *mapService) checkTransitionAvailability(ctx context.Context, db *gorm.DB, player entities.Player, transition entities.Transition) bool {
	return section.CheckConditions(ctx, db, transition.Condition, &player)
}
