package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/tconnellan/macro-tracker-backend/internal/data"
	"github.com/tconnellan/macro-tracker-backend/internal/validator"
)

func (app *application) getConsumed(w http.ResponseWriter, r *http.Request) {

	start := app.readString(r.URL.Query(), "start", "")
	end := app.readString(r.URL.Query(), "end", "")

	var consumed []*data.Consumed
	var err error

	if start == "" && end == "" {
		consumed, err = app.models.Consumed.GetAllByUserID(app.contextGetUser(r).ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	} else {
		startTime, startErr := time.Parse(time.RFC3339, start)
		endTime, endErr := time.Parse(time.RFC3339, end)
		v := validator.New()
		v.Check(startErr == nil, "start", "format must be RFC3339")
		v.Check(endErr == nil, "end", "format must be RFC3339")
		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		consumed, err = app.models.Consumed.GetAllByUserIDAndDate(app.contextGetUser(r).ID, startTime, endTime)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"consumed": consumed}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) postConsumed(w http.ResponseWriter, r *http.Request) {

	var consumed data.Consumed

	err := app.readJSON(w, r, &consumed)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	consumed.UserID = app.contextGetUser(r).ID

	err = app.models.Consumed.Insert(&consumed)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrReferencedUserDoesNotExist):
			app.foreignKeyViolationResponse(w, r, err)
		case errors.Is(err, data.ErrRecipeDoesNotExist):
			app.foreignKeyViolationResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"consumed": consumed}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateConsumed(w http.ResponseWriter, r *http.Request) {

	var consumed data.Consumed

	err := app.readJSON(w, r, &consumed)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	consumed.UserID = app.contextGetUser(r).ID

	err = app.models.Consumed.Update(&consumed)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, data.ErrReferencedUserDoesNotExist):
			app.foreignKeyViolationResponse(w, r, err)
		case errors.Is(err, data.ErrRecipeDoesNotExist):
			app.foreignKeyViolationResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"consumed": consumed}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteConsumed(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Consumed.Delete(id, app.contextGetUser(r).ID)

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
