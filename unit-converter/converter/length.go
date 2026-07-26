package converter

import "fmt"

var lengthToMeter = map[string]float64{
	"millimeter": 0.001,
	"centimeter": 0.01,
	"meter":      1,
	"kilometer":  1000,
	"inch":       0.0254,
	"foot":       0.3048,
	"yard":       0.9144,
	"mile":       1609.344,
}

func ConvertLength(value float64, from, to string) (float64, error) {
	fromFactor, ok1 := lengthToMeter[from]
	toFactor, ok2 := lengthToMeter[to]
	if !ok1 {
		return 0, fmt.Errorf("the unit is not recognized: %v", from)
	} else if !ok2 {
		return 0, fmt.Errorf("the unit is not recognized: %v", to)
	}

	toMeters := value * fromFactor
	result := toMeters / toFactor

	return result, nil
}
