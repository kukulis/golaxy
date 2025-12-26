package galaxy

const ONE_TECH_RESOURCES = 100

type Technologies struct {
	Attack  float64
	Defense float64
	Engine  float64
	Cargo   float64
}

func NewTechnologies() *Technologies {
	return &Technologies{
		Attack:  1.0,
		Defense: 1.0,
		Engine:  1.0,
		Cargo:   1.0,
	}
}

func (t *Technologies) Research(attackResources, defenseResources, engineResources, cargoResources float64) {
	t.Attack += attackResources / ONE_TECH_RESOURCES
	t.Defense += defenseResources / ONE_TECH_RESOURCES
	t.Engine += engineResources / ONE_TECH_RESOURCES
	t.Cargo += cargoResources / ONE_TECH_RESOURCES
}
