// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/gin-gonic/gin"
	_ "github.com/harmannkibue/golang_gin_clean_architecture/docs"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/controller/http/v1/blog_route"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase/blog_usecase"
	"github.com/harmannkibue/golang_gin_clean_architecture/pkg/logger"
	"net/http"
	// Swagger docs.
	swaggerFiles "github.com/swaggo/files"     // gin-swagger middleware
	ginSwagger "github.com/swaggo/gin-swagger" // swagger embed files
)

// NewRouter -.
// Swagger spec:
// @title       Gin Gonic golang Clean Architecture.
// @description Illustration of uncle Bob's clean architecture using a demo blogs api.
// @description It serves as Blog.
// @version     1.0
// @host        localhost:8089
// @BasePath    /api/v1
// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization
func NewRouter(handler *gin.Engine, l logger.Interface, u blog_usecase.BlogUseCase) {
	// Options -.
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger ui router group with basic authentication in implemented -.
	doc := handler.Group("/swagger", gin.BasicAuth(gin.Accounts{
		"admin": "admin",
	}))

	// Creating a swaggo instance -.
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	doc.GET("/*any", swaggerHandler)

	// K8s probe for kubernetes health checks -.
	handler.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "The server is up and running.Hurray blog")
	})

	// Handling a page not found endpoint -.
	handler.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "The requested page is not found.Please try later!"})
	})

	// Routers -.
	unversionedGroup := handler.Group("/api/v1")

	{
		blog_route.NewBlogRoute(unversionedGroup, u, l)
	}
}
