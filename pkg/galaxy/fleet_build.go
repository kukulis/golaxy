package galaxy

type FleetBuild struct {
	// stored to DB
	ID         string `json:"id"`
	DivisionId string `json:"division_id"`
	RaceId     string `json:"race_id"`

	AttackResources  float64 `json:"attack_resources"`
	DefenseResources float64 `json:"defense_resources"`
	EngineResources  float64 `json:"engine_resources"`
	CargoResources   float64 `json:"cargo_resources"`

	// not stored to DB
	ModelContainer *FleetBuildToShipModel
	UsedResources  float64
}
