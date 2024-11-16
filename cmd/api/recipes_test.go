package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/tconnellan/macro-tracker-backend/internal/assert"
	"github.com/tconnellan/macro-tracker-backend/internal/data"
	"github.com/tconnellan/macro-tracker-backend/internal/jsonlog"
)

func TestListRecipes(t *testing.T) {

	tests := []struct {
		Name       string
		StatusCode int
		User       *data.User
		ExpectBody string
	}{
		{
			Name:       "valid",
			StatusCode: http.StatusOK,
			User: &data.User{
				ID:       1,
				Username: "test1",
				Email:    "test1@gmail.com",
			},
			ExpectBody: `{
	"metadata": {
		"current_page": 1,
		"page_size": 1000,
		"first_page": 1,
		"last_page": 1,
		"total_records": 1
	},
	"recipes": [
		{
			"id": 1,
			"recipe_name": "Lasagne",
			"creator_id": 1,
			"created_at": "2024-01-01T10:00:00Z",
			"last_edited_at": "2024-01-01T10:00:00Z",
			"notes": "a recipe",
			"parent_recipe_id": 0,
			"is_latest": true
		}
	]
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			db, err := data.NewETETestDB(t, "recipe")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}

			app := &application{
				logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
				// models: mocks.NewTestModel(),
				models: data.NewModel(db),
			}

			ts := httptest.NewTLSServer(app.routes())
			defer ts.Close()

			ctx := app.testContextSetUser(context.Background(), tt.User)

			request := httptest.NewRequestWithContext(ctx, "GET", ts.URL+"/v1/recipes", nil)

			rr := httptest.NewRecorder()

			app.listRecipes(rr, request)

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
