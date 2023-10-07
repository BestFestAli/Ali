package data

import (
	"awesomeProject3/internal/validator"
	"time"
)

type FoodScales struct {
	Model       string  `json:"model" `
	SpecialCode int64   `json:"code"`
	Price       int32   `json:"price"`
	Year        int32   `json:"year,omitempty" `
	Dimensions  []int32 `json:"dimensions,omitempty" `
	Runtime     Runtime `json:"runtime,omitempty" `
}

func ValidateMovie(v *validator.Validator, food_scales *FoodScales) {
	v.Check(food_scales.Model != "", "brand", "must be provided")
	v.Check(len(food_scales.Model) <= 100, "brand", "must not be more than 100 bytes long")

	v.Check(food_scales.SpecialCode != 0, "code", "must be provided")

	v.Check(food_scales.Year != 0, "year", "must be provided")
	v.Check(food_scales.Year >= 2000, "year", "must be greater than 2000")
	v.Check(food_scales.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(food_scales.Runtime != 0, "runtime", "must be provided")
	v.Check(food_scales.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(food_scales.Dimensions != nil, "dimensions", "must be provided")
	v.Check(len(food_scales.Dimensions) == 3, "genres", "must contain at only 3 numbers for size")

	v.Check(food_scales.Price != 0, "price", "must be provided")
	v.Check(food_scales.Price <= 1000, "price", "must be cheaper than 1000")
}
