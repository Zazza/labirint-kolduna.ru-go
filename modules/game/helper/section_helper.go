package helper

import (
	"gamebook-backend/database/entities"
)

const MinSectionNumber = 0
const MaxSectionNumber = 155

func IsVisited(sectionNumber uint, playerSection []entities.PlayerSection) bool {
	for _, item := range GetVisitedSections(playerSection) {
		if item == sectionNumber {
			return true
		}
	}

	return false
}

func GetVisitedSections(playerSection []entities.PlayerSection) []uint {
	mapResult := make(map[uint]bool)
	for _, section := range playerSection {
		mapResult[section.Section.Number] = true
	}

	var result []uint
	for sectionNumber := range mapResult {
		if MinSectionNumber <= sectionNumber && sectionNumber <= MaxSectionNumber {
			result = append(result, sectionNumber)
		}
	}

	return result
}

func GetNotVisitedSections(playerSection []entities.PlayerSection) []uint {
	mapResult := make(map[uint]bool)
	for _, section := range playerSection {
		mapResult[section.Section.Number] = true
	}

	var result []uint
	for i := MinSectionNumber; i <= MaxSectionNumber; i++ {
		if _, ok := mapResult[uint(i)]; !ok {
			result = append(result, uint(i))
		}
	}

	return result
}

func GetAllSections() []uint {
	var result []uint
	for i := MinSectionNumber; i <= MaxSectionNumber; i++ {
		result = append(result, uint(i))
	}

	return result
}
