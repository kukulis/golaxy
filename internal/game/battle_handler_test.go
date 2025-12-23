package game

import (
	"glaktika.eu/galaktika/pkg/gamemath"
	"testing"
)

func TestEvaluateTargetIndex(t *testing.T) {
	tests := []struct {
		name         string
		targetsCount int
		randomValue  float64
		expected     int
	}{
		// Single target
		{"Single target, random 0.0", 1, 0.0, 0},
		{"Single target, random 0.5", 1, 0.5, 0},
		{"Single target, random 0.99", 1, 0.99, 0},

		// Five targets
		{"Five targets, random 0.0", 5, 0.0, 0},
		{"Five targets, random 0.1", 5, 0.1, 0},
		{"Five targets, random 0.19", 5, 0.19, 0},
		{"Five targets, random 0.2", 5, 0.2, 1},
		{"Five targets, random 0.4", 5, 0.4, 2},
		{"Five targets, random 0.6", 5, 0.6, 3},
		{"Five targets, random 0.8", 5, 0.8, 4},
		{"Five targets, random 0.99", 5, 0.99, 4},

		// Ten targets
		{"Ten targets, random 0.0", 10, 0.0, 0},
		{"Ten targets, random 0.15", 10, 0.15, 1},
		{"Ten targets, random 0.5", 10, 0.5, 5},
		{"Ten targets, random 0.95", 10, 0.95, 9},
		{"Ten targets, random 0.999", 10, 0.999, 9},

		// Two targets
		{"Two targets, random 0.0", 2, 0.0, 0},
		{"Two targets, random 0.49", 2, 0.49, 0},
		{"Two targets, random 0.5", 2, 0.5, 1},
		{"Two targets, random 0.99", 2, 0.99, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := gamemath.NewPredefinedRandomGenerator([]float64{tt.randomValue})

			result := EvaluateTargetIndex(tt.targetsCount, rng)

			if result != tt.expected {
				t.Errorf("EvaluateTargetIndex(%d, random=%v) = %d, expected %d",
					tt.targetsCount, tt.randomValue, result, tt.expected)
			}
		})
	}
}
