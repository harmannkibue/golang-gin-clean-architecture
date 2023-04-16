// Package app configures and runs application.
package app

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/harmannkibue/golang_gin_clean_architecture/config"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/controller/http/v1"
	db "github.com/harmannkibue/golang_gin_clean_architecture/internal/entity"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase/blog_usecase"
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
	handler := gin.Default()

	conn, err := postgres.New(cfg)

	if err != nil {
		fmt.Errorf("failed to connect to database %w", err)
	}

	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			panic("ERROR CLOSING POSTGRES CONNECTION")
		}
	}(conn)

	// Initializing a store for repository -.
	store := db.NewStore(conn)

	blogUsecase := blog_usecase.NewBlogUseCase(store, cfg)

	// Passing also the basic auth middleware to all  Routers -.
	v1.NewRouter(handler, l, *blogUsecase)

	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

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
