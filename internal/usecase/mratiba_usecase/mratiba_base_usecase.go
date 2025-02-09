package mratiba_usecase

import (
	"github.com/harmannkibue/spectabill_psp_connector_clean_architecture/config"
	"github.com/harmannkibue/spectabill_psp_connector_clean_architecture/internal/entity/intfaces"
)

type MRatibaUseCase struct {
	config        *config.Config
	store         intfaces.Store
	darajaFactory intfaces.DarajaFactory // New: Factory to create tenant-specific Daraja instances
}

func NewMRatibaUseCase(store intfaces.Store, config *config.Config, darajaFactory intfaces.DarajaFactory) intfaces.IntMRatibaUsecase {
	return &MRatibaUseCase{
		store:         store,
		config:        config,
		darajaFactory: darajaFactory, // Inject the factory
	}
}
