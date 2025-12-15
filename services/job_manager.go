package services

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Davidrc26/api_zinc.git/models"
)

type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

type Job struct {
	ID        string                 `json:"id"`
	Status    JobStatus              `json:"status"`
	StartTime time.Time              `json:"start_time"`
	EndTime   *time.Time             `json:"end_time,omitempty"`
	Duration  string                 `json:"duration,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Result    map[string]interface{} `json:"result,omitempty"`
}

type JobManager struct {
	jobs map[string]*Job
	mu   sync.RWMutex
}

var manager = &JobManager{
	jobs: make(map[string]*Job),
}

// CreateJob crea un nuevo job y lo registra
func CreateJob() string {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	jobID := fmt.Sprintf("job_%d", time.Now().UnixNano())
	job := &Job{
		ID:        jobID,
		Status:    JobStatusPending,
		StartTime: time.Now(),
	}
	manager.jobs[jobID] = job
	return jobID
}

// GetJob obtiene el estado de un job
func GetJob(jobID string) (*Job, bool) {
	manager.mu.RLock()
	defer manager.mu.RUnlock()

	job, exists := manager.jobs[jobID]
	return job, exists
}

// UpdateJobStatus actualiza el estado de un job
func UpdateJobStatus(jobID string, status JobStatus, errorMsg string, result map[string]interface{}) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if job, exists := manager.jobs[jobID]; exists {
		job.Status = status
		if status == JobStatusCompleted || status == JobStatusFailed {
			endTime := time.Now()
			job.EndTime = &endTime
			job.Duration = endTime.Sub(job.StartTime).String()
		}
		if errorMsg != "" {
			job.Error = errorMsg
		}
		if result != nil {
			job.Result = result
		}
	}
}

// StartIndexingAsync inicia la indexación de forma asíncrona
func StartIndexingAsync(jobID string) {
	go func() {
		// Actualizar a running
		UpdateJobStatus(jobID, JobStatusRunning, "", nil)

		// Setup CPU profiling
		cpu, err := StartCPUProfiling()
		if err != nil {
			log.Println("Error creating CPU profile:", err)
			UpdateJobStatus(jobID, JobStatusFailed, "Error al crear perfil de CPU: "+err.Error(), nil)
			return
		}
		defer StopCPUProfiling(cpu)

		// Setup log file
		logFile, err := CreateLogFile()
		if err != nil {
			log.Println("Error opening log file:", err)
			UpdateJobStatus(jobID, JobStatusFailed, "Error al abrir archivo de log: "+err.Error(), nil)
			return
		}
		defer logFile.Close()

		log.SetOutput(logFile)
		start := time.Now()
		StartIndexing()
		end := time.Now()
		elapsed := end.Sub(start)

		log.Println("Time taken:", elapsed)
		log.Println("Execution finished")

		// Memory profiling
		memoryResult := CreateMemoryProfile()
		if memoryResult.Status != 200 {
			log.Println("Error creating memory profile:", memoryResult.Message)
			UpdateJobStatus(jobID, JobStatusFailed, memoryResult.Message, nil)
			return
		}

		// Job completado exitosamente
		result := map[string]interface{}{
			"time_taken": elapsed.String(),
			"message":    "Indexación completada exitosamente",
		}
		UpdateJobStatus(jobID, JobStatusCompleted, "", result)
	}()
}

// GetJobResponse retorna una respuesta estructurada para un job
func GetJobResponse(jobID string) models.Response {
	job, exists := GetJob(jobID)
	if !exists {
		return models.Response{
			Status:  404,
			Message: "Job no encontrado",
			Result:  nil,
		}
	}

	return models.Response{
		Status:  200,
		Message: "Estado del job obtenido exitosamente",
		Result: map[string]interface{}{
			"job": job,
		},
	}
}
