package utils

import (
	"fmt"

	"github.com/google/uuid"
)

var NilUUID = uuid.Nil

func GenerateUUID() uuid.UUID {
	return uuid.New()
}

func ParseUUID(uuidStr string) (uuid.UUID, error) {
	return uuid.Parse(uuidStr)
}

func IsValidUUID(uuidStr string) bool {
	_, err := uuid.Parse(uuidStr)
	return err == nil
}

func IsNil(id uuid.UUID) bool {
	return id == uuid.Nil
}

func String(id uuid.UUID) string {
	return id.String()
}

func UUIDSliceToStringSlice(ids []uuid.UUID) []string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = id.String()
	}
	return result
}

func StringSliceToUUIDSlice(strs []string) ([]uuid.UUID, error) {
	result := make([]uuid.UUID, len(strs))
	for i, str := range strs {
		id, err := uuid.Parse(str)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID at index %d: %w", i, err)
		}
		result[i] = id
	}
	return result, nil
}

func ContainsUUID(slice []uuid.UUID, item uuid.UUID) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
