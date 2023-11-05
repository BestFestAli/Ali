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
	Foodscales FoodScaleModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Foodscales: FoodScaleModel{DB: db},
	}
}
