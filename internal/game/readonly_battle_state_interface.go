package game

import "glaktika.eu/galaktika/pkg/galaxy"

type ReadonlyBattleStateInterface interface {
	GetAliveShipCount(side Side) int
	GetAliveGunnedShipCount(side Side) int
	GetShipAt(side Side, position int) galaxy.Ship
	GetGunnedShipAt(side Side, position int) galaxy.Ship
	IsBattleOver() bool
}
