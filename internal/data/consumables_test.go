package data

import (
	"fmt"
	"testing"

	"github.com/tconnellan/macro-tracker-backend/internal/assert"
	"github.com/tconnellan/macro-tracker-backend/internal/validator"
)

func TestConsumableHelpers(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"

	tests := []struct {
		name       string
		valid      bool
		consumable Consumable
	}{
		{
			name:  "valid consumable",
			valid: true,
			consumable: Consumable{
				ID:        1,
				CreatorID: 1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Name:      "Oats",
				BrandName: "Uncle Tobys",
				Size:      100,
				Units:     "g",
				Macros: Macronutrients{
					Carbs:    40,
					Fats:     0.5,
					Proteins: 3,
					Alcohol:  0,
				},
			},
		},
		{
			name:  "valid consumable units grams",
			valid: true,
			consumable: Consumable{
				ID:        1,
				CreatorID: 1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Name:      "Oats",
				BrandName: "Uncle Tobys",
				Size:      100,
				Units:     "g",
				Macros: Macronutrients{
					Carbs:    40,
					Fats:     0.5,
					Proteins: 3,
					Alcohol:  0,
				},
			},
		},
		{
			name:  "valid consumable units ml",
			valid: true,
			consumable: Consumable{
				ID:        1,
				CreatorID: 1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Name:      "Oats",
				BrandName: "Uncle Tobys",
				Size:      100,
				Units:     "ml",
				Macros: Macronutrients{
					Carbs:    40,
					Fats:     0.5,
					Proteins: 3,
					Alcohol:  0,
				},
			},
		},
		{
			name:  "invalid consumable units",
			valid: false,
			consumable: Consumable{
				ID:        1,
				CreatorID: 1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Name:      "Oats",
				BrandName: "Uncle Tobys",
				Size:      100,
				Units:     "malarkey",
				Macros: Macronutrients{
					Carbs:    40,
					Fats:     0.5,
					Proteins: 3,
					Alcohol:  0,
				},
			},
		},
		{
			name:  "invalid consumable name empty",
			valid: false,
			consumable: Consumable{
				ID:        1,
				CreatorID: 1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Name:      "",
				BrandName: "Uncle Tobys",
				Size:      100,
				Units:     "g",
				Macros: Macronutrients{
					Carbs:    40,
					Fats:     0.5,
					Proteins: 3,
					Alcohol:  0,
				},
			},
		},
		{
			name:  "invalid consumable name too long",
			valid: false,
			consumable: Consumable{
				ID:        1,
				CreatorID: 1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Name:      "oats oats oats oats oats oats oats oats oats oats oats",
				BrandName: "Uncle Tobys",
				Size:      100,
				Units:     "g",
				Macros: Macronutrients{
					Carbs:    40,
					Fats:     0.5,
					Proteins: 3,
					Alcohol:  0,
				},
			},
		},
		{
			name:  "invalid consumable size negative",
			valid: false,
			consumable: Consumable{
				ID:        1,
				CreatorID: 1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Name:      "oats",
				BrandName: "Uncle Tobys",
				Size:      -100,
				Units:     "g",
				Macros: Macronutrients{
					Carbs:    40,
					Fats:     0.5,
					Proteins: 3,
					Alcohol:  0,
				},
			},
		},
		{
			name:  "invalid consumable size zero",
			valid: false,
			consumable: Consumable{
				ID:        1,
				CreatorID: 1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Name:      "oats",
				BrandName: "Uncle Tobys",
				Size:      0,
				Units:     "g",
				Macros: Macronutrients{
					Carbs:    40,
					Fats:     0.5,
					Proteins: 3,
					Alcohol:  0,
				},
			},
		},
		{
			name:  "invalid consumable no macros",
			valid: false,
			consumable: Consumable{
				ID:        1,
				CreatorID: 1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Name:      "oats",
				BrandName: "Uncle Tobys",
				Size:      100,
				Units:     "g",
				Macros: Macronutrients{
					Carbs:    0,
					Fats:     0,
					Proteins: 0,
					Alcohol:  0,
				},
			},
		},
		{
			name:  "invalid consumable negative carb",
			valid: false,
			consumable: Consumable{
				ID:        1,
				CreatorID: 1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Name:      "oats",
				BrandName: "Uncle Tobys",
				Size:      100,
				Units:     "g",
				Macros: Macronutrients{
					Carbs:    -1,
					Fats:     0,
					Proteins: 0,
					Alcohol:  0,
				},
			},
		},
		{
			name:  "invalid consumable negative fat",
			valid: false,
			consumable: Consumable{
				ID:        1,
				CreatorID: 1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Name:      "oats",
				BrandName: "Uncle Tobys",
				Size:      100,
				Units:     "g",
				Macros: Macronutrients{
					Carbs:    0,
					Fats:     -1,
					Proteins: 0,
					Alcohol:  0,
				},
			},
		},
		{
			name:  "invalid consumable negative protein",
			valid: false,
			consumable: Consumable{
				ID:        1,
				CreatorID: 1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Name:      "oats",
				BrandName: "Uncle Tobys",
				Size:      100,
				Units:     "g",
				Macros: Macronutrients{
					Carbs:    0,
					Fats:     0,
					Proteins: -1,
					Alcohol:  0,
				},
			},
		},
		{
			name:  "invalid consumable negative alcohol",
			valid: false,
			consumable: Consumable{
				ID:        1,
				CreatorID: 1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Name:      "oats",
				BrandName: "Uncle Tobys",
				Size:      100,
				Units:     "g",
				Macros: Macronutrients{
					Carbs:    0,
					Fats:     0,
					Proteins: 0,
					Alcohol:  -1,
				},
			},
		},
		{
			name:  "invalid consumable negative macros but positive sum",
			valid: false,
			consumable: Consumable{
				ID:        1,
				CreatorID: 1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Name:      "oats",
				BrandName: "Uncle Tobys",
				Size:      100,
				Units:     "g",
				Macros: Macronutrients{
					Carbs:    -1,
					Fats:     -1,
					Proteins: -1,
					Alcohol:  20,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			v := validator.New()
			ValidateConsumable(v, &tt.consumable)
			assert.ValidatorValid(t, v, tt.valid)
		})
	}
}

func TestConsumableModelGetByID(t *testing.T) {

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	timeFormat := "2006-01-02 15:04:05"

	tests := []struct {
		name           string
		ID             int64
		wantError      error
		wantMetadata   Metadata
		wantConsumable Consumable
	}{
		{
			name:      "get existing",
			ID:        1,
			wantError: nil,
			wantConsumable: Consumable{
				ID:        1,
				CreatorID: 1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Name:      "Oats",
				BrandName: "Uncle Tobys",
				Size:      100,
				Units:     "g",
				Macros: Macronutrients{
					Carbs:    40,
					Fats:     0.5,
					Proteins: 3,
					Alcohol:  0,
				},
			},
		},
		{
			name:           "get non-existing",
			ID:             999999,
			wantError:      ErrRecordNotFound,
			wantConsumable: Consumable{},
		},
		{
			name:           "get Bad ID",
			ID:             -1,
			wantError:      ErrRecordNotFound,
			wantConsumable: Consumable{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipes")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}

			m := ConsumableModel{db}

			cons, err := m.GetByID(tt.ID)
			assert.ExpectError(t, err, tt.wantError)
			if err != nil {
				return
			}

			assert.Equal(t, *cons, tt.wantConsumable)
		})
	}
}

func TestGetByCreatorID(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	timeFormat := "2006-01-02 15:04:05"

	tests := []struct {
		name            string
		creatorID       int64
		wantError       error
		wantConsumables []*Consumable
		wantMetadata    Metadata
		filters         ConsumableFilters
	}{
		{
			name:      "get existing",
			creatorID: 1,
			wantError: nil,
			filters: ConsumableFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch:                   "",
				BrandNameSearch:              "",
				RequireNameAndBrandNameMatch: false,
			},
			wantMetadata: Metadata{
				CurrentPage:  1,
				PageSize:     100,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 6,
			},
			wantConsumables: []*Consumable{
				{
					ID:        1,
					CreatorID: 1,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Oats",
					BrandName: "Uncle Tobys",
					Size:      100,
					Units:     "g",
					Macros: Macronutrients{
						Carbs:    40,
						Fats:     0.5,
						Proteins: 3,
						Alcohol:  0,
					},
				},
				{
					ID:        2,
					CreatorID: 1,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Cavendish Banana",
					BrandName: "Coles",
					Size:      100,
					Units:     "g",
					Macros: Macronutrients{
						Carbs:    38,
						Fats:     0.1,
						Proteins: 2,
						Alcohol:  0,
					},
				},
				{
					ID:        3,
					CreatorID: 1,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Greek Yogurt",
					BrandName: "Jalna",
					Size:      90,
					Units:     "g",
					Macros: Macronutrients{
						Carbs:    3.8,
						Fats:     5.0,
						Proteins: 9.0,
						Alcohol:  0,
					},
				},
				{
					ID:        4,
					CreatorID: 1,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Wholemeal Bread",
					BrandName: "Tip Top",
					Size:      110,
					Units:     "g",
					Macros: Macronutrients{
						Carbs:    41.8,
						Fats:     2.2,
						Proteins: 8.8,
						Alcohol:  0,
					},
				},
				{
					ID:        5,
					CreatorID: 1,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Red Apple",
					BrandName: "Aldi",
					Size:      95,
					Units:     "g",
					Macros: Macronutrients{
						Carbs:    14.0,
						Fats:     0.2,
						Proteins: 0.3,
						Alcohol:  0,
					},
				},
				{
					ID:        6,
					CreatorID: 1,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Chicken Breast",
					BrandName: "IGA",
					Size:      105,
					Units:     "g",
					Macros: Macronutrients{
						Carbs:    0,
						Fats:     2.6,
						Proteins: 22.5,
						Alcohol:  0,
					},
				},
			},
		},
		{
			name:      "get existing 2",
			creatorID: 2,
			wantError: nil,
			filters: ConsumableFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch:                   "",
				BrandNameSearch:              "",
				RequireNameAndBrandNameMatch: false,
			},
			wantMetadata: Metadata{
				CurrentPage:  1,
				PageSize:     100,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 5,
			},
			wantConsumables: []*Consumable{
				{
					ID:        7,
					CreatorID: 2,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Almond Milk",
					BrandName: "Vitasoy",
					Size:      250,
					Units:     "ml",
					Macros: Macronutrients{
						Carbs:    0.8,
						Fats:     1.2,
						Proteins: 0.5,
						Alcohol:  0,
					},
				},
				{
					ID:        8,
					CreatorID: 2,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Sweet Potato",
					BrandName: "Woolworths",
					Size:      150,
					Units:     "g",
					Macros: Macronutrients{
						Carbs:    27.5,
						Fats:     0.1,
						Proteins: 2.0,
						Alcohol:  0,
					},
				},
				{
					ID:        9,
					CreatorID: 2,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Salmon Fillet",
					BrandName: "Tassal",
					Size:      125,
					Units:     "g",
					Macros: Macronutrients{
						Carbs:    0,
						Fats:     12.5,
						Proteins: 25.0,
						Alcohol:  0,
					},
				},
				{
					ID:        10,
					CreatorID: 2,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Quinoa",
					BrandName: "Coles",
					Size:      85,
					Units:     "g",
					Macros: Macronutrients{
						Carbs:    15.6,
						Fats:     2.4,
						Proteins: 4.8,
						Alcohol:  0,
					},
				},
				{
					ID:        11,
					CreatorID: 2,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Red Wine",
					BrandName: "Penfolds",
					Size:      150,
					Units:     "ml",
					Macros: Macronutrients{
						Carbs:    3.8,
						Fats:     0,
						Proteins: 0,
						Alcohol:  13.5,
					},
				},
			},
		},
		{
			name:      "get existing paginated",
			creatorID: 1,
			wantError: nil,
			filters: ConsumableFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     2,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch:                   "",
				BrandNameSearch:              "",
				RequireNameAndBrandNameMatch: false,
			},
			wantMetadata: Metadata{
				CurrentPage:  1,
				PageSize:     2,
				FirstPage:    1,
				LastPage:     3,
				TotalRecords: 6,
			},
			wantConsumables: []*Consumable{
				{
					ID:        1,
					CreatorID: 1,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Oats",
					BrandName: "Uncle Tobys",
					Size:      100,
					Units:     "g",
					Macros: Macronutrients{
						Carbs:    40,
						Fats:     0.5,
						Proteins: 3,
						Alcohol:  0,
					},
				},
				{
					ID:        2,
					CreatorID: 1,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Cavendish Banana",
					BrandName: "Coles",
					Size:      100,
					Units:     "g",
					Macros: Macronutrients{
						Carbs:    38,
						Fats:     0.1,
						Proteins: 2,
						Alcohol:  0,
					},
				},
			},
		},
		{
			name:      "get existing paginated offset",
			creatorID: 1,
			wantError: nil,
			filters: ConsumableFilters{
				Metadata: MetadataFilters{
					Page:         2,
					PageSize:     2,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch:                   "",
				BrandNameSearch:              "",
				RequireNameAndBrandNameMatch: false,
			},
			wantMetadata: Metadata{
				CurrentPage:  2,
				PageSize:     2,
				FirstPage:    1,
				LastPage:     3,
				TotalRecords: 6,
			},
			wantConsumables: []*Consumable{
				{
					ID:        3,
					CreatorID: 1,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Greek Yogurt",
					BrandName: "Jalna",
					Size:      90,
					Units:     "g",
					Macros: Macronutrients{
						Carbs:    3.8,
						Fats:     5.0,
						Proteins: 9.0,
						Alcohol:  0,
					},
				},
				{
					ID:        4,
					CreatorID: 1,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Wholemeal Bread",
					BrandName: "Tip Top",
					Size:      110,
					Units:     "g",
					Macros: Macronutrients{
						Carbs:    41.8,
						Fats:     2.2,
						Proteins: 8.8,
						Alcohol:  0,
					},
				},
			},
		},
		{
			name:      "get not existing",
			creatorID: 999999,
			wantError: nil,
			filters: ConsumableFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch:                   "",
				BrandNameSearch:              "",
				RequireNameAndBrandNameMatch: false,
			},
			wantMetadata:    Metadata{},
			wantConsumables: []*Consumable{},
		},
		{
			name:      "get bad id",
			creatorID: -1,
			wantError: nil,
			filters: ConsumableFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch:                   "",
				BrandNameSearch:              "",
				RequireNameAndBrandNameMatch: false,
			},
			wantMetadata:    Metadata{},
			wantConsumables: []*Consumable{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipes")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}

			m := ConsumableModel{db}

			consumables, metadata, err := m.GetByCreatorID(tt.creatorID, tt.filters)
			assert.ExpectError(t, err, tt.wantError)
			if err != nil {
				return
			}

			assert.Equal(t, len(consumables), len(tt.wantConsumables))
			for i := range consumables {
				assert.Equal(t, *consumables[i], *tt.wantConsumables[i])
			}

			assert.Equal(t, metadata, tt.wantMetadata)
		})
	}
}

func TestConsumableModelSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	timeFormat := "2006-01-02 15:04:05"

	tests := []struct {
		name            string
		wantError       error
		wantConsumables []*Consumable
		wantMetadata    Metadata
		filters         ConsumableFilters
	}{
		{
			name:      "search existing by name",
			wantError: nil,
			filters: ConsumableFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch:                   "milk",
				BrandNameSearch:              "",
				RequireNameAndBrandNameMatch: true,
			},
			wantMetadata: Metadata{
				CurrentPage:  1,
				PageSize:     100,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 6,
			},
			wantConsumables: []*Consumable{
				{
					ID:        7,
					CreatorID: 2,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Almond Milk",
					BrandName: "Vitasoy",
					Size:      250,
					Units:     "ml",
					Macros: Macronutrients{
						Carbs:    0.8,
						Fats:     1.2,
						Proteins: 0.5,
						Alcohol:  0,
					},
				},
				{
					ID:        12,
					CreatorID: 3,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Full Cream Milk",
					BrandName: "Dairy Farmers",
					Size:      250,
					Units:     "ml",
					Macros: Macronutrients{
						Carbs:    12.5,
						Fats:     8.8,
						Proteins: 8.5,
						Alcohol:  0,
					},
				},
				{
					ID:        13,
					CreatorID: 3,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Skim Milk",
					BrandName: "Pauls",
					Size:      250,
					Units:     "ml",
					Macros: Macronutrients{
						Carbs:    12.8,
						Fats:     0.3,
						Proteins: 8.8,
						Alcohol:  0,
					},
				},
				{
					ID:        14,
					CreatorID: 3,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Soy Milk",
					BrandName: "Sanitarium",
					Size:      250,
					Units:     "ml",
					Macros: Macronutrients{
						Carbs:    5.5,
						Fats:     3.2,
						Proteins: 7.5,
						Alcohol:  0,
					},
				},
				{
					ID:        15,
					CreatorID: 3,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Oat Milk",
					BrandName: "Oatly",
					Size:      250,
					Units:     "ml",
					Macros: Macronutrients{
						Carbs:    16.0,
						Fats:     3.0,
						Proteins: 1.0,
						Alcohol:  0,
					},
				},
				{
					ID:        16,
					CreatorID: 3,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Coconut Milk",
					BrandName: "Pure Harvest",
					Size:      250,
					Units:     "ml",
					Macros: Macronutrients{
						Carbs:    3.0,
						Fats:     5.0,
						Proteins: 0.5,
						Alcohol:  0,
					},
				},
			},
		},
		{
			name:      "get existing by name or brand",
			wantError: nil,
			filters: ConsumableFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch:                   "milk",
				BrandNameSearch:              "Aldi",
				RequireNameAndBrandNameMatch: false,
			},
			wantMetadata: Metadata{
				CurrentPage:  1,
				PageSize:     100,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 7,
			},
			wantConsumables: []*Consumable{
				{
					ID:        5,
					CreatorID: 1,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Red Apple",
					BrandName: "Aldi",
					Size:      95,
					Units:     "g",
					Macros: Macronutrients{
						Carbs:    14.0,
						Fats:     0.2,
						Proteins: 0.3,
						Alcohol:  0,
					},
				},
				{
					ID:        7,
					CreatorID: 2,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Almond Milk",
					BrandName: "Vitasoy",
					Size:      250,
					Units:     "ml",
					Macros: Macronutrients{
						Carbs:    0.8,
						Fats:     1.2,
						Proteins: 0.5,
						Alcohol:  0,
					},
				},
				{
					ID:        12,
					CreatorID: 3,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Full Cream Milk",
					BrandName: "Dairy Farmers",
					Size:      250,
					Units:     "ml",
					Macros: Macronutrients{
						Carbs:    12.5,
						Fats:     8.8,
						Proteins: 8.5,
						Alcohol:  0,
					},
				},
				{
					ID:        13,
					CreatorID: 3,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Skim Milk",
					BrandName: "Pauls",
					Size:      250,
					Units:     "ml",
					Macros: Macronutrients{
						Carbs:    12.8,
						Fats:     0.3,
						Proteins: 8.8,
						Alcohol:  0,
					},
				},
				{
					ID:        14,
					CreatorID: 3,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Soy Milk",
					BrandName: "Sanitarium",
					Size:      250,
					Units:     "ml",
					Macros: Macronutrients{
						Carbs:    5.5,
						Fats:     3.2,
						Proteins: 7.5,
						Alcohol:  0,
					},
				},
				{
					ID:        15,
					CreatorID: 3,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Oat Milk",
					BrandName: "Oatly",
					Size:      250,
					Units:     "ml",
					Macros: Macronutrients{
						Carbs:    16.0,
						Fats:     3.0,
						Proteins: 1.0,
						Alcohol:  0,
					},
				},
				{
					ID:        16,
					CreatorID: 3,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Coconut Milk",
					BrandName: "Pure Harvest",
					Size:      250,
					Units:     "ml",
					Macros: Macronutrients{
						Carbs:    3.0,
						Fats:     5.0,
						Proteins: 0.5,
						Alcohol:  0,
					},
				},
			},
		},
		{
			name:      "get existing by name and brand",
			wantError: nil,
			filters: ConsumableFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch:                   "milk",
				BrandNameSearch:              "Pauls",
				RequireNameAndBrandNameMatch: true,
			},
			wantMetadata: Metadata{
				CurrentPage:  1,
				PageSize:     100,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 1,
			},
			wantConsumables: []*Consumable{
				{
					ID:        13,
					CreatorID: 3,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Skim Milk",
					BrandName: "Pauls",
					Size:      250,
					Units:     "ml",
					Macros: Macronutrients{
						Carbs:    12.8,
						Fats:     0.3,
						Proteins: 8.8,
						Alcohol:  0,
					},
				},
			},
		},
		{
			name:      "search paginated",
			wantError: nil,
			filters: ConsumableFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     2,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch:                   "milk",
				BrandNameSearch:              "",
				RequireNameAndBrandNameMatch: true,
			},
			wantMetadata: Metadata{
				CurrentPage:  1,
				PageSize:     2,
				FirstPage:    1,
				LastPage:     3,
				TotalRecords: 6,
			},
			wantConsumables: []*Consumable{
				{
					ID:        7,
					CreatorID: 2,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Almond Milk",
					BrandName: "Vitasoy",
					Size:      250,
					Units:     "ml",
					Macros: Macronutrients{
						Carbs:    0.8,
						Fats:     1.2,
						Proteins: 0.5,
						Alcohol:  0,
					},
				},
				{
					ID:        12,
					CreatorID: 3,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Full Cream Milk",
					BrandName: "Dairy Farmers",
					Size:      250,
					Units:     "ml",
					Macros: Macronutrients{
						Carbs:    12.5,
						Fats:     8.8,
						Proteins: 8.5,
						Alcohol:  0,
					},
				},
			},
		},
		{
			name:      "search existing paginated offset",
			wantError: nil,
			filters: ConsumableFilters{
				Metadata: MetadataFilters{
					Page:         2,
					PageSize:     2,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch:                   "milk",
				BrandNameSearch:              "",
				RequireNameAndBrandNameMatch: true,
			},
			wantMetadata: Metadata{
				CurrentPage:  2,
				PageSize:     2,
				FirstPage:    1,
				LastPage:     3,
				TotalRecords: 6,
			},
			wantConsumables: []*Consumable{
				{
					ID:        13,
					CreatorID: 3,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Skim Milk",
					BrandName: "Pauls",
					Size:      250,
					Units:     "ml",
					Macros: Macronutrients{
						Carbs:    12.8,
						Fats:     0.3,
						Proteins: 8.8,
						Alcohol:  0,
					},
				},
				{
					ID:        14,
					CreatorID: 3,
					CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
					Name:      "Soy Milk",
					BrandName: "Sanitarium",
					Size:      250,
					Units:     "ml",
					Macros: Macronutrients{
						Carbs:    5.5,
						Fats:     3.2,
						Proteins: 7.5,
						Alcohol:  0,
					},
				},
			},
		},
		{
			name:      "search not existing",
			wantError: nil,
			filters: ConsumableFilters{
				Metadata: MetadataFilters{
					Page:         1,
					PageSize:     100,
					Sort:         "ID",
					SortSafeList: []string{"ID"},
				},
				NameSearch:                   "nonamewillhavethisnameitspreposterous",
				BrandNameSearch:              "samehere",
				RequireNameAndBrandNameMatch: true,
			},
			wantMetadata:    Metadata{},
			wantConsumables: []*Consumable{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipes")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}

			m := ConsumableModel{db}

			consumables, metadata, err := m.Search(tt.filters)
			assert.ExpectError(t, err, tt.wantError)
			if err != nil {
				return
			}

			assert.Equal(t, len(consumables), len(tt.wantConsumables))
			for i := range consumables {
				assert.Equal(t, *consumables[i], *tt.wantConsumables[i])
			}

			assert.Equal(t, metadata, tt.wantMetadata)
		})
	}
}

func TestConsumableModelInsert(t *testing.T) {

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	timeFormat := "2006-01-02 15:04:05"

	tests := []struct {
		name          string
		wantError     error
		newConsumable Consumable
	}{
		{
			name:      "insert ok",
			wantError: nil,
			newConsumable: Consumable{
				ID:        0,
				CreatorID: 1,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Name:      "Rolled Oats",
				BrandName: "Uncle Tobys",
				Size:      100,
				Units:     "g",
				Macros: Macronutrients{
					Carbs:    40,
					Fats:     0.5,
					Proteins: 3,
					Alcohol:  0,
				},
			},
		},
		{
			name:      "insert creator does not exist",
			wantError: ErrReferencedUserDoesNotExist,
			newConsumable: Consumable{
				ID:        0,
				CreatorID: 9999999,
				CreatedAt: MustParse(timeFormat, "2024-01-01 10:00:00"),
				Name:      "Rolled Oats",
				BrandName: "Uncle Tobys",
				Size:      100,
				Units:     "g",
				Macros: Macronutrients{
					Carbs:    40,
					Fats:     0.5,
					Proteins: 3,
					Alcohol:  0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipes")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}

			m := ConsumableModel{db}

			insertConsumable := tt.newConsumable

			err = m.Insert(&tt.newConsumable)
			assert.ExpectError(t, err, tt.wantError)
			if err != nil {
				return
			}

			// new ID and createdAt is populated inplace of the given consumable
			assert.NotEqual(t, insertConsumable.ID, tt.newConsumable.ID)
			insertConsumable.ID = tt.newConsumable.ID
			insertConsumable.CreatedAt = tt.newConsumable.CreatedAt
			assert.Equal(t, tt.newConsumable, insertConsumable)
		})
	}
}

func TestConsumableModelDelete(t *testing.T) {

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name      string
		ID        int64
		wantError error
	}{
		{
			name:      "delete ok",
			ID:        1,
			wantError: nil,
		},
		{
			name:      "delete not exist",
			ID:        99999,
			wantError: ErrRecordNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := newTestDB(t, "recipes")
			if err != nil {
				t.Fatal(fmt.Errorf("Failed test db setup: %w", err))
			}

			m := ConsumableModel{db}

			err = m.Delete(tt.ID)
			assert.ExpectError(t, err, tt.wantError)
		})
	}
}
