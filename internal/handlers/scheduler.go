package handlers

import (
	"net/http"

	"github.com/idestis/mock/internal/models"
	"github.com/idestis/mock/internal/scheduler"
	"github.com/labstack/echo/v4"
)

// SchedulerHandler handles scheduler-related HTTP requests
type SchedulerHandler struct {
	manager *scheduler.Manager
}

// NewSchedulerHandler creates a new scheduler handler
func NewSchedulerHandler(manager *scheduler.Manager) *SchedulerHandler {
	return &SchedulerHandler{
		manager: manager,
	}
}

// TriggerJob handles POST /api/v1/system/scheduler/trigger
func (h *SchedulerHandler) TriggerJob(c echo.Context) error {
	var req models.JobTriggerRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if req.JobName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "job_name is required",
		})
	}

	job, err := h.manager.TriggerJob(req.JobName)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	response := models.JobTriggerResponse{
		JobID:   job.ID,
		Message: "Job triggered successfully",
		Status:  job.Status,
	}

	return c.JSON(http.StatusAccepted, response)
}

// GetJobStatus handles GET /api/v1/system/scheduler/status?job_id=xxx
func (h *SchedulerHandler) GetJobStatus(c echo.Context) error {
	jobID := c.QueryParam("job_id")

	if jobID == "" {
		// Return all jobs if no job_id provided
		jobs := h.manager.GetAllJobs()
		return c.JSON(http.StatusOK, map[string]interface{}{
			"jobs":  jobs,
			"count": len(jobs),
		})
	}

	job, err := h.manager.GetJobStatus(jobID)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.JobStatusResponse{
			Message: err.Error(),
		})
	}

	response := models.JobStatusResponse{
		Job: job,
	}

	return c.JSON(http.StatusOK, response)
}

// HealthCheck handles GET /health
func (h *SchedulerHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":         "healthy",
		"service":        "mock",
		"available_jobs": h.manager.GetAvailableJobs(),
	})
}
