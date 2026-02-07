package dao

import "glaktika.eu/galaktika/pkg/galaxy"

type BattleRepository struct {
	battle *galaxy.Battle
}

func NewBattleRepository() *BattleRepository {
	return &BattleRepository{
		battle: &galaxy.Battle{
			ID: "1",
			SideA: galaxy.NewFleet([]*galaxy.Ship{
				{
					ID: "A1",
					Tech: galaxy.ShipTech{
						Attack:        5,
						Guns:          2,
						Defense:       3,
						Speed:         10,
						CargoCapacity: 50,
						Mass:          100,
					},
					Destroyed: false,
					Name:      "Cruiser Alpha",
					Owner:     "race_a",
				},
				{
					ID: "A2",
					Tech: galaxy.ShipTech{
						Attack:        4,
						Guns:          2,
						Defense:       4,
						Speed:         8,
						CargoCapacity: 40,
						Mass:          90,
					},
					Destroyed: false,
					Name:      "Destroyer Alpha",
					Owner:     "race_a",
				},
			}),
			SideB: galaxy.NewFleet([]*galaxy.Ship{
				{
					ID: "B1",
					Tech: galaxy.ShipTech{
						Attack:        6,
						Guns:          3,
						Defense:       2,
						Speed:         12,
						CargoCapacity: 45,
						Mass:          95,
					},
					Destroyed: false,
					Name:      "Cruiser Beta",
					Owner:     "race_b",
				},
				{
					ID: "B1_2",
					Tech: galaxy.ShipTech{
						Attack:        6,
						Guns:          3,
						Defense:       2,
						Speed:         12,
						CargoCapacity: 45,
						Mass:          95,
					},
					Destroyed: false,
					Name:      "Cruiser Beta",
					Owner:     "race_b",
				},
				{
					ID: "B1_3",
					Tech: galaxy.ShipTech{
						Attack:        6,
						Guns:          3,
						Defense:       2,
						Speed:         12,
						CargoCapacity: 45,
						Mass:          95,
					},
					Destroyed: false,
					Name:      "Cruiser Beta",
					Owner:     "race_b",
				},
				{
					ID: "B2",
					Tech: galaxy.ShipTech{
						Attack:        5,
						Guns:          2,
						Defense:       3,
						Speed:         9,
						CargoCapacity: 55,
						Mass:          105,
					},
					Destroyed: false,
					Name:      "Destroyer Beta",
					Owner:     "race_b",
				},
			}),
			Shots: []*galaxy.Shot{
				{"A1", "B1", false},
				{"B2", "A2", false},
				{"A2", "B2", true},
				{"B1", "A1", true},
				{"A2", "B1", true},
				{"A2", "B1_2", true},
				{"A2", "B1_3", false},
			},
		},
	}
}

func (r *BattleRepository) GetBattle(battleId string) *galaxy.Battle {
	return r.battle
}
