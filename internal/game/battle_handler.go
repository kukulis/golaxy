package game

import (
	"glaktika.eu/galaktika/pkg/galaxy"
	"glaktika.eu/galaktika/pkg/util"
)

type ShipRef struct {
	id    string
	index uint
}

type BattleHandler struct {
	decisionProducer DecisionProducerInterface
	idGenerator      util.IdGenerator

	// Battle state (implements BattleState interface)
	shipsMapA   map[string]*galaxy.Ship
	shipsMapB   map[string]*galaxy.Ship
	shipsA      []*galaxy.Ship // Preserves original order
	shipsB      []*galaxy.Ship // Preserves original order
	poolA       *util.IndexMapPool
	poolB       *util.IndexMapPool
	gunnedPoolA *util.IndexMapPool
	gunnedPoolB *util.IndexMapPool
}

func (bh *BattleHandler) initializeBattleState(fleetA *galaxy.Fleet, fleetB *galaxy.Fleet) {
	shipsA := copyShips(fleetA.Ships)
	shipsB := copyShips(fleetB.Ships)
	gunnedShipsA := util.ArrayFilter(shipsA, func(ship *galaxy.Ship) bool { return ship.Tech.Guns > 0 })
	gunnedShipsB := util.ArrayFilter(shipsB, func(ship *galaxy.Ship) bool { return ship.Tech.Guns > 0 })

	bh.shipsA = shipsA
	bh.shipsB = shipsB

	bh.poolA = util.NewIndexMapPool(util.ArrayMap(shipsA, func(ship *galaxy.Ship) string { return ship.ID }))
	bh.poolB = util.NewIndexMapPool(util.ArrayMap(shipsB, func(ship *galaxy.Ship) string { return ship.ID }))
	bh.gunnedPoolA = util.NewIndexMapPool(util.ArrayMap(gunnedShipsA, func(ship *galaxy.Ship) string { return ship.ID }))
	bh.gunnedPoolB = util.NewIndexMapPool(util.ArrayMap(gunnedShipsB, func(ship *galaxy.Ship) string { return ship.ID }))

	bh.shipsMapA = make(map[string]*galaxy.Ship)
	for _, ship := range shipsA {
		bh.shipsMapA[ship.ID] = ship
	}

	bh.shipsMapB = make(map[string]*galaxy.Ship)
	for _, ship := range shipsB {
		bh.shipsMapB[ship.ID] = ship
	}
}

func (bh *BattleHandler) GetAliveShipCount(side Side) int {
	if side == 0 {
		return bh.poolA.Count()
	}

	return bh.poolB.Count()
}

func (bh *BattleHandler) GetAliveGunnedShipCount(side Side) int {
	if side == 0 {
		return bh.gunnedPoolA.Count()
	}

	return bh.gunnedPoolB.Count()
}

func (bh *BattleHandler) GetShipAt(side Side, position int) galaxy.Ship {
	if side == 0 {
		key := bh.poolA.GetKey(position)
		return *bh.shipsMapA[key]
	}

	key := bh.poolB.GetKey(position)
	return *bh.shipsMapB[key]
}

func (bh *BattleHandler) GetGunnedShipAt(side Side, position int) galaxy.Ship {
	if side == 0 {
		key := bh.gunnedPoolA.GetKey(position)
		return *bh.shipsMapA[key]
	}

	key := bh.gunnedPoolB.GetKey(position)
	return *bh.shipsMapB[key]
}

func (bh *BattleHandler) IsBattleOver() bool {
	return bh.poolA.Count() == 0 || bh.poolB.Count() == 0 || (bh.gunnedPoolA.Count() == 0 && bh.gunnedPoolB.Count() == 0)
}

func NewBattleHandler(
	idGenerator util.IdGenerator,
	decisionProducer DecisionProducerInterface,
) *BattleHandler {
	return &BattleHandler{
		decisionProducer: decisionProducer,
		idGenerator:      idGenerator,
	}
}

func (bh *BattleHandler) ExecuteBattle(fleetA *galaxy.Fleet, fleetB *galaxy.Fleet) *galaxy.Battle {

	if bh.decisionProducer == nil {
		panic("BattleHandler.ExecuteBattle: The decision producer is nil ")
	}

	bh.initializeBattleState(fleetA, fleetB)

	battle := galaxy.Battle{
		ID:    bh.idGenerator.NextId(),
		SideA: fleetA,
		SideB: fleetB,
	}

	maxShots := 1000

	for i := 0; i < maxShots; i++ {
		shotDecision := bh.decisionProducer.ProduceNextShot()
		if shotDecision == nil {
			break
		}

		shot := galaxy.Shot{
			Source:      shotDecision.ShooterId,
			Destination: shotDecision.TargetId,
			Result:      shotDecision.Destroyed,
		}

		if shotDecision.Destroyed {
			if shotDecision.Side == 0 {
				bh.shipsMapB[shotDecision.TargetId].Destroyed = true
				if err := bh.poolB.RemoveKey(shotDecision.TargetId); err != nil {
					panic("BUG: failed to remove destroyed ship from pool: " + err.Error())
				}
				if err := bh.gunnedPoolB.RemoveKey(shotDecision.TargetId); err != nil {
					// Ignore error - ship might not have guns
					_ = err
				}
			} else {
				bh.shipsMapA[shotDecision.TargetId].Destroyed = true
				if err := bh.poolA.RemoveKey(shotDecision.TargetId); err != nil {
					panic("BUG: failed to remove destroyed ship from pool: " + err.Error())
				}
				if err := bh.gunnedPoolA.RemoveKey(shotDecision.TargetId); err != nil {
					// Ignore error - ship might not have guns
					_ = err
				}
			}
		}

		battle.Shots = append(battle.Shots, &shot)

		if bh.IsBattleOver() {
			break
		}
	}

	battle.PostSideA = galaxy.NewFleet(bh.shipsA)
	battle.PostSideB = galaxy.NewFleet(bh.shipsB)

	return &battle
}

// copyShips creates a deep copy of a ship slice
func copyShips(ships []*galaxy.Ship) []*galaxy.Ship {
	copies := make([]*galaxy.Ship, len(ships))
	for i, ship := range ships {
		shipCopy := *ship
		copies[i] = &shipCopy
	}
	return copies
}
