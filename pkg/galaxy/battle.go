package galaxy

type Battle struct {
	ID    string `json:"id"`
	SideA *Fleet `json:"side_a"`
	SideB *Fleet `json:"side_b"`

	Shots     []*Shot `json:"shots"`
	PostSideA *Fleet  `json:"post_side_a"`
	PostSideB *Fleet  `json:"post_side_b"`
}

// CompareShots compares the shots of this battle with another battle's shots
func (b *Battle) CompareShots(other *Battle, logger Logger) bool {
	if b == nil || other == nil {
		return b == other
	}

	if len(b.Shots) != len(other.Shots) {
		if logger != nil {
			logger.Printf("Battle shots count mismatch: %d vs %d", len(b.Shots), len(other.Shots))
		}
		return false
	}

	for i, shot := range b.Shots {
		if !shot.Equal(other.Shots[i]) {
			if logger != nil {
				logger.Printf("Shot[%d] mismatch: source=%s->%s (result=%t) does not equal source=%s->%s (result=%t)",
					i, shot.Source, shot.Destination, shot.Result,
					other.Shots[i].Source, other.Shots[i].Destination, other.Shots[i].Result)
			}
			return false
		}
	}

	return true
}
