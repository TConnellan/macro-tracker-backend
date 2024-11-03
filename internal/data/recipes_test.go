package data

import (
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
