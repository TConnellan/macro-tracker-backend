package data

import (
	"fmt"
	"testing"
	"time"

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
			assert.ValidatorValid(t, v, tt.valid)
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

func TestUserHelpers(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	tests := []struct {
		name  string
		pass  string
		user  User
		valid bool
	}{
		{
			name:  "valid user",
			pass:  "password1",
			valid: true,
			user: User{
				ID:        1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Username:  "john doe",
				Email:     "John.Doe@gmail.com",
				Password:  password{},
				Version:   1,
			},
		},
		{
			name:  "valid user non ascii",
			pass:  "password1",
			valid: true,
			user: User{
				ID:        1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Username:  "ギギ ジジ",
				Email:     "user@email.jp",
				Password:  password{},
				Version:   1,
			},
		},
		{
			name:  "invalid username empty",
			pass:  "password1",
			valid: false,
			user: User{
				ID:        1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Username:  "",
				Email:     "John.Doe@gmail.com",
				Password:  password{},
				Version:   1,
			},
		},
		{
			name:  "valid username long",
			pass:  "password1",
			valid: true,
			user: User{
				ID:        1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Username:  "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				Email:     "John.Doe@gmail.com",
				Password:  password{},
				Version:   1,
			},
		},
		{
			name:  "invalid username long",
			pass:  "password1",
			valid: false,
			user: User{
				ID:        1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Username:  "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				Email:     "John.Doe@gmail.com",
				Password:  password{},
				Version:   1,
			},
		},
		{
			name:  "invalid username long due to runes",
			pass:  "password1",
			valid: false,
			user: User{
				ID:        1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Username:  "ギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギギ",
				Email:     "John.Doe@gmail.com",
				Password:  password{},
				Version:   1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.user.Password.Set(tt.pass)

			v := validator.New()

			ValidateUser(v, &tt.user)
			assert.ValidatorValid(t, v, tt.valid)

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
		// {
		// 	name:  "invalid consecutive periods",
		// 	email: "john.doe@gmail.com",
		// 	valid: false,
		// },
		// {
		// 	name:  "invalid no characters",
		// 	email: ".@gmail.com",
		// 	valid: false,
		// },
		// {
		// 	name:  "invalid trailing special character 1",
		// 	email: "john.doe-@gmail.com",
		// 	valid: false,
		// },
		// {
		// 	name:  "invalid trailing special character 2",
		// 	email: "john.doe_@gmail.com",
		// 	valid: false,
		// },
		// {
		// 	name:  "invalid trailing special character 3",
		// 	email: "john.doe.@gmail.com",
		// 	valid: false,
		// },
		// {
		// 	name:  "invalid no top level domain",
		// 	email: "john.doe@gmail",
		// 	valid: false,
		// },
		// {
		// 	name:  "invalid top level domain qualifier too short",
		// 	email: "john.doe@gmail.c",
		// 	valid: false,
		// },
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

	// Test cases where created based on the criteria given in first result of google search blind to the regex being used by the validator
	// cases that the regex didn't catch where commented out, might be interesting to update regex to cover these cases later, but:
	// the regex being used in the validator (not written by me) is based on https://html.spec.whatwg.org/#valid-e-mail-address which is probably more definitive that whatever site I based the cases on

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			v := validator.New()

			ValidateEmail(v, tt.email)
			assert.ValidatorValid(t, v, tt.valid)

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

func TestUserModelAuthenticate(t *testing.T) {

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name      string
		user      User
		pass      string
		trypass   string
		wantError error
	}{
		{
			name: "Valid auth",
			user: User{
				Username: "test1",
				Password: password{},
				Email:    "test1@gmail.com",
			},
			pass:      "Pass1",
			trypass:   "Pass1",
			wantError: nil,
		},
		{
			name: "invalid auth space",
			user: User{
				Username: "test2",
				Password: password{},
				Email:    "test1@email.com",
			},
			pass:      "pass2",
			trypass:   "notcorrect",
			wantError: ErrInvalidCredentials,
		},
		{
			name: "invalid auth case",
			user: User{
				Username: "test2",
				Password: password{},
				Email:    "test1@email.com",
			},
			pass:      "pass2",
			trypass:   "Pass2",
			wantError: ErrInvalidCredentials,
		},
		{
			name: "invalid auth space",
			user: User{
				Username: "test2",
				Password: password{},
				Email:    "test1@email.com",
			},
			pass:      "pass2",
			trypass:   "pass2 ",
			wantError: ErrInvalidCredentials,
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
			assert.ExpectError(t, err, nil)
			if err != nil {
				return
			}

			_, err = m.Authenticate(tt.user.Email, tt.trypass)

			assert.ExpectError(t, err, tt.wantError)
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

func TestUserModelGetForToken(t *testing.T) {

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name      string
		ID        int64
		ttl       time.Duration
		scope     string
		wantError error
	}{
		{
			name:      "user exists for token",
			ID:        1,
			ttl:       3 * time.Hour,
			scope:     "authentication",
			wantError: nil,
		},
		{
			name:      "user doesn't exist",
			ID:        999999,
			ttl:       3 * time.Hour,
			scope:     "authentication",
			wantError: ErrReferencedUserDoesNotExist,
		},
		{
			name:      "user cannot exist",
			ID:        -1,
			ttl:       3 * time.Hour,
			scope:     "authentication",
			wantError: ErrReferencedUserDoesNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "users")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}

			m := TokenModel{db}

			_, err = m.New(int64(tt.ID), tt.ttl, tt.scope)
			assert.ExpectError(t, err, tt.wantError)
			if err != nil {
				return
			}
		})
	}
}
