package mocks

import (
	"time"

	"github.com/tconnellan/macro-tracker-backend/internal/data"
)

func NewTestModel() data.Models {
	return data.Models{
		Users:            UserModelMock{},
		Tokens:           TokenModelMock{},
		Consumed:         ConsumedModelMock{},
		Consumables:      ConsumableModelMock{},
		Recipes:          RecipeModelMock{},
		RecipeComponents: RecipeComponentModelMock{},
		PantryItems:      PantryItemModelMock{},
	}
}

func MustParse(layout, value string) time.Time {
	t, _ := time.Parse(layout, value)
	return t
}
