package scheduler

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/idestis/mock/internal/models"
	"github.com/google/uuid"
)

// Manager handles job scheduling and execution
type Manager struct {
	jobs       map[string]*models.Job
	jobConfigs map[string]models.JobConfig
	mu         sync.RWMutex
}

// NewManager creates a new scheduler manager with predefined jobs
func NewManager() *Manager {
	m := &Manager{
		jobs:       make(map[string]*models.Job),
		jobConfigs: make(map[string]models.JobConfig),
	}

	// Initialize hardcoded jobs with default configurations
	m.jobConfigs["user_annonymization"] = models.JobConfig{
		Name:          "user_annonymization",
		ExecutionTime: 5 * time.Second,
		ShouldFail:    false,
	}
	m.jobConfigs["zendesk_import"] = models.JobConfig{
		Name:          "zendesk_import",
		ExecutionTime: 8 * time.Second,
		ShouldFail:    false,
	}
	m.jobConfigs["blueshift_export"] = models.JobConfig{
		Name:          "blueshift_export",
		ExecutionTime: 10 * time.Second,
		ShouldFail:    true, // This job is configured to fail by default
	}

	return m
}

// UpdateJobConfig updates the configuration for a specific job
func (m *Manager) UpdateJobConfig(jobName string, executionTime time.Duration, shouldFail bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	config, exists := m.jobConfigs[jobName]
	if !exists {
		return fmt.Errorf("job %s not found", jobName)
	}

	config.ExecutionTime = executionTime
	config.ShouldFail = shouldFail
	m.jobConfigs[jobName] = config

	return nil
}

// TriggerJob starts a job execution
func (m *Manager) TriggerJob(jobName string) (*models.Job, error) {
	m.mu.Lock()

	config, exists := m.jobConfigs[jobName]
	if !exists {
		m.mu.Unlock()
		return nil, fmt.Errorf("job %s not found", jobName)
	}

	job := &models.Job{
		ID:        uuid.New().String(),
		Name:      jobName,
		Status:    models.StatusPending,
		StartedAt: time.Now(),
	}

	m.jobs[job.ID] = job
	m.mu.Unlock()

	// Execute job in background
	go m.executeJob(job.ID, config)

	return job, nil
}

// executeJob simulates job execution
func (m *Manager) executeJob(jobID string, config models.JobConfig) {
	m.mu.Lock()
	job, exists := m.jobs[jobID]
	if !exists {
		m.mu.Unlock()
		return
	}
	job.Status = models.StatusRunning
	m.mu.Unlock()

	// Simulate work with sleep
	time.Sleep(config.ExecutionTime)

	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	job.FinishedAt = &now

	if config.ShouldFail {
		// Simulate random failure
		job.Status = models.StatusFailed
		job.Error = fmt.Sprintf("Job %s failed during execution", config.Name)
	} else {
		// Random chance of failure even for non-configured failures (10% chance)
		if rand.Float32() < 0.1 {
			job.Status = models.StatusFailed
			job.Error = fmt.Sprintf("Job %s encountered an unexpected error", config.Name)
		} else {
			job.Status = models.StatusSuccess
		}
	}
}

// GetJobStatus retrieves the status of a job by ID
func (m *Manager) GetJobStatus(jobID string) (*models.Job, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	job, exists := m.jobs[jobID]
	if !exists {
		return nil, fmt.Errorf("job with ID %s not found", jobID)
	}

	return job, nil
}

// GetAllJobs returns all jobs
func (m *Manager) GetAllJobs() []*models.Job {
	m.mu.RLock()
	defer m.mu.RUnlock()

	jobs := make([]*models.Job, 0, len(m.jobs))
	for _, job := range m.jobs {
		jobs = append(jobs, job)
	}

	return jobs
}

// GetAvailableJobs returns the list of available job names
func (m *Manager) GetAvailableJobs() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	jobNames := make([]string, 0, len(m.jobConfigs))
	for name := range m.jobConfigs {
		jobNames = append(jobNames, name)
	}

	return jobNames
}
