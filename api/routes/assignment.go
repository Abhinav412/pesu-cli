package routes

import (
	"pesu-cli/api/handlers"
	"pesu-cli/api/middleware"

	"github.com/gin-gonic/gin"
)

func AssignmentRoutes(r *gin.RouterGroup) {
	assignments := r.Group("/assignments")
	{
		// Public or Authenticated (View)
		assignments.GET("", middleware.AuthMiddleware(), handlers.GetAssignments)

		// Teacher Only (Create)
		assignments.POST("", middleware.AuthMiddleware(), middleware.RoleMiddleware("teacher", "admin"), handlers.CreateAssignment)
	}
}
