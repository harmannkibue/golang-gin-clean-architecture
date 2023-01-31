// Package app configures and runs application.
package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/harmannkibue/golang_gin_clean_architecture/config"
	v1 "github.com/harmannkibue/golang_gin_clean_architecture/internal/controller/http/v1"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase"
	_ "github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase"
	db "github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase/repositories"
	"github.com/harmannkibue/golang_gin_clean_architecture/pkg/httpserver"
	"github.com/harmannkibue/golang_gin_clean_architecture/pkg/logger"
	"github.com/harmannkibue/golang_gin_clean_architecture/pkg/postgres"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
	"syscall"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// HTTP Server -.
	router := gin.Default()

	conn, err := postgres.New(cfg)

	if err != nil {
		fmt.Errorf("failed to connect to database %w", err)
	}

	defer conn.Close()

	// Initializing a store for repository -.
	store := db.NewStore(conn)

	ledgerUsecase := usecase.NewBlogUseCase(store, cfg)

	// Passing also the basic auth middleware to all  Routers
	v1.NewRouter(router, l, *ledgerUsecase)

	//log.Fatal(router.Run(":" + cfg.Port))

	httpServer := httpserver.New(router, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal -.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
