package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type DatabaseService struct {
	poolManager        *PoolManager
	standardManager    *ConnectionManager
	migrationRunner    *MigrationRunner
	enableFallbackMode bool
}

type DatabaseInterface interface {
	GetConnection(ctx context.Context) (*RequestConnection, error)
	GetStandardDB() Service
	Health() map[string]string
	GetMetrics() *PoolMetrics
	IsHealthy() bool
	RunMigrations() error
	Close() error
}

func NewDatabaseService() (DatabaseInterface, error) {
	log.Printf("🚀 Initializing enhanced database service")

	poolManager, err := NewPoolManager()
	if err != nil {
		log.Printf("❌ Failed to initialize pool manager, falling back to standard mode: %v", err)

		// Try to create a standard manager for fallback
		standardManager, standardErr := NewConnectionManager()
		if standardErr != nil {
			log.Printf("❌ Failed to initialize standard manager as fallback: %v", standardErr)
			return nil, fmt.Errorf("both pool manager and standard manager failed: pool=%v, standard=%v", err, standardErr)
		}

		return &DatabaseService{
			standardManager:    standardManager,
			poolManager:        nil,
			enableFallbackMode: true,
		}, nil
	}

	// Only create standard manager if pool manager fails - avoid dual connection managers
	service := &DatabaseService{
		poolManager:        poolManager,
		standardManager:    nil,
		enableFallbackMode: false,
	}

	service.migrationRunner = NewMigrationRunner(poolManager.db)

	log.Printf("✅ Enhanced database service initialized successfully")
	return service, nil
}

func (ds *DatabaseService) GetConnection(ctx context.Context) (*RequestConnection, error) {
	if ds.enableFallbackMode {
		return nil, fmt.Errorf("fallback mode enabled, use GetStandardDB() instead")
	}

	if ds.poolManager == nil {
		return nil, fmt.Errorf("pool manager not available")
	}

	return ds.poolManager.GetConnection(ctx)
}

func (ds *DatabaseService) GetStandardDB() Service {
	if ds.enableFallbackMode {
		if ds.standardManager != nil {
			return &standardServiceAdapter{manager: ds.standardManager}
		}
		return New()
	}

	// For enhanced mode, provide standard access to the pool manager's underlying DB
	if ds.poolManager != nil {
		return &poolServiceAdapter{db: ds.poolManager.db}
	}

	return New()
}

func (ds *DatabaseService) Health() map[string]string {
	if ds.enableFallbackMode {
		return ds.GetStandardDB().Health()
	}

	if ds.poolManager == nil {
		return map[string]string{
			"status": "down",
			"error":  "pool manager not available",
		}
	}

	metrics := ds.poolManager.GetMetrics()

	return map[string]string{
		"status":                getHealthStatus(ds.poolManager.IsHealthy()),
		"active_connections":    fmt.Sprintf("%d", metrics.ActiveConnections),
		"idle_connections":      fmt.Sprintf("%d", metrics.IdleConnections),
		"total_connections":     fmt.Sprintf("%d", metrics.TotalConnections),
		"failed_connections":    fmt.Sprintf("%d", metrics.FailedConnections),
		"wait_count":            fmt.Sprintf("%d", metrics.WaitCount),
		"wait_duration_ms":      fmt.Sprintf("%d", metrics.WaitDuration/1000000),
		"circuit_breaker_trips": fmt.Sprintf("%d", metrics.CircuitBreakerTrips),
		"last_health_check":     metrics.LastHealthCheck.Format(time.RFC3339),
		"message":               generateHealthMessage(metrics),
	}
}

func (ds *DatabaseService) GetMetrics() *PoolMetrics {
	if ds.poolManager == nil {
		return &PoolMetrics{
			HealthStatus: false,
		}
	}

	return ds.poolManager.GetMetrics()
}

func (ds *DatabaseService) IsHealthy() bool {
	if ds.enableFallbackMode {
		standardHealth := ds.GetStandardDB().Health()
		return standardHealth["status"] == "up"
	}

	if ds.poolManager == nil {
		return false
	}

	return ds.poolManager.IsHealthy()
}

func (ds *DatabaseService) RunMigrations() error {
	if ds.enableFallbackMode {
		return ds.GetStandardDB().RunMigrations()
	}

	if ds.migrationRunner == nil {
		return fmt.Errorf("migration runner not available")
	}

	return ds.migrationRunner.RunMigrations("migrations")
}

func (ds *DatabaseService) Close() error {
	log.Printf("🔌 Closing database service")

	var errs []error

	if ds.poolManager != nil {
		if err := ds.poolManager.Close(); err != nil {
			errs = append(errs, fmt.Errorf("pool manager close error: %w", err))
		}
	}

	if ds.standardManager != nil {
		if err := ds.standardManager.Close(); err != nil {
			errs = append(errs, fmt.Errorf("standard manager close error: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during close: %v", errs)
	}

	return nil
}

func getHealthStatus(healthy bool) string {
	if healthy {
		return "up"
	}
	return "down"
}

func generateHealthMessage(metrics *PoolMetrics) string {
	if !metrics.HealthStatus {
		return "Database is unhealthy"
	}

	if metrics.FailedConnections > 10 {
		return "High number of failed connections detected"
	}

	if metrics.CircuitBreakerTrips > 0 {
		return "Circuit breaker has been triggered"
	}

	if metrics.ActiveConnections > 20 {
		return "High connection usage detected"
	}

	return "Database is healthy"
}

type standardServiceAdapter struct {
	manager *ConnectionManager
}

type poolServiceAdapter struct {
	db *sql.DB
}

func (ssa *standardServiceAdapter) Health() map[string]string {
	stats := ssa.manager.GetStats()
	result := make(map[string]string)

	for k, v := range stats {
		result[k] = fmt.Sprintf("%v", v)
	}

	if ssa.manager.IsHealthy() {
		result["status"] = "up"
		result["message"] = "It's healthy"
	} else {
		result["status"] = "down"
		result["error"] = "Standard manager unhealthy"
	}

	return result
}

func (ssa *standardServiceAdapter) Close() error {
	return ssa.manager.Close()
}

func (ssa *standardServiceAdapter) GetDB() *sql.DB {
	if ssa.manager == nil {
		return nil
	}
	return ssa.manager.GetDB()
}

func (ssa *standardServiceAdapter) RunMigrations() error {
	return fmt.Errorf("migrations not supported in standard adapter")
}

func (ssa *standardServiceAdapter) CheckConnection() error {
	if ssa.manager.IsHealthy() {
		return nil
	}
	return fmt.Errorf("connection unhealthy")
}

func (ssa *standardServiceAdapter) Reconnect() error {
	return fmt.Errorf("reconnect not supported in standard adapter")
}

// Pool service adapter methods
func (psa *poolServiceAdapter) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := psa.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("pool db down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "Pool database healthy"

	dbStats := psa.db.Stats()
	stats["open_connections"] = fmt.Sprintf("%d", dbStats.OpenConnections)
	stats["in_use"] = fmt.Sprintf("%d", dbStats.InUse)
	stats["idle"] = fmt.Sprintf("%d", dbStats.Idle)

	return stats
}

func (psa *poolServiceAdapter) Close() error {
	// Don't close the pool manager's DB directly
	return nil
}

func (psa *poolServiceAdapter) GetDB() *sql.DB {
	return psa.db
}

func (psa *poolServiceAdapter) RunMigrations() error {
	migrationRunner := NewMigrationRunner(psa.db)
	return migrationRunner.RunMigrations("migrations")
}

func (psa *poolServiceAdapter) CheckConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return psa.db.PingContext(ctx)
}

func (psa *poolServiceAdapter) Reconnect() error {
	return fmt.Errorf("reconnect not supported in pool adapter - managed by pool manager")
}
