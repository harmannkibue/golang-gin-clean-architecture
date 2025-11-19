package intfaces

import "github.com/harmannkibue/actsml-jobs-orchestrator/pkg/logger"

// Dependencies holds all injected dependencies -.
type Dependencies struct {
	Logger logger.Interface
	// Register all the usecases below for dependency injection -.
	JobUsecase IntJobUsecase
}
