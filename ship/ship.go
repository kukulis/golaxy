package ship

type Ship struct {
	tech           ShipTech
	X              float64
	Y              float64
	cargoPeople    float64
	cargoMaterials float64
}

func (ship *Ship) CargoExceeded() bool {
	totalCargo := ship.cargoPeople + ship.cargoMaterials

	return totalCargo > ship.tech.CargoCapacity
}

func (ship *Ship) AddCargoPeople(c float64) bool {

	if ship.cargoPeople+ship.cargoMaterials+c > ship.tech.CargoCapacity {
		return false
	}

	ship.cargoPeople = ship.cargoPeople + c

	return true
}

func (ship *Ship) AddCargoMaterials(c float64) bool {
	if ship.cargoPeople+ship.cargoMaterials+c > ship.tech.CargoCapacity {
		return false
	}

	ship.cargoMaterials = ship.cargoPeople + c

	return true
}
