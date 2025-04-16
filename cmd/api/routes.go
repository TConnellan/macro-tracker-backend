package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	standardMiddleware := alice.New(app.recoverPanic, app.rateLimit, app.logRequest, secureHeaders, app.enableCORS)

	dynamicMiddleware := alice.New(app.checkAuthentication)

	protectedMiddleware := dynamicMiddleware.Append(app.requireUserAuthentication) // .Append(noSurf)

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.Handler(http.MethodPost, "/api/v1/users", dynamicMiddleware.ThenFunc(app.registerUserHandler))
	router.Handler(http.MethodPost, "/api/v1/users/login", dynamicMiddleware.ThenFunc(app.UserLoginHandler))

	router.Handler(http.MethodGet, "/api/v1/consumed", protectedMiddleware.ThenFunc(app.getConsumed))
	router.Handler(http.MethodPost, "/api/v1/consumed", protectedMiddleware.ThenFunc(app.postConsumed))
	router.Handler(http.MethodPut, "/api/v1/consumed", protectedMiddleware.ThenFunc(app.updateConsumed))
	router.Handler(http.MethodDelete, "/api/v1/consumed/:id", protectedMiddleware.ThenFunc(app.deleteConsumed))
	router.Handler(http.MethodOptions, "/api/v1/consumed", standardMiddleware.Then(app.respondCors(nil)))

	router.Handler(http.MethodGet, "/api/v1/recipes", protectedMiddleware.ThenFunc(app.listRecipes))
	router.Handler(http.MethodGet, "/api/v1/recipes/:id", protectedMiddleware.ThenFunc(app.getRecipe))
	// post new recipe, creates child recipe if it already exists
	router.Handler(http.MethodPost, "/api/v1/recipes", protectedMiddleware.ThenFunc(app.createNewRecipe))
	router.Handler(http.MethodPost, "/api/v1/recipes/:id", protectedMiddleware.ThenFunc(app.createChildRecipe))
	// allow update of modifiable parts of step
	router.Handler(http.MethodPut, "/api/v1/recipes/:id/step", protectedMiddleware.ThenFunc(app.updateStep))
	router.Handler(http.MethodGet, "/api/v1/recipes/:id/ancestors", protectedMiddleware.ThenFunc(app.getAncestors))
	router.Handler(http.MethodOptions, "/api/v1/recipes", standardMiddleware.Then(app.respondCors(nil)))

	// consumables
	router.Handler(http.MethodGet, "/api/v1/consumable/personal", protectedMiddleware.ThenFunc(app.getUserConsumables))
	// router.Handler(http.MethodGet, "/api/v1/consumable/:id", protectedMiddleware.ThenFunc(app.getConsumable))
	router.Handler(http.MethodGet, "/api/v1/consumable/search", protectedMiddleware.ThenFunc(app.searchConsumables))
	router.Handler(http.MethodPost, "/api/v1/consumable", protectedMiddleware.ThenFunc(app.createConsumable))
	router.Handler(http.MethodPut, "/api/v1/consumable/:id", protectedMiddleware.ThenFunc(app.updateConsumable))
	router.Handler(http.MethodOptions, "/api/v1/consumable", standardMiddleware.Then(app.respondCors(nil)))

	// pantry items
	router.Handler(http.MethodGet, "/api/v1/pantryitems", protectedMiddleware.ThenFunc(app.getPantryItems))
	router.Handler(http.MethodPost, "/api/v1/pantryitems", protectedMiddleware.ThenFunc(app.createPantryItem))
	router.Handler(http.MethodPut, "/api/v1/pantryitems", protectedMiddleware.ThenFunc(app.updatePantryItem))
	router.Handler(http.MethodDelete, "/api/v1/pantryitems/:id", protectedMiddleware.ThenFunc(app.deletePantryItem))
	router.Handler(http.MethodOptions, "/api/v1/pantryitems", standardMiddleware.Then(app.respondCors(nil)))

	return standardMiddleware.Then(router)
}
