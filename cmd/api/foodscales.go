package main

import (
	"awesomeProject3/internal/data"
	"awesomeProject3/internal/validator"
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
	headers.Set("Location", fmt.Sprintf("/v1/foodscales/%d", foodscale.ServerID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"foodscale": foodscale}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showFoodScalesHandler(w http.ResponseWriter, r *http.Request) {
	serverID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	foodscales := &data.FoodScales{
		ServerID:    serverID,
		Model:       "Escali Primo Digital Scale",
		SpecialCode: 2204211300,
		Price:       15,
		Year:        2022,
		Dimensions:  []float32{8.5, 6, 1.5},
		Runtime:     102,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"foodscales": foodscales}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
