package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "glaktika.eu/galaktika/docs"
	"glaktika.eu/galaktika/internal/di"
	"glaktika.eu/galaktika/internal/ws"
)

// @title Galaktika API
// @version 1.0
// @description API for Galaktika galaxy game
// @host localhost:8080
// @BasePath /api
func main() {
	router := gin.Default()

	di.RegisterWebRoutes(router)

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API endpoints
	apiRoute := router.Group("/api")

	di.CreateSingletons("dev")
	di.RegisterApiRoutes(apiRoute)

	// ws messages hub
	hub := ws.NewHub()
	go hub.Run()

	router.GET("/ws", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	}))

	_ = router.Run(":8080")
}
