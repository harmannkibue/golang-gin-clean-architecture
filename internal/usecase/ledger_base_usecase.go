package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/harmannkibue/golang_gin_clean_architecture/config"
	db "github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase/repositories"
)

type LedgerUseCase struct {
	config *config.Config
	store  db.Store
	router *gin.Engine
}

func NewLedgerUseCase(store db.Store, config *config.Config) *LedgerUseCase {
	return &LedgerUseCase{
		store:  store,
		config: config,
		router: gin.New(),
	}
}
