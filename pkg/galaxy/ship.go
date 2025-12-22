package galaxy

type Ship struct {
	ID        string   `json:"id"`
	Tech      ShipTech `json:"tech"`
	Destroyed bool     `json:"destroyed"`
	Name      string   `json:"name"`
	Owner     string   `json:"owner"` // race owner id
}
