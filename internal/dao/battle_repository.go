package dao

import "glaktika.eu/galaktika/pkg/ship"

type BattleRepository struct {
	battle *ship.Battle
}

func NewBattleRepository() *BattleRepository {
	return &BattleRepository{
		battle: &ship.Battle{
			ID:    "1",
			SideA: &ship.Fleet{
				// TODO ships A1, A2
			},
			SideB: &ship.Fleet{
				// TODO ships B1, B2
			},
			Shots: []*ship.Shot{
				{"A1", "B1", false},
				{"B2", "A2", false},
				{"A2", "B2", true},
				{"B1", "A1", true},
				{"A2", "B1", true},
			},
		},
	}
}

func (r *BattleRepository) GetBattle(battleId string) *ship.Battle {
	return r.battle
}
