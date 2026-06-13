package di

import (
	"glaktika.eu/galaktika/internal/api"
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/internal/web"
	"glaktika.eu/galaktika/internal/ws"
	"glaktika.eu/galaktika/pkg/util"
)

var authenticationManagerInstance api.AuthenticationManager
var battleRepositoryInstance *dao.BattleRepository
var challengeRepositoryInstance *dao.ChallengeRepository
var challengeControllerInstance *api.ChallengeController
var battleControllerInstance *api.BattleController
var divisionRepositoryInstance *dao.DivisionRepository
var divisionControllerInstance *api.DivisionController
var fleetBuildRepositoryInstance *dao.FleetBuildRepository
var fleetBuildControllerInstance *api.FleetBuildController
var raceControllerInstance *api.RaceController
var fleetRepositoryInstance *dao.FleetRepository
var raceRepositoryInstance *dao.RaceRepository
var shipModelRepositoryInstance *dao.ShipModelRepository
var shipModelControllerInstance *api.ShipModelController
var webControllerInstance *web.WebController
var wsControllerInstance *api.WsController
var hubInstance *ws.Hub
var wsDispatcherInstance *util.Dispatcher

func GetHubInstance() *ws.Hub {
	return hubInstance
}

func GetRaceRepositoryInstance() *dao.RaceRepository {
	return raceRepositoryInstance
}

func CreateSingletons(env string) {
	// Based on env, choose repository implementation
	// Currently only in-memory repos are implemented
	switch env {
	case "test":
		raceRepositoryInstance = dao.NewRaceRepository()
		authenticationManagerInstance = NewAuthenticationManager(raceRepositoryInstance)
		battleRepositoryInstance = dao.NewBattleRepository()
		challengeRepositoryInstance = dao.NewChallengeRepository()
		divisionRepositoryInstance = dao.NewDivisionRepository()
		fleetBuildRepositoryInstance = dao.NewFleetBuildRepository()
		fleetRepositoryInstance = dao.NewFleetRepository()
		shipModelRepositoryInstance = dao.NewShipModelRepository()

	case "dev":
		raceRepositoryInstance = NewRaceRepository()
		authenticationManagerInstance = NewAuthenticationManager(raceRepositoryInstance)
		battleRepositoryInstance = dao.NewBattleRepository()
		challengeRepositoryInstance = NewChallengeRepository()
		divisionRepositoryInstance = NewDivisionRepository()
		fleetBuildRepositoryInstance = NewFleetBuildRepository()
		fleetRepositoryInstance = dao.NewFleetRepository()
		shipModelRepositoryInstance = NewShipModelRepository()

	case "prod":
		// Future: DB-backed repositories
		// Example: divisionRepositoryInstance = dao.NewDBDivisionRepository(dbConn)
		panic("prod environment not yet implemented")
	default:
		panic("unknown environment: " + env)
	}

	// Controllers are environment-agnostic
	wsDispatcherInstance = util.NewDispatcher()
	hubInstance = ws.NewHub(wsDispatcherInstance)
	webControllerInstance = web.NewWebController(fleetBuildRepositoryInstance)
	wsControllerInstance = api.NewWsController(authenticationManagerInstance, hubInstance)
	challengeControllerInstance = api.NewChallengeController(authenticationManagerInstance, challengeRepositoryInstance, divisionRepositoryInstance, raceRepositoryInstance, fleetBuildRepositoryInstance)
	raceControllerInstance = api.NewRaceController(authenticationManagerInstance, raceRepositoryInstance)
	battleControllerInstance = api.NewBattleController(authenticationManagerInstance, battleRepositoryInstance)
	divisionControllerInstance = api.NewDivisionController(authenticationManagerInstance, divisionRepositoryInstance)
	fleetBuildControllerInstance = api.NewFleetBuildController(authenticationManagerInstance, fleetBuildRepositoryInstance, fleetRepositoryInstance, shipModelRepositoryInstance, divisionRepositoryInstance)
	shipModelControllerInstance = api.NewShipModelController(authenticationManagerInstance, shipModelRepositoryInstance)
}

// ResetTestData clears all data in repositories for testing.
// This function will be used in the database tests where the test server
// is shared across multiple test cases and repositories need to be reset
// between tests to ensure data isolation.
func ResetTestData() {
	if challengeRepositoryInstance != nil {
		challengeRepositoryInstance.ResetData()
	}
	if divisionRepositoryInstance != nil {
		divisionRepositoryInstance.ResetData()
	}
	if fleetBuildRepositoryInstance != nil {
		fleetBuildRepositoryInstance.ResetData()
	}
	if raceRepositoryInstance != nil {
		raceRepositoryInstance.ResetData()
	}
	if shipModelRepositoryInstance != nil {
		shipModelRepositoryInstance.ResetData()
	}
}
