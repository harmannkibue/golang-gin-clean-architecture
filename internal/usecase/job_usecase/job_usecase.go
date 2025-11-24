package job_usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/harmannkibue/actsml-jobs-orchestrator/config"
	"github.com/harmannkibue/actsml-jobs-orchestrator/internal/entity"
	"github.com/harmannkibue/actsml-jobs-orchestrator/internal/entity/intfaces"
	"github.com/harmannkibue/actsml-jobs-orchestrator/pkg/logger"
	batchv1 "k8s.io/api/batch/v1"
)

// JobUseCase represents the job usecase with its dependencies
type JobUseCase struct {
	cfg    *config.Config
	logger logger.Interface
	k8s    intfaces.KubernetesClient
}

// NewJobUseCase creates a new instance of JobUseCase
func NewJobUseCase(cfg *config.Config, l logger.Interface, k8s intfaces.KubernetesClient) *JobUseCase {
	return &JobUseCase{
		cfg:    cfg,
		logger: l,
		k8s:    k8s,
	}
}

// CreateJob creates a new Kubernetes job from the provided payload
func (uc *JobUseCase) CreateJob(ctx context.Context, payload json.RawMessage) (*intfaces.CreateJobResult, error) {
	// Validate payload
	if len(payload) == 0 {
		return nil, entity.CreateError(entity.ErrBadRequest.Error(), "payload cannot be empty")
	}

	// Validate JSON structure
	var payloadMap map[string]interface{}
	if err := json.Unmarshal(payload, &payloadMap); err != nil {
		return nil, entity.CreateError(entity.ErrBadRequest.Error(), fmt.Sprintf("invalid JSON payload: %v", err))
	}

	// Generate unique job ID and name
	jobID := uuid.New().String()
	jobName := fmt.Sprintf("actsml-job-%s", jobID)

	// Marshal payload to string for environment variable
	payloadStr := string(payload)

	// Build environment variables
	envMap := map[string]string{
		"PAYLOAD_JSON": payloadStr,
	}

	// Get namespace from config or use default
	namespace := uc.getNamespace()

	// Get worker image from config or use default
	workerImage := uc.getWorkerImage()

	// Build Kubernetes Job manifest
	backoffLimit := int32(3)
	job := BuildK8sJob(jobName, workerImage, envMap, backoffLimit, namespace)

	// Submit job to Kubernetes
	createdJob, err := uc.k8s.CreateJob(ctx, namespace, job)
	if err != nil {
		uc.logger.Error(fmt.Errorf("failed to create job in Kubernetes: %w", err), "job_usecase - CreateJob")
		return nil, entity.CreateError(entity.ErrInternalServerError.Error(), fmt.Sprintf("failed to submit job: %v", err))
	}

	uc.logger.Info(fmt.Sprintf("Job created successfully: %s in namespace %s", jobName, namespace))

	return &intfaces.CreateJobResult{
		JobID:       jobID,
		Status:      "submitted",
		K8sJobName:  jobName,
		SubmittedAt: time.Now(),
		Namespace:   namespace,
		JobUID:      string(createdJob.UID),
		Submitted:   true,
	}, nil
}

// GetJobStatus retrieves the status of a Kubernetes job
func (uc *JobUseCase) GetJobStatus(ctx context.Context, name string) (*batchv1.Job, error) {
	if name == "" {
		return nil, entity.CreateError(entity.ErrBadRequest.Error(), "job name cannot be empty")
	}

	namespace := uc.getNamespace()
	job, err := uc.k8s.GetJob(ctx, namespace, name)
	if err != nil {
		uc.logger.Error(fmt.Errorf("failed to get job status: %w", err), "job_usecase - GetJobStatus")
		return nil, entity.CreateError(entity.ErrNotFound.Error(), fmt.Sprintf("job not found: %s", name))
	}

	return job, nil
}

// DeleteJob deletes a Kubernetes job
func (uc *JobUseCase) DeleteJob(ctx context.Context, name string) error {
	if name == "" {
		return entity.CreateError(entity.ErrBadRequest.Error(), "job name cannot be empty")
	}

	namespace := uc.getNamespace()
	err := uc.k8s.DeleteJob(ctx, namespace, name)
	if err != nil {
		uc.logger.Error(fmt.Errorf("failed to delete job: %w", err), "job_usecase - DeleteJob")
		return entity.CreateError(entity.ErrInternalServerError.Error(), fmt.Sprintf("failed to delete job: %v", err))
	}

	uc.logger.Info(fmt.Sprintf("Job deleted successfully: %s in namespace %s", name, namespace))
	return nil
}

// getNamespace returns the namespace to use for jobs
func (uc *JobUseCase) getNamespace() string {
	// Default namespace, can be overridden via config in the future
	return "default"
}

// getWorkerImage returns the worker image to use
func (uc *JobUseCase) getWorkerImage() string {
	// Default worker image, can be overridden via config in the future
	return "ghcr.io/afriqsiliconltd/actsml-base-worker:staging"
}
