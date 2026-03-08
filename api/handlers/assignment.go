package handlers

import (
	"net/http"
	"time"

	"pesu-cli/api/config"
	"pesu-cli/api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateAssignmentInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Language    string `json:"language" binding:"required,oneof=python c"`
	DueDate     string `json:"due_date" binding:"required"` // Format: RFC3339
}

func CreateAssignment(c *gin.Context) {
	var input CreateAssignmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	uid, _ := uuid.Parse(userID.(string))

	dueDate, err := time.Parse(time.RFC3339, input.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use RFC3339 (e.g., 2023-10-01T15:04:05Z)"})
		return
	}

	assignment := models.Assignment{
		Title:       input.Title,
		Description: input.Description,
		Language:    input.Language,
		DueDate:     dueDate,
		CreatedBy:   uid,
	}

	if result := config.DB.Create(&assignment); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create assignment"})
		return
	}

	c.JSON(http.StatusCreated, assignment)
}

func GetAssignments(c *gin.Context) {
	var assignments []models.Assignment
	if result := config.DB.Find(&assignments); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch assignments"})
		return
	}

	c.JSON(http.StatusOK, assignments)
}
