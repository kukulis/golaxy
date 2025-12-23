package gamemath

import (
	"math/rand/v2"
)

// RandomGenerator is an interface for generating random numbers.
type RandomGenerator interface {
	// NextRandom returns the next random float64 value in the range [0.0, 1.0).
	NextRandom() float64
}

// StdRandomGenerator generates random numbers using Go's standard library.
type StdRandomGenerator struct {
	rng *rand.Rand
}

// NewStdRandomGenerator creates a new random generator using Go's standard library.
// If seed is 0, a random seed will be used. Otherwise, the provided seed is used
// for deterministic random number generation.
func NewStdRandomGenerator(seed uint64) *StdRandomGenerator {
	var rng *rand.Rand
	if seed == 0 {
		rng = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	} else {
		rng = rand.New(rand.NewPCG(seed, seed))
	}
	return &StdRandomGenerator{rng: rng}
}

// NextRandom returns the next random float64 value in the range [0.0, 1.0).
func (g *StdRandomGenerator) NextRandom() float64 {
	return g.rng.Float64()
}

// PredefinedRandomGenerator returns predefined values from an array in sequence.
// When the array is exhausted, it wraps around to the beginning.
type PredefinedRandomGenerator struct {
	values []float64
	index  int
}

// NewPredefinedRandomGenerator creates a new generator with predefined values.
func NewPredefinedRandomGenerator(values []float64) *PredefinedRandomGenerator {
	return &PredefinedRandomGenerator{
		values: values,
		index:  0,
	}
}

// NextRandom returns the next value from the predefined array.
// Wraps around to the beginning when the end is reached.
func (g *PredefinedRandomGenerator) NextRandom() float64 {
	value := g.values[g.index]
	g.index = (g.index + 1) % len(g.values)
	return value
}
