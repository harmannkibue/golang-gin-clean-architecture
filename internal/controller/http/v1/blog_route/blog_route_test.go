package blog_route

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/entity/mocks"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase/repositories/sqlc"
	"github.com/harmannkibue/golang_gin_clean_architecture/pkg/logger"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

//https://dev.to/jacobsngoodwin/04-testing-first-gin-http-handler-9m0
func TestGetByID(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockBlogUsecase := new(mocks.BlogUsecase)

	mockBlog := sqlc.Blog{
		ID:           uuid.New(),
		Descriptions: sql.NullString{String: "Test Description ", Valid: true},
		UserRole:     sqlc.UserRolesAuthor,
		CreatedAt:    time.Now(),
		UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
	}

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	t.Run("Success", func(t *testing.T) {

		mockBlogUsecase.On("GetBlog", c, mockBlog.ID.String()).Return(&mockBlog, nil)

		req, err := http.NewRequestWithContext(c, http.MethodGet, "/blogs/"+mockBlog.ID.String(), strings.NewReader(""))
		assert.NoError(t, err)

		c.Request = req
		c.Params = append(c.Params, gin.Param{Key: "id", Value: mockBlog.ID.String()})

		handler := BlogRoute{
			u: mockBlogUsecase,
			l: logger.New("info"),
		}

		handler.blog(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockBlogUsecase.AssertExpectations(t)
	})

}
