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
	db      *sql.DB
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

		// Configure SSL mode (defaults to require for production)
		sslMode := os.Getenv("POSTGRES_SSL_MODE")
		if sslMode == "" {
			sslMode = "require"
		}

		connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s",
			username, password, host, port, database, sslMode, schema)
		log.Println("Using individual environment variables for connection")
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Configure connection pool settings based on environment
	env := os.Getenv("APP_ENV")
	environment := os.Getenv("ENVIRONMENT")
	databaseURL := os.Getenv("DATABASE_URL")

	// Improved environment detection
	isStaging := env == "staging" || environment == "staging" ||
		env == "dev" || environment == "dev" ||
		(databaseURL != "" && (env != "production" && environment != "production"))

	isProduction := env == "production" || environment == "production"

	log.Printf("🌍 Database Environment Detection: APP_ENV=%s, ENVIRONMENT=%s, DATABASE_URL=%t",
		env, environment, databaseURL != "")

	if isProduction {
		// Production settings - most conservative
		db.SetMaxOpenConns(8)
		db.SetMaxIdleConns(2)
		db.SetConnMaxLifetime(1 * time.Minute)
		db.SetConnMaxIdleTime(10 * time.Second)
		log.Printf("Applied production database pool settings (MaxOpen: 8, MaxIdle: 2)")
	} else if isStaging {
		// Staging settings - conservative for cloud
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(3)
		db.SetConnMaxLifetime(2 * time.Minute)
		db.SetConnMaxIdleTime(15 * time.Second)
		log.Printf("Applied staging database pool settings (MaxOpen: 10, MaxIdle: 3)")
	} else {
		// Development settings - more permissive
		db.SetMaxOpenConns(15)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)
		db.SetConnMaxIdleTime(30 * time.Second)
		log.Printf("Applied development database pool settings (MaxOpen: 15, MaxIdle: 5)")
	}

	// Test the connection with retries
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			log.Printf("Database ping attempt %d failed: %v", i+1, err)
			if i == maxRetries-1 {
				log.Fatal(err)
			}
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		break
	}

	log.Println("Database connection established successfully")

	return &service{
		db:      db,
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

	// Configure connection pool settings based on environment
	env := os.Getenv("APP_ENV")
	environment := os.Getenv("ENVIRONMENT")
	databaseURL := os.Getenv("DATABASE_URL")

	// Improved environment detection
	isStaging := env == "staging" || environment == "staging" ||
		env == "dev" || environment == "dev" ||
		(databaseURL != "" && (env != "production" && environment != "production"))

	isProduction := env == "production" || environment == "production"

	log.Printf("🌍 Database Environment Detection: APP_ENV=%s, ENVIRONMENT=%s, DATABASE_URL=%t",
		env, environment, databaseURL != "")

	if isProduction {
		// Production settings - most conservative
		db.SetMaxOpenConns(8)
		db.SetMaxIdleConns(2)
		db.SetConnMaxLifetime(1 * time.Minute)
		db.SetConnMaxIdleTime(10 * time.Second)
		log.Printf("Applied production database pool settings (MaxOpen: 8, MaxIdle: 2)")
	} else if isStaging {
		// Staging settings - conservative for cloud
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(3)
		db.SetConnMaxLifetime(2 * time.Minute)
		db.SetConnMaxIdleTime(15 * time.Second)
		log.Printf("Applied staging database pool settings (MaxOpen: 10, MaxIdle: 3)")
	} else {
		// Development settings - more permissive
		db.SetMaxOpenConns(15)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)
		db.SetConnMaxIdleTime(30 * time.Second)
		log.Printf("Applied development database pool settings (MaxOpen: 15, MaxIdle: 5)")
	}

	// Test the new connection with retries
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			log.Printf("Database reconnect ping attempt %d failed: %v", i+1, err)
			if i == maxRetries-1 {
				return fmt.Errorf("failed to ping database after reconnect: %w", err)
			}
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		break
	}

	s.db = db
	log.Println("Database reconnected successfully")
	return nil
}
