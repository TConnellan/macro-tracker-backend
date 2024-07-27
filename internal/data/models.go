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
	Users    UserModelInterface
	Consumed ConsumedModelInterface
}

func NewModel(db *sql.DB) Models {
	return Models{
		Users:    UserModel{DB: db},
		Consumed: ConsumedModel{DB: db},
	}
}
