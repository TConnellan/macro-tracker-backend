package mocks

import "github.com/tconnellan/macro-tracker-backend/internal/data"

type RecipeComponentModelMock struct{}

func (m RecipeComponentModelMock) Get(ID int64) (*data.RecipeComponent, error) {
	return nil, nil
}

func (m RecipeComponentModelMock) Insert(*data.RecipeComponent) error {
	return nil
}

func (m RecipeComponentModelMock) Update(*data.RecipeComponent) error {
	return nil
}
