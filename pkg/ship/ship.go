package ship

type Ship struct {
	ID        string   `json:"id"`
	Tech      ShipTech `json:"tech"`
	X         float64  `json:"x"`
	Y         float64  `json:"y"`
	Destroyed bool     `json:"destroyed"`
}
