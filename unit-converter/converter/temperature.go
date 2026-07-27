package converter

import "fmt"

func toCelsius(value float64, from string) (float64, error) {
	switch from {
	case "celsius":
		return value, nil
	case "fahrenheit":
		return (value - 32) * 5 / 9, nil
	case "kelvin":
		return value - 273.15, nil
	default:
		return 0, fmt.Errorf("the unit from is not recognized: %s", from)
	}
}

func fromCelsius(celsius float64, to string) (float64, error) {
	switch to {
	case "celsius":
		return celsius, nil
	case "fahrenheit":
		return (celsius * 9 / 5) + 32, nil
	case "kelvin":
		return celsius + 273.15, nil
	default:
		return 0, fmt.Errorf("the unit to is not recognized: %s", to)
	}
}

func ConvertTemperature(value float64, from, to string) (float64, error) {
	celsius, err := toCelsius(value, from)
	if err != nil {
		return 0, err
	}

	result, err := fromCelsius(celsius, to)
	if err != nil {
		return 0, err
	}

	return result, nil
}
