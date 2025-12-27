package di

import "github.com/gin-gonic/gin"

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
	apiRoute.GET("/fleet-builds/:id/ship-models", func(c *gin.Context) { FleetBuildControllerInstance.GetAssignedShipModels(c) })
	apiRoute.POST("/fleet-builds/:id/ship-models", func(c *gin.Context) { FleetBuildControllerInstance.AssignShipModel(c) })
	apiRoute.DELETE("/fleet-builds/:id/ship-models/:shipModelId", func(c *gin.Context) { FleetBuildControllerInstance.UnassignShipModel(c) })

	apiRoute.GET("/ship-models", func(c *gin.Context) { ShipModelControllerInstance.GetAllShipModels(c) })
	apiRoute.GET("/ship-models/:id", func(c *gin.Context) { ShipModelControllerInstance.GetShipModel(c) })
	apiRoute.POST("/ship-models", func(c *gin.Context) { ShipModelControllerInstance.CreateShipModel(c) })
	apiRoute.PUT("/ship-models/:id", func(c *gin.Context) { ShipModelControllerInstance.UpdateShipModel(c) })
	apiRoute.DELETE("/ship-models/:id", func(c *gin.Context) { ShipModelControllerInstance.DeleteShipModel(c) })
}
