package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/tconnellan/macro-tracker-backend/internal/validator"
)

type MeasurementUnit string

var (
	ValidMeasurementUnits = []MeasurementUnit{
		"g",
		"ml",
	}
)

type Consumable struct {
	ID        int64           `json:"id"`
	CreatorID int64           `json:"creator_id"`
	CreatedAt time.Time       `json:"created_at"`
	Name      string          `json:"name"`
	BrandName string          `json:"brand_name"`
	Size      float64         `json:"size"`
	Units     MeasurementUnit `json:"units"`
	Macros    Macronutrients  `json:"macros"`
}

type ConsumableSearchOptions struct {
	Name                         string
	BrandName                    string
	RequireNameAndBrandNameMatch bool
}

func (options ConsumableSearchOptions) GetKeyWord() string {
	if options.RequireNameAndBrandNameMatch {
		return "AND"
	} else {
		return "OR"
	}
}

func ValidateMeasurementUnit(v *validator.Validator, consumable *Consumable) {
	valid := false
	for _, unit := range ValidMeasurementUnits {
		if consumable.Units == unit {
			valid = true
			break
		}
	}

	v.Check(valid, "units", "must be valid")
}

func ValidateConsumable(v *validator.Validator, consumable *Consumable) {
	v.Check(consumable.Name != "", "name", "must be provided")
	v.Check(len(consumable.Name) <= 50, "name", "must be maximum 50 characters")

	v.Check(consumable.BrandName != "", "brand_name", "must be provided")
	v.Check(len(consumable.BrandName) <= 50, "brand_name", "must be maximum 50 characters")

	v.Check(consumable.Size > 0, "size", "must be positive")

	ValidateMeasurementUnit(v, consumable)

	ValidateMacroNutrients(v, consumable.Macros)
}

type ConsumableModel struct {
	DB *sql.DB
}

type ConsumableModelInterface interface {
	GetByID(int64) (*Consumable, error)
	GetByCreatorID(int64, Filters) ([]*Consumable, Metadata, error)
	Search(ConsumableSearchOptions, Filters) ([]*Consumable, Metadata, error)
	Insert(*Consumable) error
	Update(*Consumable) error
	Delete(int64) error
}

func (m *ConsumableModel) GetByID(ID int64) (*Consumable, error) {
	stmt := `SELECT id, name, creator_id, created_at, brand_name, size, units, carbs, fats, proteins, alcohol
	FROM consumables
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var consumable Consumable

	err := m.DB.QueryRowContext(ctx, stmt, ID).Scan(
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
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}

	return &consumable, nil
}

func (m *ConsumableModel) readConsumableRows(stmt string, ctx context.Context, args ...any) ([]*Consumable, int, error) {

	rows, err := m.DB.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var recordCount int = 0
	var consumables []*Consumable

	for rows.Next() {
		var consumable Consumable
		err = rows.Scan(
			&recordCount,
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
			return nil, 0, err
		}
		consumables = append(consumables, &consumable)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return consumables, recordCount, nil
}

func (m *ConsumableModel) GetByCreatorID(ID int64, filters Filters) ([]*Consumable, Metadata, error) {
	stmt := fmt.Sprintf(`
	SELECT COUNT(*) OVER(), id, creator_id, created_at, name, brand_name, size, units, carbs, fats, proteins, alcohol
	FROM consumables
	WHERE creator_id = $1
	ORDER BY %s %s, id ASC
	LIMIT $2
	OFFSET $3
	`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	consumables, recordCount, err := m.readConsumableRows(stmt, ctx, ID, filters.pageLimit(), filters.pageOffset())
	if err != nil {
		return nil, Metadata{}, err
	}

	return consumables, calculateMetadata(recordCount, filters.Page, filters.PageSize), nil
}

func (m *ConsumableModel) Search(searchOptions ConsumableSearchOptions, filters Filters) ([]*Consumable, Metadata, error) {
	stmt := fmt.Sprintf(`
	SELECT COUNT(*) OVER(), id, creator_id, created_at, name, brand_name, size, units, carbs, fats, proteins, alcohol
	FROM consumables
	WHERE ($1 = '' OR to_tsvector('simple', name) @@ plainto_tsquery('simple', $1))
	   %s ($2 = '' OR to_tsvector('simple', brand_name) @@ plainto_tsquery('simple', $2))
	ORDER BY %s %s, id ASC
	LIMIT $3
	OFFSET $4
	`, searchOptions.GetKeyWord(), filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	consumables, recordCount, err := m.readConsumableRows(stmt, ctx, searchOptions.Name, searchOptions.BrandName, filters.pageLimit(), filters.pageOffset())
	if err != nil {
		return nil, Metadata{}, err
	}

	return consumables, calculateMetadata(recordCount, filters.Page, filters.PageSize), nil
}

func (m *ConsumableModel) Insert(consumable *Consumable) error {
	stmt := `
	INSERT INTO consumables (creator_id, name, brand_name, size, units, carbs, fats, proteins, alcohol)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		consumable.CreatorID,
		consumable.Name,
		consumable.BrandName,
		consumable.Size,
		consumable.Units,
		consumable.Macros.Carbs,
		consumable.Macros.Fats,
		consumable.Macros.Proteins,
		consumable.Macros.Alcohol,
	}

	if err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&consumable.ID, &consumable.CreatedAt); err != nil {
		return err
	}

	return nil
}

func (m *ConsumableModel) Update(consumable *Consumable) error {
	stmt := `
	UPDATE consumables
	SET name = $2, brand_name = $3, size = $4, units = $5, carbs = $6, fats = $7, proteins = $8, alcohol = $9
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		consumable.Name,
		consumable.BrandName,
		consumable.Size,
		consumable.Units,
		consumable.Macros.Carbs,
		consumable.Macros.Fats,
		consumable.Macros.Proteins,
		consumable.Macros.Alcohol,
	}

	result, err := m.DB.ExecContext(ctx, stmt, args...)
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

func (m *ConsumableModel) Delete(ID int64) error {
	stmt := `
	DELETE FROM consumables
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
