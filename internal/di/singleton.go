package di

import (
	"glaktika.eu/galaktika/internal/api"
	"glaktika.eu/galaktika/internal/dao"
)

var BattleRepositoryInstance *dao.BattleRepository
var BattleControllerInstance *api.BattleController
var DivisionRepositoryInstance *dao.DivisionRepository
var DivisionControllerInstance *api.DivisionController
var FleetBuildRepositoryInstance *dao.FleetBuildRepository
var FleetBuildControllerInstance *api.FleetBuildController
var ShipModelRepositoryInstance *dao.ShipModelRepository
var ShipModelControllerInstance *api.ShipModelController

func CreateSingletons(env string) {
	// Based on env, choose repository implementation
	// Currently only in-memory repos are implemented
	switch env {
	case "test", "dev":
		BattleRepositoryInstance = dao.NewBattleRepository()
		DivisionRepositoryInstance = dao.NewDivisionRepository()
		FleetBuildRepositoryInstance = dao.NewFleetBuildRepository()
		ShipModelRepositoryInstance = dao.NewShipModelRepository()
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
	FleetBuildControllerInstance = api.NewFleetBuildController(FleetBuildRepositoryInstance, ShipModelRepositoryInstance)
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
