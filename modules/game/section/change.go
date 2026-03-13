package section

import (
	"context"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/expression"
	"gamebook-backend/modules/game/helper"
	player2 "gamebook-backend/modules/game/player"
	"gamebook-backend/modules/game/repository"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChangeDTO struct {
	Player  *entities.Player
	Message []string
}

// Change modifies player attributes based on transition rules
func Change(
	ctx context.Context,
	item entities.Transition,
	player entities.Player,
	tx *gorm.DB,
) (ChangeDTO, error) {
	var messages []string

	if item.PlayerChange != nil {
		if item.PlayerChange.Health != nil {
			playerService := player2.NewPlayerService(tx, &player)
			playerChanged, resultMessages, err := playerService.ChangeHealthFromString(ctx, *item.PlayerChange.Health)
			if err != nil {
				return ChangeDTO{}, err
			}

			player = *playerChanged
			messages = append(messages, resultMessages...)
		}

		if item.PlayerChange.Weapons != nil {
			for _, itemWeaponSection := range *item.PlayerChange.Weapons {
				for playerWeaponIndex, playerWeaponItem := range player.Weapons {
					if *itemWeaponSection.Item == playerWeaponItem.Item {
						result, err := expression.RunAndReturnRoundUint(
							fmt.Sprintf("%d %s", player.Weapons[playerWeaponIndex].Count, *itemWeaponSection.Change),
						)
						if err != nil {
							return ChangeDTO{}, err
						}

						player.Weapons[playerWeaponIndex].Count = result
					}
				}
			}
		}

		if item.PlayerChange.Meds != nil {

		}

		if item.PlayerChange.Bag != nil {
			for _, item := range *item.PlayerChange.Bag {
				player.Bag = append(player.Bag, item)
			}
		}

		if item.PlayerChange.ReturnToSection != nil {
			player.ReturnToSection = *item.PlayerChange.ReturnToSection
		}

		if item.PlayerChange.Gold != nil {
			gold := *item.PlayerChange.Gold
			result, err := expression.RunAndReturnRoundUint(
				fmt.Sprintf("%d %s", player.Gold, gold),
			)
			if err != nil {
				return ChangeDTO{}, err
			}
			player.Gold = result
		}

		if item.PlayerChange.Bonus != nil {
			for _, bonus := range *item.PlayerChange.Bonus {
				player.Bonus = append(player.Bonus, bonus)
			}
		}
	}

	return ChangeDTO{
		Player:  &player,
		Message: messages,
	}, nil
}

func CheckConditions(
	ctx context.Context,
	db *gorm.DB,
	conditions *string,
	player *entities.Player,
) bool {
	if conditions == nil {
		return true
	}

	playerBag := player.Bag
	playerSection := player.PlayerSection

	clearedBagItems := strings.ReplaceAll(*conditions, "(", "")
	clearedBagItems = strings.ReplaceAll(clearedBagItems, ")", "")

	itemsDirty := strings.Split(clearedBagItems, "||")
	var items []string
	for _, item := range itemsDirty {
		items = append(items, strings.TrimSpace(item))
	}

	checkConditionSlice := make(map[string]bool)
	for _, item := range items {
		andItemsDirty := strings.Split(item, "&&")
		if len(andItemsDirty) > 1 {
			var andItems []string
			for _, item := range andItemsDirty {
				andItems = append(andItems, strings.TrimSpace(item))

				for _, andItem := range andItems {
					checkConditionSlice[andItem] = checkCondition(ctx, db, player.ID, andItem, playerBag, playerSection)
				}
			}
		} else {
			checkConditionSlice[item] = checkCondition(ctx, db, player.ID, item, playerBag, playerSection)
		}
	}

	resultString := *conditions
	for mapKey, mapVal := range checkConditionSlice {
		resultString = strings.ReplaceAll(resultString, mapKey, fmt.Sprintf("%t", mapVal))
	}

	result, err := expression.RunAndReturnBoolean(resultString)
	if err != nil {
		return false
	}

	return result
}

func checkCondition(
	ctx context.Context,
	db *gorm.DB,
	playerID uuid.UUID,
	bagItem string,
	playerBag []entities.Bag,
	playerSection []entities.PlayerSection,
) bool {
	result := true

	findItemInPlayerBag := func(item string, array []entities.Bag) bool {
		for _, el := range array {
			if item == el.Name {
				return true
			}
		}
		return false
	}

	findEnemyInSection := func(enemyAlias string, bagItem string) bool {
		playerSectionEnemyRepository := repository.NewPlayerSectionEnemyRepository(db)
		playerSectionEnemies, _ := playerSectionEnemyRepository.GetEnemiesByPlayerIDAndAlias(ctx, db, playerID, enemyAlias)
		if playerSectionEnemies == nil {
			return true // Мы даже не были в этой секции - враг жив
		}
		for _, el := range *playerSectionEnemies {
			if enemyAlias == el.Enemy.Alias {
				if el.Health > 0 {
					return true
				}
			}
		}
		return false
	}

	findSectionInPlayerHistory := func(section string, array []string) bool {
		for _, el := range array {
			if section == el {
				return true
			}
		}
		return false
	}

	parseString := bagItem
	firstChar := string(parseString[0])

	if firstChar == "!" {
		if strings.HasPrefix(bagItem[1:], "Bag.") {
			if findItemInPlayerBag(bagItem[5:], playerBag) {
				result = false
			}
		} else if strings.HasPrefix(bagItem[1:], "Enemy.") {
			if findEnemyInSection(bagItem[7:], bagItem) {
				result = false
			}
		} else if strings.HasPrefix(bagItem[1:], "History.") {
			history := helper.GetVisitedSections(playerSection)
			var historyStrings []string
			for _, item := range history {
				historyStrings = append(historyStrings, fmt.Sprintf("%d", item))
			}
			if findSectionInPlayerHistory(bagItem[9:], historyStrings) {
				result = false
			}
		}
	} else {
		if strings.HasPrefix(bagItem, "Bag.") {
			if !findItemInPlayerBag(bagItem[4:], playerBag) {
				result = false
			}
		} else if strings.HasPrefix(bagItem, "Enemy.") {
			if !findEnemyInSection(bagItem[6:], bagItem) {
				result = false
			}
		} else if strings.HasPrefix(bagItem, "History.") {
			history := helper.GetVisitedSections(playerSection)
			var historyStrings []string
			for _, item := range history {
				historyStrings = append(historyStrings, fmt.Sprintf("%d", item))
			}
			if !findSectionInPlayerHistory(bagItem[8:], historyStrings) {
				result = false
			}
		}

	}

	return result
}
