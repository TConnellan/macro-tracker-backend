package data

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tconnellan/macro-tracker-backend/internal/validator"
)

type PantryItem struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	ConsumableId int64     `json:"consumable_id"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	LastEditedAt time.Time `json:"last_edited_at"`
}

func ValidatePantryItem(v *validator.Validator, pantryItem *PantryItem) {
	v.Check(pantryItem.Name != "", "pantry_item_name", "must not be empty")
	v.Check(len(pantryItem.Name) <= 50, "pantry_item_name", "must be less than or equal to 50 characters")
}

type PantryItemModel struct {
	DB *pgxpool.Pool
}

type IPantryItemModel interface {
	GetAllByUserID(int64) ([]*PantryItem, error)
	Get(int64) (*PantryItem, error)
	Create(*PantryItem) error
	Update(*PantryItem) error
	Delete(int64, int64) error
}

func (m PantryItemModel) GetAllByUserID(userID int64) ([]*PantryItem, error) {

	stmt := `
	SELECT id, user_id, consumable_id, name, created_at, last_edited_at
	FROM pantry_items
	WHERE user_id = $1
	`

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	rows, err := m.DB.Query(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pantryItems := []*PantryItem{}

	var pantryItem PantryItem
	for rows.Next() {
		pantryItem = PantryItem{}

		err = rows.Scan(
			&pantryItem.ID,
			&pantryItem.UserID,
			&pantryItem.ConsumableId,
			&pantryItem.Name,
			&pantryItem.CreatedAt,
			&pantryItem.LastEditedAt,
		)

		if err != nil {
			return nil, err
		}
		pantryItems = append(pantryItems, &pantryItem)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pantryItems, nil

}

func (m PantryItemModel) Get(ID int64) (*PantryItem, error) {
	return get(ID, m.DB)
}

func get(ID int64, db psqlDB) (*PantryItem, error) {

	stmt := `
	SELECT id, user_id, consumable_id, name, created_at, last_edited_at
	FROM pantry_items
	WHERE id = $1
	`

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	var pantryItem PantryItem

	err := db.QueryRow(ctx, stmt, ID).Scan(
		&pantryItem.ID,
		&pantryItem.UserID,
		&pantryItem.ConsumableId,
		&pantryItem.Name,
		&pantryItem.CreatedAt,
		&pantryItem.LastEditedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &pantryItem, nil

}

func (m PantryItemModel) Create(pantryItem *PantryItem) error {
	return create(pantryItem, m.DB)
}

func create(pantryItem *PantryItem, db psqlDB) error {
	stmt := `
	INSERT INTO pantry_items(user_id, consumable_id, name)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, last_edited_at
	`

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	err := db.QueryRow(ctx, stmt, pantryItem.UserID, pantryItem.ConsumableId, pantryItem.Name).Scan(
		&pantryItem.ID,
		&pantryItem.CreatedAt,
		&pantryItem.LastEditedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m PantryItemModel) Update(pantryItem *PantryItem) error {
	return updatepantryItem(pantryItem, m.DB)
}

func updatepantryItem(pantryItem *PantryItem, db psqlDB) error {
	stmt := `
	UPDATE TABLE pantry_items
	SET name = $2, consumable_id = $3
	WHERE id = $1
	RETURING id, user_id, consumable_id, name, created_at, last_edited_at
	`

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	err := db.QueryRow(ctx, stmt, pantryItem.ID, pantryItem.Name, pantryItem.ConsumableId).Scan(
		&pantryItem.ID,
		&pantryItem.UserID,
		&pantryItem.ConsumableId,
		&pantryItem.Name,
		&pantryItem.CreatedAt,
		&pantryItem.LastEditedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (m PantryItemModel) Delete(ID int64, userID int64) error {
	stmt := `
	DELETE FROM pantry_items
	WHERE id = $1 AND user_id = $2
	`

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	result, err := m.DB.Exec(ctx, stmt, ID, userID)
	if err != nil {
		return err
	}
	rows := result.RowsAffected()

	if rows == 0 {
		return ErrRecordNotFound
	}

	return nil
}
