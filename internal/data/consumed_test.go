package data

import (
	"fmt"
	"testing"
	"time"

	"github.com/tconnellan/macro-tracker-backend/internal/assert"
	"github.com/tconnellan/macro-tracker-backend/internal/validator"
)

func TestConsumedHelpers(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	tests := []struct {
		name     string
		valid    bool
		consumed Consumed
	}{
		{
			name:  "valid component",
			valid: true,
			consumed: Consumed{
				ID:           1,
				RecipeID:     1,
				UserID:       1,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
			},
		},
		{
			name:  "invalid component negative quantity",
			valid: false,
			consumed: Consumed{
				ID:           1,
				RecipeID:     1,
				UserID:       1,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     -1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
			},
		},
		{
			name:  "invalid component zero quantity",
			valid: false,
			consumed: Consumed{
				ID:           1,
				RecipeID:     1,
				UserID:       1,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     0,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
			},
		},
		// {
		// 	name:  "invalid component bad user ID",
		// 	valid: false,
		// 	consumed: Consumed{
		// 		RecipeID:     -1,
		// 		UserID:       1,
		// 		CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
		// 		ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
		// 		Quantity:     1,
		// 		LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
		// 		Macros: Macronutrients{
		// 			Carbs:    1,
		// 			Fats:     1,
		// 			Proteins: 1,
		// 			Alcohol:  1,
		// 		},
		// 	},
		// },
		// {
		// 	name:  "invalid component zero user ID",
		// 	valid: false,
		// 	consumed: Consumed{
		// 		RecipeID:     1,
		// 		UserID:       0,
		// 		CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
		// 		ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
		// 		Quantity:     1,
		// 		LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
		// 		Macros: Macronutrients{
		// 			Carbs:    1,
		// 			Fats:     1,
		// 			Proteins: 1,
		// 			Alcohol:  1,
		// 		},
		// 	},
		// },
		// {
		// 	name:  "invalid component no existent user ID",
		// 	valid: false,
		// 	consumed: Consumed{
		// 		RecipeID:     1,
		// 		UserID:       9999999,
		// 		CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
		// 		ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
		// 		Quantity:     1,
		// 		LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
		// 		Macros: Macronutrients{
		// 			Carbs:    1,
		// 			Fats:     1,
		// 			Proteins: 1,
		// 			Alcohol:  1,
		// 		},
		// 	},
		// },
		// {
		// 	name:  "invalid component bad recipe ID",
		// 	valid: false,
		// 	consumed: Consumed{
		// 		RecipeID:     -1,
		// 		UserID:       1,
		// 		CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
		// 		ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
		// 		Quantity:     1,
		// 		LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
		// 		Macros: Macronutrients{
		// 			Carbs:    1,
		// 			Fats:     1,
		// 			Proteins: 1,
		// 			Alcohol:  1,
		// 		},
		// 	},
		// },
		// {
		// 	name:  "invalid component zero recipe ID",
		// 	valid: false,
		// 	consumed: Consumed{
		// 		RecipeID:     0,
		// 		UserID:       1,
		// 		CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
		// 		ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
		// 		Quantity:     1,
		// 		LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
		// 		Macros: Macronutrients{
		// 			Carbs:    1,
		// 			Fats:     1,
		// 			Proteins: 1,
		// 			Alcohol:  1,
		// 		},
		// 	},
		// },
		// {
		// 	name:  "invalid component no existent recipe ID",
		// 	valid: false,
		// 	consumed: Consumed{
		// 		RecipeID:     99999,
		// 		UserID:       1,
		// 		CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
		// 		ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
		// 		Quantity:     1,
		// 		LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
		// 		Macros: Macronutrients{
		// 			Carbs:    1,
		// 			Fats:     1,
		// 			Proteins: 1,
		// 			Alcohol:  1,
		// 		},
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()

			ValidateConsumed(v, &tt.consumed)
			assert.ValidatorValid(t, v, tt.valid)
		})
	}
}

func TestConsumedModelGetByConsumedID(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name           string
		expectError    error
		ID             int64
		expectConsumed Consumed
	}{
		{
			name:        "get component existing",
			ID:          1,
			expectError: nil,
			expectConsumed: Consumed{
				ID:           1,
				RecipeID:     1,
				UserID:       1,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
				Notes: "notes",
			},
		},
		{
			name:        "get component zero id",
			ID:          0,
			expectError: ErrRecordNotFound,
		},
		{
			name:        "get component negative id",
			ID:          -1,
			expectError: ErrRecordNotFound,
		},
		{
			name:        "get component non existent id",
			ID:          99999999,
			expectError: ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipe_component")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := ConsumedModel{db}

			consumed, err := m.GetByConsumedID(tt.ID)

			assert.ExpectError(t, err, tt.expectError)
			if err != nil {
				return
			}

			assert.Equal(t, *consumed, tt.expectConsumed)
		})
	}
}

func TestConsumedModelGetAllByUserID(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name           string
		expectError    error
		ID             int64
		expectConsumed []*Consumed
	}{
		{
			name:        "get component existing",
			ID:          1,
			expectError: nil,
			expectConsumed: []*Consumed{
				{
					ID:           1,
					RecipeID:     1,
					UserID:       1,
					CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
					ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Quantity:     1,
					LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Macros: Macronutrients{
						Carbs:    1,
						Fats:     1,
						Proteins: 1,
						Alcohol:  1,
					},
					Notes: "notes",
				},
				{
					ID:           2,
					RecipeID:     2,
					UserID:       1,
					CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
					ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Quantity:     5,
					LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Macros: Macronutrients{
						Carbs:    7,
						Fats:     7,
						Proteins: 7,
						Alcohol:  1,
					},
					Notes: "notes 2",
				},
				{
					ID:           3,
					RecipeID:     2,
					UserID:       1,
					CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
					ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Quantity:     8,
					LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Macros: Macronutrients{
						Carbs:    3,
						Fats:     2,
						Proteins: 1,
						Alcohol:  0,
					},
					Notes: "notes 3",
				},
			},
		},
		{
			name:           "get component zero id",
			ID:             0,
			expectError:    nil,
			expectConsumed: []*Consumed{},
		},
		{
			name:           "get component negative id",
			ID:             -1,
			expectError:    nil,
			expectConsumed: []*Consumed{},
		},
		{
			name:           "get component non existent id",
			ID:             99999999,
			expectError:    nil,
			expectConsumed: []*Consumed{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipe_component")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := ConsumedModel{db}

			consumed, err := m.GetAllByUserID(tt.ID)

			assert.ExpectError(t, err, tt.expectError)

			assert.SliceEqual(t, consumed, tt.expectConsumed)
		})
	}
}

func TestConsumedModelGetAllByUserIDAndDate(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name           string
		expectError    error
		ID             int64
		from           time.Time
		to             time.Time
		expectConsumed []*Consumed
	}{
		{
			name:        "get component existing",
			ID:          3,
			from:        MustParse(timeFormat, "2024-01-01 10:00:00"),
			to:          MustParse(timeFormat, "2024-01-01 10:10:00"),
			expectError: nil,
			expectConsumed: []*Consumed{
				{
					ID:           5,
					RecipeID:     1,
					UserID:       3,
					CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
					ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Quantity:     1,
					LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Macros: Macronutrients{
						Carbs:    1,
						Fats:     1,
						Proteins: 1,
						Alcohol:  1,
					},
					Notes: "notes",
				},
				{
					ID:           6,
					RecipeID:     2,
					UserID:       3,
					CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
					ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:01:00"),
					Quantity:     5,
					LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Macros: Macronutrients{
						Carbs:    7,
						Fats:     7,
						Proteins: 7,
						Alcohol:  1,
					},
					Notes: "notes 2",
				},
				{
					ID:           7,
					RecipeID:     2,
					UserID:       3,
					CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
					ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:10:00"),
					Quantity:     8,
					LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Macros: Macronutrients{
						Carbs:    3,
						Fats:     2,
						Proteins: 1,
						Alcohol:  0,
					},
					Notes: "notes 3",
				},
			},
		},
		{
			name:        "get component existing limited",
			ID:          3,
			from:        MustParse(timeFormat, "2024-01-01 10:00:00"),
			to:          MustParse(timeFormat, "2024-01-01 10:09:59"),
			expectError: nil,
			expectConsumed: []*Consumed{
				{
					ID:           5,
					RecipeID:     1,
					UserID:       3,
					CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
					ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Quantity:     1,
					LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Macros: Macronutrients{
						Carbs:    1,
						Fats:     1,
						Proteins: 1,
						Alcohol:  1,
					},
					Notes: "notes",
				},
				{
					ID:           6,
					RecipeID:     2,
					UserID:       3,
					CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
					ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:01:00"),
					Quantity:     5,
					LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Macros: Macronutrients{
						Carbs:    7,
						Fats:     7,
						Proteins: 7,
						Alcohol:  1,
					},
					Notes: "notes 2",
				},
			},
		},
		{
			name:        "get component existing cross day",
			ID:          3,
			from:        MustParse(timeFormat, "2024-01-01 10:00:00"),
			to:          MustParse(timeFormat, "2024-01-02 10:10:00"),
			expectError: nil,
			expectConsumed: []*Consumed{
				{
					ID:           5,
					RecipeID:     1,
					UserID:       3,
					CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
					ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
					Quantity:     1,
					LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Macros: Macronutrients{
						Carbs:    1,
						Fats:     1,
						Proteins: 1,
						Alcohol:  1,
					},
					Notes: "notes",
				},
				{
					ID:           6,
					RecipeID:     2,
					UserID:       3,
					CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
					ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:01:00"),
					Quantity:     5,
					LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Macros: Macronutrients{
						Carbs:    7,
						Fats:     7,
						Proteins: 7,
						Alcohol:  1,
					},
					Notes: "notes 2",
				},
				{
					ID:           7,
					RecipeID:     2,
					UserID:       3,
					CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
					ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:10:00"),
					Quantity:     8,
					LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Macros: Macronutrients{
						Carbs:    3,
						Fats:     2,
						Proteins: 1,
						Alcohol:  0,
					},
					Notes: "notes 3",
				},
				{
					ID:           8,
					RecipeID:     2,
					UserID:       3,
					CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
					ConsumedAt:   MustParse(timeFormat, "2024-01-02 10:00:00"),
					Quantity:     8,
					LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Macros: Macronutrients{
						Carbs:    3,
						Fats:     2,
						Proteins: 1,
						Alcohol:  0,
					},
					Notes: "notes 3",
				},
			},
		},
		{
			name:           "get component zero id",
			ID:             0,
			from:           MustParse(timeFormat, "2024-01-01 09:00:00"),
			to:             MustParse(timeFormat, "2024-01-01 11:00:00"),
			expectError:    nil,
			expectConsumed: []*Consumed{},
		},
		{
			name:           "get component negative id",
			ID:             -1,
			from:           MustParse(timeFormat, "2024-01-01 09:00:00"),
			to:             MustParse(timeFormat, "2024-01-01 11:00:00"),
			expectError:    nil,
			expectConsumed: []*Consumed{},
		},
		{
			name:           "get component non existent id",
			ID:             99999999,
			from:           MustParse(timeFormat, "2024-01-01 09:00:00"),
			to:             MustParse(timeFormat, "2024-01-01 11:00:00"),
			expectError:    nil,
			expectConsumed: []*Consumed{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipe_component")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := ConsumedModel{db}

			consumed, err := m.GetAllByUserIDAndDate(tt.ID, tt.from, tt.to)

			assert.ExpectError(t, err, tt.expectError)
			if err != nil {
				return
			}

			assert.SliceEqual(t, consumed, tt.expectConsumed)
		})
	}
}

func TestConsumedModelInsert(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name        string
		expectError error
		consumed    Consumed
	}{
		{
			name:        "insert component existing",
			expectError: nil,
			consumed: Consumed{
				RecipeID:     3,
				UserID:       3,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     0,
					Proteins: 0,
					Alcohol:  0.5,
				},
				Notes: "notes",
			},
		},
		{
			name:        "invalid component bad user ID",
			expectError: ErrReferencedUserDoesNotExist,
			consumed: Consumed{
				RecipeID:     1,
				UserID:       -1,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
			},
		},
		{
			name:        "invalid component zero user ID",
			expectError: ErrReferencedUserDoesNotExist,
			consumed: Consumed{
				RecipeID:     1,
				UserID:       0,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
			},
		},
		{
			name:        "invalid component no existent user ID",
			expectError: ErrReferencedUserDoesNotExist,
			consumed: Consumed{
				RecipeID:     1,
				UserID:       9999999,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
			},
		},
		{
			name:        "invalid component bad recipe ID",
			expectError: ErrRecipeDoesNotExist,
			consumed: Consumed{
				RecipeID:     -1,
				UserID:       1,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
			},
		},
		{
			name:        "invalid component zero recipe ID",
			expectError: ErrRecipeDoesNotExist,
			consumed: Consumed{
				RecipeID:     0,
				UserID:       1,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
			},
		},
		{
			name:        "invalid component no existent recipe ID",
			expectError: ErrRecipeDoesNotExist,
			consumed: Consumed{
				RecipeID:     99999,
				UserID:       1,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipe_component")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := ConsumedModel{db}

			err = m.Insert(&tt.consumed)

			assert.ExpectError(t, err, tt.expectError)
			if err != nil {
				return
			}
		})
	}
}

func TestConsumedModelUpdate(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name        string
		ID          int64
		expectError error
		consumed    Consumed
	}{
		{
			name:        "update component existing",
			ID:          1,
			expectError: nil,
			consumed: Consumed{
				ID:           1,
				RecipeID:     4,
				UserID:       4,
				CreatedAt:    MustParse(timeFormat, "2024-01-02 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-04 10:00:00"),
				Quantity:     6.7,
				LastEditedAt: MustParse(timeFormat, "2024-01-06 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1.6,
					Fats:     7,
					Proteins: 1,
					Alcohol:  0.2,
				},
				Notes: "notes modified",
			},
		},
		{
			name:        "invalid bad component id",
			expectError: ErrRecordNotFound,
			consumed: Consumed{
				ID:           -1,
				RecipeID:     1,
				UserID:       1,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
			},
		},
		{
			name:        "invalid zero component id",
			expectError: ErrRecordNotFound,
			consumed: Consumed{
				ID:           0,
				RecipeID:     1,
				UserID:       1,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
			},
		},
		{
			name:        "invalid non existent component id",
			expectError: ErrRecordNotFound,
			consumed: Consumed{
				ID:           999999,
				RecipeID:     1,
				UserID:       1,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
			},
		},
		{
			name:        "invalid component bad user ID",
			expectError: ErrReferencedUserDoesNotExist,
			consumed: Consumed{
				ID:           1,
				RecipeID:     1,
				UserID:       -1,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
			},
		},
		{
			name:        "invalid component zero user ID",
			expectError: ErrReferencedUserDoesNotExist,
			consumed: Consumed{
				ID:           1,
				RecipeID:     1,
				UserID:       0,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
			},
		},
		{
			name:        "invalid component no existent user ID",
			expectError: ErrReferencedUserDoesNotExist,
			consumed: Consumed{
				ID:           1,
				RecipeID:     1,
				UserID:       9999999,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
			},
		},
		{
			name:        "invalid component bad recipe ID",
			expectError: ErrRecipeDoesNotExist,
			consumed: Consumed{
				ID:           1,
				RecipeID:     -1,
				UserID:       1,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
			},
		},
		{
			name:        "invalid component zero recipe ID",
			expectError: ErrRecipeDoesNotExist,
			consumed: Consumed{
				ID:           1,
				RecipeID:     0,
				UserID:       1,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
			},
		},
		{
			name:        "invalid component no existent recipe ID",
			expectError: ErrRecipeDoesNotExist,
			consumed: Consumed{
				ID:           1,
				RecipeID:     99999,
				UserID:       1,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipe_component")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := ConsumedModel{db}

			err = m.Update(&tt.consumed)

			assert.ExpectError(t, err, tt.expectError)
			if err != nil {
				return
			}

			updated, err := m.GetByConsumedID(tt.ID)
			assert.ExpectError(t, err, nil)
			if err != nil {
				return
			}
			tt.consumed.CreatedAt = updated.CreatedAt
			tt.consumed.LastEditedAt = updated.LastEditedAt
			assert.Equal(t, *updated, tt.consumed)

		})
	}
}

func TestConsumedModelDelete(t *testing.T) {

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name        string
		expectError error
		ID          int64
	}{
		{
			name:        "delete existing component no components",
			expectError: nil,
			ID:          1,
		},
		{
			name:        "delete existing component invalid ID",
			expectError: ErrRecordNotFound,
			ID:          -1,
		},
		{
			name:        "delete existing component non existing ID",
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
			m := ConsumedModel{db}

			err = m.Delete(tt.ID)

			assert.ExpectError(t, err, tt.expectError)
			if err != nil {
				return
			}

			_, err = m.GetByConsumedID(tt.ID)
			assert.ExpectError(t, err, ErrRecordNotFound)
		})
	}
}
