package game

import (
	"glaktika.eu/galaktika/pkg/galaxy"
	"glaktika.eu/galaktika/pkg/gamemath"
	"glaktika.eu/galaktika/pkg/util"
	"testing"
)

// testLogger wraps testing.T to implement galaxy.Logger interface
type testLogger struct {
	t *testing.T
}

func (tl *testLogger) Printf(format string, v ...interface{}) {
	tl.t.Logf(format, v...)
}

func TestExecuteBattle(t *testing.T) {
	tests := []struct {
		name            string
		fleetA          *galaxy.Fleet
		fleetB          *galaxy.Fleet
		idGenerator     util.IdGenerator
		randomGenerator gamemath.RandomGenerator
		expectedBattle  *galaxy.Battle
	}{
		{
			name: "Battle between two single-ship fleets",
			fleetA: &galaxy.Fleet{
				Ships: []*galaxy.Ship{
					{
						ID:   "ship-a1",
						Name: "Ship A1",
						Tech: galaxy.ShipTech{
							Attack:        1,
							Guns:          1,
							Defense:       1,
							Speed:         1,
							CargoCapacity: 1,
							Mass:          1,
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
						Name: "Ship B1",
						Tech: galaxy.ShipTech{
							Attack:        1,
							Guns:          1,
							Defense:       1,
							Speed:         1,
							CargoCapacity: 1,
							Mass:          1,
						},
						Owner: "race-b",
					},
				},
				Owner: "race-b",
			},
			idGenerator:     util.NewSequenceGenerator([]string{"battle-1"}),
			randomGenerator: gamemath.NewPredefinedRandomGenerator([]float64{0.5}),
			expectedBattle: &galaxy.Battle{
				ID:    "battle-1",
				Shots: nil,
				PostSideA: &galaxy.Fleet{
					Ships: []*galaxy.Ship{
						{
							ID:   "ship-a1",
							Name: "Ship A1",
							Tech: galaxy.ShipTech{
								Attack:        1,
								Guns:          1,
								Defense:       1,
								Speed:         1,
								CargoCapacity: 1,
								Mass:          1,
							},
							Owner:     "race-a",
							Destroyed: false,
						},
					},
					Owner: "race-a",
				},
				PostSideB: &galaxy.Fleet{
					Ships: []*galaxy.Ship{
						{
							ID:   "ship-b1",
							Name: "Ship B1",
							Tech: galaxy.ShipTech{
								Attack:        1,
								Guns:          1,
								Defense:       1,
								Speed:         1,
								CargoCapacity: 1,
								Mass:          1,
							},
							Owner:     "race-b",
							Destroyed: false,
						},
					},
					Owner: "race-b",
				},
			},
		},
		{
			name: "Battle between multi-ship fleets",
			fleetA: &galaxy.Fleet{
				Ships: []*galaxy.Ship{
					{
						ID:   "ship-a1",
						Name: "Ship A1",
						Tech: galaxy.ShipTech{
							Attack:        1,
							Guns:          1,
							Defense:       1,
							Speed:         1,
							CargoCapacity: 1,
							Mass:          1,
						},
						Owner: "race-a",
					},
					{
						ID:   "ship-a2",
						Name: "Ship A2",
						Tech: galaxy.ShipTech{
							Attack:        1,
							Guns:          1,
							Defense:       1,
							Speed:         1,
							CargoCapacity: 1,
							Mass:          1,
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
						Name: "Ship B1",
						Tech: galaxy.ShipTech{
							Attack:        1,
							Guns:          1,
							Defense:       1,
							Speed:         1,
							CargoCapacity: 1,
							Mass:          1,
						},
						Owner: "race-b",
					},
					{
						ID:   "ship-b2",
						Name: "Ship B2",
						Tech: galaxy.ShipTech{
							Attack:        1,
							Guns:          1,
							Defense:       1,
							Speed:         1,
							CargoCapacity: 1,
							Mass:          1,
						},
						Owner: "race-b",
					},
				},
				Owner: "race-b",
			},
			idGenerator:     util.NewSequenceGenerator([]string{"battle-2"}),
			randomGenerator: gamemath.NewPredefinedRandomGenerator([]float64{0.5}),
			expectedBattle: &galaxy.Battle{
				ID:    "battle-2",
				Shots: nil,
				PostSideA: &galaxy.Fleet{
					Ships: []*galaxy.Ship{
						{
							ID:   "ship-a1",
							Name: "Ship A1",
							Tech: galaxy.ShipTech{
								Attack:        1,
								Guns:          1,
								Defense:       1,
								Speed:         1,
								CargoCapacity: 1,
								Mass:          1,
							},
							Owner:     "race-a",
							Destroyed: false,
						},
						{
							ID:   "ship-a2",
							Name: "Ship A2",
							Tech: galaxy.ShipTech{
								Attack:        1,
								Guns:          1,
								Defense:       1,
								Speed:         1,
								CargoCapacity: 1,
								Mass:          1,
							},
							Owner:     "race-a",
							Destroyed: false,
						},
					},
					Owner: "race-a",
				},
				PostSideB: &galaxy.Fleet{
					Ships: []*galaxy.Ship{
						{
							ID:   "ship-b1",
							Name: "Ship B1",
							Tech: galaxy.ShipTech{
								Attack:        1,
								Guns:          1,
								Defense:       1,
								Speed:         1,
								CargoCapacity: 1,
								Mass:          1,
							},
							Owner:     "race-b",
							Destroyed: false,
						},
						{
							ID:   "ship-b2",
							Name: "Ship B2",
							Tech: galaxy.ShipTech{
								Attack:        1,
								Guns:          1,
								Defense:       1,
								Speed:         1,
								CargoCapacity: 1,
								Mass:          1,
							},
							Owner:     "race-b",
							Destroyed: false,
						},
					},
					Owner: "race-b",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &testLogger{t: t}
			battleHandler := NewBattleHandler(tt.idGenerator, tt.randomGenerator)
			battle := battleHandler.ExecuteBattle(tt.fleetA, tt.fleetB)

			// Assertion: Battle.SideA must be the same fleet as FleetA
			if battle.SideA != tt.fleetA {
				t.Errorf("Battle.SideA is not the same reference as FleetA")
			}

			// Assertion: Battle.SideB must be the same fleet as FleetB
			if battle.SideB != tt.fleetB {
				t.Errorf("Battle.SideB is not the same reference as FleetB")
			}

			// Assertion: Battle.PostSideA must match expected PostSideA
			if !battle.PostSideA.EqualShips(tt.expectedBattle.PostSideA, logger) {
				t.Errorf("Battle.PostSideA does not match expected PostSideA")
			}

			// Assertion: Battle.PostSideB must match expected PostSideB
			if !battle.PostSideB.EqualShips(tt.expectedBattle.PostSideB, logger) {
				t.Errorf("Battle.PostSideB does not match expected PostSideB")
			}
		})
	}
}
