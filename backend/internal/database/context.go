package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// DBContext represents a request-scoped database context
type DBContext struct {
	conn    *sql.DB
	tx      *sql.Tx
	timeout time.Duration
}

// TxOptions defines transaction options
type TxOptions struct {
	Isolation sql.IsolationLevel
	ReadOnly  bool
	Timeout   time.Duration
}

// NewDBContext creates a new request-scoped database context
func NewDBContext(conn *sql.DB, timeout time.Duration) *DBContext {
	return &DBContext{
		conn:    conn,
		timeout: timeout,
	}
}

// Query executes a query with timeout
func (dbc *DBContext) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if dbc.tx != nil {
		return dbc.tx.QueryContext(ctx, query, args...)
	}
	return dbc.conn.QueryContext(ctx, query, args...)
}

// QueryRow executes a query that returns a single row
func (dbc *DBContext) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if dbc.tx != nil {
		return dbc.tx.QueryRowContext(ctx, query, args...)
	}
	return dbc.conn.QueryRowContext(ctx, query, args...)
}

// Exec executes a query without returning rows
func (dbc *DBContext) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if dbc.tx != nil {
		return dbc.tx.ExecContext(ctx, query, args...)
	}
	return dbc.conn.ExecContext(ctx, query, args...)
}

// BeginTx starts a new transaction with options
func (dbc *DBContext) BeginTx(ctx context.Context, opts *TxOptions) error {
	if dbc.tx != nil {
		return fmt.Errorf("transaction already in progress")
	}

	var sqlOpts *sql.TxOptions
	if opts != nil {
		sqlOpts = &sql.TxOptions{
			Isolation: opts.Isolation,
			ReadOnly:  opts.ReadOnly,
		}
	}

	tx, err := dbc.conn.BeginTx(ctx, sqlOpts)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	dbc.tx = tx
	log.Printf("🔄 Transaction started for context")
	return nil
}

// Commit commits the current transaction
func (dbc *DBContext) Commit() error {
	if dbc.tx == nil {
		return fmt.Errorf("no transaction in progress")
	}

	err := dbc.tx.Commit()
	dbc.tx = nil

	if err != nil {
		log.Printf("❌ Transaction commit failed: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("✅ Transaction committed successfully")
	return nil
}

// Rollback rolls back the current transaction
func (dbc *DBContext) Rollback() error {
	if dbc.tx == nil {
		return nil // No transaction to rollback
	}

	err := dbc.tx.Rollback()
	dbc.tx = nil

	if err != nil {
		log.Printf("❌ Transaction rollback failed: %v", err)
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	log.Printf("🔄 Transaction rolled back")
	return nil
}

// InTransaction returns true if currently in a transaction
func (dbc *DBContext) InTransaction() bool {
	return dbc.tx != nil
}

// WithTransaction executes a function within a transaction
func (dbc *DBContext) WithTransaction(ctx context.Context, opts *TxOptions, fn func(*DBContext) error) error {
	if err := dbc.BeginTx(ctx, opts); err != nil {
		return err
	}

	defer func() {
		if dbc.tx != nil {
			dbc.Rollback()
		}
	}()

	if err := fn(dbc); err != nil {
		return err
	}

	return dbc.Commit()
}

// Close ensures any open transaction is rolled back
func (dbc *DBContext) Close() {
	if dbc.tx != nil {
		dbc.Rollback()
	}
}
