package main

import (
	"awesomeProject3/internal/data"
	"awesomeProject3/internal/validator"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) newFoodScalesHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Model      string       `json:"model" `
		Version    int64        `json:"version"`
		Price      float32      `json:"price"`
		Year       int32        `json:"year" `
		Dimensions []float32    `json:"dimensions" `
		Runtime    data.Runtime `json:"runtime" `
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	foodscale := &data.FoodScales{
		Model:      input.Model,
		Year:       input.Year,
		Runtime:    input.Runtime,
		Dimensions: input.Dimensions,
	}

	v := validator.New()

	if data.ValidateFoodScales(v, foodscale); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.FoodScales.Insert(foodscale)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/scales/%d", foodscale.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"foodscale": foodscale}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showFoodScalesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	foodscales, err := app.models.FoodScales.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"foodscales": foodscales}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateFoodScalesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	foodscales, err := app.models.FoodScales.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Model      *string       `json:"model" `
		Year       *int32        `json:"year" `
		Runtime    *data.Runtime `json:"runtime" `
		Dimensions []float32     `json:"dimensions" `
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Model != nil {
		foodscales.Model = *input.Model
	}
	if input.Year != nil {
		foodscales.Year = *input.Year
	}
	if input.Runtime != nil {
		foodscales.Runtime = *input.Runtime
	}
	if input.Dimensions != nil {
		foodscales.Dimensions = input.Dimensions
	}

	v := validator.New()
	if data.ValidateFoodScales(v, foodscales); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.FoodScales.Update(foodscales)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"foodscales": foodscales}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteFoodScalesHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.FoodScales.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "foodscales successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listFoodScalesHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Model string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	input.Model = app.readString(qs, "model", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "model", "year", "runtime", "-id", "-model", "-year", "-runtime"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	foodscales, metadata, err := app.models.FoodScales.GetAll(input.Model, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"foodscales": foodscales, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
