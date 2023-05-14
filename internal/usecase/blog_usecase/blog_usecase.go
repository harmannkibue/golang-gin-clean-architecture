package blog_usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/entity"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/entity/intfaces"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase/repositories/sqlc"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase/utils"
)

// GetBlog getting a single blog by id -.
func (usecase *BlogUseCase) GetBlog(ctx context.Context, id string) (*sqlc.Blog, error) {
	uuID, err := uuid.Parse(id)

	if err != nil {
		return nil, fmt.Errorf("BlogUsecase - blog - uc.usecase.GetBlog: %w", err)
	}

	blog, err := usecase.store.GetBlog(ctx, uuID)
	if err != nil {
		return nil, entity.CreateError(entity.ErrNotFound.Error(), err.Error())
	}

	return &blog, nil
}

func (usecase *BlogUseCase) CreateBlog(ctx context.Context, description string) (*sqlc.Blog, error) {

	blog, err := usecase.store.CreateBlog(ctx, sql.NullString{String: description, Valid: true})

	if err != nil {
		return nil, fmt.Errorf("BlogUsecase - uc.usecase.CreateBlog.: %w", err)
	}

	return &blog, nil
}

// ListBlogs -.
func (usecase *BlogUseCase) ListBlogs(ctx context.Context, args intfaces.ListBlogsParams) (*intfaces.ListBlogsResponse, error) {

	page, err := utils.StringToInt32(args.Page)

	if err != nil {
		return nil, errors.New("enter a valid type for the pageId query parameter")
	}

	limit, err := utils.StringToInt32(args.Limit)

	if err != nil {
		return nil, errors.New("enter a valid type for the pageSize query parameter")
	}

	Limit, Offset := utils.PaginatorParams(page, limit)

	blogs, err := usecase.store.ListBlog(ctx, sqlc.ListBlogParams{
		Limit:  Limit,
		Offset: Offset,
	})

	if err != nil {
		return nil, fmt.Errorf("BlogUsecase - bank - uc.usecase.ListBlogs: %w", err)
	}

	nextPage, previousPage := utils.PaginatorPages(ctx, page, limit, len(blogs))

	return &intfaces.ListBlogsResponse{Blog: blogs, NextPage: nextPage, PreviousPage: previousPage}, nil
}
