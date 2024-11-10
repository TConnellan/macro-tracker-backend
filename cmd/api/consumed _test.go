package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/tconnellan/macro-tracker-backend/internal/assert"
	"github.com/tconnellan/macro-tracker-backend/internal/data"
	"github.com/tconnellan/macro-tracker-backend/internal/data/mocks"
	"github.com/tconnellan/macro-tracker-backend/internal/jsonlog"
)

func TestGetAllConsumed(t *testing.T) {

	tests := []struct {
		Name       string
		ID         int64
		StatusCode int
		User       data.User
		ExpectBody string
	}{
		{
			Name:       "valid",
			ID:         1,
			StatusCode: http.StatusOK,
			User: data.User{
				ID:       1,
				Username: "test1",
				Email:    "test1@gmail.com",
			},
			ExpectBody: `{
	"consumed": [
		{
			"id": 1,
			"user_id": 1,
			"recipe_id": 1,
			"quantity": 1,
			"macros": {
				"carbs": 1,
				"fats": 1,
				"proteins": 1,
				"alcohol": 1
			},
			"consumed_at": "2024-01-01T10:00:00Z",
			"created_at": "2024-01-01T10:00:00Z",
			"last_edited_at": "2024-01-01T10:00:00Z",
			"notes": ""
		}
	]
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			app := &application{
				logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
				models: mocks.NewTestModel(),
			}

			ts := httptest.NewTLSServer(app.routes())
			defer ts.Close()

			ctx := app.testContextSetUser(context.Background(), &tt.User)

			request := httptest.NewRequestWithContext(ctx, "GET", ts.URL+"/v1/consumed", nil)

			rr := httptest.NewRecorder()

			app.getAllConsumed(rr, request)

			rs := rr.Result()
			defer rs.Body.Close()

			body, err := io.ReadAll(rs.Body)
			if err != nil {
				t.Fatal(err)
			}
			bytes.TrimSpace(body)
			assert.Equal(t, rs.StatusCode, tt.StatusCode)
			assert.Equal(t, string(body), tt.ExpectBody)
		})
	}
}
