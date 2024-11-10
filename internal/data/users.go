package data

import (
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tconnellan/macro-tracker-backend/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Version   int       `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Username != "", "name", "must be provided")
	v.Check(len(user.Username) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type UserModel struct {
	DB *pgxpool.Pool
}

type IUserModel interface {
	Insert(*User) error
	Exists(int) (bool, error)
	GetByEmail(string) (*User, error)
	Update(*User) error
	GetForToken(string, string) (*User, error)
}

func (m UserModel) Insert(user *User) error {
	query := `
INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING id, created_at, version`

	args := []any{user.Username, user.Email, user.Password.hash}

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	err := m.DB.QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case strings.HasPrefix(err.Error(), `ERROR: duplicate key value violates unique constraint "users_email_key"`):
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id =$1);"

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	err := m.DB.QueryRow(ctx, stmt, id).Scan(&exists)
	return exists, err
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
SELECT id, created_at, username, email, password_hash, version
FROM users
WHERE email = $1`
	var user User
	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()
	err := m.DB.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) Authenticate(email string, password string) (*User, error) {

	user, err := m.GetByEmail(email)
	if err != nil {
		switch {
		case errors.Is(err, ErrRecordNotFound):
			return nil, ErrInvalidCredentials
		default:
			return nil, err
		}
	}

	check, err := user.Password.Matches(password)
	if err != nil || !check {
		return nil, ErrInvalidCredentials
	}

	return user, nil

}

func (m UserModel) Update(user *User) error {
	query := `
UPDATE users
SET username = $1, email = $2, password_hash = $3, version = version + 1
WHERE id = $4 AND version = $5
RETURNING version`
	args := []any{
		user.Username,
		user.Email,
		user.Password.hash,
		user.ID,
		user.Version,
	}
	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()
	err := m.DB.QueryRow(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case strings.HasPrefix(err.Error(), `ERROR: duplicate key value violates unique constraint "users_email_key"`):
			return ErrDuplicateEmail
		case errors.Is(err, pgx.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetForToken(tokenScope string, tokenPlaintext string) (*User, error) {

	query := `
	SELECT U.id, U.created_at, U.username, U.email, U.password_hash, U.version
	FROM users U INNER JOIN tokens T ON U.id = T.user_id
	WHERE T.hash = $1 AND T.scope = $2 AND T.expiry > $3;
	`

	tokenHash := hashToken(tokenPlaintext)

	var user User

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	err := m.DB.QueryRow(ctx, query, tokenHash[:], tokenScope, time.Now()).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
