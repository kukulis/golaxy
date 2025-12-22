package ship

type Shot struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Result      bool   `json:"result"`
}
