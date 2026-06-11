package di

import "github.com/gin-gonic/gin"

func RegisterWebRoutes(router *gin.Engine) {
	router.LoadHTMLFiles(
		"./pages/division/main.gohtml",
		"./pages/division/fleet-builds.gohtml",
		"./pages/division/fleet-build/main.gohtml",
		"./pages/division/fleet-build/edit.gohtml",
		"./pages/ship-model/list.gohtml",
		"./pages/ship-model/details.gohtml",
		"./pages/ship-model/edit.gohtml",
		"./pages/division/fleet-build/ship-model-assignment/main.gohtml",
		"./pages/challenges.gohtml",
		"./pages/challenge/edit.gohtml",
	)

	router.Static("/assets", "./assets")
	router.StaticFile("/", "./pages/index.html")
	router.StaticFile("/dummy_login.html", "./pages/dummy_login.html")
	router.StaticFile("/divisions.html", "./pages/divisions.html")
	router.StaticFile("/races.html", "./pages/races.html")
	router.GET("/division/:divisionId/main.html", func(c *gin.Context) { WebControllerInstance.RenderDivision(c) })
	router.GET("/division/:divisionId/fleet-builds.html", func(c *gin.Context) { WebControllerInstance.RenderFleetBuilds(c) })
	router.GET("/fleet-build/:id/main.html", func(c *gin.Context) { WebControllerInstance.RenderFleetBuild(c) })
	router.GET("/fleet-build/:id/edit.html", func(c *gin.Context) { WebControllerInstance.RenderFleetBuildEdit(c) })
	router.GET("/ship-model-assignment/:id/main.html", func(c *gin.Context) { WebControllerInstance.RenderShipModelAssignment(c) })
	router.GET("/ship-model/list.html", func(c *gin.Context) { WebControllerInstance.RenderShipModelList(c) })
	router.GET("/ship-model/:id/details.html", func(c *gin.Context) { WebControllerInstance.RenderShipModelDetails(c) })
	router.GET("/ship-model/:id/edit.html", func(c *gin.Context) { WebControllerInstance.RenderShipModelEdit(c) })
	router.StaticFile("/test-ship-designs", "./pages/test_ship_designs.html")
	router.StaticFile("/test-ship-group-designs", "./pages/test_ship_group_designs.html")

	router.GET("/challenges.html", func(c *gin.Context) { WebControllerInstance.RenderChallenges(c) })
	router.GET("/challenge/:id/edit.html", func(c *gin.Context) { WebControllerInstance.RenderChallengeEdit(c) })
}

func RegisterApiRoutes(apiRoute *gin.RouterGroup) {
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

	apiRoute.GET("/fleet-builds", func(c *gin.Context) { FleetBuildControllerInstance.GetFleetBuilds(c) })
	apiRoute.GET("/fleet-builds/:id", func(c *gin.Context) { FleetBuildControllerInstance.GetFleetBuild(c) })
	apiRoute.POST("/fleet-builds", func(c *gin.Context) { FleetBuildControllerInstance.CreateFleetBuild(c) })
	apiRoute.PUT("/fleet-builds/:id", func(c *gin.Context) { FleetBuildControllerInstance.UpdateFleetBuild(c) })
	apiRoute.DELETE("/fleet-builds/:id", func(c *gin.Context) { FleetBuildControllerInstance.DeleteFleetBuild(c) })
	apiRoute.GET("/fleet-builds/:id/statistics", func(c *gin.Context) { FleetBuildControllerInstance.GetStatistics(c) })
	apiRoute.GET("/fleet-builds/:id/technologies", func(c *gin.Context) { FleetBuildControllerInstance.GetTechnologies(c) })
	apiRoute.GET("/fleet-builds/:id/ship-model-assignments", func(c *gin.Context) { FleetBuildControllerInstance.GetAssignedShipModels(c) })
	apiRoute.GET("/fleet-builds/:id/calculate-assignments-ship-techs", func(c *gin.Context) { FleetBuildControllerInstance.CalculateAssignmentsShipTechs(c) })
	apiRoute.GET("/ship-model-assignment/:id", func(c *gin.Context) { FleetBuildControllerInstance.GetShipModelAssignment(c) })
	apiRoute.POST("/ship-model-assignment/:id", func(c *gin.Context) { FleetBuildControllerInstance.UpdateShipModelAssignment(c) })
	apiRoute.POST("/ship-model-assignment", func(c *gin.Context) { FleetBuildControllerInstance.AddShipModelAssignment(c) })
	apiRoute.DELETE("/ship-models-assignment/:id", func(c *gin.Context) { FleetBuildControllerInstance.UnassignShipModel(c) })
	apiRoute.GET("/ship-models-assignment/:id/calculate-ship-tech", func(c *gin.Context) { FleetBuildControllerInstance.CalculateShipTech(c) })
	apiRoute.POST("/fleet-builds/:id/build", func(c *gin.Context) { FleetBuildControllerInstance.Build(c) })
	apiRoute.GET("/fleet-builds/:id/fleet", func(c *gin.Context) { FleetBuildControllerInstance.GetFleet(c) })

	apiRoute.GET("/challenges", func(c *gin.Context) { ChallengeControllerInstance.GetChallenges(c) })
	apiRoute.GET("/challenges/:id", func(c *gin.Context) { ChallengeControllerInstance.GetChallenge(c) })
	apiRoute.POST("/challenges", func(c *gin.Context) { ChallengeControllerInstance.CreateChallenge(c) })
	apiRoute.PUT("/challenges/:id", func(c *gin.Context) { ChallengeControllerInstance.UpdateChallenge(c) })
	apiRoute.DELETE("/challenges/:id", func(c *gin.Context) { ChallengeControllerInstance.DeleteChallenge(c) })

	apiRoute.GET("/ship-models", func(c *gin.Context) { ShipModelControllerInstance.GetAllShipModels(c) })
	apiRoute.GET("/ship-models/:id", func(c *gin.Context) { ShipModelControllerInstance.GetShipModel(c) })
	apiRoute.POST("/ship-models", func(c *gin.Context) { ShipModelControllerInstance.CreateShipModel(c) })
	apiRoute.PUT("/ship-models/:id", func(c *gin.Context) { ShipModelControllerInstance.UpdateShipModel(c) })
	apiRoute.POST("/ship-models/:id/calculate-ship-tech", func(c *gin.Context) { ShipModelControllerInstance.CalculateShipTech(c) })
	apiRoute.DELETE("/ship-models/:id", func(c *gin.Context) { ShipModelControllerInstance.DeleteShipModel(c) })
}
