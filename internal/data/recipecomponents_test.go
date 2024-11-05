package data

import (
	"fmt"
	"testing"

	"github.com/tconnellan/macro-tracker-backend/internal/assert"
	"github.com/tconnellan/macro-tracker-backend/internal/validator"
)

func TestRecipeComponentHelpers(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	tests := []struct {
		name            string
		valid           bool
		recipeComponent RecipeComponent
	}{
		{
			name:  "valid component",
			valid: true,
			recipeComponent: RecipeComponent{
				ID:              1,
				RecipeID:        1,
				PantryItemID:    1,
				CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:        1,
				StepNo:          1,
				StepDescription: "description",
			},
		},
		{
			name:  "invalid zero step no",
			valid: false,
			recipeComponent: RecipeComponent{
				ID:              1,
				RecipeID:        1,
				PantryItemID:    1,
				CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:        1,
				StepNo:          0,
				StepDescription: "description",
			},
		},
		{
			name:  "invalid bad step no",
			valid: false,
			recipeComponent: RecipeComponent{
				ID:              1,
				RecipeID:        1,
				PantryItemID:    1,
				CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:        1,
				StepNo:          -1,
				StepDescription: "description",
			},
		},
		{
			name:  "invalid component zero Quantity",
			valid: false,
			recipeComponent: RecipeComponent{
				ID:              1,
				RecipeID:        1,
				PantryItemID:    1,
				CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:        0,
				StepNo:          1,
				StepDescription: "description",
			},
		},
		{
			name:  "invalid component bad Quantity",
			valid: false,
			recipeComponent: RecipeComponent{
				ID:              1,
				RecipeID:        1,
				PantryItemID:    1,
				CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:        -1,
				StepNo:          1,
				StepDescription: "description",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()

			ValidateRecipeComponent(v, &tt.recipeComponent)
			assert.ValidatorValid(t, v, tt.valid)
		})
	}
}

func TestRecipeComponentModelGet(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name            string
		expectError     error
		ID              int64
		expectComponent RecipeComponent
	}{
		{
			name:        "get component existing",
			ID:          1,
			expectError: nil,
			expectComponent: RecipeComponent{
				ID:              1,
				RecipeID:        1,
				PantryItemID:    1,
				CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:        4,
				StepNo:          1,
				StepDescription: "step 1",
			},
		},
		{
			name:        "get componenent bad ID",
			ID:          -1,
			expectError: ErrRecordNotFound,
		},
		{
			name:        "get componenent zero ID",
			ID:          0,
			expectError: ErrRecordNotFound,
		},
		{
			name:        "get componenent non existing ID",
			ID:          999999999,
			expectError: ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipe_component")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := RecipeComponentModel{db}

			component, err := m.Get(tt.ID)

			assert.ExpectError(t, err, tt.expectError)
			if err != nil {
				return
			}

			assert.Equal(t, *component, tt.expectComponent)
		})
	}
}

func TestRecipeComponentModelInsert(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name        string
		expectError error
		component   RecipeComponent
	}{
		{
			name:        "insert component valid",
			expectError: nil,
			component: RecipeComponent{
				RecipeID:        4,
				PantryItemID:    1,
				CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:        4,
				StepNo:          2,
				StepDescription: "step 1",
			},
		},
		{
			name:        "insert component bad recipe id",
			expectError: ErrRecipeDoesNotExist,
			component: RecipeComponent{
				RecipeID:        -1,
				PantryItemID:    1,
				CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:        4,
				StepNo:          1,
				StepDescription: "step 1",
			},
		},
		{
			name:        "insert component zero recipe id",
			expectError: ErrRecipeDoesNotExist,
			component: RecipeComponent{
				RecipeID:        0,
				PantryItemID:    1,
				CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:        4,
				StepNo:          1,
				StepDescription: "step 1",
			},
		},
		{
			name:        "insert component non existent recipe id",
			expectError: ErrRecipeDoesNotExist,
			component: RecipeComponent{
				RecipeID:        9999,
				PantryItemID:    1,
				CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:        4,
				StepNo:          1,
				StepDescription: "step 1",
			},
		},
		{
			name:        "insert component bad pantryitem id",
			expectError: ErrPantryItemDoesNotExist,
			component: RecipeComponent{
				RecipeID:        1,
				PantryItemID:    -1,
				CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:        4,
				StepNo:          1,
				StepDescription: "step 1",
			},
		},
		{
			name:        "insert component zero pantryitem id",
			expectError: ErrPantryItemDoesNotExist,
			component: RecipeComponent{
				RecipeID:        1,
				PantryItemID:    0,
				CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:        4,
				StepNo:          1,
				StepDescription: "step 1",
			},
		},
		{
			name:        "insert component non existent pantryitem id",
			expectError: ErrPantryItemDoesNotExist,
			component: RecipeComponent{
				RecipeID:        1,
				PantryItemID:    999999,
				CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:        4,
				StepNo:          1,
				StepDescription: "step 1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipe_component")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := RecipeComponentModel{db}

			err = m.Insert(&tt.component)

			assert.ExpectError(t, err, tt.expectError)
			if err != nil {
				return
			}

			newComponent, err := m.Get(tt.component.ID)
			assert.ExpectError(t, err, nil)
			if err != nil {
				return
			}

			assert.Equal(t, *newComponent, tt.component)
		})
	}
}

func TestRecipeComponentModelUpdate(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name        string
		expectError error
		ID          int64
		component   RecipeComponent
	}{
		{
			name:        "get component existing",
			ID:          1,
			expectError: nil,
			component: RecipeComponent{
				ID:              1,
				RecipeID:        1,
				PantryItemID:    1,
				CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:        4,
				StepNo:          1,
				StepDescription: "step 1 modified",
			},
		},
		{
			name:        "get componenent bad ID",
			ID:          -1,
			expectError: ErrRecordNotFound,
			component: RecipeComponent{
				ID:              -1,
				RecipeID:        1,
				PantryItemID:    1,
				CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:        4,
				StepNo:          1,
				StepDescription: "step 1 modified",
			},
		},
		{
			name:        "get componenent zero ID",
			ID:          0,
			expectError: ErrRecordNotFound,
			component: RecipeComponent{
				ID:              0,
				RecipeID:        1,
				PantryItemID:    1,
				CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:        4,
				StepNo:          1,
				StepDescription: "step 1 modified",
			},
		},
		{
			name:        "get componenent non existing ID",
			ID:          999999999,
			expectError: ErrRecordNotFound,
			component: RecipeComponent{
				ID:              999999999,
				RecipeID:        1,
				PantryItemID:    1,
				CreatedAt:       MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:        4,
				StepNo:          1,
				StepDescription: "step 1 modified",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipe_component")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := RecipeComponentModel{db}

			err = m.Update(&tt.component)

			assert.ExpectError(t, err, tt.expectError)
			if err != nil {
				return
			}

			updatedComponent, err := m.Get(tt.ID)
			assert.ExpectError(t, err, nil)
			if err != nil {
				return
			}

			assert.Equal(t, *updatedComponent, tt.component)
		})
	}
}
