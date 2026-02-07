package main

import (
	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/di"
)

func main() {
	router := gin.Default()

	// Serve static files from assets directory
	router.Static("/assets", "./assets")
	router.StaticFile("/", "./assets/index.html")
	router.StaticFile("/test-ship-designs", "./assets/test_ship_designs.html")
	router.StaticFile("/test-ship-group-designs", "./assets/test_ship_group_designs.html")
	// API endpoints
	apiRoute := router.Group("/api")

	di.CreateSingletons("dev")
	di.RegisterRoutes(apiRoute)

	_ = router.Run(":8080")
}
