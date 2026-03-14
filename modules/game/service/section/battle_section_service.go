package section

import (
	"context"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/battle"
	"gamebook-backend/modules/game/bribe"
	"gamebook-backend/modules/game/dto"
	enemyDTO "gamebook-backend/modules/game/dto"
	gameDTO "gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/repository"
	"gamebook-backend/modules/game/sleep"
	template2 "gamebook-backend/modules/game/template"
	"gamebook-backend/modules/game/template/section"

	"gorm.io/gorm"
)

type BattleSectionService interface {
	GetActivityByPlayer(ctx context.Context, player entities.Player) (gameDTO.CurrentResponse, error)
}

type battleSectionService struct {
	sectionRepository            repository.SectionRepository
	dicesRepository              repository.DiceRepository
	battleRepository             repository.BattleRepository
	playerRepository             repository.PlayerRepository
	playerSectionEnemyRepository repository.PlayerSectionEnemyRepository
	db                           *gorm.DB
}

func NewBattleSectionService(
	sectionRepo repository.SectionRepository,
	dicesRepo repository.DiceRepository,
	battleRepo repository.BattleRepository,
	playerRepo repository.PlayerRepository,
	playerSectionEnemyRepo repository.PlayerSectionEnemyRepository,
	db *gorm.DB,
) BattleSectionService {
	return &battleSectionService{
		sectionRepository:            sectionRepo,
		dicesRepository:              dicesRepo,
		battleRepository:             battleRepo,
		playerRepository:             playerRepo,
		playerSectionEnemyRepository: playerSectionEnemyRepo,
		db:                           db,
	}
}

func (s *battleSectionService) GetActivityByPlayer(ctx context.Context, player entities.Player) (gameDTO.CurrentResponse, error) {
	var result dto.ActivityResponse
	var err error

	exit := sleep.NewExit(s.db, player)
	isExit, err := exit.IsExit(ctx)
	if err != nil {
		return gameDTO.CurrentResponse{}, err
	}

	if isExit {
		if player.Section.Type == gameDTO.SectionTypeSleepy {
			sectionResponse := section.NewSectionResponse(s.db, player)
			response, err := sectionResponse.GetResponse(ctx)
			if err != nil {
				return gameDTO.CurrentResponse{}, err
			}

			return response, nil
		}
	}

	common, err := battle.NewCommon(ctx, s.db, &player)
	if err != nil {
		return gameDTO.CurrentResponse{}, err
	}

	result, err = s.GetBattleActivityByPlayer(ctx, common)
	if err != nil {
		return gameDTO.CurrentResponse{}, err
	}

	if common.GetLastBattleLog() == nil {
		isBribe, err := bribe.IsBribe(s.db, player.ID)
		if err != nil {
			return gameDTO.CurrentResponse{}, err
		}
		if bribe.IsPossible(player.Section) && !isBribe {
			bribeTransition := bribe.GetBribeTransition(player.Section.ID)
			result.Transitions = append(result.Transitions, bribeTransition)
		}
	}

	var bonuses []gameDTO.PlayerInfoBonus
	for _, playerBonus := range player.Bonus {
		if playerBonus.Option != nil {
			for _, option := range *playerBonus.Option {
				bonuses = append(bonuses, gameDTO.PlayerInfoBonus{
					Alias:  *playerBonus.Alias,
					Name:   *playerBonus.Name,
					Option: &option,
				})
			}
		} else {
			bonuses = append(bonuses, gameDTO.PlayerInfoBonus{
				Alias: *playerBonus.Alias,
				Name:  *playerBonus.Name,
			})
		}
	}

	return gameDTO.CurrentResponse{
		Section:      player.Section.Number,
		Text:         result.Text,
		Type:         result.Type,
		Transitions:  result.Transitions,
		RollTheDices: result.RollTheDices,
		Player: gameDTO.PlayerInfo{
			Health: player.Health,
			Meds:   uint(player.Meds.Count),
			Gold:   player.Gold,
			Bonus:  bonuses,
		},
		MapAvailable: helper.HasBagItem(player.Bag, "mapIngredients"),
	}, nil
}

func (s *battleSectionService) GetBattleActivityByPlayer(ctx context.Context, common battle.Common) (dto.ActivityResponse, error) {
	player := *common.GetPlayer()

	if len(player.Section.SectionEnemies) == 0 {
		return dto.ActivityResponse{}, dto.ErrGetActivityBySectionId
	}
	if player.Section.BattleStart == nil {
		return dto.ActivityResponse{}, dto.ErrBattleStartNotDefined
	}

	var text string
	var gotoItems []dto.TransitionDTO

	common, err := battle.NewCommon(ctx, s.db, &player)
	if err != nil {
		return dto.ActivityResponse{}, err
	}

	if common.GetLastBattleLog() == nil {
		isBribe, err := bribe.IsBribe(s.db, player.ID)
		if err != nil {
			return dto.ActivityResponse{}, err
		}

		if isBribe {
			bribeResultTransition, err := bribe.GetBribeResultTransition(s.db, player)
			if err != nil {
				return dto.ActivityResponse{}, err
			}

			gotoItems = []gameDTO.TransitionDTO{bribeResultTransition}

			return dto.ActivityResponse{
				Text:        player.Section.Text,
				Transitions: gotoItems,
				Type:        dto.CHOICE,
			}, nil
		}
	}

	if common.GetLastBattleLog() == nil {
		battleDicesDTO := s.dicesRepository.FindBattleDicesByPlayerId(ctx, s.db, player.ID)
		if battleDicesDTO.Error != nil {
			return dto.ActivityResponse{}, battleDicesDTO.Error
		}

		if !battleDicesDTO.Exists && *player.Section.BattleStart == enemyDTO.BattleStartDices {
			rollTheDices := s.sectionRepository.IsDicesRequired(player.Section)
			if !rollTheDices {
				gotoItems = []gameDTO.TransitionDTO{
					{Text: "В бой!", TransitionID: player.Section.ID},
				}
			}

			return dto.ActivityResponse{
				Text:         player.Section.Text,
				Transitions:  gotoItems,
				Type:         dto.CHOICE,
				RollTheDices: rollTheDices,
			}, nil
		}

		if *player.Section.BattleStart == enemyDTO.BattleStartDices {
			if battleDicesDTO.Dices.DiceFirst == battleDicesDTO.Dices.DiceSecond {
				dicesEqualsTemplate, err := template2.GetDicesTemplate(
					ctx,
					battleDicesDTO.Dices.DiceFirst,
					battleDicesDTO.Dices.DiceSecond,
					true,
				)
				if err != nil {
					return dto.ActivityResponse{}, err
				}

				text += "<p>Брось кубики еще раз</p>" + dicesEqualsTemplate

				return dto.ActivityResponse{
					Text:         text,
					Type:         dto.CHOICE,
					Transitions:  gotoItems,
					RollTheDices: true,
				}, nil
			}
		}
	}

	result, err := s.showStartBattleInfo(ctx, player, common)
	if err != nil {
		return dto.ActivityResponse{}, err
	}

	return result, nil
}

func (s *battleSectionService) showStartBattleInfo(
	ctx context.Context,
	player entities.Player,
	common battle.Common,
) (dto.ActivityResponse, error) {
	text := player.Section.Text

	// Сообщения о прошлых ударах
	if len(*common.GetBattleLog()) > 0 {
		text += "<div style='margin: 20px 0; text-align: center;'>⚔️</div>"

		for _, item := range *common.GetBattleLog() {
			template, err := template2.GetDicesTemplate(
				ctx,
				item.Dice1,
				item.Dice2,
				false,
			)
			if err != nil {
				return dto.ActivityResponse{}, err
			}

			text += fmt.Sprintf(
				"<div class='section-log'>%s%s</div>",
				template,
				item.Description,
			)
		}
	}

	allEnemiesHealth := common.Enemy().GetAllEnemiesHealth()

	if allEnemiesHealth == 0 {
		gotoSection, err := common.GetGotoSectionWin()
		if err != nil {
			return dto.ActivityResponse{}, err
		}

		var dtoGotoSections []dto.TransitionDTO
		for _, item := range gotoSection {
			if player.Section.Type == gameDTO.StepTypeNormal {
				sectionEntity, err := s.sectionRepository.GetByID(ctx, s.db, item.TargetSectionID)
				if err != nil {
					return dto.ActivityResponse{}, err
				}

				dtoGotoSections = append(dtoGotoSections, dto.TransitionDTO{
					Text:         item.Text,
					TransitionID: item.ID,
					Visited:      helper.IsVisited(sectionEntity.Number, player.PlayerSection),
				})
			} else if player.Section.Type == gameDTO.SectionTypeSleepy {
				dtoGotoSections = append(dtoGotoSections, dto.TransitionDTO{
					Text:             item.Text,
					SleepyTransition: true,
				})
			} else {
				return dto.ActivityResponse{}, gameDTO.ErrSectionNotFound
			}
		}

		return dto.ActivityResponse{
			Text:        text + "<div align='center'>🏆 Ты победил!</div>",
			Type:        dto.CHOICE,
			Transitions: dtoGotoSections,
		}, nil
	}

	if player.Health == 0 {
		gotoSection, err := common.GetGotoSectionLose()
		if err != nil {
			return dto.ActivityResponse{}, err
		}

		sleepyTransition := false
		if player.Section.Type == gameDTO.SectionTypeSleepy {
			sleepyTransition = true
		}

		return dto.ActivityResponse{
			Text: text + " <div align='center'>💀 Тебя убили</div>",
			Type: dto.CHOICE,
			Transitions: []dto.TransitionDTO{
				{
					Text:             "Продолжить",
					TransitionID:     gotoSection,
					SleepyTransition: sleepyTransition,
				},
			},
		}, nil
	}

	isMyMove, err := common.Step().IsMyMove()
	if err != nil {
		return dto.ActivityResponse{}, err
	}

	var gotoItems []dto.TransitionDTO

	if isMyMove {
		for _, weapon := range player.Weapons {
			var gotoText string
			if weapon.Item == "hand" || weapon.Item == "sword" {
				gotoText = weapon.Name
			} else {
				if weapon.Count > 0 {
					gotoText = fmt.Sprintf("%s [x%d]", weapon.Name, weapon.Count)
				}
			}
			enemy, err := common.Enemy().GetEnemy()
			if err != nil {
				return dto.ActivityResponse{}, err
			}
			if enemy.Abstract.OnlyDiceHits != nil {
				if weapon.Item == "sword" {
					gotoItems = append(gotoItems, dto.TransitionDTO{
						Text:         gotoText,
						TransitionID: player.Section.ID,
						Weapon:       weapon.Item,
					})
				}
			} else if player.Section.Type == gameDTO.SectionTypeSleepy {
				if weapon.Item == "hand" {
					gotoItems = append(gotoItems, dto.TransitionDTO{
						Text:         gotoText,
						TransitionID: player.Section.ID,
						Weapon:       weapon.Item,
					})
				}
			} else {
				if weapon.Item == "hand" {
					continue
				}

				if weapon.Count == 0 && weapon.Item != "sword" {
					continue
				}

				gotoItems = append(gotoItems, dto.TransitionDTO{
					Text:         gotoText,
					TransitionID: player.Section.ID,
					Weapon:       weapon.Item,
				})
			}
		}
	} else {
		gotoItems = []dto.TransitionDTO{
			{Text: "Бьёт враг", TransitionID: player.Section.ID},
		}
	}

	return dto.ActivityResponse{
		Text:        text,
		Type:        dto.BATTLE,
		Transitions: gotoItems,
	}, nil
}
