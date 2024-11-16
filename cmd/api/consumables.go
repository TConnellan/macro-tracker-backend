package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tconnellan/macro-tracker-backend/internal/data"
	"github.com/tconnellan/macro-tracker-backend/internal/validator"
)

func (app *application) getConsumable(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	consumableID, err := strconv.Atoi(params.ByName("id"))
	if err != nil || consumableID < 1 {
		app.notFoundResponse(w, r)
		return
	}

	consumable, err := app.models.Consumables.GetByID(int64(consumableID))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"consumable": consumable}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getUserConsumables(w http.ResponseWriter, r *http.Request) {

	v := validator.New()

	page := app.readInt(r.URL.Query(), "page", 1, v)
	pagesize := app.readInt(r.URL.Query(), "pagesize ", 1000, v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	filters := data.ConsumableFilters{
		Metadata: data.MetadataFilters{
			Page:         page,
			PageSize:     pagesize,
			Sort:         "ID",
			SortSafeList: []string{"ID"},
		},
		NameSearch:                   "",
		BrandNameSearch:              "",
		RequireNameAndBrandNameMatch: false,
	}

	consumables, metadata, err := app.models.Consumables.GetByCreatorID(app.contextGetUser(r).ID, filters)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"consumables": consumables, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) searchConsumables(w http.ResponseWriter, r *http.Request) {

	v := validator.New()

	page := app.readInt(r.URL.Query(), "page", 1, v)
	pagesize := app.readInt(r.URL.Query(), "pagesize", 1000, v)
	name := app.readString(r.URL.Query(), "name", "")
	brandName := app.readString(r.URL.Query(), "name", "")
	bothMatch := app.readBool(r.URL.Query(), "bothmatch", false, v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	filters := data.ConsumableFilters{
		Metadata: data.MetadataFilters{
			Page:         page,
			PageSize:     pagesize,
			Sort:         "ID",
			SortSafeList: []string{"ID"},
		},
		NameSearch:                   name,
		BrandNameSearch:              brandName,
		RequireNameAndBrandNameMatch: bothMatch,
	}

	consumables, metadata, err := app.models.Consumables.Search(filters)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"consumables": consumables, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createConsumable(w http.ResponseWriter, r *http.Request) {
	var consumable data.Consumable
	err := app.readJSON(w, r, consumable)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	consumable.CreatorID = app.contextGetUser(r).ID

	v := validator.New()
	data.ValidateConsumable(v, &consumable)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Consumables.Insert(&consumable)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrReferencedUserDoesNotExist):
			app.foreignKeyViolationResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"consumable": consumable}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateConsumable(w http.ResponseWriter, r *http.Request) {
	var consumable data.Consumable
	err := app.readJSON(w, r, consumable)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateConsumable(v, &consumable)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Consumables.Update(&consumable)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrReferencedUserDoesNotExist):
			app.foreignKeyViolationResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"consumable": consumable}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
