package di

import (
	"glaktika.eu/galaktika/internal/api"
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/pkg/galaxy"
)

func NewAuthenticationManager(raceRepository *dao.RaceRepository) api.AuthenticationManager {
	return api.NewMemoryAuthenticationManager(raceRepository)
}

func NewFleetBuildRepository() *dao.FleetBuildRepository {
	r := dao.NewFleetBuildRepository()

	races := []string{"rex", "zyx", "keth"}
	for i, divisionID := range []string{"alpha", "beta", "gamma"} {
		// fleet build with assigned two ship models
		fleetBuild := &galaxy.FleetBuild{ID: divisionID + "-build-1", DivisionId: divisionID, RaceId: races[i]}
		r.Upsert(fleetBuild)
		nullShipModelAssigned := &galaxy.ShipModelAssignment{
			ID:           fleetBuild.ID + "_" + races[i] + "_",
			FleetBuildID: fleetBuild.ID,
			ShipModel:    nil,
			Amount:       0,
			ResultMass:   0,
		}
		r.AssignShipModel(nullShipModelAssigned)

		shipModelAssigned := &galaxy.ShipModelAssignment{
			ID:           fleetBuild.ID + "_" + races[i] + "-fighter",
			FleetBuildID: fleetBuild.ID,
			ShipModelID:  races[i] + "-fighter",
			Amount:       1,
			ResultMass:   5,
		}
		r.AssignShipModel(shipModelAssigned)

		// fleet build without assigned ship models
		r.Upsert(&galaxy.FleetBuild{ID: divisionID + "-build-2", DivisionId: divisionID, RaceId: races[(i+1)%len(races)]})
	}

	return r
}

func NewShipModelRepository() *dao.ShipModelRepository {
	r := dao.NewShipModelRepository()

	for _, raceID := range []string{"rex", "zyx", "keth"} {
		r.Upsert(&galaxy.ShipModel{ID: raceID + "-fighter", Name: "Fighter", OwnerId: raceID, Guns: 2, OneGunMass: 1, DefenseMass: 1, EngineMass: 1, CargoMass: 0})
		r.Upsert(&galaxy.ShipModel{ID: raceID + "-cruiser", Name: "Cruiser", OwnerId: raceID, Guns: 4, OneGunMass: 2, DefenseMass: 3, EngineMass: 2, CargoMass: 1})
		r.Upsert(&galaxy.ShipModel{ID: raceID + "-freighter", Name: "Freighter", OwnerId: raceID, Guns: 1, OneGunMass: 1, DefenseMass: 1, EngineMass: 2, CargoMass: 5})
	}

	return r
}

func NewRaceRepository() *dao.RaceRepository {
	r := dao.NewRaceRepository()

	r.Upsert(&galaxy.Race{ID: "rex", Name: "Commander Rex", Role: galaxy.RolePlayer, Token: "token-rex-001"})
	r.Upsert(&galaxy.Race{ID: "zyx", Name: "Admiral Zyx", Role: galaxy.RolePlayer, Token: "token-zyx-002"})
	r.Upsert(&galaxy.Race{ID: "keth", Name: "Warlord Keth", Role: galaxy.RolePlayer, Token: "token-keth-003"})
	r.Upsert(&galaxy.Race{ID: "admin", Name: "Administrator", Role: galaxy.RoleAdmin, Token: "token-admin-000"})

	return r
}

func NewDivisionRepository() *dao.DivisionRepository {
	divisionRepository := dao.NewDivisionRepository()

	divisionRepository.Upsert(&galaxy.Division{ID: "alpha", ResourcesAmount: 500, TechAttack: 1, TechDefense: 1, TechEngines: 1, TechCargo: 1})
	divisionRepository.Upsert(&galaxy.Division{ID: "beta", ResourcesAmount: 750, TechAttack: 1, TechDefense: 1, TechEngines: 1, TechCargo: 1})
	divisionRepository.Upsert(&galaxy.Division{ID: "gamma", ResourcesAmount: 1000, TechAttack: 1, TechDefense: 1, TechEngines: 1, TechCargo: 1})

	return divisionRepository
}
