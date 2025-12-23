package galaxy

type Battle struct {
	ID    string `json:"id"`
	SideA *Fleet `json:"side_a"`
	SideB *Fleet `json:"side_b"`

	Shots     []*Shot `json:"shots"`
	PostSideA *Fleet  `json:"post_side_a"`
	PostSideB *Fleet  `json:"post_side_b"`
}
