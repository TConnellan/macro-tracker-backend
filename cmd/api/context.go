package main

import (
	"context"
	"net/http"

	"github.com/tconnellan/macro-tracker-backend/internal/data"
)

type contextKey string

// key for getting and setting user information in the request context.
const userContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// exists for tests, need to pass context into httptest.NewRequestWithContext which creates and executes request
func (app *application) testContextSetUser(ctx context.Context, user *data.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
