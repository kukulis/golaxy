package game

import (
	"glaktika.eu/galaktika/pkg/galaxy"
	"glaktika.eu/galaktika/pkg/gamemath"
	"glaktika.eu/galaktika/pkg/util"
	"testing"
)

func TestBattleHandlerWithRuntimeDecisionProducer(t *testing.T) {
	tests := []struct {
		name            string
		fleetA          *galaxy.Fleet
		fleetB          *galaxy.Fleet
		randomValues    []float64
		expectedShots   int
		expectedSideA   func(*galaxy.Fleet) bool // Validation function
		expectedSideB   func(*galaxy.Fleet) bool
		battleOverEarly bool // True if battle ends before max shots
	}{
		{
			name: "Two fleets each with one ship without guns - no battle",
			fleetA: &galaxy.Fleet{
				Ships: []*galaxy.Ship{
					{
						ID:   "ship-a1",
						Name: "Cargo A",
						Tech: galaxy.ShipTech{
							Attack:  0,
							Guns:    0, // No guns
							Defense: 5,
						},
						Owner: "race-a",
					},
				},
				Owner: "race-a",
			},
			fleetB: &galaxy.Fleet{
				Ships: []*galaxy.Ship{
					{
						ID:   "ship-b1",
						Name: "Cargo B",
						Tech: galaxy.ShipTech{
							Attack:  0,
							Guns:    0, // No guns
							Defense: 5,
						},
						Owner: "race-b",
					},
				},
				Owner: "race-b",
			},
			randomValues:  []float64{0.5}, // Doesn't matter, no gunned ships
			expectedShots: 0,
			expectedSideA: func(f *galaxy.Fleet) bool {
				return len(f.Ships) == 1 && !f.Ships[0].Destroyed
			},
			expectedSideB: func(f *galaxy.Fleet) bool {
				return len(f.Ships) == 1 && !f.Ships[0].Destroyed
			},
			battleOverEarly: true,
		},
		{
			name: "One fleet has ship with gun, other has no guns - one-sided battle",
			fleetA: &galaxy.Fleet{
				Ships: []*galaxy.Ship{
					{
						ID:   "ship-a1",
						Name: "Fighter A",
						Tech: galaxy.ShipTech{
							Attack:  2, // Equal attack/defense for 50% destruction probability
							Guns:    2, // 2 guns for multiple shots
							Defense: 5,
						},
						Owner: "race-a",
					},
				},
				Owner: "race-a",
			},
			fleetB: &galaxy.Fleet{
				Ships: []*galaxy.Ship{
					{
						ID:   "ship-b1",
						Name: "Cargo B",
						Tech: galaxy.ShipTech{
							Attack:  0,
							Guns:    0, // No guns
							Defense: 2, // Equal to attacker's attack for 50% prob
						},
						Owner: "race-b",
					},
				},
				Owner: "race-b",
			},
			randomValues: []float64{
				// First attempt: select side B (has no guns, returns nil)
				0.6, // Side B (>= 0.5)
				// No shooter selection happens because side B has no gunned ships

				// Second attempt: select side A
				0.2,  // Side A (< 0.5)
				0.0,  // Shooter index 0 (ship-a1)
				0.0,  // Target index 0 (ship-b1)
				0.99, // Destruction: HIGH value = no destruction (ship survives)

				// Shot 2: same shooter (ship-a1 has 2 guns)
				0.0,  // Target index 0 (ship-b1)
				0.01, // Destruction: LOW value = destruction (ship destroyed)
			},
			expectedShots: 2,
			expectedSideA: func(f *galaxy.Fleet) bool {
				return len(f.Ships) == 1 && !f.Ships[0].Destroyed
			},
			expectedSideB: func(f *galaxy.Fleet) bool {
				return len(f.Ships) == 1 && f.Ships[0].Destroyed
			},
			battleOverEarly: true,
		},
		{
			name: "Two fleets with two ships each - one with guns, one without",
			fleetA: &galaxy.Fleet{
				Ships: []*galaxy.Ship{
					{
						ID:   "ship-a1",
						Name: "Fighter A",
						Tech: galaxy.ShipTech{
							Attack:  15,
							Guns:    2,
							Defense: 5,
						},
						Owner: "race-a",
					},
					{
						ID:   "ship-a2",
						Name: "Cargo A",
						Tech: galaxy.ShipTech{
							Attack:  0,
							Guns:    0, // No guns
							Defense: 3,
						},
						Owner: "race-a",
					},
				},
				Owner: "race-a",
			},
			fleetB: &galaxy.Fleet{
				Ships: []*galaxy.Ship{
					{
						ID:   "ship-b1",
						Name: "Fighter B",
						Tech: galaxy.ShipTech{
							Attack:  12,
							Guns:    2,
							Defense: 4,
						},
						Owner: "race-b",
					},
					{
						ID:   "ship-b2",
						Name: "Cargo B",
						Tech: galaxy.ShipTech{
							Attack:  0,
							Guns:    0, // No guns
							Defense: 3,
						},
						Owner: "race-b",
					},
				},
				Owner: "race-b",
			},
			randomValues: []float64{
				// Shot 1: A shoots
				0.2,  // Side A
				0.0,  // Shooter index 0 (ship-a1)
				0.8,  // Target index 1 (ship-b2) - 0.8 * 2 = 1.6 -> floor = 1
				0.01, // Destruction
				// Shot 2: A shoots again (same shooter, has 2 guns)
				0.1,  // Target index 0 (ship-b1)
				0.01, // Destruction
			},
			expectedShots: 2,
			expectedSideA: func(f *galaxy.Fleet) bool {
				return len(f.Ships) == 2 && !f.Ships[0].Destroyed && !f.Ships[1].Destroyed
			},
			expectedSideB: func(f *galaxy.Fleet) bool {
				// Both ships destroyed
				return len(f.Ships) == 2 && f.Ships[0].Destroyed && f.Ships[1].Destroyed
			},
			battleOverEarly: true,
		},
		{
			name: "High defense vs low attack - no casualties",
			fleetA: &galaxy.Fleet{
				Ships: []*galaxy.Ship{
					{
						ID:   "ship-a1",
						Name: "Armored Ship A",
						Tech: galaxy.ShipTech{
							Attack:  0,
							Guns:    0, // No guns
							Defense: 4,
						},
						Owner: "race-a",
					},
				},
				Owner: "race-a",
			},
			fleetB: &galaxy.Fleet{
				Ships: []*galaxy.Ship{
					{
						ID:   "ship-b1",
						Name: "Weak Attacker B",
						Tech: galaxy.ShipTech{
							Attack:  1,
							Guns:    5, // 5 guns - will keep firing until max shots
							Defense: 1,
						},
						Owner: "race-b",
					},
				},
				Owner: "race-b",
			},
			randomValues: []float64{
				// Initial shooter selection
				0.6, // Side B
				0.0, // Shooter index 0 (ship-b1)
				// Shots - the random generator will wrap around these values
				0.0,  // Target index 0 (ship-a1)
				0.99, // Destruction check: high value = no destruction (defense 4 / attack 1 = ratio 4 -> prob 0)
			},
			expectedShots: 1000, // Hits max shot limit - demonstrates impenetrable defense
			expectedSideA: func(f *galaxy.Fleet) bool {
				// Ship A should survive all 1000 attacks due to high defense
				return len(f.Ships) == 1 && !f.Ships[0].Destroyed
			},
			expectedSideB: func(f *galaxy.Fleet) bool {
				return len(f.Ships) == 1 && !f.Ships[0].Destroyed
			},
			battleOverEarly: false, // Battle continues until max shots
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := gamemath.NewPredefinedRandomGenerator(tt.randomValues)
			idGenerator := util.NewSequenceGenerator([]string{"battle-1"})

			battleHandler := NewBattleHandler(idGenerator, nil)
			battleHandler.initializeBattleState(tt.fleetA, tt.fleetB)

			decisionProducer := NewRuntimeDecisionProducer(rng, battleHandler)
			battleHandler.decisionProducer = decisionProducer

			battle := battleHandler.ExecuteBattle(tt.fleetA, tt.fleetB)

			// Verify number of shots
			if len(battle.Shots) != tt.expectedShots {
				t.Errorf("Expected %d shots, got %d", tt.expectedShots, len(battle.Shots))
			}

			// Verify PostSideA state
			if !tt.expectedSideA(battle.PostSideA) {
				t.Errorf("PostSideA state validation failed")
				for i, ship := range battle.PostSideA.Ships {
					t.Logf("  Ship A[%d]: ID=%s, Destroyed=%v", i, ship.ID, ship.Destroyed)
				}
			}

			// Verify PostSideB state
			if !tt.expectedSideB(battle.PostSideB) {
				t.Errorf("PostSideB state validation failed")
				for i, ship := range battle.PostSideB.Ships {
					t.Logf("  Ship B[%d]: ID=%s, Destroyed=%v", i, ship.ID, ship.Destroyed)
				}
			}

			// Log shots for debugging
			t.Logf("Battle shots:")
			for i, shot := range battle.Shots {
				t.Logf("  Shot %d: %s -> %s (destroyed=%v)", i+1, shot.Source, shot.Destination, shot.Result)
			}
		})
	}
}

func TestBattleHandlerWithRuntimeDecisionProducerStalemate(t *testing.T) {
	// Special test case: Battle reaches max shots without ending
	fleetA := &galaxy.Fleet{
		Ships: []*galaxy.Ship{
			{
				ID:   "ship-a1",
				Name: "Eternal Fighter A",
				Tech: galaxy.ShipTech{
					Attack:  1,
					Guns:    1000, // Many guns to keep shooting
					Defense: 10,   // High defense
				},
				Owner: "race-a",
			},
		},
		Owner: "race-a",
	}

	fleetB := &galaxy.Fleet{
		Ships: []*galaxy.Ship{
			{
				ID:   "ship-b1",
				Name: "Eternal Fighter B",
				Tech: galaxy.ShipTech{
					Attack:  1,
					Guns:    1000, // Many guns to keep shooting
					Defense: 10,   // High defense
				},
				Owner: "race-b",
			},
		},
		Owner: "race-b",
	}

	// Generate enough random values for many shots
	// Each shot needs: side selection (if new shooter), shooter selection (if new), target, destruction
	randomValues := make([]float64, 0, 4000)
	for i := 0; i < 1000; i++ {
		if i%1000 == 0 {
			// New shooter selection
			randomValues = append(randomValues,
				0.3, // Side A
				0.0, // Shooter index 0
			)
		}
		randomValues = append(randomValues,
			0.0,  // Target index 0
			0.99, // No destruction (high defense vs low attack)
		)
	}

	rng := gamemath.NewPredefinedRandomGenerator(randomValues)
	idGenerator := util.NewSequenceGenerator([]string{"battle-stalemate"})

	battleHandler := NewBattleHandler(idGenerator, nil)
	battleHandler.initializeBattleState(fleetA, fleetB)

	decisionProducer := NewRuntimeDecisionProducer(rng, battleHandler)
	battleHandler.decisionProducer = decisionProducer

	battle := battleHandler.ExecuteBattle(fleetA, fleetB)

	// Should hit max shots limit (1000)
	if len(battle.Shots) != 1000 {
		t.Errorf("Expected battle to hit max shots limit of 1000, got %d shots", len(battle.Shots))
	}

	// Both ships should still be alive
	if battle.PostSideA.Ships[0].Destroyed {
		t.Error("Ship A should not be destroyed in stalemate")
	}
	if battle.PostSideB.Ships[0].Destroyed {
		t.Error("Ship B should not be destroyed in stalemate")
	}

	t.Logf("Stalemate battle ended with %d shots, both ships alive", len(battle.Shots))
}
