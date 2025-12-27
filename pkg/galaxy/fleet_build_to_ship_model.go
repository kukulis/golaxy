package galaxy

type FleetBuildToShipModel struct {
	// stored to db

	FleetBuildID string  `json:"fleet_build_id"`
	ShipModelID  string  `json:"ship_model_id"`
	Amount       int     `json:"amount"`
	ResultMass   float64 `json:"result_mass"`

	// not stored to DB directly
	ShipModel *ShipModel
}

func (c *FleetBuildToShipModel) CalculateResultMass() float64 {
	return c.ShipModel.CalculateTotalMass() * float64(c.Amount)
}
