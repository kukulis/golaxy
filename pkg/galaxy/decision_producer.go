package galaxy

import (
	"glaktika.eu/galaktika/pkg/gamemath"
	"glaktika.eu/galaktika/pkg/util"
)

// Side represents which side of the battle
type Side int

const (
	SideA Side = iota
	SideB
)

// DecisionProducer produces battle turn decisions
type DecisionProducer interface {
	// ProduceNextTurn produces the next turn decision based on current fleet state
	// Returns nil when battle should end
	ProduceNextTurn(postSideA, postSideB *Fleet) *TurnDecision
}

// TurnDecision represents a single turn in the battle
type TurnDecision struct {
	Side      Side
	ShooterID string
	Shots     []ShotDecision
}

// ShotDecision represents a single shot decision
type ShotDecision struct {
	TargetID string
	Destroys bool
}

// PredefinedDecisionProducer returns predefined turns for testing
type PredefinedDecisionProducer struct {
	turns []TurnDecision
	index int
}

// NewPredefinedDecisionProducer creates a new predefined decision producer
func NewPredefinedDecisionProducer(turns []TurnDecision) *PredefinedDecisionProducer {
	return &PredefinedDecisionProducer{
		turns: turns,
		index: 0,
	}
}

// ProduceNextTurn returns the next predefined turn
func (p *PredefinedDecisionProducer) ProduceNextTurn(postSideA, postSideB *Fleet) *TurnDecision {
	if p.index >= len(p.turns) {
		return nil
	}
	turn := p.turns[p.index]
	p.index++
	return &turn
}

// RuntimeDecisionProducer produces battle decisions using randomness
type RuntimeDecisionProducer struct {
	randomGenerator     gamemath.RandomGenerator
	destructionFunction *gamemath.ConfigurableFunction
}

// NewRuntimeDecisionProducer creates a new runtime decision producer
func NewRuntimeDecisionProducer(rng gamemath.RandomGenerator) *RuntimeDecisionProducer {
	f, err := gamemath.NewConfigurableFunction([]float64{0.25, 1, 4}, []float64{1, 0.5, 0})
	if err != nil {
		panic(err)
	}
	return &RuntimeDecisionProducer{
		randomGenerator:     rng,
		destructionFunction: f,
	}
}

// ProduceNextTurn produces a turn decision based on current fleet state and randomness
func (r *RuntimeDecisionProducer) ProduceNextTurn(postSideA, postSideB *Fleet) *TurnDecision {
	// Create index pools for both sides
	poolA := r.createAliveShipsPool(postSideA.Ships)
	poolB := r.createAliveShipsPool(postSideB.Ships)

	// Battle ends if either side has no alive ships
	if poolA.Count() == 0 || poolB.Count() == 0 {
		return nil
	}

	// Select which side fires (50/50)
	side := r.selectSide()

	// Determine active and enemy fleets and pools
	var activeFleet, enemyFleet *Fleet
	var activePool, enemyPool *util.IndexPool
	if side == SideA {
		activeFleet = postSideA
		enemyFleet = postSideB
		activePool = poolA
		enemyPool = poolB
	} else {
		activeFleet = postSideB
		enemyFleet = postSideA
		activePool = poolB
		enemyPool = poolA
	}

	// Select shooter from active side (evenly distributed)
	shooterIdx := activePool.GetRandom(r.randomGenerator.NextRandom())
	shooter := activeFleet.Ships[shooterIdx]

	// Generate shots (one per gun)
	numGuns := int(shooter.Tech.Guns)
	shots := make([]ShotDecision, numGuns)

	for i := 0; i < numGuns; i++ {
		// Select target (evenly distributed)
		targetIdx := enemyPool.GetRandom(r.randomGenerator.NextRandom())
		target := enemyFleet.Ships[targetIdx]

		// Calculate if target is destroyed based on attack/defense ratio
		probability := r.destructionFunction.CalculateRatio(
			target.Tech.Defense,
			shooter.Tech.Attack,
		)
		destroys := r.randomGenerator.NextRandom() < probability

		shots[i] = ShotDecision{
			TargetID: target.ID,
			Destroys: destroys,
		}
	}

	return &TurnDecision{
		Side:      side,
		ShooterID: shooter.ID,
		Shots:     shots,
	}
}

// selectSide randomly selects which side fires (50/50 probability)
func (r *RuntimeDecisionProducer) selectSide() Side {
	if r.randomGenerator.NextRandom() < 0.5 {
		return SideA
	}
	return SideB
}

// createAliveShipsPool creates an IndexPool with only alive ships
func (r *RuntimeDecisionProducer) createAliveShipsPool(ships []*Ship) *util.IndexPool {
	pool := util.NewIndexPool(len(ships))
	for i, ship := range ships {
		if ship.Destroyed {
			pool.Remove(i)
		}
	}
	return pool
}
