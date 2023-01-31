package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store This interface helps to mock the database during testing -.
type Store interface {
	Querier
}

// SqlStore provides all functions to execute db queries as well as transactions
type SqlStore struct {
	*Queries
	db *sql.DB
}

// NewStore SqlStore creates a new SqlStore
func NewStore(db *sql.DB) Store {
	return &SqlStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a callback function within a single transaction
func (store *SqlStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// The options part can be used to set up database isolation level.If nil then default will be
	//used which is read commited in postgres
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// calling the new with a sql transaction not sql.queries object
	q := New(tx)

	err = fn(q)
	if err != nil {
		// If there is an error in rollbacks
		if rbEr := tx.Rollback(); rbEr != nil {
			return fmt.Errorf("fn execTx rollback error: %v transaction error: %v", rbEr, err)
		}
		return fmt.Errorf("transaction error: %v", err)
	}
	return tx.Commit()
}
