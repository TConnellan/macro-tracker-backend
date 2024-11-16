package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tconnellan/macro-tracker-backend/internal/data"
	"github.com/tconnellan/macro-tracker-backend/internal/validator"
)

func (app *application) getPantryItems(w http.ResponseWriter, r *http.Request) {

	pantryItems, err := app.models.PantryItems.GetAllByUserID(app.contextGetUser(r).ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"pantryitems": pantryItems}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createPantryItem(w http.ResponseWriter, r *http.Request) {
	var pantryItem data.PantryItem

	err := app.readJSON(w, r, &pantryItem)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	pantryItem.ID = app.contextGetUser(r).ID

	v := validator.New()
	data.ValidatePantryItem(v, &pantryItem)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.PantryItems.Create(&pantryItem)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"pantryitem": pantryItem}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updatePantryItem(w http.ResponseWriter, r *http.Request) {
	var pantryItem data.PantryItem

	err := app.readJSON(w, r, &pantryItem)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	pantryItem.ID = app.contextGetUser(r).ID

	v := validator.New()
	data.ValidatePantryItem(v, &pantryItem)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.PantryItems.Update(&pantryItem)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"pantryitem": pantryItem}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deletePantryItem(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())
	pantryItemID, err := strconv.Atoi(params.ByName("id"))
	if err != nil || pantryItemID < 1 {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.PantryItems.Delete(int64(pantryItemID), app.contextGetUser(r).ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
