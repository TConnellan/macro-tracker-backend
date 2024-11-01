package data

import (
	"fmt"
	"testing"

	"github.com/tconnellan/macro-tracker-backend/internal/assert"
	"github.com/tconnellan/macro-tracker-backend/internal/validator"
)

func TestPasswordHelpers(t *testing.T) {

	tests := []struct {
		name           string
		pass           string
		valid          bool
		expectPassword password
		checkAgainst   string
		matches        bool
		matchesError   error
	}{
		{
			name:         "valid password",
			pass:         "AbCdEfGh123!@",
			valid:        true,
			checkAgainst: "AbCdEfGh123!@",
			matches:      true,
			matchesError: nil,
		},
		{
			name:         "valid password mismatch compare not equal",
			pass:         "AbCdEfGh123!@",
			valid:        true,
			checkAgainst: "something else",
			matches:      false,
			matchesError: nil,
		},
		{
			name:         "valid password mismatch compare upper case",
			pass:         "AbCdEfGh123!@",
			valid:        true,
			checkAgainst: "ABCDEFGH123!@",
			matches:      false,
			matchesError: nil,
		},
		{
			name:         "valid password mismatch compare lower case",
			pass:         "AbCdEfGh123!@",
			valid:        true,
			checkAgainst: "abcdefgh123!@",
			matches:      false,
			matchesError: nil,
		},
		{
			name:  "invalid password empty",
			pass:  "",
			valid: false,
		},
		{
			name:  "invalid password short",
			pass:  "AbCd",
			valid: false,
		},
		{
			name:  "invalid password long",
			pass:  "AbCdEfGh123!@AbCdEfGh123!@AbCdEfGh123!@AbCdEfGh123!@AbCdEfGh123!@AbCdEfGh123!@",
			valid: false,
		},
	}

	// of interest, also possible for bcrypt to return other errors on comparison, if the hash
	// is too short/long/malformatted or the cost is not in a valid range (4 - 31)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			v := validator.New()

			ValidatePasswordPlaintext(v, tt.pass)
			if v.Valid() != tt.valid {
				assert.Equal(t, v.Valid(), tt.valid)
			}
			if !v.Valid() {
				return
			}

			actualPassword := password{}
			err := actualPassword.Set(tt.pass)
			assert.ExpectError(t, err, nil)
			if err != nil {
				return
			}

			assert.Equal(t, *actualPassword.plaintext, tt.pass)

			matches, err := actualPassword.Matches(tt.checkAgainst)
			assert.ExpectError(t, err, tt.matchesError)
			assert.Equal(t, matches, tt.matches)

		})
	}
}

func TestEmailHelpers(t *testing.T) {

	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{
			name:  "valid email",
			email: "john.doe@gmail.com",
			valid: true,
		},
		{
			name:  "valid email 2",
			email: "john.doe-second@gmail.com",
			valid: true,
		},
		{
			name:  "valid email 3",
			email: "john_doe@gmail.com",
			valid: true,
		},
		{
			name:  "valid email 4",
			email: "john.doe@gmail-test.com",
			valid: true,
		},
		{
			name:  "invalid email",
			email: "john.doe",
			valid: false,
		},
		{
			name:  "invalid empty email",
			email: "",
			valid: false,
		},
		{
			name:  "invalid consecutive periods",
			email: "john.doe@gmail.com",
			valid: false,
		},
		{
			name:  "invalid no characters",
			email: ".@gmail.com",
			valid: false,
		},
		{
			name:  "invalid trailing special character 1",
			email: "john.doe-@gmail.com",
			valid: false,
		},
		{
			name:  "invalid trailing special character 2",
			email: "john.doe_@gmail.com",
			valid: false,
		},
		{
			name:  "invalid trailing special character 3",
			email: "john.doe.@gmail.com",
			valid: false,
		},
		{
			name:  "invalid no top level domain",
			email: "john.doe@gmail",
			valid: false,
		},
		{
			name:  "invalid top level domain qualifier two short",
			email: "john.doe@gmail.c",
			valid: false,
		},
		{
			name:  "invalid no email domain",
			email: "john.doe@.com",
			valid: false,
		},
		{
			name:  "invalid email domain",
			email: "john.doe@-gmail.com",
			valid: false,
		},
		{
			name:  "invalid email domain 2",
			email: "john.doe@gmail.com-",
			valid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			v := validator.New()

			ValidateEmail(v, tt.email)
			if v.Valid() != tt.valid {
				assert.ValidatorValid(t, v)
			}

		})
	}
}

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
