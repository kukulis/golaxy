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

	shipMap map[string]*Ship
}

func NewFleet(ships []*Ship) *Fleet {
	shipMap := make(map[string]*Ship)
	for _, ship := range ships {
		shipMap[ship.ID] = ship
	}
	return &Fleet{
		Ships:   ships,
		shipMap: shipMap,
	}
}

func (fleet *Fleet) Speed() float64 {
	speed := math.MaxFloat64

	for _, ship := range fleet.Ships {
		speed = min(speed, ship.Tech.Speed)
	}

	return speed
}

func (fleet *Fleet) Speed2() float64 {
	if len(fleet.Ships) == 0 {
		return 0
	}

	speeds := util.ArrayMap(fleet.Ships, func(ship *Ship) float64 { return ship.Tech.Speed })

	return util.ArrayReduce(speeds, math.MaxFloat64, func(a float64, b float64) float64 { return min(a, b) })
}

// EqualShips compares two fleets for equality based on their ships
func (fleet *Fleet) EqualShips(other *Fleet, logger Logger) bool {
	if fleet == nil || other == nil {
		return fleet == other
	}

	if len(fleet.Ships) != len(other.Ships) {
		if logger != nil {
			logger.Printf("Fleet ship count mismatch: %d vs %d", len(fleet.Ships), len(other.Ships))
		}
		return false
	}

	for i, ship := range fleet.Ships {
		if !ship.EqualFields(other.Ships[i]) {
			if logger != nil {
				logger.Printf("Ship[%d] mismatch: ship ID=%s does not equal ship ID=%s", i, ship.ID, other.Ships[i].ID)
			}
			return false
		}
	}

	return true
}

func (fleet *Fleet) GetShipById(id string) *Ship {
	return fleet.shipMap[id]
}
