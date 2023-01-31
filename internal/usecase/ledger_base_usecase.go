package usecase

import (
	"github.com/harmannkibue/golang_gin_clean_architecture/config"
	db "github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase/repositories"
)

type BlogUseCase struct {
	config *config.Config
	store  db.Store
}

func NewBlogUseCase(store db.Store, config *config.Config) *BlogUseCase {
	return &BlogUseCase{
		store:  store,
		config: config,
	}
}
