package main

import (
	"pesu-cli/api/config"
	"pesu-cli/api/middleware"
	"pesu-cli/api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Configurations
	config.ConnectDB()
	config.ConnectRedis()

	r := gin.Default()

	// Default Route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// API Group
	api := r.Group("/api")
	{
		routes.AuthRoutes(api)
		routes.AssignmentRoutes(api)
		routes.SubmissionRoutes(api)
	}

	// Protected Routes Example
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Add protected routes here
	}

	r.Run()
}
