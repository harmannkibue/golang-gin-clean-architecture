package intfaces

import (
	"context"
)

type IntJobUsecase interface {
	CreateJob(ctx context.Context, description string) error
}
