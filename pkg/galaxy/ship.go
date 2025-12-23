package galaxy

type Ship struct {
	ID        string   `json:"id"`
	Tech      ShipTech `json:"tech"`
	Destroyed bool     `json:"destroyed"`
	Name      string   `json:"name"`
	Owner     string   `json:"owner"` // race owner id
}

// EqualsWithoutDamage compares two ships for equality, ignoring the Destroyed field
func (s *Ship) EqualsWithoutDamage(other *Ship) bool {
	if s == nil || other == nil {
		return s == other
	}

	return s.ID == other.ID &&
		s.Name == other.Name &&
		s.Owner == other.Owner &&
		s.Tech == other.Tech
}
