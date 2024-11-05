package data

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound             = errors.New("record not found")
	ErrEditConflict               = errors.New("edit conflict")
	ErrReferencedUserDoesNotExist = errors.New("user doesn't exist")
	ErrParentRecipeDoesNotExist   = errors.New("parent recipe doesn't exist")
	ErrPantryItemDoesNotExist     = errors.New("pantry item does not exist")
	ErrChildRecipeExists          = errors.New("child recipe exists")
	ErrRecipeDoesNotExist         = errors.New("recipe does not exists")
)

type Models struct {
	Users            IUserModel
	Consumed         IConsumedModel
	Consumables      IConsumableModel
	Recipes          IRecipeModel
	RecipeComponents IRecipeComponentModel
	PantryItems      IPantryItemModel
}

func NewModel(db *pgxpool.Pool) Models {
	return Models{
		Users:            UserModel{DB: db},
		Consumed:         ConsumedModel{DB: db},
		Consumables:      ConsumableModel{DB: db},
		Recipes:          RecipeModel{DB: db},
		RecipeComponents: RecipeComponentModel{DB: db},
		PantryItems:      PantryItemModel{DB: db},
	}
}

// use as input to helper functions to allow passing of transactions in more complicated actions
type psqlDB interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...any) pgx.Row
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
}
