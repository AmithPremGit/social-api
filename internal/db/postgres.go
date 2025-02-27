package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver

	"social-api/internal/config"
)

// New creates a new database connection pool
func New(cfg config.DBConfig) (*sql.DB, error) {
	// Build connection string
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
	)

	// Create connection pool
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Set connection pool configuration
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.MaxIdleTime)

	// Create a context with a 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test the connection
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Transaction represents a database transaction
type Transaction struct {
	*sql.Tx
}

// ExecuteTx executes a function within a database transaction
func ExecuteTx(ctx context.Context, db *sql.DB, fn func(*Transaction) error) error {
	// Begin transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Wrap in Transaction struct
	txWrapper := &Transaction{tx}

	// Execute function
	err = fn(txWrapper)
	if err != nil {
		// Attempt to roll back on error
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	// Commit transaction
	return tx.Commit()
}
