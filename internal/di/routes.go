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
	router.GET("/division/:divisionId/main.html", func(c *gin.Context) { webControllerInstance.RenderDivision(c) })
	router.GET("/division/:divisionId/fleet-builds.html", func(c *gin.Context) { webControllerInstance.RenderFleetBuilds(c) })
	router.GET("/fleet-build/:id/main.html", func(c *gin.Context) { webControllerInstance.RenderFleetBuild(c) })
	router.GET("/fleet-build/:id/edit.html", func(c *gin.Context) { webControllerInstance.RenderFleetBuildEdit(c) })
	router.GET("/ship-model-assignment/:id/main.html", func(c *gin.Context) { webControllerInstance.RenderShipModelAssignment(c) })
	router.GET("/ship-model/list.html", func(c *gin.Context) { webControllerInstance.RenderShipModelList(c) })
	router.GET("/ship-model/:id/details.html", func(c *gin.Context) { webControllerInstance.RenderShipModelDetails(c) })
	router.GET("/ship-model/:id/edit.html", func(c *gin.Context) { webControllerInstance.RenderShipModelEdit(c) })
	router.StaticFile("/test-ship-designs", "./pages/test_ship_designs.html")
	router.StaticFile("/test-ship-group-designs", "./pages/test_ship_group_designs.html")

	router.GET("/challenges.html", func(c *gin.Context) { webControllerInstance.RenderChallenges(c) })
	router.GET("/challenge/:id/edit.html", func(c *gin.Context) { webControllerInstance.RenderChallengeEdit(c) })
}

func RegisterApiRoutes(apiRoute *gin.RouterGroup) {
	apiRoute.GET("/current-race", func(c *gin.Context) { raceControllerInstance.GetCurrentRace(c) })
	apiRoute.GET("/races", func(c *gin.Context) { raceControllerInstance.GetAllRaces(c) })
	apiRoute.POST("/races", func(c *gin.Context) { raceControllerInstance.CreateRace(c) })
	apiRoute.GET("/races/:id", func(c *gin.Context) { raceControllerInstance.GetRace(c) })
	apiRoute.DELETE("/races/:id", func(c *gin.Context) { raceControllerInstance.DeleteRace(c) })

	apiRoute.GET("/battle", func(c *gin.Context) { battleControllerInstance.GetBattle(c) })

	apiRoute.GET("/divisions", func(c *gin.Context) { divisionControllerInstance.GetAllDivisions(c) })
	apiRoute.GET("/divisions/:id", func(c *gin.Context) { divisionControllerInstance.GetDivision(c) })
	apiRoute.POST("/divisions", func(c *gin.Context) { divisionControllerInstance.CreateDivision(c) })
	apiRoute.PUT("/divisions/:id", func(c *gin.Context) { divisionControllerInstance.UpdateDivision(c) })
	apiRoute.DELETE("/divisions/:id", func(c *gin.Context) { divisionControllerInstance.DeleteDivision(c) })

	apiRoute.GET("/fleet-builds", func(c *gin.Context) { fleetBuildControllerInstance.GetFleetBuilds(c) })
	apiRoute.GET("/fleet-builds/:id", func(c *gin.Context) { fleetBuildControllerInstance.GetFleetBuild(c) })
	apiRoute.POST("/fleet-builds", func(c *gin.Context) { fleetBuildControllerInstance.CreateFleetBuild(c) })
	apiRoute.PUT("/fleet-builds/:id", func(c *gin.Context) { fleetBuildControllerInstance.UpdateFleetBuild(c) })
	apiRoute.DELETE("/fleet-builds/:id", func(c *gin.Context) { fleetBuildControllerInstance.DeleteFleetBuild(c) })
	apiRoute.GET("/fleet-builds/:id/statistics", func(c *gin.Context) { fleetBuildControllerInstance.GetStatistics(c) })
	apiRoute.GET("/fleet-builds/:id/technologies", func(c *gin.Context) { fleetBuildControllerInstance.GetTechnologies(c) })
	apiRoute.GET("/fleet-builds/:id/ship-model-assignments", func(c *gin.Context) { fleetBuildControllerInstance.GetAssignedShipModels(c) })
	apiRoute.GET("/fleet-builds/:id/calculate-assignments-ship-techs", func(c *gin.Context) { fleetBuildControllerInstance.CalculateAssignmentsShipTechs(c) })
	apiRoute.GET("/ship-model-assignment/:id", func(c *gin.Context) { fleetBuildControllerInstance.GetShipModelAssignment(c) })
	apiRoute.POST("/ship-model-assignment/:id", func(c *gin.Context) { fleetBuildControllerInstance.UpdateShipModelAssignment(c) })
	apiRoute.POST("/ship-model-assignment", func(c *gin.Context) { fleetBuildControllerInstance.AddShipModelAssignment(c) })
	apiRoute.DELETE("/ship-models-assignment/:id", func(c *gin.Context) { fleetBuildControllerInstance.UnassignShipModel(c) })
	apiRoute.GET("/ship-models-assignment/:id/calculate-ship-tech", func(c *gin.Context) { fleetBuildControllerInstance.CalculateShipTech(c) })
	apiRoute.POST("/fleet-builds/:id/build", func(c *gin.Context) { fleetBuildControllerInstance.Build(c) })
	apiRoute.GET("/fleet-builds/:id/fleet", func(c *gin.Context) { fleetBuildControllerInstance.GetFleet(c) })

	apiRoute.GET("/challenges", func(c *gin.Context) { challengeControllerInstance.GetChallenges(c) })
	apiRoute.GET("/challenges/:id", func(c *gin.Context) { challengeControllerInstance.GetChallenge(c) })
	apiRoute.POST("/challenges", func(c *gin.Context) { challengeControllerInstance.CreateChallenge(c) })
	apiRoute.PUT("/challenges/:id", func(c *gin.Context) { challengeControllerInstance.UpdateChallenge(c) })
	apiRoute.DELETE("/challenges/:id", func(c *gin.Context) { challengeControllerInstance.DeleteChallenge(c) })

	apiRoute.GET("/ship-models", func(c *gin.Context) { shipModelControllerInstance.GetAllShipModels(c) })
	apiRoute.GET("/ship-models/:id", func(c *gin.Context) { shipModelControllerInstance.GetShipModel(c) })
	apiRoute.POST("/ship-models", func(c *gin.Context) { shipModelControllerInstance.CreateShipModel(c) })
	apiRoute.PUT("/ship-models/:id", func(c *gin.Context) { shipModelControllerInstance.UpdateShipModel(c) })
	apiRoute.POST("/ship-models/:id/calculate-ship-tech", func(c *gin.Context) { shipModelControllerInstance.CalculateShipTech(c) })
	apiRoute.DELETE("/ship-models/:id", func(c *gin.Context) { shipModelControllerInstance.DeleteShipModel(c) })
}

func RegisterWsRoutes(router *gin.Engine) {
	router.GET("/ws", func(c *gin.Context) { wsControllerInstance.ServeWs(c) })
}
