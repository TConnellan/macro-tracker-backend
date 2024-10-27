package data

import (
	"errors"
	"fmt"
	"testing"

	"github.com/tconnellan/macro-tracker-backend/internal/assert"
)

func TestUserModelExists(t *testing.T) {

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "users")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}

			m := UserModel{db}

			exists, err := m.Exists(tt.userID)
			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}

func TestUserModelInsert(t *testing.T) {

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name      string
		user      User
		pass      string
		wantError error
	}{
		{
			name: "Valid Insert",
			user: User{
				Username: "test1",
				Password: password{},
				Email:    "test1@gmail.com",
			},
			pass:      "Pass1",
			wantError: nil,
		},
		{
			name: "Email already exists",
			user: User{
				Username: "test2",
				Password: password{},
				Email:    "John@email.com",
			},
			pass:      "pass2",
			wantError: ErrDuplicateEmail,
		},
		{
			name: "Non valid email",
			user: User{
				Username: "test3",
				Password: password{},
				Email:    "notanemail",
			},
			pass:      "pass3",
			wantError: errors.New("error"),
		},
	}

	for _, tst := range tests {
		err := tst.user.Password.Set(tst.pass)
		if err != nil {
			assert.FailWithError(t, err)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "users")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}

			m := UserModel{db}

			err = m.Insert(&tt.user)

			assert.ExpectError(t, err, tt.wantError)
		})
	}
}
