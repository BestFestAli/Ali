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

	foodscales := &data.FoodScales{
		Model:       input.Model,
		SpecialCode: input.SpecialCode,
		Price:       input.Price,
		Year:        input.Year,
		Dimensions:  input.Dimensions,
		Runtime:     input.Runtime,
	}

	v := validator.New()

	if data.ValidateFoodScales(v, foodscales); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
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