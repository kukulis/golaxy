package galaxy

const RolePlayer = "player"
const RoleAdmin = "admin"

type Race struct {
	ID    string
	Name  string
	Role  string
	Token string
}
