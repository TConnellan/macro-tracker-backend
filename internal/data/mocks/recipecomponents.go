package mocks

import "github.com/tconnellan/macro-tracker-backend/internal/data"

type RecipeComponentModelMock struct{}

func (m RecipeComponentModelMock) Get(ID int64) (*data.RecipeComponent, error) {
	return nil, nil
}

func (m RecipeComponentModelMock) Insert(recipeComponent *data.RecipeComponent) error {
	return nil
}

func (m RecipeComponentModelMock) Update(recipeComponent *data.RecipeComponent, userId int64) error {
	return nil
}
