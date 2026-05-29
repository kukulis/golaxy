package di

import "github.com/gin-gonic/gin"

func RegisterRoutes(apiRoute *gin.RouterGroup) {
	apiRoute.GET("/current-race", func(c *gin.Context) { RaceControllerInstance.GetCurrentRace(c) })
	apiRoute.GET("/races", func(c *gin.Context) { RaceControllerInstance.GetAllRaces(c) })
	apiRoute.POST("/races", func(c *gin.Context) { RaceControllerInstance.CreateRace(c) })
	apiRoute.GET("/races/:id", func(c *gin.Context) { RaceControllerInstance.GetRace(c) })
	apiRoute.DELETE("/races/:id", func(c *gin.Context) { RaceControllerInstance.DeleteRace(c) })

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
	apiRoute.GET("/fleet-builds/:id/statistics", func(c *gin.Context) { FleetBuildControllerInstance.GetStatistics(c) })
	apiRoute.GET("/fleet-builds/:id/technologies", func(c *gin.Context) { FleetBuildControllerInstance.GetTechnologies(c) })
	apiRoute.GET("/fleet-builds/:id/ship-model-assignments", func(c *gin.Context) { FleetBuildControllerInstance.GetAssignedShipModels(c) })
	apiRoute.GET("/ship-model-assignment/:id", func(c *gin.Context) { FleetBuildControllerInstance.GetShipModelAssignment(c) })
	apiRoute.POST("/ship-model-assignment/:id", func(c *gin.Context) { FleetBuildControllerInstance.UpdateShipModelAssignment(c) })
	apiRoute.POST("/ship-model-assignment", func(c *gin.Context) { FleetBuildControllerInstance.AddShipModelAssignment(c) })
	apiRoute.DELETE("/ship-models-assignment/:id", func(c *gin.Context) { FleetBuildControllerInstance.UnassignShipModel(c) })
	apiRoute.POST("/fleet-builds/:id/build", func(c *gin.Context) { FleetBuildControllerInstance.Build(c) })
	apiRoute.GET("/fleet-builds/:id/fleet", func(c *gin.Context) { FleetBuildControllerInstance.GetFleet(c) })
	apiRoute.GET("/fleet-builds/:id/ship-models/:shipModelId/calculate-ship-tech", func(c *gin.Context) { FleetBuildControllerInstance.CalculateShipTech(c) })

	apiRoute.GET("/ship-models", func(c *gin.Context) { ShipModelControllerInstance.GetAllShipModels(c) })
	apiRoute.GET("/ship-models/:id", func(c *gin.Context) { ShipModelControllerInstance.GetShipModel(c) })
	apiRoute.POST("/ship-models", func(c *gin.Context) { ShipModelControllerInstance.CreateShipModel(c) })
	apiRoute.PUT("/ship-models/:id", func(c *gin.Context) { ShipModelControllerInstance.UpdateShipModel(c) })
	apiRoute.POST("/ship-models/:id/calculate-ship-tech", func(c *gin.Context) { ShipModelControllerInstance.CalculateShipTech(c) })
	apiRoute.DELETE("/ship-models/:id", func(c *gin.Context) { ShipModelControllerInstance.DeleteShipModel(c) })
}
