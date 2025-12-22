package galaxy

import (
	"glaktika.eu/galaktika/pkg/util"
	"math"
)

type Fleet struct {
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
