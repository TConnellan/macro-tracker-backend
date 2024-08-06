package main

import "net/http"

func (app *application) createConsumedEvent(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserId int64 `json:"user_id"`
	}
}
