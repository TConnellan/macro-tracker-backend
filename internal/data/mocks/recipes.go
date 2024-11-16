package mocks

import "github.com/tconnellan/macro-tracker-backend/internal/data"

type RecipeModelMock struct{}

func (m RecipeModelMock) Get(ID int64) (*data.Recipe, error) {
	return nil, nil
}

func (m RecipeModelMock) GetByCreatorID(creatorID int64, filters data.RecipeFilters) ([]*data.Recipe, data.Metadata, error) {

	switch creatorID {
	case 1:
		return []*data.Recipe{}, data.Metadata{}, nil
	default:
		return []*data.Recipe{}, data.Metadata{}, nil
	}
}

func (m RecipeModelMock) GetLatestByCreatorID(int64, data.RecipeFilters) ([]*data.Recipe, data.Metadata, error) {
	return nil, data.Metadata{}, nil
}

func (m RecipeModelMock) GetFullRecipe(int64, int64) (*data.FullRecipe, error) {
	return nil, nil
}

func (m RecipeModelMock) Insert(*data.Recipe) error {
	return nil
}

func (m RecipeModelMock) InsertFullRecipe(*data.FullRecipe) error {
	return nil
}

func (m RecipeModelMock) Update(*data.Recipe) error {
	return nil
}

func (m RecipeModelMock) UpdateFullRecipe(*data.FullRecipe) error {
	return nil
}

func (m RecipeModelMock) Delete(int64) error {
	return nil
}

func (m RecipeModelMock) GetParentRecipe(*data.Recipe) (*data.Recipe, error) {
	return nil, nil
}

func (m RecipeModelMock) GetAllAncestors(*data.Recipe, data.RecipeFilters) ([]*data.Recipe, data.Metadata, error) {
	return nil, data.Metadata{}, nil
}
