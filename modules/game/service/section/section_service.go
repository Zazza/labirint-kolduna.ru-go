package section

import (
	"context"
	"errors"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	diceDTO "gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/repository"
	"gamebook-backend/modules/game/service"
	"gamebook-backend/modules/game/service/section/custom"
	template2 "gamebook-backend/modules/game/template"
	section2 "gamebook-backend/modules/game/template/section"

	"gorm.io/gorm"
)

type SectionService interface {
	GetSection(ctx context.Context, db *gorm.DB, player entities.Player) (dto.CurrentResponse, error)
	GetProfile(ctx context.Context, db *gorm.DB, player entities.Player) (dto.ProfileResponse, error)
}

type sectionService struct {
	sectionRepository    repository.SectionRepository
	transitionService    service.TransitionService
	bonusService         service.BonusService
	playerProfileService service.PlayerProfileService
	diceService          service.DiceService
	bribeService         service.BribeService
}

func NewSectionService(
	sectionRepo repository.SectionRepository,
	transitionSvc service.TransitionService,
	bonusSvc service.BonusService,
	playerProfileSvc service.PlayerProfileService,
	diceSvc service.DiceService,
	bribeSvc service.BribeService,
) SectionService {
	return &sectionService{
		sectionRepository:    sectionRepo,
		transitionService:    transitionSvc,
		bonusService:         bonusSvc,
		playerProfileService: playerProfileSvc,
		diceService:          diceSvc,
		bribeService:         bribeSvc,
	}
}

func (s *sectionService) GetSection(ctx context.Context, db *gorm.DB, player entities.Player) (dto.CurrentResponse, error) {
	sectionResponse := section2.NewSectionResponse(db, player)
	sectionText, err := sectionResponse.GetSectionText(ctx)
	if err != nil {
		return dto.CurrentResponse{}, err
	}

	if custom.IsCustom(player.Section.Number) {
		return s.handleCustomSection(ctx, db, player, sectionText)
	}

	return s.handleRegularSection(ctx, db, player, sectionText)
}

func (s *sectionService) handleCustomSection(ctx context.Context, db *gorm.DB, player entities.Player, sectionText string) (dto.CurrentResponse, error) {
	customSection, err := custom.GetSection(db, player, player.Section.Number)
	if err != nil {
		return dto.CurrentResponse{}, err
	}

	customSectionDTO, err := customSection.Handle(ctx)
	if err != nil {
		return dto.CurrentResponse{}, err
	}

	return s.buildCustomSectionResponse(player, customSectionDTO)
}

func (s *sectionService) buildCustomSectionResponse(player entities.Player, customSectionDTO dto.CustomSectionDTO) (dto.CurrentResponse, error) {
	bonuses := s.bonusService.ConvertBonusesToDTO(player.Bonus)
	return BuildCurrentResponse(player, customSectionDTO.SectionText, customSectionDTO.GotoItems, entities.Choice{}, bonuses, false, dto.SectionTypeChoice), nil
}

func (s *sectionService) handleRegularSection(ctx context.Context, db *gorm.DB, player entities.Player, sectionText string) (dto.CurrentResponse, error) {
	rollTheDicesNeeded, gotoItems, err := s.getTransitions(ctx, db, &player, sectionText)
	if err != nil {
		return dto.CurrentResponse{}, err
	}

	choice := s.getChoice(player)
	bonuses := s.bonusService.ConvertBonusesToDTO(player.Bonus)

	if s.bribeService.IsPossible(player.Section) {
		return s.handleBribeSection(ctx, db, player, sectionText, gotoItems, choice, bonuses, rollTheDicesNeeded, dto.SectionTypeChoice)
	}

	return s.buildResponse(player, sectionText, gotoItems, choice, bonuses, rollTheDicesNeeded, dto.SectionTypeChoice), nil
}

func (s *sectionService) getTransitions(ctx context.Context, db *gorm.DB, player *entities.Player, sectionText string) (bool, []dto.TransitionDTO, error) {
	rollTheDicesNeeded := false
	var gotoItems []dto.TransitionDTO
	var err error

	if s.sectionRepository.IsDicesRequired(player.Section) {
		diceEntity, err := s.diceService.GetLastDice(ctx, db, player.ID, dto.ReasonChoice)
		if errors.Is(err, diceDTO.MessageDicesNotDefined) {
			rollTheDicesNeeded = true
		} else if err != nil {
			return false, nil, err
		}

		if !rollTheDicesNeeded {
			gotoItems, sectionText, err = s.processTransitionsWithDice(ctx, db, player, diceEntity, sectionText)
			if err != nil {
				return false, nil, err
			}
		}
	} else {
		gotoItems, err = s.transitionService.GetAvailableTransitions(ctx, db, player, nil)
		if err != nil {
			return false, nil, err
		}
	}

	return rollTheDicesNeeded, gotoItems, nil
}

func (s *sectionService) getChoice(player entities.Player) entities.Choice {
	if player.Section.Choice != nil {
		return *player.Section.Choice
	}
	return entities.Choice{}
}

func (s *sectionService) buildResponse(player entities.Player, sectionText string, gotoItems []dto.TransitionDTO, choice entities.Choice, bonuses []dto.PlayerInfoBonus, rollTheDicesNeeded bool, dataType entities.SectionType) dto.CurrentResponse {
	return BuildCurrentResponse(player, sectionText, gotoItems, choice, bonuses, rollTheDicesNeeded, dataType)
}

func (s *sectionService) processTransitionsWithDice(ctx context.Context, db *gorm.DB, player *entities.Player, diceEntity *entities.Dice, sectionText string) ([]dto.TransitionDTO, string, error) {
	var gotoItems []dto.TransitionDTO
	ctxWithDices := context.WithValue(ctx, "dices", diceEntity)

	for _, item := range player.Section.Transitions {
		template, err := s.getTransitionTemplate(ctx, &item, diceEntity)
		if err != nil {
			return nil, "", err
		}

		if template != "" {
			sectionText = player.Section.Text + fmt.Sprintf(
				"<p style='text-align: center;'>%s</p>",
				template,
			)
		}

		result, err := s.transitionService.Check(ctxWithDices, db, &item, player)
		if err != nil {
			return nil, "", err
		}

		if result {
			sectionEntity, err := s.sectionRepository.GetByID(ctx, db, item.TargetSectionID)
			if err != nil {
				return nil, "", err
			}

			gotoItems = append(gotoItems, dto.TransitionDTO{
				Text:         item.Text,
				TransitionID: item.ID,
				Visited:      helper.IsVisited(sectionEntity.Number, player.PlayerSection),
			})
		}
	}

	return gotoItems, sectionText, nil
}

func (s *sectionService) getTransitionTemplate(ctx context.Context, transition *entities.Transition, dices *entities.Dice) (string, error) {
	if transition.Dice != nil {
		template, err := template2.GetDiceTemplate(
			ctx,
			dices.DiceFirst,
			true,
		)
		if err != nil {
			return "", err
		}
		return template, nil
	} else if transition.Dices != nil {
		template, err := template2.GetDicesTemplate(
			ctx,
			dices.DiceFirst,
			dices.DiceSecond,
			true,
		)
		if err != nil {
			return "", err
		}
		return template, nil
	}
	return "", nil
}

func (s *sectionService) handleBribeSection(ctx context.Context, db *gorm.DB, player entities.Player, sectionText string, gotoItems []dto.TransitionDTO, choice entities.Choice, bonuses []dto.PlayerInfoBonus, rollTheDicesNeeded bool, dataType entities.SectionType) (dto.CurrentResponse, error) {
	isBribe, err := s.bribeService.IsBribe(db, player.ID)
	if err != nil {
		return dto.CurrentResponse{}, err
	}

	if isBribe {
		bribeResultTransition, err := s.bribeService.GetBribeResultTransition(db, player)
		if err != nil {
			return dto.CurrentResponse{}, err
		}

		gotoItems = []dto.TransitionDTO{bribeResultTransition}
	} else {
		bribeTransition := s.bribeService.GetBribeTransition(player.Section.ID)
		gotoItems = append(gotoItems, bribeTransition)
	}

	return s.buildResponse(player, sectionText, gotoItems, choice, bonuses, rollTheDicesNeeded, dataType), nil
}

func (s *sectionService) GetProfile(ctx context.Context, db *gorm.DB, player entities.Player) (dto.ProfileResponse, error) {
	return s.playerProfileService.GetProfile(ctx, db, player)
}
