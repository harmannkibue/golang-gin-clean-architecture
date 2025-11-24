package intfaces

import (
	"context"
	"encoding/json"
	"time"

	batchv1 "k8s.io/api/batch/v1"
)

// CreateJobResult represents the result of creating a job
type CreateJobResult struct {
	JobID       string    `json:"job_id"`
	Status      string    `json:"status"`
	K8sJobName  string    `json:"k8s_job_name"`
	SubmittedAt time.Time `json:"submitted_at"`
	Namespace   string    `json:"namespace"`
	JobUID      string    `json:"job_uid"`
	Submitted   bool      `json:"submitted"`
}

type IntJobUsecase interface {
	CreateJob(ctx context.Context, payload json.RawMessage) (*CreateJobResult, error)
	GetJobStatus(ctx context.Context, name string) (*batchv1.Job, error)
	DeleteJob(ctx context.Context, name string) error
}
