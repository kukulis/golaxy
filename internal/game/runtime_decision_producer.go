package game

import (
	"glaktika.eu/galaktika/pkg/galaxy"
	"glaktika.eu/galaktika/pkg/gamemath"
	"math"
)

// RuntimeDecisionProducer produces battle decisions using randomness
type RuntimeDecisionProducer struct {
	randomGenerator     gamemath.RandomGenerator
	destructionFunction *gamemath.ConfigurableFunction
	battleState         ReadonlyBattleStateInterface

	// State for current shooter
	currentSide Side
	shooter     galaxy.Ship
	target      galaxy.Ship
	shotsMade   int
}

// NewRuntimeDecisionProducer creates a new runtime decision producer
func NewRuntimeDecisionProducer(
	rng gamemath.RandomGenerator,
	battleState ReadonlyBattleStateInterface,
) *RuntimeDecisionProducer {

	f, err := gamemath.NewConfigurableFunction([]float64{0.25, 1, 4}, []float64{1, 0.5, 0})
	if err != nil {
		panic(err)
	}
	return &RuntimeDecisionProducer{
		randomGenerator:     rng,
		destructionFunction: f,
		currentSide:         0,
		battleState:         battleState,
	}
}

// ProduceNextShot produces the next shot decision based on current fleet state
func (r *RuntimeDecisionProducer) ProduceNextShot() *ShotDecision {
	if r.shooter.Tech.Guns <= r.shotsMade {
		// select new shooter
		r.currentSide = 1
		if r.randomGenerator.NextRandom() < 0.5 {
			r.currentSide = 0
		}

		shooterIndex := int(math.Floor(r.randomGenerator.NextRandom() * float64(r.battleState.GetAliveGunnedShipCount(r.currentSide))))
		r.shooter = r.battleState.GetGunnedShipAt(r.currentSide, shooterIndex)
		r.shotsMade = 0
	}

	targetIndex := int(math.Floor(r.randomGenerator.NextRandom() * float64(r.battleState.GetAliveShipCount(r.currentSide.Flip()))))
	r.target = r.battleState.GetShipAt(r.currentSide, targetIndex)

	destroyed := r.randomGenerator.NextRandom() < r.destructionFunction.CalculateRatio(r.target.Tech.Defense, r.shooter.Tech.Attack)

	shotDecision := &ShotDecision{
		Side:      r.currentSide,
		ShooterId: r.shooter.ID,
		TargetId:  r.target.ID,
		Destroyed: destroyed,
	}
	r.shotsMade++

	return shotDecision
}
