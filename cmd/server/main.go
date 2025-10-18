package main

import (
	"fmt"
	"os"
	"time"

	"github.com/idestis/mock/internal/handlers"
	"github.com/idestis/mock/internal/scheduler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize scheduler manager
	manager := scheduler.NewManager()

	// Optional: Configure jobs via environment variables
	configureJobsFromEnv(manager)

	// Initialize handlers
	schedulerHandler := handlers.NewSchedulerHandler(manager)

	// Routes
	e.GET("/health", schedulerHandler.HealthCheck)
	e.POST("/api/v1/system/scheduler/trigger", schedulerHandler.TriggerJob)
	e.GET("/api/v1/system/scheduler/status", schedulerHandler.GetJobStatus)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}

// configureJobsFromEnv allows configuring job behavior via environment variables
func configureJobsFromEnv(manager *scheduler.Manager) {
	// Example: USER_ANNONYMIZATION_DURATION=10s
	// Example: USER_ANNONYMIZATION_SHOULD_FAIL=true

	jobs := []string{"user_annonymization", "zendesk_import", "blueshift_export"}

	for _, jobName := range jobs {
		durationKey := fmt.Sprintf("%s_DURATION", getEnvKey(jobName))
		failKey := fmt.Sprintf("%s_SHOULD_FAIL", getEnvKey(jobName))

		duration := 5 * time.Second
		if durationStr := os.Getenv(durationKey); durationStr != "" {
			if d, err := time.ParseDuration(durationStr); err == nil {
				duration = d
			}
		}

		shouldFail := false
		if failStr := os.Getenv(failKey); failStr == "true" {
			shouldFail = true
		}

		// Only update if environment variables are set
		if os.Getenv(durationKey) != "" || os.Getenv(failKey) != "" {
			manager.UpdateJobConfig(jobName, duration, shouldFail)
		}
	}
}

// getEnvKey converts job name to environment variable key format
func getEnvKey(jobName string) string {
	// Convert "user_annonymization" to "USER_ANNONYMIZATION"
	key := ""
	for _, c := range jobName {
		if c == '_' {
			key += "_"
		} else if c >= 'a' && c <= 'z' {
			key += string(c - 32)
		} else {
			key += string(c)
		}
	}
	return key
}
