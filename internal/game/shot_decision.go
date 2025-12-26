package game

// ShotDecision represents a single shot decision
type ShotDecision struct {
	Side      Side
	ShooterId string
	TargetId  string
	Destroyed bool
}
