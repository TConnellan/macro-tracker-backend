package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/tconnellan/macro-tracker-backend/internal/validator"
)

type Consumed struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	RecipeID     int64     `json:"recipe_id"`
	Quantity     float64   `json:"quantity"`
	Carbs        float64   `json:"carbs"`
	Fats         float64   `json:"fats"`
	Proteins     float64   `json:"proteins"`
	Alcohol      float64   `json:"alcohol"`
	ConsumedAt   time.Time `json:"consumed_at"`
	CreatedAt    time.Time `json:"created_at"`
	LastEditedAt time.Time `json:"last_edited_at"`
	Notes        string    `json:"notes"`
}

func ValidateMacroNutrients(v *validator.Validator, consumed *Consumed) {
	v.Check(consumed.Carbs >= 0, "carbs", "must be non-negative")
	v.Check(consumed.Fats >= 0, "fats", "must be non-negative")
	v.Check(consumed.Proteins >= 0, "proteins", "must be non-negative")
	v.Check(consumed.Alcohol >= 0, "alcohol", "must be non-negative")
	v.Check(consumed.Carbs+consumed.Fats+consumed.Proteins+consumed.Alcohol > 0, "macros", "one macronutrient must be non-negative")
}

func ValidateConsumed(v *validator.Validator, consumed *Consumed) {
	v.Check(consumed.Quantity > 0, "quantity", "quantity must be positive")
	ValidateMacroNutrients(v, consumed)
}

type ConsumedModel struct {
	DB *sql.DB
}

type ConsumerModelInterface interface {
	GetByConsumedID(int64) (*Consumed, error)
	GetAllByUserID(int64) ([]*Consumed, error)
	GetAllByUserIDAndDate(int64, time.Time, time.Time) ([]*Consumed, error)
	Insert(*Consumed) error
	Update(*Consumed) error
}

func (m ConsumedModel) GetByConsumedID(ConsumedID int64) (*Consumed, error) {
	stmt := `SELECT id, user_id, recipe_id, quantity, carbs, fats, proteins, alcohol, consumed_at, created_at, last_edited_at, notes
	FROM consumed
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	consumed := &Consumed{}

	err := m.DB.QueryRowContext(ctx, stmt, ConsumedID).Scan(
		&consumed.ID,
		&consumed.UserID,
		&consumed.RecipeID,
		&consumed.Quantity,
		&consumed.Carbs,
		&consumed.Fats,
		&consumed.Proteins,
		&consumed.Alcohol,
		&consumed.ConsumedAt,
		&consumed.CreatedAt,
		&consumed.LastEditedAt,
		&consumed.Notes,
	)
	if err != nil {
		return nil, err
	}

	return consumed, nil
}

func (m ConsumedModel) GetAllByUserID(userID int64) ([]*Consumed, error) {
	stmt := `SELECT id, user_id, recipe_id, quantity, carbs, fats, proteins, alcohol, consumed_at, created_at, last_edited_at, notes
	FROM consumed
	WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}

	var allConsumed []*Consumed

	for rows.Next() {
		var consumed Consumed
		err = rows.Scan(
			&consumed.ID,
			&consumed.UserID,
			&consumed.RecipeID,
			&consumed.Quantity,
			&consumed.Carbs,
			&consumed.Fats,
			&consumed.Proteins,
			&consumed.Alcohol,
			&consumed.ConsumedAt,
			&consumed.CreatedAt,
			&consumed.LastEditedAt,
			&consumed.Notes,
		)
		if err != nil {
			break
		}
		allConsumed = append(allConsumed, &consumed)
	}

	if closeErr := rows.Close(); closeErr != nil {
		return nil, closeErr
	}

	if err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return allConsumed, nil
}

func (m ConsumedModel) GetAllByUserIDAndDate(userID int64, from time.Time, to time.Time) ([]*Consumed, error) {
	stmt := `SELECT id, user_id, recipe_id, quantity, carbs, fats, proteins, alcohol, consumed_at, created_at, last_edited_at, notes
	FROM consumed
	WHERE user_id = $1 AND consumed_at >= $2 and consumed_at <= $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, stmt, userID, from, to)
	if err != nil {
		return nil, err
	}

	var allConsumed []*Consumed

	for rows.Next() {
		var consumed Consumed
		err = rows.Scan(
			&consumed.ID,
			&consumed.UserID,
			&consumed.RecipeID,
			&consumed.Quantity,
			&consumed.Carbs,
			&consumed.Fats,
			&consumed.Proteins,
			&consumed.Alcohol,
			&consumed.ConsumedAt,
			&consumed.CreatedAt,
			&consumed.LastEditedAt,
			&consumed.Notes,
		)
		if err != nil {
			break
		}
		allConsumed = append(allConsumed, &consumed)
	}

	if closeErr := rows.Close(); closeErr != nil {
		return nil, closeErr
	}

	if err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return allConsumed, nil
}

func (m ConsumedModel) Insert(consumed *Consumed) error {
	stmt := `INSERT INTO consumed (user_id, recipe_id, quantity, carbs, fats, proteins, alcohol, consumed_at, notes)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id, created_at, last_edited_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		consumed.UserID,
		consumed.RecipeID,
		consumed.Quantity,
		consumed.Carbs,
		consumed.Fats,
		consumed.Proteins,
		consumed.Alcohol,
		consumed.ConsumedAt,
		consumed.LastEditedAt,
	}

	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(
		&consumed.ID,
		&consumed.CreatedAt,
		&consumed.LastEditedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m ConsumedModel) Update(consumed *Consumed) error {
	stmt := `UPDATE consumed 
	SET user_id = $1, recipe_id = $2, quantity = $3, carbs = $4, fats = $5, proteins = $6, alcohol = $7, consumed_at = $8, last_edited_at = current_timestamp, notes=$9
	WHERE id = $10`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		consumed.UserID,
		consumed.RecipeID,
		consumed.Quantity,
		consumed.Carbs,
		consumed.Fats,
		consumed.Proteins,
		consumed.Alcohol,
		consumed.ConsumedAt,
		consumed.Notes,
		consumed.ID,
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
