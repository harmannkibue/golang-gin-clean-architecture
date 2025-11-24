package job_usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/harmannkibue/actsml-jobs-orchestrator/config"
	"github.com/harmannkibue/actsml-jobs-orchestrator/internal/entity"
	"github.com/harmannkibue/actsml-jobs-orchestrator/internal/entity/intfaces"
	"github.com/harmannkibue/actsml-jobs-orchestrator/pkg/logger"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	// Validate payload against ACTSML v1 specification
	validatedPayload, err := ValidatePayload(payload)
	if err != nil {
		uc.logger.Error(fmt.Errorf("payload validation failed: %w", err), "job_usecase - CreateJob")
		return nil, entity.CreateError(entity.ErrBadRequest.Error(), fmt.Sprintf("payload validation failed: %v", err))
	}

	uc.logger.Info(fmt.Sprintf("Validated payload for project_id=%s, experiment_id=%s, algorithm=%s",
		validatedPayload.ProjectID, validatedPayload.ExperimentID, validatedPayload.Problem.Algorithm))

	// Generate unique job ID and name
	jobID := uuid.New().String()
	jobName := fmt.Sprintf("actsml-job-%s", jobID)

	// Get namespace from config or use default
	namespace := uc.getNamespace()

	// Marshal payload to JSON string
	// This ensures proper JSON encoding (handles special characters, etc.)
	payloadBytes, err := json.Marshal(validatedPayload)
	if err != nil {
		uc.logger.Error(fmt.Errorf("failed to marshal validated payload: %w", err), "job_usecase - CreateJob")
		return nil, entity.CreateError(entity.ErrInternalServerError.Error(), "failed to process payload")
	}
	payloadStr := string(payloadBytes)

	// Create ConfigMap with payload (worker expects PAYLOAD_FILE pointing to a file)
	configMapName := fmt.Sprintf("%s-payload", jobName)
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: namespace,
			Labels: map[string]string{
				"managed-by": "actsml-orchestrator",
				"job-name":   jobName,
			},
		},
		Data: map[string]string{
			"payload.json": payloadStr,
		},
	}

	// Create ConfigMap before creating the job
	_, err = uc.k8s.CreateConfigMap(ctx, namespace, configMap)
	if err != nil {
		uc.logger.Error(fmt.Errorf("failed to create ConfigMap for payload: %w", err), "job_usecase - CreateJob")
		return nil, entity.CreateError(entity.ErrInternalServerError.Error(), fmt.Sprintf("failed to create payload ConfigMap: %v", err))
	}

	// Build environment variables for the worker
	envMap := map[string]string{
		"PAYLOAD_FILE": "/workspace/payload.json", // Worker expects file path
		// MinIO configuration (can be overridden via config in the future)
		"MINIO_ENDPOINT":  uc.getMinIOEndpoint(),
		"MINIO_ACCESS_KEY": uc.getMinIOAccessKey(),
		"MINIO_SECRET_KEY": uc.getMinIOSecretKey(),
		"MINIO_SECURE":     uc.getMinIOSecure(),
	}

	// Get worker image from config or use default
	workerImage := uc.getWorkerImage()

	// Get image pull secret name
	imagePullSecretName := uc.getImagePullSecretName()

	// Build Kubernetes Job manifest with compute resources from payload
	// Pass ConfigMap name for volume mounting
	backoffLimit := int32(3)
	job := BuildK8sJob(jobName, workerImage, envMap, backoffLimit, namespace, &validatedPayload.Compute, imagePullSecretName, configMapName)
	
	// Add project and experiment labels for better tracking
	if job.Labels == nil {
		job.Labels = make(map[string]string)
	}
	job.Labels["project-id"] = validatedPayload.ProjectID
	job.Labels["experiment-id"] = validatedPayload.ExperimentID
	job.Spec.Template.Labels["project-id"] = validatedPayload.ProjectID
	job.Spec.Template.Labels["experiment-id"] = validatedPayload.ExperimentID

	// Submit job to Kubernetes
	createdJob, err := uc.k8s.CreateJob(ctx, namespace, job)
	if err != nil {
		uc.logger.Error(fmt.Errorf("failed to create job in Kubernetes: %w", err), "job_usecase - CreateJob")
		return nil, entity.CreateError(entity.ErrInternalServerError.Error(), fmt.Sprintf("failed to submit job: %v", err))
	}

	uc.logger.Info(fmt.Sprintf("Job created successfully: %s in namespace %s", jobName, namespace))

	return &intfaces.CreateJobResult{
		JobID:        jobID,
		Status:       "submitted",
		K8sJobName:   jobName,
		SubmittedAt:  time.Now(),
		Namespace:    namespace,
		JobUID:       string(createdJob.UID),
		Submitted:    true,
		ProjectID:    validatedPayload.ProjectID,
		ExperimentID: validatedPayload.ExperimentID,
	}, nil
}

// GetJobStatus retrieves the status of a Kubernetes job -.
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
	// Can be overridden via environment variable or config
	if ns := os.Getenv("K8S_NAMESPACE"); ns != "" {
		return ns
	}
	// Default to actsml namespace for ACTSML workspace
	return "actsml"
}

// getWorkerImage returns the worker image to use
func (uc *JobUseCase) getWorkerImage() string {
	// Default worker image, can be overridden via config in the future
	return "ghcr.io/afriqsiliconltd/actsml-worker-image:staging"
}

// getMinIOEndpoint returns the MinIO endpoint
func (uc *JobUseCase) getMinIOEndpoint() string {
	// Can be overridden via config or environment variable
	if endpoint := os.Getenv("MINIO_ENDPOINT"); endpoint != "" {
		return endpoint
	}
	return "http://api.staging.minio.actsml.com"
}

// getMinIOAccessKey returns the MinIO access key
func (uc *JobUseCase) getMinIOAccessKey() string {
	if key := os.Getenv("MINIO_ACCESS_KEY"); key != "" {
		return key
	}
	return "actsMl"
}

// getMinIOSecretKey returns the MinIO secret key
func (uc *JobUseCase) getMinIOSecretKey() string {
	if key := os.Getenv("MINIO_SECRET_KEY"); key != "" {
		return key
	}
	return "6m35xip2UX50SpKh"
}

// getMinIOSecure returns whether MinIO uses secure connection
func (uc *JobUseCase) getMinIOSecure() string {
	if secure := os.Getenv("MINIO_SECURE"); secure != "" {
		return secure
	}
	return "false"
}

// getImagePullSecretName returns the name of the image pull secret to use
func (uc *JobUseCase) getImagePullSecretName() string {
	// Can be overridden via environment variable
	if secretName := os.Getenv("IMAGE_PULL_SECRET_NAME"); secretName != "" {
		return secretName
	}
	// Default to ghcr-pull-secret for GHCR authentication
	return "ghcr-pull-secret"
}
