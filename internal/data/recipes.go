package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/tconnellan/macro-tracker-backend/internal/validator"
)

type Recipe struct {
	ID           int64     `json:"id"`
	Name         string    `json:"recipe_name"`
	CreatorID    int64     `json:"creator_id"`
	CreatedAt    time.Time `json:"created_at"`
	LastEditedAt time.Time `json:"last_edited_at"`
	Notes        string    `json:"notes"`
}

func ValidateRecipe(v *validator.Validator, recipe *Recipe) {
	v.Check(recipe.Name != "", "recipe_name", "Cannot be empty")
	v.Check(len(recipe.Name) <= 50, "recipe_name", "Must be at most 50 characters")
}

type FullRecipe struct {
	Recipe           Recipe
	RecipeComponents []*RecipeComponent
	Consumables      []*Consumable
}

func ValidateComponentConsumableList(v *validator.Validator, recipeID int64, recipeComponents []*RecipeComponent, consumables []*Consumable) {
	// same length, more than zero
	v.Check(len(recipeComponents) == len(consumables), "recipe_steps", "must have same number of components as steps")
	v.Check(len(recipeComponents) > 0, "recipe_steps", "must have at least one step")
	v.Check(len(consumables) > 0, "recipe_steps", "must have at least one step")

	if len(recipeComponents) == len(consumables) {
		n := int64(len(recipeComponents))
		i := int64(0)

		for i < n {
			v.Check(recipeComponents[i].ConsumableID == consumables[i].ID, "recipe_steps", fmt.Sprintf("consumable ids of step %d must match", recipeComponents[i].ConsumableID))
			v.Check(recipeComponents[i].RecipeID == recipeID, "recipe_id", "must be the same for all steps")
		}

		j := int64(0)

		// sort recipeComponents based on StepNo, mirroring the alterations
		// in consumables. Check that step numbers are valid
		for j < n {
			for recipeComponents[j].StepNo != j+1 {

				if recipeComponents[j].StepNo < 1 || recipeComponents[j].StepNo > n {
					v.Check(false, "step_numbers", "must be in range of 1..num_steps")
					break
				}

				actualIndex := recipeComponents[j].StepNo - 1

				if recipeComponents[j].StepNo == recipeComponents[actualIndex].StepNo {
					v.Check(false, "step_numbers", "must be unique")
					break
				}

				recipeComponents[j], recipeComponents[actualIndex] = recipeComponents[actualIndex], recipeComponents[j]
				consumables[j], consumables[actualIndex] = consumables[actualIndex], consumables[j]
			}

			j += 1
		}
	}
}

func ValidateFullRecipe(v *validator.Validator, fullRecipe *FullRecipe) {
	ValidateRecipe(v, &fullRecipe.Recipe)
	ValidateComponentConsumableList(v, fullRecipe.Recipe.ID, fullRecipe.RecipeComponents, fullRecipe.Consumables)
	for _, recipeComponent := range fullRecipe.RecipeComponents {
		ValidateRecipeComponent(v, recipeComponent)
	}
	for _, consumable := range fullRecipe.Consumables {
		ValidateConsumable(v, consumable)
	}
}

type RecipeFilters struct {
	Metadata   MetadataFilters
	NameSearch string
}

type IRecipeModel interface {
	Get(int64) (*Recipe, error)
	GetByCreatorID(int64, RecipeFilters) ([]*Recipe, Metadata, error)
	// GetFullRecipe(int64) (*FullRecipe, error)
	Insert(*Recipe) error
	// InsertFullRecipe(*FullRecipe) error
	Update(*Recipe) error
	// UpdateFullRecipe(*FullRecipe) error
	Delete(int64) error
}

type RecipeModel struct {
	DB *sql.DB
}

func (m RecipeModel) Get(ID int64) (*Recipe, error) {
	stmt := `
	SELECT id, recipe_name, creator_id, created_at, last_edited_at, notes
	FROM recipes
	WHERE id = $1
	`

	var recipe Recipe

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, stmt, ID).Scan(
		&recipe.ID,
		&recipe.Name,
		&recipe.CreatorID,
		&recipe.CreatedAt,
		&recipe.LastEditedAt,
		&recipe.Notes,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &recipe, nil
}

func (m RecipeModel) GetByCreatorID(ID int64, filters RecipeFilters) ([]*Recipe, Metadata, error) {
	stmt := fmt.Sprintf(`
	SELECT COUNT(*) OVER(), id, recipe_name, creator_id, created_at, last_edited_at, notes
	FROM recipes
	WHERE ID = $1
	  AND $2 = "" or recipe_name LIKE $2
	ORDER BY %s %s, id ASC
	LIMIT $3
	OFFSET $4
	`, filters.Metadata.sortColumn(), filters.Metadata.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, stmt, ID, filters.NameSearch, filters.Metadata.pageLimit(), filters.Metadata.pageOffset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	var recordCount int = 0
	var recipes []*Recipe

	for rows.Next() {
		var recipe Recipe
		err = rows.Scan(
			&recordCount,
			&recipe.ID,
			&recipe.Name,
			&recipe.CreatorID,
			&recipe.CreatedAt,
			&recipe.LastEditedAt,
			&recipe.Notes,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		recipes = append(recipes, &recipe)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	return recipes, calculateMetadata(recordCount, filters.Metadata.Page, filters.Metadata.PageSize), nil
}

func (m RecipeModel) GetFullRecipe(ID int64) (*FullRecipe, error) {
	// join recipe on componets first, then join components on consumables
	stmtRecipe := `
	SELECT id, recipe_name, creator_id, created_at, last_edited_at, notes
	FROM recipes
	WHERE id = $1
	`

	stmtComponents := `
	SELECT RC.id, RC.recipe_id, RC.consumable_id, RC.created-at, RC.quantity, RC.step_no, RC.step_description, C.id, C.creator_id, C.created_at, C.name, C.brand_name, C.size, C.units, C.carbs, c.fats, C.proteins, C.alcohol
	FROM recipe_components RC INNER JOIN consumables C ON RC.consumable_id = C.id
	WHERE RC.recipe_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var recipe Recipe

	if err := m.DB.QueryRowContext(ctx, stmtRecipe, ID).Scan(
		&recipe.ID,
		&recipe.Name,
		&recipe.CreatorID,
		&recipe.CreatedAt,
		&recipe.LastEditedAt,
		&recipe.Notes,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	rows, err := m.DB.QueryContext(ctx, stmtComponents, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	components := []*RecipeComponent{}
	consumables := []*Consumable{}

	for rows.Next() {
		var component RecipeComponent
		var consumable Consumable
		err = rows.Scan(
			&component.ID,
			&component.RecipeID,
			&component.ConsumableID,
			&component.CreatedAt,
			&component.Quantity,
			&component.StepNo,
			&component.StepDescription,
			&consumable.ID,
			&consumable.CreatorID,
			&consumable.CreatedAt,
			&consumable.Name,
			&consumable.BrandName,
			&consumable.Size,
			&consumable.Units,
			&consumable.Macros.Carbs,
			&consumable.Macros.Fats,
			&consumable.Macros.Proteins,
			&consumable.Macros.Alcohol,
		)
		if err != nil {
			return nil, err
		}
		components = append(components, &component)
		consumables = append(consumables, &consumable)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &FullRecipe{Recipe: recipe, RecipeComponents: components, Consumables: consumables}, nil

}

func (m RecipeModel) Insert(recipe *Recipe) error {
	stmt := `
	INSERT INTO recipes (recipe_name, creator_id, notes)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, last_edited_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, stmt, recipe.Name, recipe.CreatorID, recipe.Notes).Scan(&recipe.ID, &recipe.CreatedAt, &recipe.LastEditedAt)

	if err != nil {
		return err
	}

	return nil
}

func (m RecipeModel) Update(recipe *Recipe) error {
	stmt := `
	UPDATE recipes
	SET recipe_name = $2, last_edited_at = current_timestamp, notes = $3
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, stmt, recipe.ID, recipe.Name, recipe.Notes)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m RecipeModel) Delete(ID int64) error {

	stmtRecipeComponent := `
	DELETE FROM recipe_components 
	WHERE recipe_id = $1
	`

	stmtRecipe := `
	DELETE FROM recipes
	WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	txn, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer txn.Rollback()

	_, err = txn.ExecContext(ctx, stmtRecipeComponent, ID)
	if err != nil {
		return err
	}

	result, err := txn.ExecContext(ctx, stmtRecipe, ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrRecordNotFound
	}

	txn.Commit()

	return nil
}
