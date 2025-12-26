package game

import (
	"glaktika.eu/galaktika/pkg/galaxy"
	"glaktika.eu/galaktika/pkg/gamemath"
	"testing"
)

func createTestShip(id string, attack float64, guns int, defense float64) galaxy.Ship {
	return galaxy.Ship{
		ID: id,
		Tech: galaxy.ShipTech{
			Attack:  attack,
			Guns:    guns,
			Defense: defense,
		},
	}
}

func TestNewRuntimeDecisionProducer(t *testing.T) {
	rng := gamemath.NewPredefinedRandomGenerator([]float64{0.5})
	battleState := MockReadonlyBattleState{
		AliveShipCount:       []int{2, 2},
		AliveGunnedShipCount: []int{1, 1},
	}

	producer := NewRuntimeDecisionProducer(rng, battleState)

	if producer == nil {
		t.Fatal("NewRuntimeDecisionProducer returned nil")
	}

	if producer.randomGenerator == nil {
		t.Error("randomGenerator not set")
	}

	if producer.destructionFunction == nil {
		t.Error("destructionFunction not set")
	}

	if producer.battleState == nil {
		t.Error("battleState not set")
	}

	if producer.currentSide != SideA {
		t.Errorf("Expected initial side to be SideA (0), got %d", producer.currentSide)
	}

	if producer.shotsMade != 0 {
		t.Errorf("Expected shotsMade to be 0, got %d", producer.shotsMade)
	}
}

func TestProduceNextShotFirstShot(t *testing.T) {
	// Random values: [0.3, 0.6, 0.2, 0.1]
	// 0.3 < 0.5 -> SideA is shooter
	// 0.6 * 2 = 1.2 -> floor = 1 -> shooter at index 1
	// 0.2 * 3 = 0.6 -> floor = 0 -> target at index 0
	// 0.1 -> destruction check
	rng := gamemath.NewPredefinedRandomGenerator([]float64{0.3, 0.6, 0.2, 0.1})

	shooterA1 := createTestShip("ship-a1", 10, 2, 5)
	shooterA2 := createTestShip("ship-a2", 15, 3, 5)
	targetB1 := createTestShip("ship-b1", 5, 1, 8)
	targetB2 := createTestShip("ship-b2", 5, 1, 8)
	targetB3 := createTestShip("ship-b3", 5, 1, 8)

	battleState := MockReadonlyBattleState{
		AliveShipCount:       []int{2, 3},
		AliveGunnedShipCount: []int{2, 2},
		AliveShips: [][]galaxy.Ship{
			{shooterA1, shooterA2},
			{targetB1, targetB2, targetB3},
		},
		AliveGunnedShips: [][]galaxy.Ship{
			{shooterA1, shooterA2},
			{targetB1, targetB2},
		},
	}

	producer := NewRuntimeDecisionProducer(rng, battleState)
	shot := producer.ProduceNextShot()

	if shot == nil {
		t.Fatal("ProduceNextShot returned nil")
	}

	if shot.Side != SideA {
		t.Errorf("Expected side SideA, got %d", shot.Side)
	}

	if shot.ShooterId != "ship-a2" {
		t.Errorf("Expected shooter 'ship-a2', got '%s'", shot.ShooterId)
	}

	if shot.TargetId != "ship-b1" {
		t.Errorf("Expected target 'ship-b1', got '%s'", shot.TargetId)
	}

	if producer.shotsMade != 1 {
		t.Errorf("Expected shotsMade to be 1, got %d", producer.shotsMade)
	}
}

func TestProduceNextShotSideB(t *testing.T) {
	// Random values: [0.6, 0.4, 0.5, 0.2]
	// 0.6 >= 0.5 -> SideB is shooter
	// 0.4 * 2 = 0.8 -> floor = 0 -> shooter at index 0
	// 0.5 * 2 = 1.0 -> floor = 1 -> target at index 1
	// 0.2 -> destruction check
	rng := gamemath.NewPredefinedRandomGenerator([]float64{0.6, 0.4, 0.5, 0.2})

	targetA1 := createTestShip("ship-a1", 5, 1, 5)
	targetA2 := createTestShip("ship-a2", 5, 1, 5)
	shooterB1 := createTestShip("ship-b1", 12, 2, 8)
	shooterB2 := createTestShip("ship-b2", 10, 2, 8)

	battleState := MockReadonlyBattleState{
		AliveShipCount:       []int{2, 2},
		AliveGunnedShipCount: []int{1, 2},
		AliveShips: [][]galaxy.Ship{
			{targetA1, targetA2},
			{shooterB1, shooterB2},
		},
		AliveGunnedShips: [][]galaxy.Ship{
			{targetA1},
			{shooterB1, shooterB2},
		},
	}

	producer := NewRuntimeDecisionProducer(rng, battleState)
	shot := producer.ProduceNextShot()

	if shot.Side != SideB {
		t.Errorf("Expected side SideB, got %d", shot.Side)
	}

	if shot.ShooterId != "ship-b1" {
		t.Errorf("Expected shooter 'ship-b1', got '%s'", shot.ShooterId)
	}

	if shot.TargetId != "ship-a2" {
		t.Errorf("Expected target 'ship-a2', got '%s'", shot.TargetId)
	}
}

func TestProduceNextShotMultipleShotsFromSameShooter(t *testing.T) {
	// Random values for multiple shots:
	// Shot 1: [0.2, 0.0, 0.1, 0.5] -> SideA, shooter[0], target[0], destruction check
	// Shot 2: [0.3, 0.9] -> target[2], destruction check (same shooter)
	// Shot 3: [0.4, 0.8] -> target[1], destruction check (same shooter, exhausted after this)
	rng := gamemath.NewPredefinedRandomGenerator([]float64{
		0.2, 0.0, 0.1, 0.5,
		0.3, 0.9,
		0.4, 0.8,
	})

	shooterA1 := createTestShip("ship-a1", 10, 3, 5) // 3 guns
	targetB1 := createTestShip("ship-b1", 5, 1, 8)
	targetB2 := createTestShip("ship-b2", 5, 1, 8)
	targetB3 := createTestShip("ship-b3", 5, 1, 8)

	battleState := MockReadonlyBattleState{
		AliveShipCount:       []int{1, 3},
		AliveGunnedShipCount: []int{1, 1},
		AliveShips: [][]galaxy.Ship{
			{shooterA1},
			{targetB1, targetB2, targetB3},
		},
		AliveGunnedShips: [][]galaxy.Ship{
			{shooterA1},
			{targetB1},
		},
	}

	producer := NewRuntimeDecisionProducer(rng, battleState)

	// First shot
	shot1 := producer.ProduceNextShot()
	if shot1.ShooterId != "ship-a1" {
		t.Errorf("Shot 1: Expected shooter 'ship-a1', got '%s'", shot1.ShooterId)
	}
	if shot1.TargetId != "ship-b1" {
		t.Errorf("Shot 1: Expected target 'ship-b1', got '%s'", shot1.TargetId)
	}
	if producer.shotsMade != 1 {
		t.Errorf("Shot 1: Expected shotsMade=1, got %d", producer.shotsMade)
	}

	// Second shot - same shooter
	shot2 := producer.ProduceNextShot()
	if shot2.ShooterId != "ship-a1" {
		t.Errorf("Shot 2: Expected shooter 'ship-a1', got '%s'", shot2.ShooterId)
	}
	if shot2.TargetId != "ship-b1" {
		t.Errorf("Shot 2: Expected target 'ship-b1', got '%s'", shot2.TargetId)
	}
	if producer.shotsMade != 2 {
		t.Errorf("Shot 2: Expected shotsMade=2, got %d", producer.shotsMade)
	}

	// Third shot - same shooter
	shot3 := producer.ProduceNextShot()
	if shot3.ShooterId != "ship-a1" {
		t.Errorf("Shot 3: Expected shooter 'ship-a1', got '%s'", shot3.ShooterId)
	}
	if shot3.TargetId != "ship-b2" {
		t.Errorf("Shot 3: Expected target 'ship-b2', got '%s'", shot3.TargetId)
	}
	if producer.shotsMade != 3 {
		t.Errorf("Shot 3: Expected shotsMade=3, got %d", producer.shotsMade)
	}
}

func TestProduceNextShotShooterExhausted(t *testing.T) {
	// Random values:
	// Shot 1: [0.2, 0.0, 0.1, 0.5] -> SideA, shooter[0], target[0]
	// Shot 2: [0.3, 0.9] -> target[0]
	// Shot 3: [0.6, 0.5, 0.2, 0.3] -> new shooter: SideB, shooter[1], target[0]
	rng := gamemath.NewPredefinedRandomGenerator([]float64{
		0.2, 0.0, 0.1, 0.5, // Shot 1
		0.3, 0.9, // Shot 2
		0.6, 0.5, 0.2, 0.3, // Shot 3 - new shooter
	})

	shooterA1 := createTestShip("ship-a1", 10, 2, 5) // 2 guns
	targetA2 := createTestShip("ship-a2", 5, 1, 5)
	shooterB1 := createTestShip("ship-b1", 8, 2, 8)
	shooterB2 := createTestShip("ship-b2", 12, 2, 8)
	targetB3 := createTestShip("ship-b3", 5, 1, 8)

	battleState := MockReadonlyBattleState{
		AliveShipCount:       []int{2, 3},
		AliveGunnedShipCount: []int{1, 2},
		AliveShips: [][]galaxy.Ship{
			{shooterA1, targetA2},
			{shooterB1, shooterB2, targetB3},
		},
		AliveGunnedShips: [][]galaxy.Ship{
			{shooterA1},
			{shooterB1, shooterB2},
		},
	}

	producer := NewRuntimeDecisionProducer(rng, battleState)

	// First shot
	shot1 := producer.ProduceNextShot()
	if shot1.ShooterId != "ship-a1" {
		t.Errorf("Shot 1: Expected shooter 'ship-a1', got '%s'", shot1.ShooterId)
	}

	// Second shot - same shooter
	shot2 := producer.ProduceNextShot()
	if shot2.ShooterId != "ship-a1" {
		t.Errorf("Shot 2: Expected shooter 'ship-a1', got '%s'", shot2.ShooterId)
	}
	if producer.shotsMade != 2 {
		t.Errorf("Shot 2: Expected shotsMade=2, got %d", producer.shotsMade)
	}

	// Third shot - should select new shooter
	shot3 := producer.ProduceNextShot()
	if shot3.ShooterId == "ship-a1" {
		t.Error("Shot 3: Expected different shooter than 'ship-a1'")
	}
	if shot3.Side != SideB {
		t.Errorf("Shot 3: Expected side SideB, got %d", shot3.Side)
	}
	if shot3.ShooterId != "ship-b2" {
		t.Errorf("Shot 3: Expected shooter 'ship-b2', got '%s'", shot3.ShooterId)
	}
	if producer.shotsMade != 1 {
		t.Errorf("Shot 3: Expected shotsMade reset to 1, got %d", producer.shotsMade)
	}
}

func TestProduceNextShotDestructionTrue(t *testing.T) {
	// Random values: [0.3, 0.0, 0.0, 0.01]
	// Destruction check: 0.01 is very low, should result in destruction
	rng := gamemath.NewPredefinedRandomGenerator([]float64{0.3, 0.0, 0.0, 0.01})

	shooterA1 := createTestShip("ship-a1", 10, 2, 5)
	targetB1 := createTestShip("ship-b1", 5, 1, 2) // Low defense

	battleState := MockReadonlyBattleState{
		AliveShipCount:       []int{1, 1},
		AliveGunnedShipCount: []int{1, 1},
		AliveShips: [][]galaxy.Ship{
			{shooterA1},
			{targetB1},
		},
		AliveGunnedShips: [][]galaxy.Ship{
			{shooterA1},
			{targetB1},
		},
	}

	producer := NewRuntimeDecisionProducer(rng, battleState)
	shot := producer.ProduceNextShot()

	if !shot.Destroyed {
		t.Error("Expected shot to result in destruction with high attack vs low defense")
	}
}

func TestProduceNextShotDestructionFalse(t *testing.T) {
	// Random values: [0.3, 0.0, 0.0, 0.99]
	// Destruction check: 0.99 is very high, should result in no destruction
	rng := gamemath.NewPredefinedRandomGenerator([]float64{0.3, 0.0, 0.0, 0.99})

	shooterA1 := createTestShip("ship-a1", 2, 2, 5) // Low attack
	targetB1 := createTestShip("ship-b1", 5, 1, 10) // High defense

	battleState := MockReadonlyBattleState{
		AliveShipCount:       []int{1, 1},
		AliveGunnedShipCount: []int{1, 1},
		AliveShips: [][]galaxy.Ship{
			{shooterA1},
			{targetB1},
		},
		AliveGunnedShips: [][]galaxy.Ship{
			{shooterA1},
			{targetB1},
		},
	}

	producer := NewRuntimeDecisionProducer(rng, battleState)
	shot := producer.ProduceNextShot()

	if shot.Destroyed {
		t.Error("Expected shot to not result in destruction with low attack vs high defense")
	}
}

func TestProduceNextShotTargetSelection(t *testing.T) {
	// Test that target index correctly wraps with floor calculation
	// Random values: [0.2, 0.0, 0.99, 0.5]
	// 0.99 * 3 = 2.97 -> floor = 2 -> target at index 2
	rng := gamemath.NewPredefinedRandomGenerator([]float64{0.2, 0.0, 0.99, 0.5})

	shooterA1 := createTestShip("ship-a1", 10, 2, 5)
	targetB1 := createTestShip("ship-b1", 5, 1, 8)
	targetB2 := createTestShip("ship-b2", 5, 1, 8)
	targetB3 := createTestShip("ship-b3", 5, 1, 8)

	battleState := MockReadonlyBattleState{
		AliveShipCount:       []int{1, 3},
		AliveGunnedShipCount: []int{1, 1},
		AliveShips: [][]galaxy.Ship{
			{shooterA1},
			{targetB1, targetB2, targetB3},
		},
		AliveGunnedShips: [][]galaxy.Ship{
			{shooterA1},
			{targetB1},
		},
	}

	producer := NewRuntimeDecisionProducer(rng, battleState)
	shot := producer.ProduceNextShot()

	if shot.TargetId != "ship-b3" {
		t.Errorf("Expected target 'ship-b3', got '%s'", shot.TargetId)
	}
}

func TestProduceNextShotShooterWithOneGun(t *testing.T) {
	// Test shooter with only 1 gun - should change after each shot
	// Random values:
	// Shot 1: [0.2, 0.0, 0.0, 0.5] -> SideA, shooter[0], target[0]
	// Shot 2: [0.7, 0.0, 0.0, 0.5] -> SideB, shooter[0], target[0]
	rng := gamemath.NewPredefinedRandomGenerator([]float64{
		0.2, 0.0, 0.0, 0.5, // Shot 1
		0.7, 0.0, 0.0, 0.5, // Shot 2
	})

	shooterA1 := createTestShip("ship-a1", 10, 1, 5) // Only 1 gun
	shooterB1 := createTestShip("ship-b1", 10, 1, 8) // Only 1 gun

	battleState := MockReadonlyBattleState{
		AliveShipCount:       []int{1, 1},
		AliveGunnedShipCount: []int{1, 1},
		AliveShips: [][]galaxy.Ship{
			{shooterA1},
			{shooterB1},
		},
		AliveGunnedShips: [][]galaxy.Ship{
			{shooterA1},
			{shooterB1},
		},
	}

	producer := NewRuntimeDecisionProducer(rng, battleState)

	// First shot
	shot1 := producer.ProduceNextShot()
	if shot1.ShooterId != "ship-a1" {
		t.Errorf("Shot 1: Expected shooter 'ship-a1', got '%s'", shot1.ShooterId)
	}
	if producer.shotsMade != 1 {
		t.Errorf("Shot 1: Expected shotsMade=1, got %d", producer.shotsMade)
	}

	// Second shot - should select new shooter (gun exhausted)
	shot2 := producer.ProduceNextShot()
	if shot2.ShooterId != "ship-b1" {
		t.Errorf("Shot 2: Expected shooter 'ship-b1', got '%s'", shot2.ShooterId)
	}
	if producer.shotsMade != 1 {
		t.Errorf("Shot 2: Expected shotsMade reset to 1, got %d", producer.shotsMade)
	}
}
