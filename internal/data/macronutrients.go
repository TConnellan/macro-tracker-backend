package data

import "github.com/tconnellan/macro-tracker-backend/internal/validator"

type Macronutrients struct {
	Carbs    float64 `json:"carbs"`
	Fats     float64 `json:"fats"`
	Proteins float64 `json:"proteins"`
	Alcohol  float64 `json:"alcohol"`
}

func ValidateMacroNutrients(v *validator.Validator, macros Macronutrients) {
	v.Check(macros.Carbs >= 0, "carbs", "must be non-negative")
	v.Check(macros.Fats >= 0, "fats", "must be non-negative")
	v.Check(macros.Proteins >= 0, "proteins", "must be non-negative")
	v.Check(macros.Alcohol >= 0, "alcohol", "must be non-negative")
	v.Check(macros.Carbs+macros.Fats+macros.Proteins+macros.Alcohol > 0, "macronutrients", "one macronutrient must be non-negative")
}

func (macros *Macronutrients) CalculateKJ() float64 {
	return 16.7*macros.Carbs + 37.7*macros.Fats + 16.7*macros.Proteins + 29*macros.Alcohol
}
