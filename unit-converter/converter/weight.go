package converter

import (
	"fmt"
	"slices"
)

var weightToGram = map[string]float64{
	"milligram": 0.001,
	"gram":      1,
	"kilogram":  1000,
	"ounce":     28.3495,
	"pound":     453.592,
}

func ConvertWeight(value float64, from, to string) (float64, error) {
	fromFactor, ok1 := weightToGram[from]
	toFactor, ok2 := weightToGram[to]
	if !ok1 {
		return 0, fmt.Errorf("the unit is not recognized: %v", from)
	} else if !ok2 {
		return 0, fmt.Errorf("the unit is not recognized: %v", to)
	}

	toGrams := fromFactor * value
	result := toGrams / toFactor
	return result, nil
}

func WeightUnits() []string {
	var units []string
	for key := range weightToGram {
		units = append(units, key)
	}
	slices.Sort(units)
	return units
}
