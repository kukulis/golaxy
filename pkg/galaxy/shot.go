package galaxy

type Shot struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Result      bool   `json:"result"`
}

// Equal compares two shots for equality
func (s *Shot) Equal(other *Shot) bool {
	if s == nil || other == nil {
		return s == other
	}

	return s.Source == other.Source &&
		s.Destination == other.Destination &&
		s.Result == other.Result
}
