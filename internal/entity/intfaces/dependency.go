package intfaces

import "github.com/harmannkibue/golang_gin_clean_architecture/pkg/logger"

// Dependencies holds all injected dependencies -.
type Dependencies struct {
	Logger logger.Interface
	// Register all the usecases below for dependency injection -.
	BlogUsecase IntBlogUsecase
}
