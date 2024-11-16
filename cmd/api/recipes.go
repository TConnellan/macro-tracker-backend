package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tconnellan/macro-tracker-backend/internal/data"
	"github.com/tconnellan/macro-tracker-backend/internal/validator"
)

func (app *application) listRecipes(w http.ResponseWriter, r *http.Request) {

	v := validator.New()

	latest := app.readBool(r.URL.Query(), "latest", false, v)
	page := app.readInt(r.URL.Query(), "page", 1, v)
	pagesize := app.readInt(r.URL.Query(), "pagesize ", 1000, v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	var recipes []*data.Recipe
	filters := data.RecipeFilters{
		Metadata: data.MetadataFilters{
			Page:         page,
			PageSize:     pagesize,
			Sort:         "ID",
			SortSafeList: []string{"ID"},
		},
		NameSearch: "",
	}
	var metadata data.Metadata
	var err error

	if latest {
		recipes, metadata, err = app.models.Recipes.GetLatestByCreatorID(app.contextGetUser(r).ID, filters)
	} else {
		recipes, metadata, err = app.models.Recipes.GetByCreatorID(app.contextGetUser(r).ID, filters)
	}
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"recipes": recipes, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// opportunity for caching here
func (app *application) getRecipe(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	recipeID, err := strconv.Atoi(params.ByName("id"))
	if err != nil || recipeID < 1 {
		app.notFoundResponse(w, r)
		return
	}

	fullRecipes, err := app.models.Recipes.GetFullRecipe(int64(recipeID), app.contextGetUser(r).ID)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.badRequestResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"fullrecipes": fullRecipes}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createChildRecipe(w http.ResponseWriter, r *http.Request) {
	var fullRecipe data.FullRecipe

	params := httprouter.ParamsFromContext(r.Context())
	parentId, err := strconv.Atoi(params.ByName("id"))
	if err != nil || parentId < 1 {
		app.notFoundResponse(w, r)
		return
	}

	err = app.readJSON(w, r, &fullRecipe)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	fullRecipe.Recipe.CreatorID = app.contextGetUser(r).ID
	fullRecipe.Recipe.ParentRecipeID = int64(parentId)

	err = app.models.Recipes.UpdateFullRecipe(&fullRecipe)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrReferencedUserDoesNotExist):
			app.foreignKeyViolationResponse(w, r, err)
		case errors.Is(err, data.ErrRecipeDoesNotExist):
			app.foreignKeyViolationResponse(w, r, err)
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, data.ErrPantryItemDoesNotExist):
			app.foreignKeyViolationResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"fullRecipe": fullRecipe}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createNewRecipe(w http.ResponseWriter, r *http.Request) {
	var fullRecipe data.FullRecipe

	err := app.readJSON(w, r, &fullRecipe)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	fullRecipe.Recipe.CreatorID = app.contextGetUser(r).ID
	fullRecipe.Recipe.ParentRecipeID = 0

	err = app.models.Recipes.InsertFullRecipe(&fullRecipe)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrReferencedUserDoesNotExist):
			app.foreignKeyViolationResponse(w, r, err)
		case errors.Is(err, data.ErrRecipeDoesNotExist):
			app.foreignKeyViolationResponse(w, r, err)
		case errors.Is(err, data.ErrPantryItemDoesNotExist):
			app.foreignKeyViolationResponse(w, r, err)
		case errors.Is(err, data.ErrParentRecipeDoesNotExist):
			app.foreignKeyViolationResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"fullRecipe": fullRecipe}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateStep(w http.ResponseWriter, r *http.Request) {
	var recipeComponent data.RecipeComponent

	err := app.readJSON(w, r, &recipeComponent)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	err = app.models.RecipeComponents.Update(&recipeComponent, app.contextGetUser(r).ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, data.ErrRecipeDoesNotExist):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"recipecomponent": recipeComponent}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAncestors(w http.ResponseWriter, r *http.Request) {
	v := validator.New()

	params := httprouter.ParamsFromContext(r.Context())
	recipeID, err := strconv.Atoi(params.ByName("id"))
	if err != nil || recipeID < 1 {
		app.notFoundResponse(w, r)
		return
	}

	page := app.readInt(r.URL.Query(), "page", 1, v)
	pagesize := app.readInt(r.URL.Query(), "pagesize ", 1000, v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	filters := data.RecipeFilters{
		Metadata: data.MetadataFilters{
			Page:         page,
			PageSize:     pagesize,
			Sort:         "ID",
			SortSafeList: []string{"ID"},
		},
		NameSearch: "",
	}

	ancestors, metadata, err := app.models.Recipes.GetAllAncestors(&data.Recipe{ID: int64(recipeID)}, filters)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if ancestors[len(ancestors)-1].CreatorID != app.contextGetUser(r).ID {
		app.logger.PrintInfo(fmt.Sprintf("got: %d, wanted %d", ancestors[len(ancestors)-1].CreatorID, app.contextGetUser(r).ID), nil)
		for _, x := range ancestors {
			app.logger.PrintInfo(fmt.Sprintf("%#v", x), nil)
		}
		app.notFoundResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"ancestors": ancestors, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
