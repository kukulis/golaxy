package di

import (
	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/api"
	"glaktika.eu/galaktika/internal/dao"
)

var BattleRepositoryInstance *dao.BattleRepository
var BattleControllerInstance *api.BattleController
var DivisionRepositoryInstance *dao.DivisionRepository
var DivisionControllerInstance *api.DivisionController
var FleetBuildRepositoryInstance *dao.FleetBuildRepository
var ShipModelRepositoryInstance *dao.ShipModelRepository

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
}

// ResetTestData clears all data in repositories for testing
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

func RegisterRoutes(apiRoute *gin.RouterGroup) {
	apiRoute.GET("/battle", func(c *gin.Context) { BattleControllerInstance.GetBattle(c) })
	apiRoute.GET("/divisions", func(c *gin.Context) { DivisionControllerInstance.GetAllDivisions(c) })
	apiRoute.GET("/divisions/:id", func(c *gin.Context) { DivisionControllerInstance.GetDivision(c) })
	apiRoute.POST("/divisions", func(c *gin.Context) { DivisionControllerInstance.CreateDivision(c) })
	apiRoute.PUT("/divisions/:id", func(c *gin.Context) { DivisionControllerInstance.UpdateDivision(c) })
	apiRoute.DELETE("/divisions/:id", func(c *gin.Context) { DivisionControllerInstance.DeleteDivision(c) })
}
