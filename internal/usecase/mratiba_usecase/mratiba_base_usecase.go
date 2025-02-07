package mratiba_usecase

import (
	"github.com/harmannkibue/golang_gin_clean_architecture/config"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/entity/intfaces"
)

type MRatibaUseCase struct {
	config *config.Config
	store  intfaces.Store
}

func NewMRatibaUseCase(store intfaces.Store, config *config.Config) intfaces.IntMRatibaUsecase {
	return &MRatibaUseCase{
		store:  store,
		config: config,
	}
}
