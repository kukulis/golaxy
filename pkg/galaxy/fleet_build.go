package galaxy

import "fmt"

type ShipModelAssignment struct {
	ShipModel ShipModel
	Amount    int
}

// FleetBuild consists of some ship models and some resources spent on technologies research.
type FleetBuild struct {
	// stored to DB
	ID         string `json:"id"`
	DivisionId string `json:"division_id"`
	RaceId     string `json:"race_id"`

	// Resources for research

	AttackResources  float64 `json:"attack_resources"`
	DefenseResources float64 `json:"defense_resources"`
	EngineResources  float64 `json:"engine_resources"`
	CargoResources   float64 `json:"cargo_resources"`

	// not stored to DB directly

	AssignedShipModels []ShipModelAssignment
	UsedResources      float64
}

func (fleetBuild *FleetBuild) CalculateShipTech(shipModel *ShipModel) ShipTech {
	tech := NewTechnologies()
	tech.Research(fleetBuild.AttackResources, fleetBuild.DefenseResources, fleetBuild.EngineResources, fleetBuild.CargoResources)

	return shipModel.CalculateShipTech(tech)
}

func (fleetBuild *FleetBuild) CalculateAllShipTechs() []*ShipTech {
	tech := NewTechnologies()
	tech.Research(fleetBuild.AttackResources, fleetBuild.DefenseResources, fleetBuild.EngineResources, fleetBuild.CargoResources)

	var rez = []*ShipTech{}

	for _, assignedShipModel := range fleetBuild.AssignedShipModels {
		fmt.Printf("Ship model: %v\n", assignedShipModel)

		shipTech := assignedShipModel.ShipModel.CalculateShipTech(tech)
		rez = append(rez, &shipTech)
	}

	return rez
}
