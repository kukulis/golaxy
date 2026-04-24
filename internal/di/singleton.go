package di

import (
	"glaktika.eu/galaktika/internal/api"
	"glaktika.eu/galaktika/internal/dao"
)

var AuthenticationManagerInstance api.AuthenticationManager
var BattleRepositoryInstance *dao.BattleRepository
var BattleControllerInstance *api.BattleController
var DivisionRepositoryInstance *dao.DivisionRepository
var DivisionControllerInstance *api.DivisionController
var FleetBuildRepositoryInstance *dao.FleetBuildRepository
var FleetBuildControllerInstance *api.FleetBuildController
var FleetRepositoryInstance *dao.FleetRepository
var ShipModelRepositoryInstance *dao.ShipModelRepository
var ShipModelControllerInstance *api.ShipModelController

func CreateSingletons(env string) {
	// Based on env, choose repository implementation
	// Currently only in-memory repos are implemented
	switch env {
	case "test", "dev":
		AuthenticationManagerInstance = api.NewMemoryAuthenticationManager()
		BattleRepositoryInstance = dao.NewBattleRepository()
		DivisionRepositoryInstance = NewDivisionRepository()

		FleetBuildRepositoryInstance = dao.NewFleetBuildRepository()
		FleetRepositoryInstance = dao.NewFleetRepository()
		ShipModelRepositoryInstance = dao.NewShipModelRepository()

		AuthenticationManagerInstance.AddToken("user1", "user1")
		AuthenticationManagerInstance.AddToken("user2", "user2")

	case "prod":
		// Future: DB-backed repositories
		// Example: DivisionRepositoryInstance = dao.NewDBDivisionRepository(dbConn)
		panic("prod environment not yet implemented")
	default:
		panic("unknown environment: " + env)
	}

	// Controllers are environment-agnostic
	BattleControllerInstance = api.NewBattleController(BattleRepositoryInstance)
	DivisionControllerInstance = api.NewDivisionController(DivisionRepositoryInstance)
	FleetBuildControllerInstance = api.NewFleetBuildController(AuthenticationManagerInstance, FleetBuildRepositoryInstance, FleetRepositoryInstance, ShipModelRepositoryInstance, DivisionRepositoryInstance)
	ShipModelControllerInstance = api.NewShipModelController(ShipModelRepositoryInstance)
}

// ResetTestData clears all data in repositories for testing.
// This function will be used in the database tests where the test server
// is shared across multiple test cases and repositories need to be reset
// between tests to ensure data isolation.
func ResetTestData() {
	// TODO: Add BattleRepositoryInstance.ResetData() when implemented
	if DivisionRepositoryInstance != nil {
		DivisionRepositoryInstance.ResetData()
	}
	if FleetBuildRepositoryInstance != nil {
		FleetBuildRepositoryInstance.ResetData()
	}
	if ShipModelRepositoryInstance != nil {
		ShipModelRepositoryInstance.ResetData()
	}
}
