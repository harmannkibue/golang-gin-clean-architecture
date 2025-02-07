// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/controller/http/v1/mpesa_ratiba_route"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/entity/intfaces"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"

	// Swagger docs -.
	_ "github.com/harmannkibue/golang_gin_clean_architecture/docs"
)

// NewRouter -.
// Swagger spec:
// @title       Spectacular SAS Billing.
// @description Spectacular Billing Golang Payment Services Connector Clean Architecture using Uncle Bob's clean architecture principles.
// @description It serves as Blog.
// @version     1.0
// @host        localhost:8080
// @BasePath    /api/v1
// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization
func NewRouter(handler *gin.Engine, dependecies intfaces.Dependencies) {
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
	versionedGroup := handler.Group("/api/v1")

	// Feature-based route registration
	RegisterRoutes(versionedGroup, dependecies)
}

// RegisterRoutes dynamically adds routes based on available dependencies.
func RegisterRoutes(api *gin.RouterGroup, deps intfaces.Dependencies) {
	// Register M-Pesa Ratiba routes
	if deps.MRatibaUsecase != nil {
		mpesa_ratiba_route.NewMRatibaRoute(api, deps.MRatibaUsecase, deps.Logger)
	}
}
