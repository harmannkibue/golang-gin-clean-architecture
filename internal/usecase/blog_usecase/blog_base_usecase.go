package blog_usecase

import (
	"github.com/harmannkibue/golang_gin_clean_architecture/config"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/entity/intfaces"
)

type BlogUseCase struct {
	config *config.Config
	store  intfaces.Store
}

func NewBlogUseCase(store intfaces.Store, config *config.Config) intfaces.IntBlogUsecase {
	return &BlogUseCase{
		store:  store,
		config: config,
	}
}
