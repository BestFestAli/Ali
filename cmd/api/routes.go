package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/scales", app.listFoodScalesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/scales", app.newFoodScalesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/scales/:serverID", app.showFoodScalesHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/scales/:serverID", app.updateFoodScalesHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/scales/:serverID", app.deleteFoodScalesHandler)

	return app.recoverPanic(app.rateLimit(router))

}
