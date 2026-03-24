package galaxy

import "fmt"

type ShipModelAssignment struct {
	ShipModel ShipModel
	Amount    int
}

type FleetBuildStatistics struct {
	MaxResources                 int `json:"max_resources"`
	UsedResources                int `json:"used_resources"`
	UsedResourcesForShips        int `json:"used_resources_for_ships"`
	UsedResourcesForTechnologies int `json:"used_resources_for_technologies"`
	RemainingResources           int `json:"remaining_resources"`
	ExceedingResources           int `json:"exceeding_resources"`
}

func (fleetBuild *FleetBuild) CalculateStatistics(maxResources int) FleetBuildStatistics {
	usedForShips := 0
	for _, assignment := range fleetBuild.AssignedShipModels {
		usedForShips += int(assignment.ShipModel.CalculateTotalMass()) * assignment.Amount
	}

	usedForTech := int(fleetBuild.AttackResources + fleetBuild.DefenseResources + fleetBuild.EngineResources + fleetBuild.CargoResources)

	usedResources := usedForShips + usedForTech
	remaining := maxResources - usedResources
	exceeding := 0
	if remaining < 0 {
		exceeding = -remaining
		remaining = 0
	}

	return FleetBuildStatistics{
		MaxResources:                 maxResources,
		UsedResources:                usedResources,
		UsedResourcesForShips:        usedForShips,
		UsedResourcesForTechnologies: usedForTech,
		RemainingResources:           remaining,
		ExceedingResources:           exceeding,
	}
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
