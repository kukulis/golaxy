package game

// Side represents which side of the battle
type Side int

const (
	SideA Side = iota
	SideB
)

func (s Side) Flip() Side {
	return s ^ 1
}
