package mocks

import (
	"time"

	"github.com/tconnellan/macro-tracker-backend/internal/data"
)

type ConsumedModelMock struct{}

func (m ConsumedModelMock) GetByConsumedID(ID int64) (*data.Consumed, error) {
	return nil, nil
}

func (m ConsumedModelMock) GetAllByUserID(userID int64) ([]*data.Consumed, error) {

	timeFormat := "2006-01-02 15:04:05"

	switch userID {
	case 1:
		return []*data.Consumed{
			{
				ID:           1,
				RecipeID:     1,
				UserID:       1,
				CreatedAt:    MustParse(timeFormat, "2024-01-01 10:00:00"),
				ConsumedAt:   MustParse(timeFormat, "2024-01-01 10:00:00"),
				Quantity:     1,
				LastEditedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Macros: data.Macronutrients{
					Carbs:    1,
					Fats:     1,
					Proteins: 1,
					Alcohol:  1,
				},
			},
		}, nil
	}

	return nil, nil
}

func (m ConsumedModelMock) GetAllByUserIDAndDate(userID int64, start time.Time, end time.Time) ([]*data.Consumed, error) {
	return nil, nil
}

func (m ConsumedModelMock) Insert(consumed *data.Consumed) error {
	return nil
}

func (m ConsumedModelMock) Update(consumed *data.Consumed) error {
	return nil
}

func (m ConsumedModelMock) Delete(ID int64) error {
	return nil
}
