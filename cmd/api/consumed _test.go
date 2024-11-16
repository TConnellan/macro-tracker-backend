package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
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
		User       *data.User
		ExpectBody string
	}{
		{
			Name:       "valid",
			ID:         1,
			StatusCode: http.StatusOK,
			User: &data.User{
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

			ctx := app.testContextSetUser(context.Background(), tt.User)

			request := httptest.NewRequestWithContext(ctx, "GET", ts.URL+"/v1/consumed", nil)

			rr := httptest.NewRecorder()

			app.getConsumed(rr, request)

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

func TestGetAllConsumedByDate(t *testing.T) {

	tests := []struct {
		Name       string
		ID         int64
		StatusCode int
		Start      string
		End        string
		User       *data.User
		ExpectBody string
	}{
		{
			Name:       "valid",
			ID:         1,
			StatusCode: http.StatusOK,
			Start:      "2024-01-01T10:00:00+10:00",
			End:        "2024-01-01T10:00:00+10:00",
			User: &data.User{
				ID:       1,
				Username: "test1",
				Email:    "test1@gmail.com",
			},
			ExpectBody: `{
[
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
		{
			Name:       "invalid",
			ID:         99999,
			StatusCode: http.StatusOK,
			Start:      "2024-01-01T10:00:00+10:00",
			End:        "2024-01-01T10:00:00+10:00",
			User: &data.User{
				ID:       99999,
				Username: "test1",
				Email:    "test1@gmail.com",
			},
			ExpectBody: `{

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

			ctx := app.testContextSetUser(context.Background(), tt.User)

			request := httptest.NewRequestWithContext(ctx, "GET", ts.URL+"/v1/consumed?start="+url.QueryEscape(tt.Start)+"&end="+url.QueryEscape(tt.End), nil)

			rr := httptest.NewRecorder()

			app.getConsumed(rr, request)

			rs := rr.Result()
			defer rs.Body.Close()

			body, err := io.ReadAll(rs.Body)
			if err != nil {
				t.Fatal(err)
			}
			bytes.TrimSpace(body)
			assert.Equal(t, rs.StatusCode, tt.StatusCode)

			var expectConsumed data.Consumed
			json.Unmarshal([]byte(tt.ExpectBody), &expectConsumed)
			var actualConsumed data.Consumed
			json.Unmarshal([]byte(body), &actualConsumed)
			assert.Equal(t, actualConsumed, expectConsumed)
		})
	}
}

func TestPostConsumed(t *testing.T) {

	tests := []struct {
		Name       string
		StatusCode int
		User       *data.User
		Body       string
	}{
		{
			Name:       "valid",
			StatusCode: http.StatusCreated,
			User: &data.User{
				ID:       1,
				Username: "test1",
				Email:    "test1@gmail.com",
			},
			Body: `{
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

			ctx := app.testContextSetUser(context.Background(), tt.User)

			request := httptest.NewRequestWithContext(ctx, "POST", ts.URL+"/v1/consumed", strings.NewReader(tt.Body))

			rr := httptest.NewRecorder()

			app.postConsumed(rr, request)

			rs := rr.Result()
			defer rs.Body.Close()

			body, err := io.ReadAll(rs.Body)
			if err != nil {
				t.Fatal(err)
			}
			bytes.TrimSpace(body)
			assert.Equal(t, rs.StatusCode, tt.StatusCode)
		})
	}
}

// func TestPutConsumed(t *testing.T)
