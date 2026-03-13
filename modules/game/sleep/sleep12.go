package sleep

import (
	"context"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dice"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"
	template2 "gamebook-backend/modules/game/template"

	"gorm.io/gorm"
)

type sleep12 struct {
	db     *gorm.DB
	player entities.Player
}

func NewSleep12(db *gorm.DB, player entities.Player) dto.Sleep {
	return &sleep12{
		db:     db,
		player: player,
	}
}

func (s *sleep12) Execute(
	ctx context.Context,
	_ uint,
	_ uint,
) (dto.SleepyKingdomDTO, error) {
	enemyHealth := uint(40)
	enemyDamage := uint(20)

	playerHealth := s.player.Health
	playerArrows := uint(7)
	playerArrowDamage := uint(10)

	rollTheDices := dice.NewRollTheDices(s.db, &s.player)

	fullArrowsDamage := uint(0)
	for i := 0; i <= int(playerArrows); i++ {
		diceBow1, diceBow2, err := rollTheDices.RollTheDices(ctx, s.player)
		if err != nil {
			return dto.SleepyKingdomDTO{}, err
		}

		templateArrows, err := template2.GetDicesTemplate(ctx, *diceBow1, *diceBow2, false)
		if err != nil {
			return dto.SleepyKingdomDTO{}, err
		}

		if 6 <= *diceBow1+*diceBow2 {
			s.writeToChannel(fmt.Sprintf("<p>%s ты попал стрелой!</p>", templateArrows))
			fullArrowsDamage += playerArrowDamage
		} else {
			s.writeToChannel(fmt.Sprintf("<p>%s ты промазал!</p>", templateArrows))
		}
	}

	if fullArrowsDamage >= enemyHealth {
		s.writeToChannel("<p>Людоед победил</p>")

		return dto.SleepyKingdomDTO{
			Exit:    true,
			Death:   false,
			NextTry: false,
		}, nil
	}

	death := false
	exit := false

	step := "enemy"
	for {
		if playerHealth == 0 {
			death = true

			s.writeToChannel("<p>Людоед победил</p>")

			break
		}
		if enemyHealth == 0 {
			exit = true

			s.writeToChannel("<p>Ты победил людоеда!</p>")

			break
		}

		if step == "enemy" {
			if enemyDamage >= playerHealth {
				playerHealth = 0
			} else {
				playerHealth -= enemyDamage
			}

			s.writeToChannel(fmt.Sprintf("<p>Людоед наносит урон в %d жизненных сил!</p>", enemyDamage))

			step = "player"
		} else {
			dicePlayerDamage1, dicePlayerDamage2, err := rollTheDices.RollTheDices(ctx, s.player)
			if err != nil {
				return dto.SleepyKingdomDTO{}, err
			}

			var playerDamage uint
			if *dicePlayerDamage1+*dicePlayerDamage2 <= 6 {
				playerDamage = 0
			} else {
				playerDamage = *dicePlayerDamage1 + *dicePlayerDamage2 - 6
			}

			if playerDamage >= enemyHealth {
				enemyHealth = 0
			} else {
				enemyHealth -= playerDamage
			}

			templateDamage, err := template2.GetDicesTemplate(ctx, *dicePlayerDamage1, *dicePlayerDamage2, false)
			if err != nil {
				return dto.SleepyKingdomDTO{}, err
			}

			s.writeToChannel(fmt.Sprintf(
				"<p>%s Ты наносишь урон руками %d HP!</p>",
				templateDamage,
				playerDamage,
			))

			step = "enemy"
		}

	}

	return dto.SleepyKingdomDTO{
		Exit:    exit,
		Death:   death,
		NextTry: false,
	}, nil
}

func (s *sleep12) writeToChannel(message string) {
	helper.DescriptionMessage(
		s.player.ID,
		message,
	)
}
