package expression

import (
	"gamebook-backend/modules/game/dto"
	"math"

	"github.com/expr-lang/expr"
)

func Run(expressionString string) (any, error) {
	program, err := expr.Compile(expressionString)
	if err != nil {
		return false, err
	}

	output, err := expr.Run(program, nil)
	if err != nil {
		return false, err
	}

	return output, nil
}

func RunAndReturnRoundUint(expressionString string) (uint, error) {
	result, err := Run(expressionString)
	if err != nil {
		return 0, err
	}

	roundFuncInt := func(result any) int {
		switch val := result.(type) {
		case uint:
			return int(val)
		case int:
			return val
		case float64:
			return int(math.Round(val))
		case float32:
			return int(math.Round(float64(val)))
		default:
			return 0
		}
	}

	intResult := roundFuncInt(result)
	if intResult < 0 {
		return 0, nil
	}

	return uint(intResult), nil
}

func RunAndReturnBoolean(expressionString string) (bool, error) {
	result, err := Run(expressionString)
	if err != nil {
		return false, err
	}

	resultBoolean, ok := result.(bool)
	if !ok {
		return false, dto.MessageExpectedBoolean
	}
	return resultBoolean, nil
}
