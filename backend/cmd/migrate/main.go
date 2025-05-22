package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/baldybuilds/creatorsync/internal/database"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	log.Println("Starting database migrations...")

	// Initialize database connection
	db := database.New()
	defer db.Close()

	// Test database connection
	health := db.Health()
	if health["status"] != "up" {
		log.Fatalf("Database connection failed: %v", health)
	}

	log.Println("Database connection successful")

	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	// Get list of migration files
	migrations, err := getMigrationFiles("migrations")
	if err != nil {
		log.Fatalf("Failed to get migration files: %v", err)
	}

	if len(migrations) == 0 {
		log.Println("No migration files found")
		return
	}

	log.Printf("Found %d migration files", len(migrations))

	// Run migrations
	for _, migration := range migrations {
		if err := runMigration(db, migration); err != nil {
			log.Fatalf("Failed to run migration %s: %v", migration, err)
		}
	}

	log.Println("All migrations completed successfully!")
}

func createMigrationsTable(db database.Service) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			filename VARCHAR(255) NOT NULL UNIQUE,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := db.GetDB().Exec(query)
	return err
}

func getMigrationFiles(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var migrations []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrations = append(migrations, file.Name())
		}
	}

	// Sort migrations alphabetically to ensure order
	sort.Strings(migrations)
	return migrations, nil
}

func runMigration(db database.Service, filename string) error {
	// Check if migration has already been run
	var count int
	err := db.GetDB().QueryRow("SELECT COUNT(*) FROM migrations WHERE filename = $1", filename).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check migration status: %w", err)
	}

	if count > 0 {
		log.Printf("Migration %s already executed, skipping", filename)
		return nil
	}

	log.Printf("Running migration: %s", filename)

	// Read migration file
	content, err := os.ReadFile(filepath.Join("migrations", filename))
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute migration
	_, err = db.GetDB().Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Record migration as executed
	_, err = db.GetDB().Exec("INSERT INTO migrations (filename) VALUES ($1)", filename)
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	log.Printf("Migration %s completed successfully", filename)
	return nil
}
