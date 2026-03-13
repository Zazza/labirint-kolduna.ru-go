package helper

import (
	"gamebook-backend/database/entities"
)

var BagItems = map[string]string{
	"wings":          "Канат длиною тридцать локтей (пятнадцать метров)",
	"wedges":         "Дюжину клиньев для скалолазания",
	"torches":        "Шесть факелов",
	"flashlight":     "Фонарь и четыре фляги масла к нему",
	"breadAndMeat":   "Хлеб и мясо",
	"apples":         "Два яблока",
	"mapIngredients": "Пергамент, гусиное перо и чернила, чтобы сделать карту",
	"garlic":         "Несколько головок чеснока",
	"hammer":         "Молоток, гвозди и пилу",
	"fire":           "Огниво",
}

func HasBagItem(bag []entities.Bag, itemName string) bool {
	for _, item := range bag {
		if item.Name == itemName {
			return true
		}
	}
	return false
}

func GetFullBagItems(items []string) []entities.Bag {
	var result []entities.Bag
	for _, item := range items {
		result = append(result, entities.Bag{
			Name:        item,
			Description: BagItems[item],
		})
	}

	return result
}
