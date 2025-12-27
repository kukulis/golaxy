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
	FleetBuildControllerInstance = api.NewFleetBuildController(FleetBuildRepositoryInstance)
	ShipModelControllerInstance = api.NewShipModelController(ShipModelRepositoryInstance)
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

	apiRoute.GET("/fleet-builds", func(c *gin.Context) { FleetBuildControllerInstance.GetAllFleetBuilds(c) })
	apiRoute.GET("/fleet-builds/:id", func(c *gin.Context) { FleetBuildControllerInstance.GetFleetBuild(c) })
	apiRoute.POST("/fleet-builds", func(c *gin.Context) { FleetBuildControllerInstance.CreateFleetBuild(c) })
	apiRoute.PUT("/fleet-builds/:id", func(c *gin.Context) { FleetBuildControllerInstance.UpdateFleetBuild(c) })
	apiRoute.DELETE("/fleet-builds/:id", func(c *gin.Context) { FleetBuildControllerInstance.DeleteFleetBuild(c) })

	apiRoute.GET("/ship-models", func(c *gin.Context) { ShipModelControllerInstance.GetAllShipModels(c) })
	apiRoute.GET("/ship-models/:id", func(c *gin.Context) { ShipModelControllerInstance.GetShipModel(c) })
	apiRoute.POST("/ship-models", func(c *gin.Context) { ShipModelControllerInstance.CreateShipModel(c) })
	apiRoute.PUT("/ship-models/:id", func(c *gin.Context) { ShipModelControllerInstance.UpdateShipModel(c) })
	apiRoute.DELETE("/ship-models/:id", func(c *gin.Context) { ShipModelControllerInstance.DeleteShipModel(c) })
}
