package data

import (
	"awesomeProject3/internal/validator"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"time"
)

type FoodScales struct {
	Model      string    `json:"model" `
	ID         int64     `json:"id"`
	Price      float32   `json:"price"`
	Year       int32     `json:"year,omitempty" `
	Dimensions []float32 `json:"dimensions,omitempty" `
	Runtime    Runtime   `json:"runtime,omitempty" `
	Version    int32     `json:"version"`
}

func ValidateFoodScales(v *validator.Validator, foodscale *FoodScales) {
	v.Check(foodscale.Model != "", "brand", "must be provided")
	v.Check(len(foodscale.Model) <= 100, "brand", "must not be more than 100 bytes long")

	v.Check(foodscale.Version != 0, "code", "must be provided")

	v.Check(foodscale.Year != 0, "year", "must be provided")
	v.Check(foodscale.Year >= 2000, "year", "must be greater than 2000")
	v.Check(foodscale.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(foodscale.Runtime != 0, "runtime", "must be provided")
	v.Check(foodscale.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(foodscale.Dimensions != nil, "dimensions", "must be provided")
	v.Check(len(foodscale.Dimensions) == 3, "genres", "must contain at only 3 numbers for size")

	v.Check(foodscale.Price != 0, "price", "must be provided")
	v.Check(foodscale.Price <= 1000, "price", "must be cheaper than 1000")
}

type FoodScaleModel struct {
	DB *sql.DB
}

func (m FoodScaleModel) Insert(foodscale *FoodScales) error {
	query := `
 		INSERT INTO "FoodScales" (model, year, runtime, dimensions) 
		VALUES ($1, $2, $3, $4)
 		RETURNING id, price, version`

	args := []interface{}{foodscale.Model, foodscale.Year, foodscale.Runtime, pq.Array(foodscale.Dimensions)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&foodscale.ID, &foodscale.Version, &foodscale.Price)

}

func (m FoodScaleModel) Get(id int64) (*FoodScales, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
 		SELECT id, model, year, runtime, dimensions, price, version
 		FROM "FoodScales"
 		WHERE id = $1 `

	var foodscales FoodScales

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := m.DB.QueryRowContext(ctx, query).Scan(
		&foodscales.ID,
		&foodscales.Model,
		&foodscales.Year,
		&foodscales.Runtime,
		pq.Array(&foodscales.Dimensions),
		&foodscales.Price,
		&foodscales.Version,
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
 		UPDATE "FoodScales" 
 		SET model = $1, year = $2, runtime = $3, dimensions = $4, version = version + 1
 		WHERE id = $5 AND version = $6
 		RETURNING version `

	args := []interface{}{
		foodscales.Model,
		foodscales.Year,
		foodscales.Runtime,
		pq.Array(foodscales.Dimensions),
		foodscales.ID,
		foodscales.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&foodscales.Version)
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

func (m FoodScaleModel) Delete(ID int64) error {
	if ID < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM "FoodScales"
 		WHERE id = $1 `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, ID)
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

func (m FoodScaleModel) GetAll(model string, filters Filters) ([]*FoodScales, Metadata, error) {
	query := fmt.Sprintf(`
 		SELECT count(*) OVER(), id, version, model, year, runtime, dimensions, price
 		FROM "FoodScales"
 		WHERE (to_tsvector('simple', model) @@ plainto_tsquery('simple', $1) OR $1 = '') 
 		ORDER BY %s %s, id ASC
 		LIMIT $3 OFFSET $4 `, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{model, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	foodscales := []*FoodScales{}

	for rows.Next() {
		var foodscale FoodScales
		err := rows.Scan(
			&totalRecords,
			&foodscale.ID,
			&foodscale.Model,
			&foodscale.Year,
			&foodscale.Runtime,
			pq.Array(&foodscale.Dimensions),
			&foodscale.Price,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		foodscales = append(foodscales, &foodscale)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return foodscales, metadata, nil
}
