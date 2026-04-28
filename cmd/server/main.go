package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "glaktika.eu/galaktika/docs"
	"glaktika.eu/galaktika/internal/di"
)

// @title Galaktika API
// @version 1.0
// @description API for Galaktika galaxy game
// @host localhost:8080
// @BasePath /api
func main() {
	router := gin.Default()

	router.Static("/assets", "./assets")
	router.StaticFile("/", "./pages/index.html")
	router.StaticFile("/dummy_login.html", "./pages/dummy_login.html")
	router.StaticFile("/divisions.html", "./pages/divisions.html")
	router.GET("/division/:divisionId/main.html", func(c *gin.Context) { c.File("./pages/division/main.html") })
	router.GET("/division/:divisionId/fleet-builds.html", func(c *gin.Context) { c.File("./pages/division/fleet-builds.html") })
	router.GET("/fleet-build/:id/main.html", func(c *gin.Context) { c.File("./pages/division/fleet-build/main.html") })
	router.StaticFile("/ship-model/list.html", "./pages/ship-model/list.html")
	router.GET("/ship-model/:id/details.html", func(c *gin.Context) { c.File("./pages/ship-model/details.html") })
	router.StaticFile("/test-ship-designs", "./pages/test_ship_designs.html")
	router.StaticFile("/test-ship-group-designs", "./pages/test_ship_group_designs.html")

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API endpoints
	apiRoute := router.Group("/api")

	di.CreateSingletons("dev")
	di.RegisterRoutes(apiRoute)

	_ = router.Run(":8080")
}
