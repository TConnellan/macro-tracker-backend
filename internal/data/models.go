package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users       IUserModel
	Consumed    IConsumedModel
	Consumables IConsumableModel
	Recipes     IRecipeModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Users:       UserModel{DB: db},
		Consumed:    ConsumedModel{DB: db},
		Consumables: ConsumableModel{DB: db},
		Recipes:     RecipeModel{DB: db},
	}
}
