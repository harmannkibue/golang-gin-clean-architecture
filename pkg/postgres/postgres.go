// Package postgres implements postgres connection.
package postgres

import (
	"database/sql"
	"github.com/harmannkibue/golang_gin_clean_architecture/config"
	"log"
)

const (
	_maxLifeTime            = 0
	_maxIdleConnections     = 50
	_maximumOpenConnections = 50
)

// New -.
func New(cfg *config.Config) (*sql.DB, error) {
	dbSource := BuildPostgresUrl(cfg)

	if len(dbSource) == 0 {
		log.Fatal("Postgres database url configs not passed")
	}

	// Opening a driver typically will not attempt to connect to the database.
	pool, err := sql.Open(cfg.DatabaseDriver, dbSource)

	if err != nil {
		// This will not be a connection error, but a dbSource parse error or
		// another initialization error.
		log.Fatal("unable to use data source name", err)
	}

	// Should we have the db closed after implementation -.
	//defer pool.Close()

	pool.SetConnMaxLifetime(_maxLifeTime)
	pool.SetMaxIdleConns(_maxIdleConnections)
	pool.SetMaxOpenConns(_maximumOpenConnections)

	return pool, nil
}
