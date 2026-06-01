package di

import (
	"glaktika.eu/galaktika/internal/api"
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/internal/web"
)

var AuthenticationManagerInstance api.AuthenticationManager
var BattleRepositoryInstance *dao.BattleRepository
var ChallengeRepositoryInstance *dao.ChallengeRepository
var ChallengeControllerInstance *api.ChallengeController
var BattleControllerInstance *api.BattleController
var DivisionRepositoryInstance *dao.DivisionRepository
var DivisionControllerInstance *api.DivisionController
var FleetBuildRepositoryInstance *dao.FleetBuildRepository
var FleetBuildControllerInstance *api.FleetBuildController
var RaceControllerInstance *api.RaceController
var FleetRepositoryInstance *dao.FleetRepository
var RaceRepositoryInstance *dao.RaceRepository
var ShipModelRepositoryInstance *dao.ShipModelRepository
var ShipModelControllerInstance *api.ShipModelController
var WebControllerInstance *web.WebController

func CreateSingletons(env string) {
	// Based on env, choose repository implementation
	// Currently only in-memory repos are implemented
	switch env {
	case "test":
		RaceRepositoryInstance = dao.NewRaceRepository()
		AuthenticationManagerInstance = NewAuthenticationManager(RaceRepositoryInstance)
		BattleRepositoryInstance = dao.NewBattleRepository()
		ChallengeRepositoryInstance = dao.NewChallengeRepository()
		DivisionRepositoryInstance = dao.NewDivisionRepository()
		FleetBuildRepositoryInstance = dao.NewFleetBuildRepository()
		FleetRepositoryInstance = dao.NewFleetRepository()
		ShipModelRepositoryInstance = dao.NewShipModelRepository()

	case "dev":
		RaceRepositoryInstance = NewRaceRepository()
		AuthenticationManagerInstance = NewAuthenticationManager(RaceRepositoryInstance)
		BattleRepositoryInstance = dao.NewBattleRepository()
		ChallengeRepositoryInstance = dao.NewChallengeRepository()
		DivisionRepositoryInstance = NewDivisionRepository()
		FleetBuildRepositoryInstance = NewFleetBuildRepository()
		FleetRepositoryInstance = dao.NewFleetRepository()
		ShipModelRepositoryInstance = NewShipModelRepository()

	case "prod":
		// Future: DB-backed repositories
		// Example: DivisionRepositoryInstance = dao.NewDBDivisionRepository(dbConn)
		panic("prod environment not yet implemented")
	default:
		panic("unknown environment: " + env)
	}

	// Controllers are environment-agnostic
	WebControllerInstance = web.NewWebController(FleetBuildRepositoryInstance)
	ChallengeControllerInstance = api.NewChallengeController(AuthenticationManagerInstance, ChallengeRepositoryInstance)
	RaceControllerInstance = api.NewRaceController(AuthenticationManagerInstance, RaceRepositoryInstance)
	BattleControllerInstance = api.NewBattleController(AuthenticationManagerInstance, BattleRepositoryInstance)
	DivisionControllerInstance = api.NewDivisionController(AuthenticationManagerInstance, DivisionRepositoryInstance)
	FleetBuildControllerInstance = api.NewFleetBuildController(AuthenticationManagerInstance, FleetBuildRepositoryInstance, FleetRepositoryInstance, ShipModelRepositoryInstance, DivisionRepositoryInstance)
	ShipModelControllerInstance = api.NewShipModelController(AuthenticationManagerInstance, ShipModelRepositoryInstance)
}

// ResetTestData clears all data in repositories for testing.
// This function will be used in the database tests where the test server
// is shared across multiple test cases and repositories need to be reset
// between tests to ensure data isolation.
func ResetTestData() {
	// TODO: Add BattleRepositoryInstance.ResetData() when implemented
	if ChallengeRepositoryInstance != nil {
		ChallengeRepositoryInstance.ResetData()
	}
	if DivisionRepositoryInstance != nil {
		DivisionRepositoryInstance.ResetData()
	}
	if FleetBuildRepositoryInstance != nil {
		FleetBuildRepositoryInstance.ResetData()
	}
	if RaceRepositoryInstance != nil {
		RaceRepositoryInstance.ResetData()
	}
	if ShipModelRepositoryInstance != nil {
		ShipModelRepositoryInstance.ResetData()
	}
}
