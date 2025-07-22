package ship

import "math"

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
