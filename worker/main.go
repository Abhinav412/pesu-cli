package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pesu-cli/api/config"
	"pesu-cli/api/models"
	"pesu-cli/worker/processor"
)

type SubmissionJob struct {
	SubmissionID string `json:"submission_id"`
	AssignmentID string `json:"assignment_id"`
	Language     string `json:"language"`
	CodeBundle   []byte `json:"code_bundle"`
}

func main() {
	// Initialize Logic
	config.ConnectDB()
	config.ConnectRedis()

	fmt.Println("Worker started. Listening for submissions...")

	// Graceful Shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("Shutting down worker...")
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Pop from Redis (Blocking for 2 seconds)
			result, err := config.RDB.BRPop(ctx, 2*time.Second, "submissions").Result()
			if err != nil {
				// Timeout or error (Redis closed, etc)
				continue
			}

			// result[0] is key, result[1] is value
			payload := result[1]
			var job SubmissionJob
			if err := json.Unmarshal([]byte(payload), &job); err != nil {
				log.Printf("Error processing job parsing: %v", err)
				continue
			}

			fmt.Printf("Processing submission: %s\n", job.SubmissionID)

			// Update status to processing
			config.DB.Model(&models.Submission{}).Where("id = ?", job.SubmissionID).Update("status", "processing")

			// Execute Job
			output, score, status := processor.Execute(job.Language, job.CodeBundle)

			// Update Result
			config.DB.Model(&models.Submission{}).Where("id = ?", job.SubmissionID).Updates(map[string]interface{}{
				"status":       status,
				"score":        score,
				"feedback":     output,
				"processed_at": time.Now(),
			})

			fmt.Printf("Finished submission: %s (Status: %s)\n", job.SubmissionID, status)
		}
	}
}
