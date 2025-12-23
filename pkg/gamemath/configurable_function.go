package gamemath

import "fmt"

type ConfigurableFunction struct {
	keys   []float64
	values []float64
}

func NewConfigurableFunction(keys []float64, values []float64) (*ConfigurableFunction, error) {
	// Check for nil arrays
	if keys == nil || values == nil {
		return nil, fmt.Errorf("keys and values cannot be nil")
	}

	// Check for empty arrays
	if len(keys) == 0 || len(values) == 0 {
		return nil, fmt.Errorf("keys and values cannot be empty")
	}

	// Check for mismatched lengths
	if len(keys) != len(values) {
		return nil, fmt.Errorf("keys and values must have the same length")
	}

	// Check that keys are strictly increasing (prevents division by zero and ensures proper interpolation)
	for i := 1; i < len(keys); i++ {
		if keys[i] <= keys[i-1] {
			if keys[i] == keys[i-1] {
				return nil, fmt.Errorf("keys must be strictly increasing (duplicate key at index %d)", i)
			}
			return nil, fmt.Errorf("keys must be strictly increasing (key at index %d is not greater than previous)", i)
		}
	}

	return &ConfigurableFunction{keys: keys, values: values}, nil
}

func (cf *ConfigurableFunction) Calculate(x float64) float64 {
	if x <= cf.keys[0] {
		return cf.values[0]
	}

	for i := 1; i < len(cf.keys); i++ {
		if x <= cf.keys[i] {
			deltaX := x - cf.keys[i-1]
			tangent := (cf.values[i] - cf.values[i-1]) / (cf.keys[i] - cf.keys[i-1])
			deltaY := tangent * deltaX

			return cf.values[i-1] + deltaY
		}
	}

	return cf.values[len(cf.values)-1]
}

// CalculateRatio calculates the function value at the ratio of numerator/denominator.
// Special case: if both numerator and denominator are 0, the ratio is treated as 1.
func (cf *ConfigurableFunction) CalculateRatio(numerator, denominator float64) float64 {
	var ratio float64

	// Special case: 0/0 treated as 1
	if numerator == 0 && denominator == 0 {
		ratio = 1.0
	} else {
		ratio = numerator / denominator
	}

	return cf.Calculate(ratio)
}
