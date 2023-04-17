// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/controller/http/v1/blog_route"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase/blog_usecase"
	"github.com/harmannkibue/golang_gin_clean_architecture/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"

	// Swagger docs -.
	_ "github.com/harmannkibue/golang_gin_clean_architecture/docs"
)

// NewRouter -.
// Swagger spec:
// @title       Gin Gonic golang Clean Architecture.
// @description Illustration of uncle Bob's clean architecture using a demo blogs api.
// @description It serves as Blog.
// @version     1.0
// @host        localhost:8080
// @BasePath    /api/v1
// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization
func NewRouter(handler *gin.Engine, l logger.Interface, u *blog_usecase.BlogUseCase) {
	// Options -.
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	//// Swagger ui router group with basic authentication in implemented -.
	//doc := handler.Group("/swagger", gin.BasicAuth(gin.Accounts{
	//	"admin": "admin",
	//}))
	//
	//// Creating a swaggo instance -.
	//swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	//
	//doc.GET("/*any", swaggerHandler)

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe for kubernetes health checks -.
	handler.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "The server is up and running.Hurray Blog")
	})

	// Handling a page not found endpoint -.
	handler.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "The requested page is not found.Please try later!"})
	})

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers -.
	unversionedGroup := handler.Group("/api/v1")

	{
		blog_route.NewBlogRoute(unversionedGroup, u, l)
	}
}
