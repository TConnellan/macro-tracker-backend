package data

import (
	"fmt"
	"testing"

	"github.com/tconnellan/macro-tracker-backend/internal/assert"
	"github.com/tconnellan/macro-tracker-backend/internal/validator"
)

func TestRecipeHelpers(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	tests := []struct {
		name   string
		valid  bool
		recipe Recipe
	}{
		{
			name:  "valid recipe",
			valid: true,
			recipe: Recipe{
				ID:             1,
				Name:           "Lasagne",
				CreatorID:      1,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "some notes",
				ParentRecipeID: 0,
				IsLatest:       true,
			},
		},
		{
			name:  "invalid recipe name empty",
			valid: false,
			recipe: Recipe{
				ID:             1,
				Name:           "",
				CreatorID:      1,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "some notes",
				ParentRecipeID: 0,
				IsLatest:       true,
			},
		},
		{
			name:  "invalid recipe name too long",
			valid: false,
			recipe: Recipe{
				ID:             1,
				Name:           "LasagneLasagneLasagneLasagneLasagneLasagneLasagneLasagne",
				CreatorID:      1,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "some notes",
				ParentRecipeID: 0,
				IsLatest:       true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			v := validator.New()

			ValidateRecipe(v, &tt.recipe)
			assert.ValidatorValid(t, v, tt.valid)
		})
	}
}

func TestFullRecipeHelpers(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	tests := []struct {
		name       string
		valid      bool
		recipe     Recipe
		fullRecipe FullRecipe
	}{
		{
			name:  "valid recipe",
			valid: true,
			fullRecipe: FullRecipe{
				Recipe: Recipe{
					ID:             1,
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "some notes",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				RecipeComponents: []*RecipeComponent{
					{
						ID:              1,
						RecipeID:        1,
						PantryItemID:    1,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        4,
						StepNo:          1,
						StepDescription: "step 1",
					},
				},
				PantryItems: []*PantryItem{
					{
						ID:           1,
						UserID:       1,
						ConsumableId: 1,
						Name:         "lasagne sheet",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
				},
				Consumables: []*Consumable{
					{
						ID:        1,
						CreatorID: 1,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "Lasagne Past Large",
						BrandName: "San Remo",
						Size:      62.5,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    46.6,
							Fats:     0.9,
							Proteins: 7.9,
							Alcohol:  0,
						},
					},
				},
			},
		},
		{
			name:  "invalid recipe empty",
			valid: false,
			fullRecipe: FullRecipe{
				Recipe: Recipe{
					ID:             1,
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "some notes",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				RecipeComponents: []*RecipeComponent{},
				PantryItems:      []*PantryItem{},
				Consumables:      []*Consumable{},
			},
		},
		{
			name:  "invalid recipe non matching recipe id",
			valid: false,
			fullRecipe: FullRecipe{
				Recipe: Recipe{
					ID:             1,
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "some notes",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				RecipeComponents: []*RecipeComponent{
					{
						ID:              1,
						RecipeID:        9999,
						PantryItemID:    1,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        4,
						StepNo:          1,
						StepDescription: "step 1",
					},
				},
				PantryItems: []*PantryItem{
					{
						ID:           1,
						UserID:       1,
						ConsumableId: 1,
						Name:         "lasagne sheet",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
				},
				Consumables: []*Consumable{
					{
						ID:        1,
						CreatorID: 1,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "Lasagne Past Large",
						BrandName: "San Remo",
						Size:      62.5,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    46.6,
							Fats:     0.9,
							Proteins: 7.9,
							Alcohol:  0,
						},
					},
				},
			},
		},
		{
			name:  "valid recipe multipart",
			valid: true,
			fullRecipe: FullRecipe{
				Recipe: Recipe{
					ID:             1,
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "some notes",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				RecipeComponents: []*RecipeComponent{
					{
						ID:              1,
						RecipeID:        1,
						PantryItemID:    1,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        4,
						StepNo:          1,
						StepDescription: "step 1",
					},
					{
						ID:              1,
						RecipeID:        1,
						PantryItemID:    2,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        5,
						StepNo:          2,
						StepDescription: "step 2",
					},
				},
				PantryItems: []*PantryItem{
					{
						ID:           1,
						UserID:       1,
						ConsumableId: 1,
						Name:         "lasagne sheet",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
					{
						ID:           2,
						UserID:       1,
						ConsumableId: 2,
						Name:         "Minced Beef",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
				},
				Consumables: []*Consumable{
					{
						ID:        1,
						CreatorID: 1,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "Lasagne Pasta Large",
						BrandName: "San Remo",
						Size:      62.5,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    46.6,
							Fats:     0.9,
							Proteins: 7.9,
							Alcohol:  0,
						},
					},
					{
						ID:        2,
						CreatorID: 1,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "No Added Hormone Beef 5 Star Extra Trim Mince",
						BrandName: "Coles",
						Size:      100,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    .5,
							Fats:     2,
							Proteins: 21.3,
							Alcohol:  0,
						},
					},
				},
			},
		},
		{
			name:  "invalid recipe duplicate step no",
			valid: false,
			fullRecipe: FullRecipe{
				Recipe: Recipe{
					ID:             1,
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "some notes",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				RecipeComponents: []*RecipeComponent{
					{
						ID:              1,
						RecipeID:        1,
						PantryItemID:    1,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        4,
						StepNo:          1,
						StepDescription: "step 1",
					},
					{
						ID:              1,
						RecipeID:        1,
						PantryItemID:    1,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        5,
						StepNo:          1,
						StepDescription: "step 2",
					},
				},
				PantryItems: []*PantryItem{
					{
						ID:           1,
						UserID:       1,
						ConsumableId: 1,
						Name:         "lasagne sheet",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
					{
						ID:           2,
						UserID:       1,
						ConsumableId: 2,
						Name:         "Minced Beef",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
				},
				Consumables: []*Consumable{
					{
						ID:        1,
						CreatorID: 1,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "Lasagne Pasta Large",
						BrandName: "San Remo",
						Size:      62.5,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    46.6,
							Fats:     0.9,
							Proteins: 7.9,
							Alcohol:  0,
						},
					},
					{
						ID:        2,
						CreatorID: 1,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "No Added Hormone Beef 5 Star Extra Trim Mince",
						BrandName: "Coles",
						Size:      100,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    .5,
							Fats:     2,
							Proteins: 21.3,
							Alcohol:  0,
						},
					},
				},
			},
		},
		{
			name:  "invalid recipe step no out of range",
			valid: false,
			fullRecipe: FullRecipe{
				Recipe: Recipe{
					ID:             1,
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "some notes",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				RecipeComponents: []*RecipeComponent{
					{
						ID:              1,
						RecipeID:        1,
						PantryItemID:    1,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        4,
						StepNo:          1,
						StepDescription: "step 1",
					},
					{
						ID:              1,
						RecipeID:        1,
						PantryItemID:    1,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        5,
						StepNo:          3,
						StepDescription: "step 2",
					},
				},
				PantryItems: []*PantryItem{
					{
						ID:           1,
						UserID:       1,
						ConsumableId: 1,
						Name:         "lasagne sheet",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
					{
						ID:           2,
						UserID:       1,
						ConsumableId: 2,
						Name:         "Minced Beef",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
				},
				Consumables: []*Consumable{
					{
						ID:        1,
						CreatorID: 1,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "Lasagne Pasta Large",
						BrandName: "San Remo",
						Size:      62.5,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    46.6,
							Fats:     0.9,
							Proteins: 7.9,
							Alcohol:  0,
						},
					},
					{
						ID:        2,
						CreatorID: 1,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "No Added Hormone Beef 5 Star Extra Trim Mince",
						BrandName: "Coles",
						Size:      100,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    .5,
							Fats:     2,
							Proteins: 21.3,
							Alcohol:  0,
						},
					},
				},
			},
		},
		{
			name:  "invalid recipe non matching length items",
			valid: false,
			fullRecipe: FullRecipe{
				Recipe: Recipe{
					ID:             1,
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "some notes",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				RecipeComponents: []*RecipeComponent{
					{
						ID:              1,
						RecipeID:        1,
						PantryItemID:    1,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        4,
						StepNo:          1,
						StepDescription: "step 1",
					},
					{
						ID:              1,
						RecipeID:        1,
						PantryItemID:    1,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        5,
						StepNo:          2,
						StepDescription: "step 2",
					},
				},
				PantryItems: []*PantryItem{
					{
						ID:           1,
						UserID:       1,
						ConsumableId: 1,
						Name:         "lasagne sheet",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
				},
				Consumables: []*Consumable{
					{
						ID:        1,
						CreatorID: 1,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "Lasagne Pasta Large",
						BrandName: "San Remo",
						Size:      62.5,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    46.6,
							Fats:     0.9,
							Proteins: 7.9,
							Alcohol:  0,
						},
					},
					{
						ID:        2,
						CreatorID: 1,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "No Added Hormone Beef 5 Star Extra Trim Mince",
						BrandName: "Coles",
						Size:      100,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    .5,
							Fats:     2,
							Proteins: 21.3,
							Alcohol:  0,
						},
					},
				},
			},
		},
		{
			name:  "invalid recipe non matching length consumables",
			valid: false,
			fullRecipe: FullRecipe{
				Recipe: Recipe{
					ID:             1,
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "some notes",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				RecipeComponents: []*RecipeComponent{
					{
						ID:              1,
						RecipeID:        1,
						PantryItemID:    1,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        4,
						StepNo:          1,
						StepDescription: "step 1",
					},
					{
						ID:              1,
						RecipeID:        1,
						PantryItemID:    1,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        5,
						StepNo:          2,
						StepDescription: "step 2",
					},
				},
				PantryItems: []*PantryItem{
					{
						ID:           1,
						UserID:       1,
						ConsumableId: 1,
						Name:         "lasagne sheet",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
					{
						ID:           2,
						UserID:       1,
						ConsumableId: 2,
						Name:         "Minced Beef",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
				},
				Consumables: []*Consumable{
					{
						ID:        1,
						CreatorID: 1,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "Lasagne Pasta Large",
						BrandName: "San Remo",
						Size:      62.5,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    46.6,
							Fats:     0.9,
							Proteins: 7.9,
							Alcohol:  0,
						},
					},
				},
			},
		},
		{
			name:  "invalid recipe nonmatching length recipe components",
			valid: false,
			fullRecipe: FullRecipe{
				Recipe: Recipe{
					ID:             1,
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "some notes",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				RecipeComponents: []*RecipeComponent{
					{
						ID:              1,
						RecipeID:        1,
						PantryItemID:    1,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        4,
						StepNo:          1,
						StepDescription: "step 1",
					},
				},
				PantryItems: []*PantryItem{
					{
						ID:           1,
						UserID:       1,
						ConsumableId: 1,
						Name:         "lasagne sheet",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
					{
						ID:           2,
						UserID:       1,
						ConsumableId: 2,
						Name:         "Minced Beef",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
				},
				Consumables: []*Consumable{
					{
						ID:        1,
						CreatorID: 1,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "Lasagne Pasta Large",
						BrandName: "San Remo",
						Size:      62.5,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    46.6,
							Fats:     0.9,
							Proteins: 7.9,
							Alcohol:  0,
						},
					},
					{
						ID:        2,
						CreatorID: 1,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "No Added Hormone Beef 5 Star Extra Trim Mince",
						BrandName: "Coles",
						Size:      100,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    .5,
							Fats:     2,
							Proteins: 21.3,
							Alcohol:  0,
						},
					},
				},
			},
		},
		{
			name:  "invalid recipe nonmatching pantry ids",
			valid: false,
			fullRecipe: FullRecipe{
				Recipe: Recipe{
					ID:             1,
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "some notes",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				RecipeComponents: []*RecipeComponent{
					{
						ID:              1,
						RecipeID:        1,
						PantryItemID:    1,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        4,
						StepNo:          1,
						StepDescription: "step 1",
					},
					{
						ID:              1,
						RecipeID:        1,
						PantryItemID:    99999,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        5,
						StepNo:          2,
						StepDescription: "step 2",
					},
				},
				PantryItems: []*PantryItem{
					{
						ID:           1,
						UserID:       1,
						ConsumableId: 1,
						Name:         "lasagne sheet",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
					{
						ID:           2,
						UserID:       1,
						ConsumableId: 2,
						Name:         "Minced Beef",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
				},
				Consumables: []*Consumable{
					{
						ID:        1,
						CreatorID: 1,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "Lasagne Pasta Large",
						BrandName: "San Remo",
						Size:      62.5,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    46.6,
							Fats:     0.9,
							Proteins: 7.9,
							Alcohol:  0,
						},
					},
					{
						ID:        2,
						CreatorID: 1,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "No Added Hormone Beef 5 Star Extra Trim Mince",
						BrandName: "Coles",
						Size:      100,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    .5,
							Fats:     2,
							Proteins: 21.3,
							Alcohol:  0,
						},
					},
				},
			},
		},
		{
			name:  "invalid recipe nonmatching consumable ids",
			valid: false,
			fullRecipe: FullRecipe{
				Recipe: Recipe{
					ID:             1,
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "some notes",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				RecipeComponents: []*RecipeComponent{
					{
						ID:              1,
						RecipeID:        1,
						PantryItemID:    1,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        4,
						StepNo:          1,
						StepDescription: "step 1",
					},
					{
						ID:              1,
						RecipeID:        1,
						PantryItemID:    2,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        5,
						StepNo:          2,
						StepDescription: "step 2",
					},
				},
				PantryItems: []*PantryItem{
					{
						ID:           1,
						UserID:       1,
						ConsumableId: 1,
						Name:         "lasagne sheet",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
					{
						ID:           2,
						UserID:       1,
						ConsumableId: 2,
						Name:         "Minced Beef",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
				},
				Consumables: []*Consumable{
					{
						ID:        1,
						CreatorID: 1,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "Lasagne Pasta Large",
						BrandName: "San Remo",
						Size:      62.5,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    46.6,
							Fats:     0.9,
							Proteins: 7.9,
							Alcohol:  0,
						},
					},
					{
						ID:        99999,
						CreatorID: 1,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "No Added Hormone Beef 5 Star Extra Trim Mince",
						BrandName: "Coles",
						Size:      100,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    .5,
							Fats:     2,
							Proteins: 21.3,
							Alcohol:  0,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			v := validator.New()

			ValidateFullRecipe(v, &tt.fullRecipe)
			assert.ValidatorValid(t, v, tt.valid)

			if v.Valid() {
				for i, recipeComponent := range tt.fullRecipe.RecipeComponents {
					if i+1 != int(recipeComponent.StepNo) {
						t.Errorf("Component at index: %d, has stepNo %d, expect: %d", i, int(recipeComponent.StepNo), i+1)
					}
				}
			}
		})
	}
}

func TestRecipeModelGet(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name         string
		ID           int64
		expectError  error
		expectRecipe Recipe
	}{
		{
			name:        "existing recipe",
			ID:          1,
			expectError: nil,
			expectRecipe: Recipe{
				ID:             1,
				Name:           "Lasagne",
				CreatorID:      1,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "a recipe",
				ParentRecipeID: 0,
				IsLatest:       true,
			},
		},
		{
			name:        "non-existing recipe Bad ID",
			ID:          -1,
			expectError: ErrRecordNotFound,
		},
		{
			name:        "non-existing recipe Zero ID",
			ID:          0,
			expectError: ErrRecordNotFound,
		},
		{
			name:        "non-existing recipe",
			ID:          99999,
			expectError: ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipe")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := RecipeModel{db}

			recipe, err := m.Get(tt.ID)

			assert.ExpectError(t, err, tt.expectError)
			if err != nil {
				return
			}
			assert.Equal(t, *recipe, tt.expectRecipe)

		})
	}
}

func TestRecipeModelGetByCreatorID(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name           string
		creatorID      int64
		filters        RecipeFilters
		expectError    error
		expectRecipes  []*Recipe
		expectMetadata Metadata
	}{
		{
			name:      "existing recipe",
			creatorID: 1,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectError: nil,
			expectMetadata: Metadata{
				CurrentPage:  1,
				PageSize:     100,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 1,
			},
			expectRecipes: []*Recipe{
				{
					ID:             1,
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
			},
		},
		{
			name:      "existing recipes multiple",
			creatorID: 2,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectError: nil,
			expectMetadata: Metadata{
				CurrentPage:  1,
				PageSize:     100,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 4,
			},
			expectRecipes: []*Recipe{
				{
					ID:             2,
					Name:           "recipe2",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 2",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				{
					ID:             3,
					Name:           "recipe3",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 3",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				{
					ID:             5,
					Name:           "Recipe5",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 5",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				{
					ID:             6,
					Name:           "doesntmatchsearch",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "not matching",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
			},
		},
		{
			name:      "existing recipes multiple search",
			creatorID: 2,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "recipe",
			},
			expectError: nil,
			expectMetadata: Metadata{
				CurrentPage:  1,
				PageSize:     100,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 3,
			},
			expectRecipes: []*Recipe{
				{
					ID:             2,
					Name:           "recipe2",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 2",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				{
					ID:             3,
					Name:           "recipe3",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 3",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				{
					ID:             5,
					Name:           "Recipe5",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 5",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
			},
		},
		{
			name:      "existing recipes multiple limit",
			creatorID: 2,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     2,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectError: nil,
			expectMetadata: Metadata{
				CurrentPage:  1,
				PageSize:     2,
				FirstPage:    1,
				LastPage:     2,
				TotalRecords: 4,
			},
			expectRecipes: []*Recipe{
				{
					ID:             2,
					Name:           "recipe2",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 2",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				{
					ID:             3,
					Name:           "recipe3",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 3",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
			},
		},
		{
			name:      "existing recipes multiple limit and offset",
			creatorID: 2,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         2,
					PageSize:     2,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectError: nil,
			expectMetadata: Metadata{
				CurrentPage:  2,
				PageSize:     2,
				FirstPage:    1,
				LastPage:     2,
				TotalRecords: 4,
			},
			expectRecipes: []*Recipe{
				{
					ID:             5,
					Name:           "Recipe5",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 5",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				{
					ID:             6,
					Name:           "doesntmatchsearch",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "not matching",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
			},
		},
		{
			name:      "existing recipe with parent",
			creatorID: 3,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectError: nil,
			expectMetadata: Metadata{
				CurrentPage:  1,
				PageSize:     100,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 1,
			},
			expectRecipes: []*Recipe{
				{
					ID:             4,
					Name:           "recipe4",
					CreatorID:      3,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 4",
					ParentRecipeID: 1,
					IsLatest:       true,
				},
			},
		},
		{
			name:      "non-existing recipes Bad Creator ID",
			creatorID: -1,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectError: nil,
		},
		{
			name:      "non-existing recipes Zero Creator ID",
			creatorID: 0,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectError: nil,
		},
		{
			name:      "non-existing recipes",
			creatorID: 99999,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipe")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := RecipeModel{db}

			recipes, metadata, err := m.GetByCreatorID(tt.creatorID, tt.filters)

			assert.ExpectError(t, err, tt.expectError)
			if err != nil {
				return
			}
			assert.SliceEqual(t, recipes, tt.expectRecipes)
			assert.Equal(t, metadata, tt.expectMetadata)

		})
	}
}

func TestRecipeModelGetLatestByCreatorID(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name           string
		creatorID      int64
		filters        RecipeFilters
		expectError    error
		expectRecipes  []*Recipe
		expectMetadata Metadata
	}{
		{
			name:      "existing recipe",
			creatorID: 1,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectError: nil,
			expectMetadata: Metadata{
				CurrentPage:  1,
				PageSize:     100,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 1,
			},
			expectRecipes: []*Recipe{
				{
					ID:             1,
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
			},
		},
		{
			name:      "existing recipes multiple",
			creatorID: 2,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectError: nil,
			expectMetadata: Metadata{
				CurrentPage:  1,
				PageSize:     100,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 4,
			},
			expectRecipes: []*Recipe{
				{
					ID:             2,
					Name:           "recipe2",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 2",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				{
					ID:             3,
					Name:           "recipe3",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 3",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				{
					ID:             5,
					Name:           "Recipe5",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 5",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				{
					ID:             6,
					Name:           "doesntmatchsearch",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "not matching",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
			},
		},
		{
			name:      "existing recipes multiple search",
			creatorID: 2,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "recipe",
			},
			expectError: nil,
			expectMetadata: Metadata{
				CurrentPage:  1,
				PageSize:     100,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 3,
			},
			expectRecipes: []*Recipe{
				{
					ID:             2,
					Name:           "recipe2",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 2",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				{
					ID:             3,
					Name:           "recipe3",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 3",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				{
					ID:             5,
					Name:           "Recipe5",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 5",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
			},
		},
		{
			name:      "existing recipes multiple limit",
			creatorID: 2,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     2,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectError: nil,
			expectMetadata: Metadata{
				CurrentPage:  1,
				PageSize:     2,
				FirstPage:    1,
				LastPage:     2,
				TotalRecords: 4,
			},
			expectRecipes: []*Recipe{
				{
					ID:             2,
					Name:           "recipe2",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 2",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				{
					ID:             3,
					Name:           "recipe3",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 3",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
			},
		},
		{
			name:      "existing recipes multiple limit and offset",
			creatorID: 2,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         2,
					PageSize:     2,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectError: nil,
			expectMetadata: Metadata{
				CurrentPage:  2,
				PageSize:     2,
				FirstPage:    1,
				LastPage:     2,
				TotalRecords: 4,
			},
			expectRecipes: []*Recipe{
				{
					ID:             5,
					Name:           "Recipe5",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 5",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				{
					ID:             6,
					Name:           "doesntmatchsearch",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "not matching",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
			},
		},
		{
			name:      "existing recipe with parent",
			creatorID: 3,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectError: nil,
			expectMetadata: Metadata{
				CurrentPage:  1,
				PageSize:     100,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 1,
			},
			expectRecipes: []*Recipe{
				{
					ID:             4,
					Name:           "recipe4",
					CreatorID:      3,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 4",
					ParentRecipeID: 1,
					IsLatest:       true,
				},
			},
		},
		{
			name:      "non-existing recipes Bad Creator ID",
			creatorID: -1,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectError: nil,
		},
		{
			name:      "non-existing recipes Zero Creator ID",
			creatorID: 0,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectError: nil,
		},
		{
			name:      "non-existing recipes",
			creatorID: 99999,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipe")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := RecipeModel{db}

			recipes, metadata, err := m.GetLatestByCreatorID(tt.creatorID, tt.filters)

			assert.ExpectError(t, err, tt.expectError)
			if err != nil {
				return
			}
			assert.SliceEqual(t, recipes, tt.expectRecipes)
			assert.Equal(t, metadata, tt.expectMetadata)

		})
	}
}

func TestGetFullRecipe(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	tests := []struct {
		name             string
		ID               int64
		expectError      error
		expectFullRecipe FullRecipe
	}{
		{
			name:        "valid fullrecipe",
			ID:          1,
			expectError: nil,
			expectFullRecipe: FullRecipe{
				Recipe: Recipe{
					ID:             1,
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				RecipeComponents: []*RecipeComponent{
					{
						ID:              1,
						RecipeID:        1,
						PantryItemID:    1,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        4,
						StepNo:          1,
						StepDescription: "step 1",
					},
					{
						ID:              2,
						RecipeID:        1,
						PantryItemID:    2,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        5,
						StepNo:          2,
						StepDescription: "step 2",
					},
				},
				PantryItems: []*PantryItem{
					{
						ID:           1,
						UserID:       1,
						ConsumableId: 17,
						Name:         "lasagne sheet",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
					{
						ID:           2,
						UserID:       1,
						ConsumableId: 18,
						Name:         "Minced Beef",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
				},
				Consumables: []*Consumable{
					{
						ID:        17,
						CreatorID: 4,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "Lasagne Pasta Large",
						BrandName: "San Remo",
						Size:      62.5,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    46.6,
							Fats:     0.9,
							Proteins: 7.9,
							Alcohol:  0,
						},
					},
					{
						ID:        18,
						CreatorID: 4,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "No Added Hormone Beef 5 Star Extra Trim Mince",
						BrandName: "Coles",
						Size:      100,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    .5,
							Fats:     2,
							Proteins: 21.3,
							Alcohol:  0,
						},
					},
				},
			},
		},
		{
			name:        "invalid fullrecipe zero id",
			ID:          0,
			expectError: ErrRecordNotFound,
		},
		{
			name:        "invalid fullrecipe bad id",
			ID:          -1,
			expectError: ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipe")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := RecipeModel{db}

			fullRecipe, err := m.GetFullRecipe(tt.ID)

			assert.ExpectError(t, err, tt.expectError)
			if err != nil {
				return
			}
			assert.Equal(t, fullRecipe.Recipe, tt.expectFullRecipe.Recipe)
			assert.SliceEqual(t, fullRecipe.RecipeComponents, tt.expectFullRecipe.RecipeComponents)
			assert.SliceEqual(t, fullRecipe.PantryItems, tt.expectFullRecipe.PantryItems)
			assert.SliceEqual(t, fullRecipe.Consumables, tt.expectFullRecipe.Consumables)
		})
	}
}

func TestInsert(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	tests := []struct {
		name        string
		expectError error
		recipe      Recipe
	}{
		{
			name:        "valid recipe",
			expectError: nil,
			recipe: Recipe{
				Name:           "test recipe",
				CreatorID:      4,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "some notes",
				ParentRecipeID: 0,
				IsLatest:       true,
			},
		},
		{
			name:        "valid recipe parent exists",
			expectError: nil,
			recipe: Recipe{
				Name:           "test recipe 2",
				CreatorID:      4,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "some notes",
				ParentRecipeID: 1,
				IsLatest:       true,
			},
		},
		{
			name:        "invalid recipe bad Creator ID",
			expectError: ErrReferencedUserDoesNotExist,
			recipe: Recipe{
				Name:           "test recipe 2",
				CreatorID:      9999,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "some notes",
				ParentRecipeID: 1,
				IsLatest:       true,
			},
		},
		{
			name:        "invalid recipe bad parent id",
			expectError: ErrParentRecipeDoesNotExist,
			recipe: Recipe{
				Name:           "test recipe 2",
				CreatorID:      4,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "some notes",
				ParentRecipeID: -1,
				IsLatest:       true,
			},
		},
		{
			name:        "invalid recipe non existent parent id",
			expectError: ErrParentRecipeDoesNotExist,
			recipe: Recipe{
				Name:           "test recipe 2",
				CreatorID:      4,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "some notes",
				ParentRecipeID: 9999999,
				IsLatest:       true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipe")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := RecipeModel{db}

			err = m.Insert(&tt.recipe)

			assert.ExpectError(t, err, tt.expectError)
			if err != nil {
				return
			}
		})
	}
}

func TestInsertFullRecipe(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	tests := []struct {
		name        string
		expectError error
		fullRecipe  FullRecipe
	}{
		{
			name:        "insert valid fullrecipe",
			expectError: nil,
			fullRecipe: FullRecipe{
				Recipe: Recipe{
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				RecipeComponents: []*RecipeComponent{
					{
						RecipeID:        1,
						PantryItemID:    1,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        4,
						StepNo:          1,
						StepDescription: "step 1",
					},
					{
						RecipeID:        1,
						PantryItemID:    2,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        5,
						StepNo:          2,
						StepDescription: "step 2",
					},
				},
				PantryItems: []*PantryItem{
					{
						UserID:       1,
						ConsumableId: 17,
						Name:         "lasagne sheet",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
					{
						UserID:       1,
						ConsumableId: 18,
						Name:         "Minced Beef",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
				},
				Consumables: []*Consumable{
					{
						CreatorID: 4,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "Lasagne Pasta Large",
						BrandName: "San Remo",
						Size:      62.5,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    46.6,
							Fats:     0.9,
							Proteins: 7.9,
							Alcohol:  0,
						},
					},
					{
						CreatorID: 4,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "No Added Hormone Beef 5 Star Extra Trim Mince",
						BrandName: "Coles",
						Size:      100,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    .5,
							Fats:     2,
							Proteins: 21.3,
							Alcohol:  0,
						},
					},
				},
			},
		},
		{
			name:        "insert invalid fullrecipe bad recipe user",
			expectError: ErrReferencedUserDoesNotExist,
			fullRecipe: FullRecipe{
				Recipe: Recipe{
					Name:           "Lasagne",
					CreatorID:      -1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				RecipeComponents: []*RecipeComponent{
					{
						RecipeID:        1,
						PantryItemID:    1,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        4,
						StepNo:          1,
						StepDescription: "step 1",
					},
					{
						RecipeID:        1,
						PantryItemID:    2,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        5,
						StepNo:          2,
						StepDescription: "step 2",
					},
				},
				PantryItems: []*PantryItem{
					{
						UserID:       1,
						ConsumableId: 17,
						Name:         "lasagne sheet",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
					{
						UserID:       1,
						ConsumableId: 18,
						Name:         "Minced Beef",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
				},
				Consumables: []*Consumable{
					{
						CreatorID: 4,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "Lasagne Pasta Large",
						BrandName: "San Remo",
						Size:      62.5,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    46.6,
							Fats:     0.9,
							Proteins: 7.9,
							Alcohol:  0,
						},
					},
					{
						CreatorID: 4,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "No Added Hormone Beef 5 Star Extra Trim Mince",
						BrandName: "Coles",
						Size:      100,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    .5,
							Fats:     2,
							Proteins: 21.3,
							Alcohol:  0,
						},
					},
				},
			},
		},
		{
			name:        "insert invalid fullrecipe pantry item does not exist",
			expectError: ErrPantryItemDoesNotExist,
			fullRecipe: FullRecipe{
				Recipe: Recipe{
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				RecipeComponents: []*RecipeComponent{
					{
						RecipeID:        1,
						PantryItemID:    99999,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        4,
						StepNo:          1,
						StepDescription: "step 1",
					},
					{
						RecipeID:        1,
						PantryItemID:    2,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        5,
						StepNo:          2,
						StepDescription: "step 2",
					},
				},
				PantryItems: []*PantryItem{
					{
						UserID:       -1,
						ConsumableId: 17,
						Name:         "lasagne sheet",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
					{
						UserID:       1,
						ConsumableId: 18,
						Name:         "Minced Beef",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
				},
				Consumables: []*Consumable{
					{
						CreatorID: 4,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "Lasagne Pasta Large",
						BrandName: "San Remo",
						Size:      62.5,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    46.6,
							Fats:     0.9,
							Proteins: 7.9,
							Alcohol:  0,
						},
					},
					{
						CreatorID: 4,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "No Added Hormone Beef 5 Star Extra Trim Mince",
						BrandName: "Coles",
						Size:      100,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    .5,
							Fats:     2,
							Proteins: 21.3,
							Alcohol:  0,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipe")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := RecipeModel{db}

			err = m.InsertFullRecipe(&tt.fullRecipe)

			assert.ExpectError(t, err, tt.expectError)
			if err != nil {
				return
			}

			createdRecipe, err := m.GetFullRecipe(tt.fullRecipe.Recipe.ID)
			assert.ExpectError(t, err, nil)
			if err != nil {
				return
			}
			assert.Equal(t, createdRecipe.Recipe, tt.fullRecipe.Recipe)
			assert.SliceEqual(t, createdRecipe.RecipeComponents, tt.fullRecipe.RecipeComponents)
			assert.SliceEqual(t, createdRecipe.PantryItems, tt.fullRecipe.PantryItems)
			assert.SliceEqual(t, createdRecipe.Consumables, tt.fullRecipe.Consumables)
		})
	}
}

func TestUpdate(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	tests := []struct {
		name        string
		ID          int64
		expectError error
		newRecipe   Recipe
	}{
		{
			name:        "update valid recipe",
			ID:          1,
			expectError: nil,
			newRecipe: Recipe{
				ID:             1,
				Name:           "Lasagne",
				CreatorID:      1,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "a recipe with modification",
				ParentRecipeID: 0,
				IsLatest:       true,
			},
		},
		{
			name:        "update invalid recipe",
			ID:          -1,
			expectError: ErrRecordNotFound,
			newRecipe: Recipe{
				ID:             -1,
				Name:           "Lasagne",
				CreatorID:      1,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "a recipe with modification",
				ParentRecipeID: 0,
				IsLatest:       true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipe")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := RecipeModel{db}

			err = m.Update(&tt.newRecipe)
			assert.ExpectError(t, err, tt.expectError)
			if err != nil {
				return
			}
			updatedRecipe, err := m.Get(tt.ID)
			assert.ExpectError(t, err, nil)
			if err != nil {
				return
			}

			// last edited at timestamp will not be predicable
			tt.newRecipe.LastEditedAt = updatedRecipe.LastEditedAt
			assert.Equal(t, *updatedRecipe, tt.newRecipe)

		})
	}
}

func TestUpdateFullRecipe(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	tests := []struct {
		name          string
		expectError   error
		newFullRecipe FullRecipe
	}{
		{
			name:        "update valid fullrecipe",
			expectError: nil,
			newFullRecipe: FullRecipe{
				Recipe: Recipe{
					ID:             1,
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe with modification",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				RecipeComponents: []*RecipeComponent{
					{
						RecipeID:        1,
						PantryItemID:    1,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        4,
						StepNo:          1,
						StepDescription: "step 1",
					},
					{
						RecipeID:        1,
						PantryItemID:    2,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        5,
						StepNo:          2,
						StepDescription: "step 2",
					},
				},
				PantryItems: []*PantryItem{
					{
						UserID:       1,
						ConsumableId: 17,
						Name:         "lasagne sheet",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
					{
						UserID:       1,
						ConsumableId: 18,
						Name:         "Minced Beef",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
				},
				Consumables: []*Consumable{
					{
						CreatorID: 4,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "Lasagne Pasta Large",
						BrandName: "San Remo",
						Size:      62.5,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    46.6,
							Fats:     0.9,
							Proteins: 7.9,
							Alcohol:  0,
						},
					},
					{
						CreatorID: 4,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "No Added Hormone Beef 5 Star Extra Trim Mince",
						BrandName: "Coles",
						Size:      100,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    .5,
							Fats:     2,
							Proteins: 21.3,
							Alcohol:  0,
						},
					},
				},
			},
		},
		{
			name:        "update invalid fullrecipe bad recipe user",
			expectError: ErrReferencedUserDoesNotExist,
			newFullRecipe: FullRecipe{
				Recipe: Recipe{
					ID:             1,
					Name:           "Lasagne",
					CreatorID:      -1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				RecipeComponents: []*RecipeComponent{
					{
						RecipeID:        1,
						PantryItemID:    1,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        4,
						StepNo:          1,
						StepDescription: "step 1",
					},
					{
						RecipeID:        1,
						PantryItemID:    2,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        5,
						StepNo:          2,
						StepDescription: "step 2",
					},
				},
				PantryItems: []*PantryItem{
					{
						UserID:       1,
						ConsumableId: 17,
						Name:         "lasagne sheet",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
					{
						UserID:       1,
						ConsumableId: 18,
						Name:         "Minced Beef",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
				},
				Consumables: []*Consumable{
					{
						CreatorID: 4,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "Lasagne Pasta Large",
						BrandName: "San Remo",
						Size:      62.5,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    46.6,
							Fats:     0.9,
							Proteins: 7.9,
							Alcohol:  0,
						},
					},
					{
						CreatorID: 4,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "No Added Hormone Beef 5 Star Extra Trim Mince",
						BrandName: "Coles",
						Size:      100,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    .5,
							Fats:     2,
							Proteins: 21.3,
							Alcohol:  0,
						},
					},
				},
			},
		},
		{
			name:        "update invalid fullrecipe pantry item does not exist",
			expectError: ErrPantryItemDoesNotExist,
			newFullRecipe: FullRecipe{
				Recipe: Recipe{
					ID:             1,
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				RecipeComponents: []*RecipeComponent{
					{
						RecipeID:        1,
						PantryItemID:    99999,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        4,
						StepNo:          1,
						StepDescription: "step 1",
					},
					{
						RecipeID:        1,
						PantryItemID:    2,
						CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
						Quantity:        5,
						StepNo:          2,
						StepDescription: "step 2",
					},
				},
				PantryItems: []*PantryItem{
					{
						UserID:       -1,
						ConsumableId: 17,
						Name:         "lasagne sheet",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
					{
						UserID:       1,
						ConsumableId: 18,
						Name:         "Minced Beef",
						CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
						LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					},
				},
				Consumables: []*Consumable{
					{
						CreatorID: 4,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "Lasagne Pasta Large",
						BrandName: "San Remo",
						Size:      62.5,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    46.6,
							Fats:     0.9,
							Proteins: 7.9,
							Alcohol:  0,
						},
					},
					{
						CreatorID: 4,
						CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
						Name:      "No Added Hormone Beef 5 Star Extra Trim Mince",
						BrandName: "Coles",
						Size:      100,
						Units:     "g",
						Macros: Macronutrients{
							Carbs:    .5,
							Fats:     2,
							Proteins: 21.3,
							Alcohol:  0,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipe")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := RecipeModel{db}

			err = m.UpdateFullRecipe(&tt.newFullRecipe)

			assert.ExpectError(t, err, tt.expectError)
			if err != nil {
				return
			}

			createdRecipe, err := m.GetFullRecipe(tt.newFullRecipe.Recipe.ID)
			assert.ExpectError(t, err, nil)
			if err != nil {
				return
			}
			assert.Equal(t, createdRecipe.Recipe, tt.newFullRecipe.Recipe)
			assert.SliceEqual(t, createdRecipe.RecipeComponents, tt.newFullRecipe.RecipeComponents)
			assert.SliceEqual(t, createdRecipe.PantryItems, tt.newFullRecipe.PantryItems)
			assert.SliceEqual(t, createdRecipe.Consumables, tt.newFullRecipe.Consumables)
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name        string
		expectError error
		ID          int64
	}{
		{
			name:        "delete existing recipe no components",
			expectError: nil,
			ID:          5,
		},
		{
			name:        "delete existing recipe valid with component",
			expectError: nil,
			ID:          7,
		},
		{
			name:        "delete existing recipe invalid has child",
			expectError: ErrChildRecipeExists,
			ID:          1,
		},
		{
			name:        "delete existing recipe invalid ID",
			expectError: ErrRecordNotFound,
			ID:          -1,
		},
		{
			name:        "delete existing recipe non existing ID",
			expectError: ErrRecordNotFound,
			ID:          999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipe")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := RecipeModel{db}

			err = m.Delete(tt.ID)

			assert.ExpectError(t, err, tt.expectError)
			if err != nil {
				return
			}

			_, err = m.Get(tt.ID)
			assert.ExpectError(t, err, ErrRecordNotFound)
		})
	}
}

func TestGetParentRecipe(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	tests := []struct {
		name         string
		expectError  error
		childRecipe  Recipe
		expectParent Recipe
	}{
		{
			name:        "get parent existing recipe no parent",
			expectError: ErrRecordNotFound,
			childRecipe: Recipe{
				ID:             1,
				Name:           "Lasagne",
				CreatorID:      1,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "a recipe 1",
				ParentRecipeID: 0,
				IsLatest:       true,
			},
		},
		{
			name:        "get parent existing recipe with parent",
			expectError: nil,
			childRecipe: Recipe{
				ID:             8,
				Name:           "recipe8",
				CreatorID:      2,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "a recipe 8",
				ParentRecipeID: 1,
				IsLatest:       true,
			},
			expectParent: Recipe{
				ID:             1,
				Name:           "Lasagne",
				CreatorID:      1,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "a recipe",
				ParentRecipeID: 0,
				IsLatest:       true,
			},
		},
		{
			name:        "get parent invalid doesn't exist",
			expectError: ErrRecordNotFound,
			childRecipe: Recipe{
				ID:             999999,
				Name:           "Lasagne",
				CreatorID:      1,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "a recipe",
				ParentRecipeID: 0,
				IsLatest:       true,
			},
		},
		{
			name:        "get parent invalid ID doesn't exist",
			expectError: ErrRecordNotFound,
			childRecipe: Recipe{
				ID:             -1,
				Name:           "Lasagne",
				CreatorID:      1,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "a recipe",
				ParentRecipeID: 0,
				IsLatest:       true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipe")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := RecipeModel{db}

			parent, err := m.GetParentRecipe(&tt.childRecipe)

			assert.ExpectError(t, err, tt.expectError)
			if err != nil {
				return
			}

			assert.Equal(t, *parent, tt.expectParent)
		})
	}
}

func TestGetAllAncestors(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	tests := []struct {
		name            string
		expectError     error
		childRecipe     Recipe
		filters         RecipeFilters
		expectMetadata  Metadata
		expectAncestors []*Recipe
	}{
		{
			name:        "get parent existing recipe",
			expectError: nil,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectMetadata: Metadata{
				CurrentPage:  1,
				PageSize:     100,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 2,
			},
			expectAncestors: []*Recipe{
				{
					ID:             1,
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				{
					ID:             4,
					Name:           "recipe4",
					CreatorID:      3,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 4",
					ParentRecipeID: 1,
					IsLatest:       true,
				},
			},
			childRecipe: Recipe{
				ID:             4,
				Name:           "recipe4",
				CreatorID:      3,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "a recipe 4",
				ParentRecipeID: 1,
				IsLatest:       true,
			},
		},
		{
			name:        "get parent existing recipe with many parent",
			expectError: nil,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectMetadata: Metadata{
				CurrentPage:  1,
				PageSize:     100,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 4,
			},
			childRecipe: Recipe{
				ID:             10,
				Name:           "recipe10",
				CreatorID:      4,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "a recipe 10",
				ParentRecipeID: 9,
				IsLatest:       true,
			},
			expectAncestors: []*Recipe{
				{
					ID:             1,
					Name:           "Lasagne",
					CreatorID:      1,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe",
					ParentRecipeID: 0,
					IsLatest:       true,
				},
				{
					ID:             8,
					Name:           "recipe8",
					CreatorID:      2,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 8",
					ParentRecipeID: 1,
					IsLatest:       true,
				},
				{
					ID:             9,
					Name:           "recipe9",
					CreatorID:      4,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 9",
					ParentRecipeID: 8,
					IsLatest:       true,
				},
				{
					ID:             10,
					Name:           "recipe10",
					CreatorID:      4,
					CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
					LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Notes:          "a recipe 10",
					ParentRecipeID: 9,
					IsLatest:       true,
				},
			},
		},
		{
			name:        "get parent invalid doesn't exist",
			expectError: ErrRecordNotFound,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectMetadata: Metadata{},
			childRecipe: Recipe{
				ID:             999999,
				Name:           "Lasagne",
				CreatorID:      1,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "a recipe",
				ParentRecipeID: 0,
				IsLatest:       true,
			},
		},
		{
			name:        "get parent invalid ID doesn't exist",
			expectError: ErrRecordNotFound,
			filters: RecipeFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch: "",
			},
			expectMetadata: Metadata{},
			childRecipe: Recipe{
				ID:             -1,
				Name:           "Lasagne",
				CreatorID:      1,
				CreatedAt:      MustParse(timeFormat, "2024-01-01 10:00:00"),
				LastEditedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Notes:          "a recipe",
				ParentRecipeID: 0,
				IsLatest:       true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipe")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := RecipeModel{db}

			ancestors, metadata, err := m.GetAllAncestors(&tt.childRecipe, tt.filters)

			assert.ExpectError(t, err, tt.expectError)
			if err != nil {
				return
			}

			assert.SliceEqual(t, ancestors, tt.expectAncestors)
			assert.Equal(t, metadata, tt.expectMetadata)
		})
	}
}
