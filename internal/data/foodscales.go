package data

import (
	"awesomeProject3/internal/validator"
	"database/sql"
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

func (m FoodScaleModel) Insert(foodscales *FoodScales) error {
	return nil
}

func (m FoodScaleModel) Get(id int64) (*FoodScales, error) {
	return nil, nil
}

func (m FoodScaleModel) Update(foodscales *FoodScales) error {
	return nil
}

func (m FoodScaleModel) Delete(id int64) error {
	return nil
}
