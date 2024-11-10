package data

import (
	"fmt"
	"testing"
	"time"

	"github.com/tconnellan/macro-tracker-backend/internal/assert"
	"github.com/tconnellan/macro-tracker-backend/internal/validator"
)

func TestTokenHelpers(t *testing.T) {

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

func TestTokenModelNew(t *testing.T) {

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name      string
		userID    int
		ttl       time.Duration
		scope     string
		wantError error
	}{
		{
			name:      "Valid user",
			userID:    1,
			ttl:       3 * time.Hour,
			scope:     ScopeAuthentication,
			wantError: nil,
		},
		{
			name:      "non existent user",
			userID:    99999,
			ttl:       3 * time.Hour,
			scope:     ScopeAuthentication,
			wantError: ErrReferencedUserDoesNotExist,
		},
		{
			name:      "non existent user",
			userID:    0,
			ttl:       3 * time.Hour,
			scope:     ScopeAuthentication,
			wantError: ErrReferencedUserDoesNotExist,
		},
		{
			name:      "non existent user",
			userID:    -1,
			ttl:       3 * time.Hour,
			scope:     ScopeAuthentication,
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

			_, err = m.New(int64(tt.userID), tt.ttl, tt.scope)
			assert.ExpectError(t, err, tt.wantError)
			if err != nil {
				return
			}

		})
	}
}
