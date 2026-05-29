package galaxy

type ShipModelAssignment struct {
	// stored to db
	ID           string `json:"id"`
	FleetBuildID string `json:"fleet_build_id"`
	ShipModelID  string `json:"ship_model_id"`
	Amount       int    `json:"amount"`
	// not stored to DB directly
	ResultMass float64    `json:"result_mass"`
	ShipModel  *ShipModel `json:"shipModel"`
}

func (c *ShipModelAssignment) CalculateResultMass() float64 {
	return c.ShipModel.CalculateTotalMass() * float64(c.Amount)
}
