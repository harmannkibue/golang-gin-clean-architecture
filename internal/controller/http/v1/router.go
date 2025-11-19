// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/harmannkibue/actsml-jobs-orchestrator/internal/controller/http/v1/job_route"
	"github.com/harmannkibue/actsml-jobs-orchestrator/internal/entity/intfaces"
	"github.com/harmannkibue/actsml-jobs-orchestrator/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"

	// Swagger docs -.
	_ "github.com/harmannkibue/actsml-jobs-orchestrator/docs"
)

// NewRouter -.
// Swagger spec:
// @title       ACTSML Gin Gonic golang Clean Architecture Orchestrator.
// @description Illustration of uncle Bob's clean architecture using a demo blogs api.
// @description It serves as Blog.
// @version     1.0
// @host        localhost:8080
// @BasePath    /api/v1
// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization
func NewRouter(handler *gin.Engine, l logger.Interface, u intfaces.Dependencies) {
	// Options -.
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	//// Swagger ui router group with basic authentication in implemented -.
	doc := handler.Group("/swagger", gin.BasicAuth(gin.Accounts{
		"admin": "admin",
	}))

	// Creating a swaggo instance -.
	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)

	doc.GET("/*any", swaggerHandler)

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
		job_route.NewJobRoute(unversionedGroup, u.JobUsecase, l)
	}
}
