package mocks

import "github.com/tconnellan/macro-tracker-backend/internal/data"

type ConsumableModelMock struct{}

func (m ConsumableModelMock) GetByID(ID int64) (*data.Consumable, error) {
	return nil, nil
}

func (m ConsumableModelMock) GetByCreatorID(ID int64, filters data.ConsumableFilters) ([]*data.Consumable, data.Metadata, error) {
	return nil, data.Metadata{}, nil
}

func (m ConsumableModelMock) Search(data.ConsumableFilters) ([]*data.Consumable, data.Metadata, error) {
	return nil, data.Metadata{}, nil
}

func (m ConsumableModelMock) Insert(*data.Consumable) error {
	return nil
}

func (m ConsumableModelMock) Update(*data.Consumable) error {
	return nil
}

func (m ConsumableModelMock) Delete(int64) error {
	return nil
}
