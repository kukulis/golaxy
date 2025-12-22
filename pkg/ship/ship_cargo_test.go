package ship

import "testing"

func TestShipCargo(t *testing.T) {
	ship := Ship{tech: ShipTech{CargoCapacity: 10}}

	ship.AddCargoPeople(2)

	if ship.cargoPeople != 2 {
		t.Errorf("Ship cargo peaple %f; want 2", ship.cargoPeople)
	}

}
