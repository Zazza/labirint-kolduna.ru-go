package utils

import (
	"errors"
	"fmt"
)

var (
	ErrEmptySlice    = errors.New("slice is empty")
	ErrEmptyMap      = errors.New("map is empty")
	ErrEmptyString   = errors.New("string is empty")
	ErrInvalidLength = errors.New("invalid length")
	ErrInvalidRange  = errors.New("value out of range")
	ErrInvalidFormat = errors.New("invalid format")
)

func IsEmptySlice[T any](slice []T) bool {
	return len(slice) == 0
}

func IsNotEmptySlice[T any](slice []T) bool {
	return len(slice) > 0
}

func IsEmptyMap[K comparable, V any](m map[K]V) bool {
	return len(m) == 0
}

func IsNotEmptyMap[K comparable, V any](m map[K]V) bool {
	return len(m) > 0
}

func IsEmptyString(s string) bool {
	return len(s) == 0
}

func IsNotEmptyString(s string) bool {
	return len(s) > 0
}

func IsNilOrEmpty[T any](slice []T) bool {
	return slice == nil || len(slice) == 0
}

func IsNotNilOrEmpty[T any](slice []T) bool {
	return slice != nil && len(slice) > 0
}

func RequireNotEmpty[T any](slice []T, fieldName string) error {
	if len(slice) == 0 {
		return NewAppError(ErrCodeValidation, fieldName+" cannot be empty", ErrEmptySlice)
	}
	return nil
}

func RequireNotEmptyString(s string, fieldName string) error {
	if len(s) == 0 {
		return NewAppError(ErrCodeValidation, fieldName+" cannot be empty", ErrEmptyString)
	}
	return nil
}

func InRange(value, min, max int) bool {
	return value >= min && value <= max
}

func RequireInRange(value, min, max int, fieldName string) error {
	if value < min || value > max {
		return NewAppError(ErrCodeValidation, fmt.Sprintf("%s out of range [%d-%d]", fieldName, min, max), ErrInvalidRange)
	}
	return nil
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func SliceContains[T comparable](slice []T, item T) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func MapContains[K comparable, V any](m map[K]V, key K) bool {
	_, exists := m[key]
	return exists
}

func IntSliceSum(slice []int) int {
	sum := 0
	for _, v := range slice {
		sum += v
	}
	return sum
}

func IntSliceMax(slice []int) int {
	if len(slice) == 0 {
		return 0
	}
	max := slice[0]
	for _, v := range slice[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func UintSliceSum(slice []uint) uint {
	sum := uint(0)
	for _, v := range slice {
		sum += v
	}
	return sum
}

func UintSliceMax(slice []uint) uint {
	if len(slice) == 0 {
		return 0
	}
	max := slice[0]
	for _, v := range slice[1:] {
		if v > max {
			max = v
		}
	}
	return max
}
