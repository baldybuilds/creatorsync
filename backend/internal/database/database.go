package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Health() map[string]string
	Close() error
	GetDB() *sql.DB
	RunMigrations() error
	CheckConnection() error
	Reconnect() error
}

type service struct {
	db     *sql.DB
	connStr string
}

func New() Service {
	var connStr string

	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		connStr = databaseURL
		log.Println("Using DATABASE_URL for connection")
	} else {
		database := os.Getenv("POSTGRES_DB_DATABASE")
		password := os.Getenv("POSTGRES_DB_PASSWORD")
		username := os.Getenv("POSTGRES_DB_USERNAME")
		port := os.Getenv("POSTGRES_DB_PORT")
		host := os.Getenv("POSTGRES_DB_HOST")
		schema := os.Getenv("POSTGRES_DB_SCHEMA")

		if schema == "" {
			schema = "public"
		}

		connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require&search_path=%s",
			username, password, host, port, database, schema)
		log.Println("Using individual environment variables for connection")
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Configure connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(30 * time.Second)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		log.Printf("Failed to ping database: %v", err)
		log.Fatal(err)
	}

	log.Println("Database connection established successfully")

	return &service{
		db:     db,
		connStr: connStr,
	}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	if dbStats.OpenConnections > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

func (s *service) Close() error {
	log.Printf("Disconnected from database")
	return s.db.Close()
}

func (s *service) GetDB() *sql.DB {
	// Check if connection is healthy
	if err := s.CheckConnection(); err != nil {
		log.Printf("Database connection unhealthy: %v, attempting reconnect...", err)
		if reconnectErr := s.Reconnect(); reconnectErr != nil {
			log.Printf("Failed to reconnect to database: %v", reconnectErr)
			// Return the existing connection anyway - let the caller handle the error
		}
	}
	return s.db
}

func (s *service) RunMigrations() error {
	migrationRunner := NewMigrationRunner(s.db)
	return migrationRunner.RunMigrations("migrations")
}

func (s *service) CheckConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	return s.db.PingContext(ctx)
}

func (s *service) Reconnect() error {
	log.Println("Attempting to reconnect to database...")
	
	// Close the existing connection
	if s.db != nil {
		s.db.Close()
	}
	
	// Create a new connection
	db, err := sql.Open("pgx", s.connStr)
	if err != nil {
		return fmt.Errorf("failed to reconnect to database: %w", err)
	}
	
	// Configure connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(30 * time.Second)
	
	// Test the new connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database after reconnect: %w", err)
	}
	
	s.db = db
	log.Println("Database reconnected successfully")
	return nil
}
