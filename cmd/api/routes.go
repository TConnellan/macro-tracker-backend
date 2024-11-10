package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/login", app.UserLoginHandler)
	router.HandlerFunc(http.MethodGet, "/v1/consumed", app.getAllConsumed)

	return app.recoverPanic(app.rateLimit(app.checkAuthentication(router)))
}
