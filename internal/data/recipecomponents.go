package data

import (
	"database/sql"
	"time"

	"github.com/tconnellan/macro-tracker-backend/internal/validator"
)

type RecipeComponent struct {
	ID              int64     `json:"id"`
	RecipeID        int64     `json:"recipe_id"`
	ConsumableID    int64     `json:"consumable_id"`
	CreatedAt       time.Time `json:"created_at"`
	Quantity        int64     `json:"quantity"`
	StepNo          int64     `json:"step_no"`
	StepDescription string    `json:"step_description"`
}

func ValidateRecipeComponent(v *validator.Validator, recipeComponent *RecipeComponent) {
	v.Check(recipeComponent.Quantity > 0, "quantity", "must be positive")
	v.Check(recipeComponent.StepNo > 0, "step_no", "must be positive")
	v.Check(len(recipeComponent.StepDescription) <= 1000, "step_description", "must be at must 1000")
}

type RecipeComponentModel struct {
	DB *sql.DB
}

type IRecipeComponent interface {
	Get(int64) (*RecipeComponent, error)
	Insert(*RecipeComponent) error
	Update(*RecipeComponent) error
	Delete(int64) error
}
