package game

import (
	"glaktika.eu/galaktika/pkg/galaxy"
	"glaktika.eu/galaktika/pkg/util"
	"maps"
	"slices"
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

	bh.poolA = util.NewIndexMapPool(util.ArrayMap(shipsA, func(ship *galaxy.Ship) string { return ship.ID }))
	bh.poolB = util.NewIndexMapPool(util.ArrayMap(shipsB, func(ship *galaxy.Ship) string { return ship.ID }))
	bh.gunnedPoolA = util.NewIndexMapPool(util.ArrayMap(gunnedShipsA, func(ship *galaxy.Ship) string { return ship.ID }))
	bh.gunnedPoolB = util.NewIndexMapPool(util.ArrayMap(gunnedShipsB, func(ship *galaxy.Ship) string { return ship.ID }))

	bh.shipsMapA = make(map[string]*galaxy.Ship)
	for _, ship := range fleetA.Ships {
		bh.shipsMapA[ship.ID] = ship
	}

	bh.shipsMapB = make(map[string]*galaxy.Ship)
	for _, ship := range fleetA.Ships {
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

		shot := galaxy.Shot{
			Source:      shotDecision.ShooterId,
			Destination: shotDecision.TargetId,
			Result:      shotDecision.Destroyed,
		}

		if shotDecision.Destroyed {
			if shotDecision.Side == 0 {
				bh.shipsMapB[shotDecision.TargetId].Destroyed = true
				bh.poolB.RemoveKey(shotDecision.TargetId)
				bh.gunnedPoolB.RemoveKey(shotDecision.TargetId)
			} else {
				bh.shipsMapA[shotDecision.TargetId].Destroyed = true
				bh.poolA.RemoveKey(shotDecision.TargetId)
				bh.gunnedPoolA.RemoveKey(shotDecision.TargetId)
			}
		}

		battle.Shots = append(battle.Shots, &shot)

		if bh.IsBattleOver() {
			break
		}
	}

	battle.PostSideA = galaxy.NewFleet(slices.Collect(maps.Values(bh.shipsMapA)))
	battle.PostSideB = galaxy.NewFleet(slices.Collect(maps.Values(bh.shipsMapB)))

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
