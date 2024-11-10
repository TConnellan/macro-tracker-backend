package mocks

import "github.com/tconnellan/macro-tracker-backend/internal/data"

type PantryItemModelMock struct{}

func (m PantryItemModelMock) Get(int64) (*data.PantryItem, error) {
	return nil, nil
}

func (m PantryItemModelMock) Create(*data.PantryItem) error {
	return nil
}

func (m PantryItemModelMock) Update(*data.PantryItem) error {
	return nil
}

func (m PantryItemModelMock) Delete(int64) error {
	return nil
}
