package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"pesu-cli/api/config"
	"pesu-cli/api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubmissionJob struct {
	SubmissionID string `json:"submission_id"`
	AssignmentID string `json:"assignment_id"`
	Language     string `json:"language"`
	CodeBundle   []byte `json:"code_bundle"`
}

func SubmitAssignment(c *gin.Context) {
	// 1. Parse Multipart Form
	file, _, err := c.Request.FormFile("bundle")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File 'bundle' is required"})
		return
	}
	defer file.Close()

	// 2. Read File Content
	content, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	assignmentIDStr := c.PostForm("assignment_id")
	if assignmentIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "assignment_id is required"})
		return
	}
	assignmentID, err := uuid.Parse(assignmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment_id"})
		return
	}

	userID, _ := c.Get("userID")
	uid, _ := uuid.Parse(userID.(string))

	// 3. Create Submission Record
	submission := models.Submission{
		UserID:       uid,
		AssignmentID: assignmentID,
		Status:       "queued",
	}

	if result := config.DB.Create(&submission); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create submission record"})
		return
	}

	// 4. Fetch Assignment Details (for Language)
	var assignment models.Assignment
	if result := config.DB.First(&assignment, assignmentID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	// 5. Push to Redis
	job := SubmissionJob{
		SubmissionID: submission.ID.String(),
		AssignmentID: assignment.ID.String(),
		Language:     assignment.Language,
		CodeBundle:   content,
	}

	jobBytes, err := json.Marshal(job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize job"})
		return
	}

	// Push to "submissions" queue
	err = config.RDB.LPush(context.Background(), "submissions", jobBytes).Err()
	if err != nil {
		// Rollback DB creation if Redis fails? Or just fail?
		// For simplicity, we just fail. Ideally we'd retry or use a transaction.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue submission: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Submission queued",
		"submission_id": submission.ID,
		"status":        submission.Status,
	})
}

func GetSubmissions(c *gin.Context) {
	// Filter by User or Assignment if needed
	userID, _ := c.Get("userID")

	var submissions []models.Submission
	if result := config.DB.Where("user_id = ?", userID).Find(&submissions); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch submissions"})
		return
	}

	c.JSON(http.StatusOK, submissions)
}
