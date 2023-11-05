package data

import (
	"awesomeProject3/internal/validator"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"time"
)

type FoodScales struct {
	ServerID    int64     `json:"serverID"`
	Model       string    `json:"model" `
	SpecialCode int64     `json:"code"`
	Price       float32   `json:"price"`
	Year        int32     `json:"year,omitempty" `
	Dimensions  []float32 `json:"dimensions,omitempty" `
	Runtime     Runtime   `json:"runtime,omitempty" `
}

func ValidateFoodScales(v *validator.Validator, foodscales *FoodScales) {
	v.Check(foodscales.Model != "", "brand", "must be provided")
	v.Check(len(foodscales.Model) <= 100, "brand", "must not be more than 100 bytes long")

	v.Check(foodscales.SpecialCode != 0, "code", "must be provided")

	v.Check(foodscales.Year != 0, "year", "must be provided")
	v.Check(foodscales.Year >= 2000, "year", "must be greater than 2000")
	v.Check(foodscales.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(foodscales.Runtime != 0, "runtime", "must be provided")
	v.Check(foodscales.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(foodscales.Dimensions != nil, "dimensions", "must be provided")
	v.Check(len(foodscales.Dimensions) == 3, "genres", "must contain at only 3 numbers for size")

	v.Check(foodscales.Price != 0, "price", "must be provided")
	v.Check(foodscales.Price <= 1000, "price", "must be cheaper than 1000")
}

type FoodScaleModel struct {
	DB *sql.DB
}

func (m FoodScaleModel) Insert(foodscale *FoodScales) error {
	query := `
 		INSERT INTO FoodScales (model, year, runtime, dimensions) 
		VALUES ($1, $2, $3, $4)
 		RETURNING id, code, price `

	args := []interface{}{foodscale.Model, foodscale.Year, foodscale.Runtime, pq.Array(foodscale.Dimensions)}

	return m.DB.QueryRow(query, args...).Scan(&foodscale.ServerID, &foodscale.SpecialCode, &foodscale.Price)

}

func (m FoodScaleModel) Get(id int64) (*FoodScales, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
 		SELECT id, code, model, year, runtime, dimensions, price
 		FROM FoodScales
 		WHERE id = $1 `

	var foodscales FoodScales

	err := m.DB.QueryRow(query, id).Scan(
		&foodscales.ServerID,
		&foodscales.SpecialCode,
		&foodscales.Model,
		&foodscales.Year,
		&foodscales.Runtime,
		pq.Array(&foodscales.Dimensions),
		&foodscales.Price,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Otherwise, return a pointer to the Movie struct.
	return &foodscales, nil

}

func (m FoodScaleModel) Update(foodscales *FoodScales) error {
	return nil
}

func (m FoodScaleModel) Delete(id int64) error {
	return nil
}
