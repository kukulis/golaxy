package galaxy

type ShipModelContainer struct {
	ShipModel  *ShipModel `json:"shipModel"`
	Amount     int        `json:"amount"`
	ResultMass float64    `json:"resultMass"`
}

func (c *ShipModelContainer) CalculateResultMass() float64 {
	return c.ShipModel.CalculateTotalMass() * float64(c.Amount)
}
