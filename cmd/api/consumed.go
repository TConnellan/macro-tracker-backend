package main

import (
	"net/http"
)

func (app *application) getAllConsumed(w http.ResponseWriter, r *http.Request) {

	consumed, err := app.models.Consumed.GetAllByUserID(app.contextGetUser(r).ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"consumed": consumed}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
