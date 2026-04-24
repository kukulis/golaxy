package di

import (
	"glaktika.eu/galaktika/internal/api"
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/pkg/galaxy"
)

func NewAuthenticationManager() api.AuthenticationManager {
	am := api.NewMemoryAuthenticationManager()

	am.AddToken("token-rex-001", &galaxy.Race{ID: "rex", Name: "Commander Rex", Role: "commander"})
	am.AddToken("token-zyx-002", &galaxy.Race{ID: "zyx", Name: "Admiral Zyx", Role: "admiral"})
	am.AddToken("token-keth-003", &galaxy.Race{ID: "keth", Name: "Warlord Keth", Role: "warlord"})

	return am
}

func NewFleetBuildRepository() *dao.FleetBuildRepository {
	r := dao.NewFleetBuildRepository()

	races := []string{"rex", "zyx", "keth"}
	for i, divisionID := range []string{"alpha", "beta", "gamma"} {
		r.Upsert(&galaxy.FleetBuild{ID: divisionID + "-build-1", DivisionId: divisionID, RaceId: races[i]})
		r.Upsert(&galaxy.FleetBuild{ID: divisionID + "-build-2", DivisionId: divisionID, RaceId: races[(i+1)%len(races)]})
	}

	return r
}

func NewDivisionRepository() *dao.DivisionRepository {
	divisionRepository := dao.NewDivisionRepository()

	divisionRepository.Upsert(&galaxy.Division{ID: "alpha", ResourcesAmount: 500, TechAttack: 1, TechDefense: 1, TechEngines: 1, TechCargo: 1})
	divisionRepository.Upsert(&galaxy.Division{ID: "beta", ResourcesAmount: 750, TechAttack: 1, TechDefense: 1, TechEngines: 1, TechCargo: 1})
	divisionRepository.Upsert(&galaxy.Division{ID: "gamma", ResourcesAmount: 1000, TechAttack: 1, TechDefense: 1, TechEngines: 1, TechCargo: 1})

	return divisionRepository
}
