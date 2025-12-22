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

	// Future: API endpoints will go here
	apiRoute := router.Group("/api")
	// {
	//     api.GET("/fleet", getFleet)
	//     api.POST("/fleet", createFleet)
	// }

	di.CreateSingletons()
	apiRoute.GET("/battle", func(c *gin.Context) { di.BattleControllerInstance.GetBattle(c) })

	router.Run(":8080")
}
