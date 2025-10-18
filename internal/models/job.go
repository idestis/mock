package models

import (
	"time"
)

// JobStatus represents the current status of a job
type JobStatus string

const (
	StatusPending JobStatus = "PENDING"
	StatusRunning JobStatus = "RUNNING"
	StatusSuccess JobStatus = "SUCCESS"
	StatusFailed  JobStatus = "FAILED"
)

// JobConfig holds the configuration for a job
type JobConfig struct {
	Name          string        `json:"name"`
	ExecutionTime time.Duration `json:"execution_time"`
	ShouldFail    bool          `json:"should_fail"`
}

// Job represents a job execution instance
type Job struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Status     JobStatus  `json:"status"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	Error      string     `json:"error,omitempty"`
}

// JobTriggerRequest represents the request to trigger a job
type JobTriggerRequest struct {
	JobName string `json:"job_name" validate:"required"`
}

// JobTriggerResponse represents the response after triggering a job
type JobTriggerResponse struct {
	JobID   string    `json:"job_id"`
	Message string    `json:"message"`
	Status  JobStatus `json:"status"`
}

// JobStatusResponse represents the response for job status
type JobStatusResponse struct {
	Job     *Job   `json:"job,omitempty"`
	Message string `json:"message,omitempty"`
}
