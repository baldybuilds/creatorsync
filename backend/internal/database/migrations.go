package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

// MigrationRunner handles database migrations
type MigrationRunner struct {
	db *sql.DB
}

// NewMigrationRunner creates a new migration runner
func NewMigrationRunner(db *sql.DB) *MigrationRunner {
	return &MigrationRunner{db: db}
}

// RunMigrations executes all pending migrations
func (mr *MigrationRunner) RunMigrations(migrationsDir string) error {
	// Create migrations table if it doesn't exist
	if err := mr.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	files, err := mr.getMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Get applied migrations
	applied, err := mr.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Execute pending migrations
	for _, file := range files {
		filename := filepath.Base(file)
		if applied[filename] {
			log.Printf("Migration %s already applied, skipping", filename)
			continue
		}

		log.Printf("Applying migration: %s", filename)
		if err := mr.executeMigration(file, filename); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}
		log.Printf("Successfully applied migration: %s", filename)
	}

	return nil
}

// createMigrationsTable creates the migrations tracking table
func (mr *MigrationRunner) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`
	_, err := mr.db.Exec(query)
	return err
}

// getMigrationFiles returns sorted list of migration files
func (mr *MigrationRunner) getMigrationFiles(dir string) ([]string, error) {
	pattern := filepath.Join(dir, "*.sql")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	// Sort files to ensure they're executed in order
	sort.Strings(files)
	return files, nil
}

// getAppliedMigrations returns a map of already applied migrations
func (mr *MigrationRunner) getAppliedMigrations() (map[string]bool, error) {
	query := "SELECT filename FROM schema_migrations"
	rows, err := mr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, err
		}
		applied[filename] = true
	}

	return applied, rows.Err()
}

// executeMigration executes a single migration file
func (mr *MigrationRunner) executeMigration(filePath, filename string) error {
	// Read migration file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute migration in a transaction
	tx, err := mr.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Split content by semicolons and execute each statement
	statements := strings.Split(string(content), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute statement '%s': %w", stmt, err)
		}
	}

	// Record migration as applied
	if _, err := tx.Exec("INSERT INTO schema_migrations (filename) VALUES ($1)", filename); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return tx.Commit()
}
