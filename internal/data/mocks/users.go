package mocks

import "github.com/tconnellan/macro-tracker-backend/internal/data"

type UserModelMock struct{}

func (m UserModelMock) Insert(*data.User) error {
	return nil
}
func (m UserModelMock) Exists(int) (bool, error) {
	return false, nil
}
func (m UserModelMock) GetByEmail(string) (*data.User, error) {
	return nil, nil
}
func (m UserModelMock) Update(*data.User) error {
	return nil
}
func (m UserModelMock) GetForToken(string, string) (*data.User, error) {
	return nil, nil
}
