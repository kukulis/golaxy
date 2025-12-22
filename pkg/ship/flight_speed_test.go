package ship

import "testing"

func TestFlightSpeed0(t *testing.T) {

	flight := Fleet{ships: []*Ship{}}

	rez := flight.Speed2()

	if rez != 0 {
		t.Errorf("Fleet speed %f; want 0", rez)
	}
}

func TestFlightSpeed1(t *testing.T) {

	flight := Fleet{ships: []*Ship{
		{Tech: ShipTech{CargoCapacity: 0, Speed: 10}},
	}}

	rez := flight.Speed2()

	if rez != 10 {
		t.Errorf("Fleet speed %f; want 10", rez)
	}
}
