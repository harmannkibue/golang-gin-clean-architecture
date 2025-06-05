package intfaces

import (
	"context"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase/repository/sqlc"
)

type IntBlogUsecase interface {
	GetBlog(ctx context.Context, id string) (*sqlc.Blog, error)
	CreateBlog(ctx context.Context, description string) (*sqlc.Blog, error)
	ListBlogs(ctx context.Context, args ListBlogsParams) (*ListBlogsResponse, error)
}

type ListBlogsParams struct {
	Page  string `json:"page"`
	Limit string `json:"limit"`
}

type ListBlogsResponse struct {
	Blog         []sqlc.Blog `json:"blogs"`
	NextPage     string      `json:"next_page"`
	PreviousPage string      `json:"previous_page"`
}
