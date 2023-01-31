// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/gin-gonic/gin"
	_ "github.com/harmannkibue/golang_gin_clean_architecture/docs"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase"
	"github.com/harmannkibue/golang_gin_clean_architecture/pkg/logger"
	"net/http"
	// Swagger docs.
	swaggerFiles "github.com/swaggo/files"     // gin-swagger middleware
	ginSwagger "github.com/swaggo/gin-swagger" // swagger embed files
)

// NewRouter -.
// Swagger spec:
// @title       Ledgers Clean Architecture.
// @description Ledgers api endpoints for the Accounts Ledgers.
// @description It serves as Accounts Ledgers.
// @version     1.0
// @host        localhost:8089
// @BasePath    /api/v1
// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization
func NewRouter(handler *gin.Engine, l logger.Interface, u usecase.LedgerUseCase) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger ui router group with basic authentication in implemented -.
	doc := handler.Group("/swagger", gin.BasicAuth(gin.Accounts{
		"churpyLedger": "admin",
	}))

	// Creating a swaggo instance -.
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	doc.GET("/*any", swaggerHandler)

	// K8s probe for kubernetes health checks -.
	handler.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "The virtual accounts server is up and running!Hurray Churpy")
	})

	// Handling a page not found endpoint -.
	handler.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "The requested page is not found.Please try later!"})
	})

	// Routers -.
	unversionedGroup := handler.Group("/api/v1")
	{
		newVirtualAccountsRoute(unversionedGroup, u, l)
	}
}
