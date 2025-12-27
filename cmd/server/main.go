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

	// API endpoints
	apiRoute := router.Group("/api")

	di.CreateSingletons("dev")
	di.RegisterRoutes(apiRoute)

	router.Run(":8080")
}
