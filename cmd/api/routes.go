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

	router.HandlerFunc(http.MethodGet, "/v1/scales", app.requirePermission("scales:read", app.listFoodScalesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/scales", app.requirePermission("scales:write", app.newFoodScalesHandler))
	router.HandlerFunc(http.MethodGet, "/v1/scales/:id", app.requirePermission("scales:read", app.showFoodScalesHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/scales/:id", app.requirePermission("scales:write", app.updateFoodScalesHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/scales/:serverID", app.requirePermission("scales:write", app.deleteFoodScalesHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))

}
