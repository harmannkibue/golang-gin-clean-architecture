package tests

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockDb "github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase/mock"
	db "github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase/repositories"
	"net/http/httptest"
	testing2 "testing"
	"time"
)

func TestListingBlog(t *testing2.T) {

	blog := randomBlog()
	fmt.Println("THE id siiiis s", blog.ID)
	ctr := gomock.NewController(t)
	defer ctr.Finish()

	store := mockDb.NewMockStore(ctr)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	// Todo: Failing by missing calls
	//	 Build the stubs -.
	store.EXPECT().
		GetBlog(c, gomock.Eq(blog)).
		AnyTimes().
		Return(blog, nil)

	//Todo: This is the main issue I am facing
	//////	start the test server and send a http request for testing -.
	//server := httpserver.New(gin.Default(), httpserver.Port("8080"))
	//url := fmt.Sprintf("/api/v1/blogs/:%s", blog.ID)
	//request, err := http.NewRequest(http.MethodGet, url, nil)
	//require.NoError(t, err)

}

func randomBlog() db.Blog {
	return db.Blog{
		ID:           uuid.New(),
		Descriptions: sql.NullString{String: "Mock Test", Valid: true},
		UserRole:     db.UserRolesAuthor,
		CreatedAt:    time.Now(),
	}
}
