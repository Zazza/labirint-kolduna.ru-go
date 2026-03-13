package template

import (
	"context"
	"errors"
	"fmt"
)

func GetDiceTemplate(ctx context.Context, number uint, isBig bool) (string, error) {
	if number < 1 || 9 < number {
		return "", errors.New("invalid number")
	}

	classBig := ""
	if isBig {
		classBig = "dice-big"
	}

	if number > 6 {
		number1 := 9 - number
		number2 := 7 - number1

		return fmt.Sprintf("<span class='dice %s dice-%d'></span>"+
			"<span class='dice %s dice-%d'></span>", classBig, number1, classBig, number2), nil
	}

	return fmt.Sprintf("<span class='dice %s dice-%d'></span>", classBig, number), nil
}

func GetDicesTemplate(ctx context.Context, number1, number2 uint, isBig bool) (string, error) {
	if number1 < 1 || 9 < number1 {
		return "", errors.New("invalid number1")
	}
	if number2 < 1 || 9 < number2 {
		return "", errors.New("invalid number2")
	}

	result1, err := GetDiceTemplate(ctx, number1, isBig)
	if err != nil {
		return "", err
	}
	result2, err := GetDiceTemplate(ctx, number2, isBig)
	if err != nil {
		return "", err
	}

	return result1 + result2, nil
}
