package game

import (
	"glaktika.eu/galaktika/pkg/galaxy"
	"glaktika.eu/galaktika/pkg/gamemath"
	"glaktika.eu/galaktika/pkg/util"
	"math"
)

type BattleHandler struct {
	destructionFunction *gamemath.ConfigurableFunction
	idGenerator         util.IdGenerator
	randomGenerator     gamemath.RandomGenerator
}

func NewBattleHandler(
	idGenerator util.IdGenerator,
	randomGenerator gamemath.RandomGenerator,
) *BattleHandler {
	f, err := gamemath.NewConfigurableFunction([]float64{0.25, 1, 4}, []float64{1, 0.5, 0})
	if err != nil {
		panic(err)
	}
	return &BattleHandler{
		destructionFunction: f,
		idGenerator:         idGenerator,
		randomGenerator:     randomGenerator,
	}
}

func (bh *BattleHandler) EvaluateShotResult(defence float64, attack float64, rGenerator gamemath.RandomGenerator) bool {
	probability := bh.destructionFunction.CalculateRatio(defence, attack)
	random := rGenerator.NextRandom()

	return probability > random
}

func EvaluateTargetIndex(targetsCount int, randomGenerator gamemath.RandomGenerator) int {
	return int(math.Floor(randomGenerator.NextRandom() * float64(targetsCount)))
}

func (bh *BattleHandler) ExecuteBattle(fleetA *galaxy.Fleet, fleetB *galaxy.Fleet) *galaxy.Battle {

	battle := galaxy.Battle{
		ID:    bh.idGenerator.NextId(),
		SideA: fleetA,
		SideB: fleetB,
	}

	// Create PostSideA with copies of all ships from SideA
	battle.PostSideA = &galaxy.Fleet{
		Ships: copyShips(fleetA.Ships),
		Owner: fleetA.Owner,
	}

	// Create PostSideB with copies of all ships from SideB
	battle.PostSideB = &galaxy.Fleet{
		Ships: copyShips(fleetB.Ships),
		Owner: fleetB.Owner,
	}

	// TODO generate shots
	// TODO use IndexPool

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
