package game

import (
	"glaktika.eu/galaktika/pkg/gamemath"
	"math"
)

type BattleHandler struct {
	destructionFunction *gamemath.ConfigurableFunction
}

func NewBattleHandler() *BattleHandler {
	f, err := gamemath.NewConfigurableFunction([]float64{0.25, 1, 4}, []float64{1, 0.5, 0})
	if err != nil {
		panic(err)
	}
	return &BattleHandler{
		destructionFunction: f,
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
