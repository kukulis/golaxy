package galaxy

type Division struct {
	ID              string `json:"id"`
	ResourcesAmount int    `json:"resources_amount"`
	TechAttack      int    `json:"tech_attack"`
	TechDefense     int    `json:"tech_defense"`
	TechEngines     int    `json:"tech_engines"`
	TechCargo       int    `json:"tech_cargo"`
}
