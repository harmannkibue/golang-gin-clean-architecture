// Package app configures and runs the application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/harmannkibue/actsml-jobs-orchestrator/config"
	"github.com/harmannkibue/actsml-jobs-orchestrator/internal/controller/http/v1"
	"github.com/harmannkibue/actsml-jobs-orchestrator/internal/entity/intfaces"
	"github.com/harmannkibue/actsml-jobs-orchestrator/internal/usecase/job_usecase"
	"github.com/harmannkibue/actsml-jobs-orchestrator/internal/usecase/microservices"

	"github.com/harmannkibue/actsml-jobs-orchestrator/pkg/httpserver"
	"github.com/harmannkibue/actsml-jobs-orchestrator/pkg/logger"
)

func Run(cfg *config.Config) {
	// ---------- LOGGER ----------
	l := logger.New(cfg.Log.Level)
	l.Info("Starting ACTSML Job Orchestrator...")

	// ---------- HTTP HANDLER ----------
	gin.SetMode(cfg.HTTP.Mode)
	handler := gin.New()
	handler.Use(gin.Logger(), gin.Recovery())

	// ---------- K8S CLIENT ----------
	k8sClient, err := microservices.NewK8sClient()
	if err != nil {
		l.Fatal(fmt.Errorf("failed to initialize Kubernetes client: %w", err))
	}

	// ---------- USE CASES ----------
	jobUC := job_usecase.NewJobUseCase(cfg, l, k8sClient)

	// ---------- DEPENDENCY CONTAINER ----------
	deps := intfaces.Dependencies{
		Logger:     l,
		JobUsecase: jobUC,
		K8sClient:  k8sClient,
	}

	// ---------- ROUTES ----------
	v1.NewRouter(handler, l, deps)

	// ---------- SERVER ----------
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// ---------- GRACEFUL SHUTDOWN ----------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		l.Info("Received shutdown signal: " + sig.String())

	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("server error: %w", err))
	}

	// ---------- FINAL SHUTDOWN ----------
	if err := httpServer.Shutdown(); err != nil {
		l.Error(fmt.Errorf("http server shutdown error: %w", err))
	}

	l.Info("ACTSML Orchestrator stopped gracefully.")
}
