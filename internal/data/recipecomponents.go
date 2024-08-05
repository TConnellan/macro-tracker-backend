package data

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tconnellan/macro-tracker-backend/internal/validator"
)

type RecipeComponent struct {
	ID              int64     `json:"id"`
	RecipeID        int64     `json:"recipe_id"`
	PantryItemID    int64     `json:"pantry_item_id"`
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
	DB *pgxpool.Pool
}

type IRecipeComponentModel interface {
	Get(int64) (*RecipeComponent, error)
	Insert(*RecipeComponent) error
	Update(*RecipeComponent) error
}

func (m RecipeComponentModel) Get(ID int64) (*RecipeComponent, error) {
	stmt := `
	SELECT id, recipe_id, pantry_item_id, created_at, quantity, step_no, step_description
	FROM recipe_components
	WHERE id = $1
	`

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	var recipeComponent RecipeComponent

	err := m.DB.QueryRow(ctx, stmt, ID).Scan(
		&recipeComponent.ID,
		&recipeComponent.RecipeID,
		&recipeComponent.PantryItemID,
		&recipeComponent.CreatedAt,
		&recipeComponent.Quantity,
		&recipeComponent.StepNo,
		&recipeComponent.StepDescription,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &recipeComponent, nil

}

func (m RecipeComponentModel) Insert(recipeComponent *RecipeComponent) error {
	stmt := `
	INSERT INTO recipe_components(recipe_id, pantry_item_id, quantity, step_no, step_description)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at
	`

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	args := []any{
		&recipeComponent.RecipeID,
		&recipeComponent.PantryItemID,
		&recipeComponent.Quantity,
		&recipeComponent.StepNo,
		&recipeComponent.StepDescription,
	}

	err := m.DB.QueryRow(ctx, stmt, args...).Scan(
		&recipeComponent.ID,
		&recipeComponent.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m RecipeComponentModel) Update(recipeComponent *RecipeComponent) error {
	//components should be functionally immutable except for description which we will allow to change
	stmt := `
	UPDATE recipe_components
	SET step_description = $2
	WHERE id = $1
	`

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	result, err := m.DB.Exec(ctx, stmt, recipeComponent.ID, recipeComponent.StepDescription)
	if err != nil {
		return err
	}

	rows := result.RowsAffected()

	if rows == 0 {
		return ErrRecordNotFound
	}

	return nil
}
