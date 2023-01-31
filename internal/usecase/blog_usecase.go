package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	db "github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase/repositories"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase/utils"
)

// GetBlog getting a single blog by id -.
func (usecase *LedgerUseCase) GetBlog(ctx context.Context, id string) (*db.Blog, error) {
	uuID, err := uuid.Parse(id)

	if err != nil {
		return nil, fmt.Errorf("BlogUsecase - blog - uc.usecase.GetBlog: %w", err)
	}

	blog, err := usecase.store.GetBlog(ctx, uuID)
	if err != nil {
		return nil, fmt.Errorf("BlogUsecase - blog - uc.repo.GetBlog: %w", err)
	}

	return &blog, nil
}

func (usecase *LedgerUseCase) CreateBlog(ctx context.Context, description string) (*db.Blog, error) {

	blog, err := usecase.store.CreateBlog(ctx, sql.NullString{String: description, Valid: true})

	if err != nil {
		return nil, fmt.Errorf("BlogUsecase - store - uc.usecase.CreateBlog.: %w", err)
	}

	return &blog, nil
}

type ListBlogsParams struct {
	Page  string `json:"page"`
	Limit string `json:"limit"`
}

type listBlogsResponse struct {
	Blog         []db.Blog `json:"blogs"`
	NextPage     string    `json:"next_page"`
	PreviousPage string    `json:"previous_page"`
}

// ListBlogs -.
func (usecase *LedgerUseCase) ListBlogs(ctx context.Context, args ListBlogsParams) (*listBlogsResponse, error) {

	page, err := utils.StringToInt32(args.Page)

	if err != nil {
		return nil, errors.New("enter a valid type for the pageId query parameter")
	}

	limit, err := utils.StringToInt32(args.Limit)

	if err != nil {
		return nil, errors.New("enter a valid type for the pageSize query parameter")
	}

	Limit, Offset := utils.PaginatorParams(page, limit)

	blogs, err := usecase.store.ListBlog(ctx, db.ListBlogParams{
		Limit:  Limit,
		Offset: Offset,
	})

	if err != nil {
		return nil, fmt.Errorf("BankUsecase - bank - uc.usecase.GetBanks: %w", err)
	}

	nextPage, previousPage := utils.PaginatorPages(ctx, page, limit, len(blogs))

	return &listBlogsResponse{blogs, nextPage, previousPage}, nil
}
