package data

import (
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
			userID: 99999,
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
			name: "Email already exists case insensitive",
			user: User{
				Username: "test3",
				Password: password{},
				Email:    "john@email.com",
			},
			pass:      "pass3",
			wantError: ErrDuplicateEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Password.Set(tt.pass)
			if err != nil {
				assert.FailWithError(t, err)
			}

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

func TestUserModelGetByEmail(t *testing.T) {

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name      string
		email     string
		ID        int64
		wantError error
	}{
		{
			name:      "email exists",
			email:     "John@email.com",
			ID:        1,
			wantError: nil,
		},
		{
			name:      "email exists case insensitive",
			email:     "john@email.com",
			ID:        1,
			wantError: nil,
		},
		{
			name:      "email doesnt exist",
			email:     "notexists",
			ID:        -1,
			wantError: ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "users")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := UserModel{db}

			foundUser, err := m.GetByEmail(tt.email)

			assert.ExpectError(t, err, tt.wantError)

			if tt.ID != -1 {
				assert.Equal(t, tt.ID, foundUser.ID)
			}
		})
	}
}

func TestUserModelUpdate(t *testing.T) {

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name        string
		email       string
		newUsername string
		newEmail    string
		newPassword string
		wantError   error
	}{
		{
			name:        "update user",
			email:       "John@email.com",
			newUsername: "Jane Doe",
			newEmail:    "Jane@email.com",
			newPassword: "newpass",
			wantError:   nil,
		},
		{
			name:        "update user email conflict",
			email:       "John@email.com",
			newUsername: "Jack Brabham",
			newEmail:    "jack@email.com",
			newPassword: "newpass",
			wantError:   ErrDuplicateEmail,
		},
		{
			name:        "update user allow conflicting usernames",
			email:       "John@email.com",
			newUsername: "Jack Brabham",
			newEmail:    "newjack@email.com",
			newPassword: "newpass",
			wantError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "users")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}
			m := UserModel{db}

			foundUser, err := m.GetByEmail(tt.email)
			if err != nil {
				t.Fatal(err)
			}

			foundUser.Username = tt.newUsername
			foundUser.Email = tt.newEmail
			foundUser.Password.Set(tt.newPassword)

			err = m.Update(foundUser)
			assert.ExpectError(t, err, tt.wantError)
			if err != nil {
				return
			}
			_, err = m.GetByEmail(tt.email)
			assert.ExpectError(t, err, ErrRecordNotFound)

			updatedUser, err := m.GetByEmail(tt.newEmail)
			assert.ExpectError(t, err, nil)
			if err != nil {
				return
			}

			assert.Equal(t, updatedUser.Email, tt.newEmail)
			assert.Equal(t, updatedUser.Username, tt.newUsername)
		})
	}
}
