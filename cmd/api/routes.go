package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	standardMiddleware := alice.New(app.recoverPanic, app.rateLimit, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.checkAuthentication) //, noSurf)

	protectedMiddleware := dynamicMiddleware.Append(app.requireUserAuthentication)

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.Handler(http.MethodPost, "/v1/users", dynamicMiddleware.ThenFunc(app.registerUserHandler))
	router.Handler(http.MethodPost, "/v1/users/login", dynamicMiddleware.ThenFunc(app.UserLoginHandler))

	router.Handler(http.MethodGet, "/v1/consumed", protectedMiddleware.ThenFunc(app.getConsumed))
	router.Handler(http.MethodPost, "/v1/consumed", protectedMiddleware.ThenFunc(app.postConsumed))
	router.Handler(http.MethodPut, "/v1/consumed", protectedMiddleware.ThenFunc(app.updateConsumed))
	router.Handler(http.MethodDelete, "/v1/consumed", protectedMiddleware.ThenFunc(app.deleteConsumed))

	router.Handler(http.MethodGet, "/v1/recipes", protectedMiddleware.ThenFunc(app.listRecipes))
	router.Handler(http.MethodGet, "/v1/recipes/:id", protectedMiddleware.ThenFunc(app.getRecipe))
	// post new recipe, creates child recipe if it already exists
	router.Handler(http.MethodPost, "/v1/recipes", protectedMiddleware.ThenFunc(app.createNewRecipe))
	router.Handler(http.MethodPost, "/v1/recipes/:id", protectedMiddleware.ThenFunc(app.createChildRecipe))
	// allow update of modifiable parts of step
	router.Handler(http.MethodPut, "/v1/recipes/:id/step", protectedMiddleware.ThenFunc(app.updateStep))
	router.Handler(http.MethodGet, "/v1/recipes/:id/ancestors", protectedMiddleware.ThenFunc(app.getAncestors))

	// return app.recoverPanic(app.rateLimit(app.checkAuthentication(router)))
	return standardMiddleware.Then(router)
}
