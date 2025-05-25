package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type RequestConnection struct {
	db             *sql.DB
	tx             *sql.Tx
	defaultTimeout time.Duration
	inTransaction  bool
	committed      bool
	rolledBack     bool
}

type TransactionOptions struct {
	Isolation sql.IsolationLevel
	ReadOnly  bool
	Timeout   time.Duration
}

func NewRequestConnection(db *sql.DB, defaultTimeout time.Duration) *RequestConnection {
	return &RequestConnection{
		db:             db,
		defaultTimeout: defaultTimeout,
	}
}

func (rc *RequestConnection) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if rc.tx != nil {
		return rc.tx.QueryContext(ctx, query, args...)
	}
	return rc.db.QueryContext(ctx, query, args...)
}

func (rc *RequestConnection) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if rc.tx != nil {
		return rc.tx.QueryRowContext(ctx, query, args...)
	}
	return rc.db.QueryRowContext(ctx, query, args...)
}

func (rc *RequestConnection) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if rc.tx != nil {
		return rc.tx.ExecContext(ctx, query, args...)
	}
	return rc.db.ExecContext(ctx, query, args...)
}

func (rc *RequestConnection) BeginTx(ctx context.Context, opts *TransactionOptions) error {
	if rc.inTransaction {
		return fmt.Errorf("transaction already in progress")
	}

	var sqlOpts *sql.TxOptions
	if opts != nil {
		sqlOpts = &sql.TxOptions{
			Isolation: opts.Isolation,
			ReadOnly:  opts.ReadOnly,
		}

		if opts.Timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
			defer cancel()
		}
	}

	tx, err := rc.db.BeginTx(ctx, sqlOpts)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	rc.tx = tx
	rc.inTransaction = true
	log.Printf("🔄 Transaction started for request")
	return nil
}

func (rc *RequestConnection) Commit() error {
	if !rc.inTransaction {
		return fmt.Errorf("no transaction in progress")
	}

	if rc.committed {
		return fmt.Errorf("transaction already committed")
	}

	if rc.rolledBack {
		return fmt.Errorf("transaction already rolled back")
	}

	err := rc.tx.Commit()
	if err != nil {
		log.Printf("❌ Transaction commit failed: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	rc.committed = true
	rc.inTransaction = false
	rc.tx = nil
	log.Printf("✅ Transaction committed successfully")
	return nil
}

func (rc *RequestConnection) Rollback() error {
	if !rc.inTransaction {
		return nil
	}

	if rc.committed {
		return fmt.Errorf("cannot rollback committed transaction")
	}

	if rc.rolledBack {
		return nil
	}

	err := rc.tx.Rollback()
	if err != nil {
		log.Printf("❌ Transaction rollback failed: %v", err)
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	rc.rolledBack = true
	rc.inTransaction = false
	rc.tx = nil
	log.Printf("🔄 Transaction rolled back")
	return nil
}

func (rc *RequestConnection) IsInTransaction() bool {
	return rc.inTransaction
}

func (rc *RequestConnection) WithTransaction(ctx context.Context, opts *TransactionOptions, fn func(*RequestConnection) error) error {
	if err := rc.BeginTx(ctx, opts); err != nil {
		return err
	}

	defer func() {
		if rc.inTransaction && !rc.committed {
			rc.Rollback()
		}
	}()

	if err := fn(rc); err != nil {
		return err
	}

	return rc.Commit()
}

func (rc *RequestConnection) Close() error {
	if rc.inTransaction {
		if err := rc.Rollback(); err != nil {
			log.Printf("⚠️ Failed to rollback transaction during close: %v", err)
		}
	}

	log.Printf("🔌 Request connection closed")
	return nil
}

func (rc *RequestConnection) Ping(ctx context.Context) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), rc.defaultTimeout)
		defer cancel()
	}

	return rc.db.PingContext(ctx)
}
