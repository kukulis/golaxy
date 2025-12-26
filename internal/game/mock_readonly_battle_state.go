package game

import "glaktika.eu/galaktika/pkg/galaxy"

type MockReadonlyBattleState struct {
	AliveShipCount       []int
	AliveGunnedShipCount []int
	AliveShips           [][]galaxy.Ship
	AliveGunnedShips     [][]galaxy.Ship
	BattleOver           bool
}

func (m MockReadonlyBattleState) GetAliveShipCount(side Side) int {
	return m.AliveShipCount[side]
}

func (m MockReadonlyBattleState) GetAliveGunnedShipCount(side Side) int {
	return m.AliveGunnedShipCount[side]
}

func (m MockReadonlyBattleState) GetShipAt(side Side, position int) galaxy.Ship {
	return m.AliveShips[side][position]
}

func (m MockReadonlyBattleState) GetGunnedShipAt(side Side, position int) galaxy.Ship {
	return m.AliveGunnedShips[side][position]
}

func (m MockReadonlyBattleState) IsBattleOver() bool {
	return m.BattleOver
}
