package data

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
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

// func NewModel(db *sql.DB) Models {
// 	return Models{
// 		Users:       UserModel{DB: db},
// 		Consumed:    ConsumedModel{DB: db},
// 		Consumables: ConsumableModel{DB: db},
// 		Recipes:     RecipeModel{DB: db},
// 	}
// }

func NewModel(db *pgxpool.Pool) Models {
	return Models{
		Users:       UserModel{DB: db},
		Consumed:    ConsumedModel{DB: db},
		Consumables: ConsumableModel{DB: db},
		Recipes:     RecipeModel{DB: db},
	}
}
