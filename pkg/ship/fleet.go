package ship

import (
	"glaktika.eu/galaktika/pkg/util"
	"math"
)

type Fleet struct {
	ships []*Ship
}

func (flight Fleet) Speed() float64 {
	speed := math.MaxFloat64

	for _, ship := range flight.ships {
		speed = min(speed, ship.Tech.Speed)
	}

	return speed
}

func (flight *Fleet) Speed2() float64 {
	if len(flight.ships) == 0 {
		return 0
	}

	speeds := util.ArrayMap(flight.ships, func(ship *Ship) float64 { return ship.Tech.Speed })

	return util.ArrayReduce(speeds, math.MaxFloat64, func(a float64, b float64) float64 { return min(a, b) })
}
