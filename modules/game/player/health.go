package player

import (
	"context"
	"fmt"
	"gamebook-backend/database/entities"
	dice2 "gamebook-backend/modules/game/dice"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/expression"
	template2 "gamebook-backend/modules/game/template"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type Service interface {
	ChangeHealthFromString(
		ctx context.Context,
		healthString string,
	) (*entities.Player, []string, error)
}

type playerService struct {
	db     *gorm.DB
	player *entities.Player
}

func NewPlayerService(db *gorm.DB, player *entities.Player) Service {
	return &playerService{
		db:     db,
		player: player,
	}
}

func (s *playerService) ChangeHealthFromString(
	ctx context.Context,
	healthString string,
) (*entities.Player, []string, error) {
	var messages []string

	healthChange := healthString

	simpleChange := false
	if healthChange[1:] == "dice" || healthChange[1:] == "dices" {
		firstChar := string(healthChange[0])
		if firstChar != "-" && firstChar != "+" {
			return nil, []string{}, dto.ErrPlayerChangesDicesFirstChar
		}

		if healthChange[1:] == "dice" {
			rollTheDices := dice2.NewRollTheDices(s.db, s.player)
			diceFirst, err := rollTheDices.RollTheDice(context.Background(), *s.player)
			if err != nil {
				return nil, []string{}, err
			}

			template, err := template2.GetDiceTemplate(context.Background(), *diceFirst, true)
			if err != nil {
				return nil, []string{}, err
			}

			messages = append(messages, fmt.Sprintf("%s", template))

			healthChange = firstChar + strconv.Itoa(int(*diceFirst))
		} else if healthChange[1:] == "dices" {
			rollTheDices := dice2.NewRollTheDices(s.db, s.player)
			diceFirst, diceSecond, err := rollTheDices.RollTheDices(context.Background(), *s.player)
			if err != nil {
				return nil, []string{}, err
			}

			template, err := template2.GetDicesTemplate(context.Background(), *diceFirst, *diceSecond, true)
			if err != nil {
				return nil, []string{}, err
			}

			messages = append(messages, fmt.Sprintf("%s", template))

			healthChange = firstChar + strconv.Itoa(int(*diceFirst+*diceSecond))
		}
	} else if strings.Contains(healthChange, "max") {
		healthChange = strings.ReplaceAll(
			healthChange,
			"max",
			fmt.Sprintf("%d", s.player.HealthMax),
		)

		messages = append(messages, "Восстанавливает все твои недостающие Жизненные Силы, но и придает тебе 25 Дополнительных Жизненных Сил")
	} else {
		simpleChange = true
	}

	result, err := expression.RunAndReturnRoundUint(
		fmt.Sprintf("%d %s", s.player.Health, healthChange),
	)
	if err != nil {
		return nil, []string{}, err
	}

	if simpleChange {
		if s.player.Health < result {
			messageResult := result - s.player.Health
			messages = append(messages, fmt.Sprintf("<p>+%d HP</p>", messageResult))
		} else if result <= s.player.Health {
			messageResult := s.player.Health - result
			messages = append(messages, fmt.Sprintf("<p>-%d HP</p>", messageResult))
		}
	}

	if result <= 0 {
		s.player.Health = 0
	} else {
		s.player.Health = result
	}

	return s.player, messages, nil
}
