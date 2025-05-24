package database

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func mustStartPostgresContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	var (
		dbName = "database"
		dbPwd  = "password"
		dbUser = "user"
	)

	dbContainer, err := postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	os.Setenv("POSTGRES_DB_DATABASE", dbName)
	os.Setenv("POSTGRES_DB_PASSWORD", dbPwd)
	os.Setenv("POSTGRES_DB_USERNAME", dbUser)

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	os.Setenv("POSTGRES_DB_HOST", dbHost)
	os.Setenv("POSTGRES_DB_PORT", dbPort.Port())
	// Disable SSL for test containers
	os.Setenv("POSTGRES_SSL_MODE", "disable")

	return dbContainer.Terminate, err
}

func TestMain(m *testing.M) {
	// Check if we should skip integration tests by looking for short flag
	for _, arg := range os.Args {
		if strings.Contains(arg, "short") {
			log.Println("Skipping database integration tests in short mode")
			os.Exit(0)
		}
	}

	teardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	result := m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown postgres container: %v", err)
	}

	os.Exit(result)
}

func TestNew(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database integration test in short mode")
	}

	srv := New()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database integration test in short mode")
	}

	srv := New()

	stats := srv.Health()

	if stats["status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["status"])
	}

	if _, ok := stats["error"]; ok {
		t.Fatalf("expected error not to be present")
	}

	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}

func TestClose(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database integration test in short mode")
	}

	srv := New()

	if srv.Close() != nil {
		t.Fatalf("expected Close() to return nil")
	}
}
