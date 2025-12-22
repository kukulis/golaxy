package galaxy

type ShipTech struct {
	Attack        float64 `json:"attack"`
	Guns          float64 `json:"guns"`
	Defense       float64 `json:"defense"`
	Speed         float64 `json:"speed"`
	CargoCapacity float64 `json:"cargo_capacity"`
	Mass          int     `json:"mass"`
}
