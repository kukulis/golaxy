package game

import (
	"glaktika.eu/galaktika/pkg/galaxy"
	"glaktika.eu/galaktika/pkg/gamemath"
	"glaktika.eu/galaktika/pkg/util"
	"testing"
)

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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			// Assertion: Battle.PostSideA must have all ships copies from SideA
			if battle.PostSideA == nil {
				t.Errorf("Battle.PostSideA is nil, expected fleet with %d ships", len(battle.SideA.Ships))
			} else {
				if len(battle.PostSideA.Ships) != len(battle.SideA.Ships) {
					t.Errorf("PostSideA has %d ships, SideA has %d ships",
						len(battle.PostSideA.Ships), len(battle.SideA.Ships))
				}

				for i, postShip := range battle.PostSideA.Ships {
					originalShip := battle.SideA.Ships[i]
					if !postShip.EqualsWithoutDamage(originalShip) {
						t.Errorf("PostSideA ship[%d] does not match SideA ship[%d] (ignoring Destroyed field)", i, i)
					}
				}
			}

			// Assertion: Battle.PostSideB must have all ships copies from SideB
			if battle.PostSideB == nil {
				t.Errorf("Battle.PostSideB is nil, expected fleet with %d ships", len(battle.SideB.Ships))
			} else {
				if len(battle.PostSideB.Ships) != len(battle.SideB.Ships) {
					t.Errorf("PostSideB has %d ships, SideB has %d ships",
						len(battle.PostSideB.Ships), len(battle.SideB.Ships))
				}

				for i, postShip := range battle.PostSideB.Ships {
					originalShip := battle.SideB.Ships[i]
					if !postShip.EqualsWithoutDamage(originalShip) {
						t.Errorf("PostSideB ship[%d] does not match SideB ship[%d] (ignoring Destroyed field)", i, i)
					}
				}
			}
		})
	}
}
