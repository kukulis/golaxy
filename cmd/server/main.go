package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Serve static files from assets directory
	router.Static("/assets", "./assets")
	router.StaticFile("/", "./assets/index.html")

	// Future: API endpoints will go here
	// api := router.Group("/api")
	// {
	//     api.GET("/fleet", getFleet)
	//     api.POST("/fleet", createFleet)
	// }

	router.Run(":8080")
}
