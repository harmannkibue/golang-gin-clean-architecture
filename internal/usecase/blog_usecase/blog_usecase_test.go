package blog_usecase

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/harmannkibue/golang_gin_clean_architecture/config"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/entity/mocks"
	db "github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase/repositories/sqlc"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	testing2 "testing"
	"time"
)

func TestMockGettingBlog(t *testing2.T) {
	mockStore := new(mocks.Store)
	blogId, _ := uuid.Parse("93979a30-a3a9-4910-aa20-3fd5f14b69f9")
	mockBlog := db.Blog{
		ID:           blogId,
		Descriptions: sql.NullString{String: "Blog description", Valid: true},
		UserRole:     "author",
		CreatedAt:    time.Now(),
		UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
	}

	t.Run("success", func(t *testing2.T) {
		res := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(res)
		mockStore.On("GetBlog", ctx, blogId).Return(mockBlog, nil).Once()
		blogUsecase := NewBlogUseCase(mockStore, &config.Config{})

		blog, err := blogUsecase.GetBlog(ctx, blogId.String())

		assert.NoError(t, err)
		assert.NotNil(t, blog)

		mockStore.AssertExpectations(t)
	})
}
