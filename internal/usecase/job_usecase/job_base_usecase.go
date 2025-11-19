package job_usecase

import (
	"github.com/harmannkibue/actsml-jobs-orchestrator/config"
	"github.com/harmannkibue/actsml-jobs-orchestrator/internal/entity/intfaces"
	"github.com/harmannkibue/actsml-jobs-orchestrator/pkg/logger"
)

type JobUseCase struct {
	config *config.Config
	logger *logger.Logger
}

func NewJobUseCase(config *config.Config, l *logger.Logger) intfaces.IntJobUsecase {
	return &JobUseCase{
		config: config,
		logger: l,
	}
}
