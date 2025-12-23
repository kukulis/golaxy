package galaxy

import (
	"glaktika.eu/galaktika/pkg/gamemath"
	"glaktika.eu/galaktika/pkg/util"
)

// @deprecated make implementation, using interfaces; example at decision_producer_proposition.txt

// Side represents which side of the battle
type Side int

const (
	SideA Side = iota
	SideB
)

// DecisionProducer produces battle shot decisions
type DecisionProducer interface {
	// ProduceNextShot produces the next shot decision based on current fleet state
	// Returns nil when battle should end
	ProduceNextShot(postSideA, postSideB *Fleet) *ShotDecision
}

// ShotDecision represents a single shot decision
type ShotDecision struct {
	Side      Side
	ShooterID string
	TargetID  string
	Destroys  bool
}

// PredefinedDecisionProducer returns predefined shots for testing
type PredefinedDecisionProducer struct {
	shots []ShotDecision
	index int
}

// NewPredefinedDecisionProducer creates a new predefined decision producer
func NewPredefinedDecisionProducer(shots []ShotDecision) *PredefinedDecisionProducer {
	return &PredefinedDecisionProducer{
		shots: shots,
		index: 0,
	}
}

// ProduceNextShot returns the next predefined shot
func (p *PredefinedDecisionProducer) ProduceNextShot(postSideA, postSideB *Fleet) *ShotDecision {
	if p.index >= len(p.shots) {
		return nil
	}
	shot := p.shots[p.index]
	p.index++
	return &shot
}

// RuntimeDecisionProducer produces battle decisions using randomness
type RuntimeDecisionProducer struct {
	randomGenerator     gamemath.RandomGenerator
	destructionFunction *gamemath.ConfigurableFunction

	// State for current shooter
	currentSide      *Side // nil means: select new shooter
	currentShooterID string
	remainingShots   int
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
		currentSide:         nil, // Start with no shooter selected
	}
}

// ProduceNextShot produces the next shot decision based on current fleet state
func (r *RuntimeDecisionProducer) ProduceNextShot(postSideA, postSideB *Fleet) *ShotDecision {
	// Create index pools for both sides
	poolA := r.createAliveShipsPool(postSideA.Ships)
	poolB := r.createAliveShipsPool(postSideB.Ships)

	// Battle ends if either side has no alive ships
	if poolA.Count() == 0 || poolB.Count() == 0 {
		return nil
	}

	// If no current shooter, select a new one
	if r.currentSide == nil {
		r.selectNewShooter(postSideA, postSideB, poolA, poolB)
	}

	// Determine active and enemy fleets and pools
	var activeFleet, enemyFleet *Fleet
	var enemyPool *util.IndexPool
	if *r.currentSide == SideA {
		activeFleet = postSideA
		enemyFleet = postSideB
		enemyPool = poolB
	} else {
		activeFleet = postSideB
		enemyFleet = postSideA
		enemyPool = poolA
	}

	// Find shooter
	shooter := findShipByID(activeFleet.Ships, r.currentShooterID)
	if shooter == nil || shooter.Destroyed {
		// Shooter was destroyed, select new shooter
		r.currentSide = nil
		return r.ProduceNextShot(postSideA, postSideB)
	}

	// Select target (evenly distributed among alive ships)
	targetIdx := enemyPool.GetRandom(r.randomGenerator.NextRandom())
	target := enemyFleet.Ships[targetIdx]

	// Calculate if target is destroyed based on attack/defense ratio
	probability := r.destructionFunction.CalculateRatio(
		target.Tech.Defense,
		shooter.Tech.Attack,
	)
	destroys := r.randomGenerator.NextRandom() < probability

	// Decrement remaining shots
	r.remainingShots--
	if r.remainingShots <= 0 {
		// Shooter is done, reset for next shooter selection
		r.currentSide = nil
	}

	return &ShotDecision{
		Side:      *r.currentSide,
		ShooterID: shooter.ID,
		TargetID:  target.ID,
		Destroys:  destroys,
	}
}

// selectNewShooter selects a new shooter and initializes state
func (r *RuntimeDecisionProducer) selectNewShooter(postSideA, postSideB *Fleet, poolA, poolB *util.IndexPool) {
	// Select which side fires (50/50)
	side := r.selectSide()
	r.currentSide = &side

	// Determine active fleet and pool
	var activeFleet *Fleet
	var activePool *util.IndexPool
	if side == SideA {
		activeFleet = postSideA
		activePool = poolA
	} else {
		activeFleet = postSideB
		activePool = poolB
	}

	// Select shooter from active side (evenly distributed)
	shooterIdx := activePool.GetRandom(r.randomGenerator.NextRandom())
	shooter := activeFleet.Ships[shooterIdx]

	r.currentShooterID = shooter.ID
	r.remainingShots = int(shooter.Tech.Guns)
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

// findShipByID finds a ship in a slice by its ID
func findShipByID(ships []*Ship, id string) *Ship {
	for _, ship := range ships {
		if ship.ID == id {
			return ship
		}
	}
	return nil
}
