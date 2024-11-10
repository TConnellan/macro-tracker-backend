package mocks

import (
	"time"

	"github.com/tconnellan/macro-tracker-backend/internal/data"
)

type TokenModelMock struct{}

func (m TokenModelMock) New(userID int64, ttl time.Duration, scope string) (*data.Token, error) {
	return nil, nil
}
func (m TokenModelMock) Insert(token *data.Token) error {
	return nil
}
func (m TokenModelMock) DeleteAllForUser(string, int64) error {
	return nil
}
