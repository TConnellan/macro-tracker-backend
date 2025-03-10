package main

import (
	"context"
	"net/http"

	"github.com/tconnellan/macro-tracker-backend/internal/data"
)

type contextKey string

// key for getting and setting user information in the request context.
const userContextKey = contextKey("user")
const bearerTokenContextKey = contextKey("token")

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

func (app *application) getSessionToken(r *http.Request) string {
	cookie, err := r.Cookie(string(bearerTokenContextKey)) // Make sure this name matches what you set
	if err != nil {
		return ""
	}

	// Get token from cookie
	return cookie.Value
}

func (app *application) responseSetTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     string(bearerTokenContextKey),
		Value:    token,
		Path:     "/api",
		MaxAge:   24 * 60 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (app *application) responseClearTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     string(bearerTokenContextKey),
		Value:    "",
		Path:     "/api",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}
