package galaxy

import (
	"glaktika.eu/galaktika/pkg/util"
	"math"
)

// Logger interface for logging ship comparison mismatches
type Logger interface {
	Printf(format string, v ...interface{})
}

type Fleet struct {
	ID    string  `json:"id"`
	Ships []*Ship `json:"ships"`
	Owner string  `json:"owner"` // owner race id
}

func NewFleet(ships []*Ship) *Fleet {
	return &Fleet{
		Ships: ships,
	}
}

func (flight Fleet) Speed() float64 {
	speed := math.MaxFloat64

	for _, ship := range flight.Ships {
		speed = min(speed, ship.Tech.Speed)
	}

	return speed
}

func (flight *Fleet) Speed2() float64 {
	if len(flight.Ships) == 0 {
		return 0
	}

	speeds := util.ArrayMap(flight.Ships, func(ship *Ship) float64 { return ship.Tech.Speed })

	return util.ArrayReduce(speeds, math.MaxFloat64, func(a float64, b float64) float64 { return min(a, b) })
}

// EqualShips compares two fleets for equality based on their ships
func (f *Fleet) EqualShips(other *Fleet, logger Logger) bool {
	if f == nil || other == nil {
		return f == other
	}

	if len(f.Ships) != len(other.Ships) {
		if logger != nil {
			logger.Printf("Fleet ship count mismatch: %d vs %d", len(f.Ships), len(other.Ships))
		}
		return false
	}

	for i, ship := range f.Ships {
		if !ship.EqualFields(other.Ships[i]) {
			if logger != nil {
				logger.Printf("Ship[%d] mismatch: ship ID=%s does not equal ship ID=%s", i, ship.ID, other.Ships[i].ID)
			}
			return false
		}
	}

	return true
}
