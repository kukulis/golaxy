package gamemath

import (
	"math"
	"testing"
)

// TestNewConfigurableFunctionValidation tests that NewConfigurableFunction
// returns errors for invalid parameters.
func TestNewConfigurableFunctionValidation(t *testing.T) {
	tests := []struct {
		name      string
		keys      []float64
		values    []float64
		wantError bool
		errorMsg  string
	}{
		{
			name:      "Empty arrays",
			keys:      []float64{},
			values:    []float64{},
			wantError: true,
			errorMsg:  "keys and values cannot be empty",
		},
		{
			name:      "Nil keys",
			keys:      nil,
			values:    []float64{1, 2, 3},
			wantError: true,
			errorMsg:  "keys and values cannot be nil",
		},
		{
			name:      "Nil values",
			keys:      []float64{1, 2, 3},
			values:    nil,
			wantError: true,
			errorMsg:  "keys and values cannot be nil",
		},
		{
			name:      "Mismatched lengths - more keys",
			keys:      []float64{0, 10, 20},
			values:    []float64{0, 100},
			wantError: true,
			errorMsg:  "keys and values must have the same length",
		},
		{
			name:      "Mismatched lengths - more values",
			keys:      []float64{0, 10},
			values:    []float64{0, 100, 200},
			wantError: true,
			errorMsg:  "keys and values must have the same length",
		},
		{
			name:      "Duplicate keys",
			keys:      []float64{0, 10, 10, 20},
			values:    []float64{0, 100, 150, 200},
			wantError: true,
			errorMsg:  "keys must be strictly increasing (duplicate key at index 2)",
		},
		{
			name:      "Unsorted keys",
			keys:      []float64{0, 20, 10, 30},
			values:    []float64{0, 100, 150, 200},
			wantError: true,
			errorMsg:  "keys must be strictly increasing (key at index 2 is not greater than previous)",
		},
		{
			name:      "Valid single point",
			keys:      []float64{10},
			values:    []float64{100},
			wantError: false,
		},
		{
			name:      "Valid multiple points",
			keys:      []float64{0, 10, 20, 30},
			values:    []float64{0, 100, 200, 300},
			wantError: false,
		},
		{
			name:      "Valid non-uniform spacing",
			keys:      []float64{0, 1, 10, 100},
			values:    []float64{0, 10, 20, 30},
			wantError: false,
		},
		{
			name:      "Valid negative keys",
			keys:      []float64{-100, -50, 0, 50},
			values:    []float64{10, 20, 30, 40},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cf, err := NewConfigurableFunction(tt.keys, tt.values)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else {
					t.Logf("Got expected error: %v", err)
				}
				if cf != nil {
					t.Errorf("Expected nil ConfigurableFunction when error occurs, got %v", cf)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if cf == nil {
					t.Errorf("Expected valid ConfigurableFunction but got nil")
				}
			}
		})
	}
}

func TestCalculateRatio(t *testing.T) {
	// Function that maps strength ratios:
	// - denominator >= 4*numerator (ratio <= 0.25) -> 1
	// - denominator == numerator (ratio == 1) -> 0.5
	// - numerator >= 4*denominator (ratio >= 4) -> 0
	keys := []float64{0.25, 1, 4}
	values := []float64{1, 0.5, 0}

	cf, err := NewConfigurableFunction(keys, values)
	if err != nil {
		t.Fatalf("Failed to create ConfigurableFunction: %v", err)
	}

	tests := []struct {
		name        string
		numerator   float64
		denominator float64
		expected    float64
	}{
		// === ZERO CASES ===
		{"0/0", 0, 0, 0.5},
		{"0/1", 0, 1, 1.0},
		{"0/10", 0, 10, 1.0},
		{"0/100", 0, 100, 1.0},
		{"0/0.001", 0, 0.001, 1.0},
		{"0/1e-100", 0, 1e-100, 1.0},
		{"1/0", 1, 0, 0.0},
		{"10/0", 10, 0, 0.0},
		{"100/0", 100, 0, 0.0},
		{"0.001/0", 0.001, 0, 0.0},
		{"1e-100/0", 1e-100, 0, 0.0},

		// === BOUNDARY CASES ===
		{"10/10 (equal)", 10, 10, 0.5},
		{"10/40 (4x weaker)", 10, 40, 1.0},
		{"40/10 (4x stronger)", 40, 10, 0.0},
		{"1/4 (min key)", 1, 4, 1.0},
		{"4/1 (max key)", 4, 1, 0.0},

		// === INTERPOLATION CASES ===
		{"20/10 (2x stronger)", 20, 10, 0.3333333333333333},
		{"10/20 (2x weaker)", 10, 20, 0.8333333333333334},
		{"5/10 (between 0.25 and 1)", 5, 10, 0.8333333333333334},

		// === OUT OF RANGE CASES ===
		{"1000/1 (very large ratio)", 1000, 1, 0.0},
		{"1/1000 (very small ratio)", 1, 1000, 1.0},

		// === EDGE EQUALITY CASES ===
		{"0.5/0.5 (equal small)", 0.5, 0.5, 0.5},
		{"1000/1000 (equal large)", 1000, 1000, 0.5},
		{"1e-50/1e-50 (equal tiny)", 1e-50, 1e-50, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cf.CalculateRatio(tt.numerator, tt.denominator)

			if math.IsNaN(result) {
				t.Errorf("Result is NaN (unexpected)")
			} else if math.Abs(result-tt.expected) > 1e-10 {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}
