package intfaces

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
)

type KubernetesClient interface {
	CreateJob(ctx context.Context, namespace string, job *batchv1.Job) (*batchv1.Job, error)
	DeleteJob(ctx context.Context, namespace, name string) error
	GetJob(ctx context.Context, namespace, name string) (*batchv1.Job, error)
}
