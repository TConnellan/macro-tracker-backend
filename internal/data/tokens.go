package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tconnellan/macro-tracker-backend/internal/validator"
)

const (
	ScopeAuthentication = "authentication"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {

	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil
}

func hashToken(tokenPlaintext string) [32]byte {
	return sha256.Sum256([]byte(tokenPlaintext))

}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

type ITokenModel interface {
	New(int64, time.Duration, string) (*Token, error)
	Insert(*Token) error
	DeleteAllForUser(string, int64) error
}

type TokenModel struct {
	DB *pgxpool.Pool
}

func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = m.Insert(token)
	if err != nil {
		switch {
		case strings.HasPrefix(err.Error(), "ERROR: insert or update on table \"tokens\" violates foreign key constraint \"tokens_user_id_fkey\""):
			return token, ErrReferencedUserDoesNotExist
		default:
			return token, err
		}
	}
	return token, err
}

func (m TokenModel) Insert(token *Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)`
	args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}
	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()
	_, err := m.DB.Exec(ctx, query, args...)
	return err
}

func (m TokenModel) DeleteAllForUser(scope string, userID int64) error {
	query := `
	DELETE FROM tokens
	WHERE scope = $1 AND user_id = $2`
	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()
	_, err := m.DB.Exec(ctx, query, scope, userID)
	return err
}
