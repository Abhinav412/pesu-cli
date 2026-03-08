package routes

import (
	"pesu-cli/api/handlers"
	"pesu-cli/api/middleware"

	"github.com/gin-gonic/gin"
)

func SubmissionRoutes(r *gin.RouterGroup) {
	submissions := r.Group("/submissions")
	{
		submissions.POST("/", middleware.AuthMiddleware(), handlers.SubmitAssignment)
		submissions.GET("/", middleware.AuthMiddleware(), handlers.GetSubmissions)
	}
}
