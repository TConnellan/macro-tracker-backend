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

type RecipeFilters struct {
	Metadata   MetadataFilters
	NameSearch string
}

type IRecipeModel interface {
	Get(int64) (*Recipe, error)
	GetByCreatorID(int64, RecipeFilters) ([]*Recipe, Metadata, error)
	Insert(*Recipe) error
	Update(*Recipe) error
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
	stmt := `
	DELETE FROM recipes
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, stmt, ID)
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
