package data

import (
	"context"
	"database/sql"
	"errors"
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
	Name      string          `json:"name"`
	BrandName string          `json:"brand_name"`
	Size      float64         `json:"size"`
	Units     MeasurementUnit `json:"units"`
	Macros    Macronutrients  `json:"macros"`
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
	GetByName(string) (*Consumable, error)
	GetByCreatorID(int64) ([]*Consumable, error)
	GetAllByBrandName(string) ([]*Consumable, error)
	SearchByName(string) ([]*Consumable, error)
	SearchByBrandName(string) ([]*Consumable, error)
	Insert(*Consumable) error
	Update(*Consumable) error
	Delete(int64) error
}

func (m *ConsumableModel) GetByID(ID int64) (*Consumable, error) {
	stmt := `SELECT id, name, brand_name, size, units, carbs, fats, proteins, alcohol
	FROM consumables
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var consumable Consumable

	err := m.DB.QueryRowContext(ctx, stmt, ID).Scan(
		&consumable.ID,
		&consumable.Name,
		&consumable.BrandName,
		&consumable.Size,
		&consumable.Units,
		&consumable.Macros.Carbs,
		&consumable.Macros.Fats,
		&consumable.Macros.Proteins,
		&consumable.Macros.Proteins,
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
