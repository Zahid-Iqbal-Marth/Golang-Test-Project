package routes

import (
	"io"

	"github.com/Zahid-Iqbal-Marth/Golang-Test-Project/handlers"
	"github.com/Zahid-Iqbal-Marth/Golang-Test-Project/utils"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the router with all routes and middleware
func SetupRouter(processor *utils.BatchProcessor) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	router := gin.New()
	router.Use(gin.Recovery())

	// Define routes
	router.GET("/healthz", handlers.HealthzHandler)
	router.POST("/log", handlers.LogHandler(processor))

	return router
}
