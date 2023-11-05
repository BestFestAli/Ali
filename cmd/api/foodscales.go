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
		Model       string       `json:"model" `
		SpecialCode int64        `json:"code"`
		Price       float32      `json:"price"`
		Year        int32        `json:"year" `
		Dimensions  []float32    `json:"dimensions" `
		Runtime     data.Runtime `json:"runtime" `
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
	err = app.models.Foodscales.Insert(foodscale)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/scales/%d", foodscale.ServerID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"foodscale": foodscale}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showFoodScalesHandler(w http.ResponseWriter, r *http.Request) {
	ServerID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	foodscales, err := app.models.Foodscales.Get(ServerID)
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
	ServerID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	foodscales, err := app.models.Foodscales.Get(ServerID)
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
		Model      string       `json:"model" `
		Year       int32        `json:"year" `
		Runtime    data.Runtime `json:"runtime" `
		Dimensions []float32    `json:"dimensions" `
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	foodscales.Model = input.Model
	foodscales.Year = input.Year
	foodscales.Runtime = input.Runtime
	foodscales.Dimensions = input.Dimensions

	v := validator.New()
	if data.ValidateFoodScales(v, foodscales); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Foodscales.Update(foodscales)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"foodscales": foodscales}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteFoodScalesHandler(w http.ResponseWriter, r *http.Request) {

	serverID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Delete the movie from the database, sending a 404 Not Found response to the
	// client if there isn't a matching record.
	err = app.models.Foodscales.Delete(serverID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return a 200 OK status code along with a success message.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "foodscales successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
