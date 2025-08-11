package ship

import (
	"galaktika.eu/util"
	"math"
)

type Flight struct {
	ships []*Ship
}

func (flight Flight) Speed() float64 {
	speed := math.MaxFloat64

	for _, ship := range flight.ships {
		speed = min(speed, ship.tech.Speed)
	}

	return speed
}

func (flight *Flight) Speed2() float64 {
	if len(flight.ships) == 0 {
		return 0
	}

	speeds := util.ArrayMap(flight.ships, func(ship *Ship) float64 { return ship.tech.Speed })

	return util.ArrayReduce(speeds, math.MaxFloat64, func(a float64, b float64) float64 { return min(a, b) })
}
