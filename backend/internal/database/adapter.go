package database

import (
	"context"
	"database/sql"
	"log"
)

// ServiceAdapter provides backward compatibility by wrapping the new database system
type ServiceAdapter struct {
	dbService DatabaseInterface
}

// NewServiceAdapter creates a new service adapter that uses the modern database system
func NewServiceAdapter() (Service, error) {
	dbService, err := NewDatabaseService()
	if err != nil {
		log.Printf("⚠️ Failed to create new database service, creating standard fallback: %v", err)
		return createStandardService(), nil
	}

	return &ServiceAdapter{dbService: dbService}, nil
}

// Health returns health information - delegates to new service
func (a *ServiceAdapter) Health() map[string]string {
	return a.dbService.Health()
}

// Close closes the database service - delegates to new service
func (a *ServiceAdapter) Close() error {
	return a.dbService.Close()
}

// GetDB returns a raw database connection (legacy support)
// Deprecated: This method is deprecated, use the new context-based approach
func (a *ServiceAdapter) GetDB() *sql.DB {
	log.Printf("⚠️ GetDB() called via adapter - this method is deprecated")
	standard := a.dbService.GetStandardDB()
	if standard != nil {
		return standard.GetDB()
	}
	return nil
}

// RunMigrations runs database migrations - delegates to new service
func (a *ServiceAdapter) RunMigrations() error {
	return a.dbService.RunMigrations()
}

// CheckConnection checks if the database is reachable - delegates to new service
func (a *ServiceAdapter) CheckConnection() error {
	standard := a.dbService.GetStandardDB()
	if standard != nil {
		return standard.CheckConnection()
	}
	return nil
}

// Reconnect is no longer needed with the new connection manager
// Deprecated: Connection management is now automatic
func (a *ServiceAdapter) Reconnect() error {
	log.Printf("⚠️ Reconnect() called via adapter - connection management is now automatic")
	standard := a.dbService.GetStandardDB()
	if standard != nil {
		return standard.Reconnect()
	}
	return nil
}

// GetConnection returns a request-scoped database connection (new functionality)
func (a *ServiceAdapter) GetConnection(ctx context.Context) (*RequestConnection, error) {
	return a.dbService.GetConnection(ctx)
}

// IsHealthy returns the current health status (new functionality)
func (a *ServiceAdapter) IsHealthy() bool {
	return a.dbService.IsHealthy()
}

// GetMetrics returns detailed connection pool statistics (new functionality)
func (a *ServiceAdapter) GetMetrics() *PoolMetrics {
	return a.dbService.GetMetrics()
}
