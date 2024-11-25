package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/tconnellan/macro-tracker-backend/internal/data"
	"github.com/tconnellan/macro-tracker-backend/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidatePasswordPlaintext(v, input.Password)
	data.ValidateEmail(v, input.Email)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user := &data.User{
		Username: input.Username,
		Email:    input.Email,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// set up smtp sever before re-enabling
	// app.background(func() {
	// 	err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
	// 	if err != nil {
	// 		app.logger.PrintError(err, nil)
	// 	}
	// })

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		if errors.Is(err, data.ErrInvalidCredentials) {
			app.credentialsInvalid(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.responseSetTokenCookie(w, token.Plaintext)

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
