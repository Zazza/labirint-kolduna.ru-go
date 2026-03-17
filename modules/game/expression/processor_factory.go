package expression

import (
	battleDTO "gamebook-backend/modules/game/dto"
	"math"
	// "github.com/google/uuid"
	// "github.com/com/expr-lang/expr" - Temporarily disabled due to dependency issues
)

type ExpressionResultProcessor interface {
	Process(result any) (any, error)
}

type ProcessorFactory interface {
	GetProcessor(resultType string) ExpressionResultProcessor
}

type expressionProcessorFactory struct {
	processors map[string]ExpressionResultProcessor
}

func NewProcessorFactory() ProcessorFactory {
	factory := &expressionProcessorFactory{
		processors: make(map[string]ExpressionResultProcessor),
	}

	factory.registerDefaults()
	return factory
}

func (f *expressionProcessorFactory) registerDefaults() {
	f.processors["uint"] = &UintProcessor{}
	f.processors["int"] = &IntProcessor{}
	f.processors["float64"] = &Float64Processor{}
	f.processors["float32"] = &Float32Processor{}
	f.processors["bool"] = &BoolProcessor{}
	f.processors["string"] = &StringProcessor{}
	f.processors["default"] = &DefaultProcessor{}
}

func (f *expressionProcessorFactory) GetProcessor(resultType string) ExpressionResultProcessor {
	if processor, exists := f.processors[resultType]; exists {
		return processor
	}
	return &DefaultProcessor{}
}

type UintProcessor struct{}

func (p *UintProcessor) Process(result any) (any, error) {
	uintResult, ok := result.(uint)
	if !ok {
		return result, nil
	}
	return uintResult, nil
}

type IntProcessor struct{}

func (p *IntProcessor) Process(result any) (any, error) {
	intResult, ok := result.(int)
	if !ok {
		return result, nil
	}
	return intResult, nil
}

type Float64Processor struct{}

func (p *Float64Processor) Process(result any) (any, error) {
	floatResult, ok := result.(float64)
	if !ok {
		return result, nil
	}
	return int(math.Round(floatResult)), nil
}

type Float32Processor struct{}

func (p *Float32Processor) Process(result any) (any, error) {
	floatResult, ok := result.(float32)
	if !ok {
		return result, nil
	}
	return int(math.Round(float64(floatResult))), nil
}

type BoolProcessor struct{}

func (p *BoolProcessor) Process(result any) (any, error) {
	boolResult, ok := result.(bool)
	if !ok {
		return false, battleDTO.MessageExpectedBoolean
	}
	return boolResult, nil
}

type StringProcessor struct{}

func (p *StringProcessor) Process(result any) (any, error) {
	stringResult, ok := result.(string)
	if !ok {
		return result, nil
	}
	return stringResult, nil
}

type DefaultProcessor struct{}

func (p *DefaultProcessor) Process(result any) (any, error) {
	return result, nil
}
