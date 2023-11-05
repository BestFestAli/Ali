package data

import (
	"awesomeProject3/internal/validator"
	"context"
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&foodscale.ServerID, &foodscale.SpecialCode, &foodscale.Price)

}

func (m FoodScaleModel) Get(serverID int64) (*FoodScales, error) {
	if serverID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
 		SELECT id, code, model, year, runtime, dimensions, price
 		FROM FoodScales
 		WHERE id = $1 `

	var foodscales FoodScales

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, serverID).Scan(
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

	return &foodscales, nil

}

func (m FoodScaleModel) Update(foodscales *FoodScales) error {
	query := `
 		UPDATE movies 
 		SET model = $1, year = $2, runtime = $3, dimensions = $4, specialcode = specialcode + 101
 		WHERE id = $5 AND specialcode = $6
 		RETURNING specialcode `

	args := []interface{}{
		foodscales.Model,
		foodscales.Year,
		foodscales.Runtime,
		pq.Array(foodscales.Dimensions),
		foodscales.ServerID,
		foodscales.SpecialCode,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&foodscales.SpecialCode)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil

}

func (m FoodScaleModel) Delete(serverID int64) error {
	if serverID < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM foodscales
 		WHERE id = $1 `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, serverID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m FoodScaleModel) GetAll(model string, filters Filters) ([]*FoodScales, error) {
	query := `
 		SELECT id, code, model, year, runtime, dimensions, price
 		FROM foodscales
 		WHERE (LOWER(model) = LOWER($1) OR $1 = '')
 		ORDER BY id `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, model)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	foodscales := []*FoodScales{}

	for rows.Next() {
		var foodscale FoodScales
		err := rows.Scan(
			&foodscale.ServerID,
			&foodscale.SpecialCode,
			&foodscale.Model,
			&foodscale.Year,
			&foodscale.Runtime,
			pq.Array(&foodscale.Dimensions),
			&foodscale.Price,
		)
		if err != nil {
			return nil, err
		}

		foodscales = append(foodscales, &foodscale)
	}
	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK, then return the slice of movies.
	return foodscales, nil
}
