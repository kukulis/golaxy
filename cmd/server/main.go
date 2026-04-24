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

	// Serve static files from assets directory
	router.Static("/assets", "./assets")
	router.StaticFile("/", "./assets/index.html")
	router.StaticFile("/divisions.html", "./assets/divisions.html")
	router.StaticFile("/test-ship-designs", "./assets/test_ship_designs.html")
	router.StaticFile("/test-ship-group-designs", "./assets/test_ship_group_designs.html")

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API endpoints
	apiRoute := router.Group("/api")

	di.CreateSingletons("dev")
	di.RegisterRoutes(apiRoute)

	_ = router.Run(":8080")
}
