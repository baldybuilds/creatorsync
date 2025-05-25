package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// ConnectionManager manages database connections and health
type ConnectionManager struct {
	db        *sql.DB
	connStr   string
	config    *PoolConfig
	healthMux sync.RWMutex
	isHealthy bool
	lastCheck time.Time
	checkFreq time.Duration
	done      chan struct{}
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() (*ConnectionManager, error) {
	connStr := buildConnectionString()
	config := getPoolConfig()

	log.Printf("🔗 Initializing connection manager for environment: %s", config.Environment)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	log.Printf("✅ Applied %s pool settings (MaxOpen: %d, MaxIdle: %d)",
		config.Environment, config.MaxOpenConns, config.MaxIdleConns)

	cm := &ConnectionManager{
		db:        db,
		connStr:   connStr,
		config:    config,
		checkFreq: 30 * time.Second, // Health check every 30 seconds
		done:      make(chan struct{}),
	}

	// Initial health check
	if err := cm.performHealthCheck(context.Background()); err != nil {
		return nil, fmt.Errorf("initial health check failed: %w", err)
	}

	// Start background health monitoring
	go cm.startHealthMonitoring()

	return cm, nil
}

// GetContext returns a new request-scoped database context
func (cm *ConnectionManager) GetContext(timeout time.Duration) *DBContext {
	return NewDBContext(cm.db, timeout)
}

// IsHealthy returns the current health status (non-blocking)
func (cm *ConnectionManager) IsHealthy() bool {
	cm.healthMux.RLock()
	defer cm.healthMux.RUnlock()
	return cm.isHealthy
}

// GetStats returns connection pool statistics
func (cm *ConnectionManager) GetStats() map[string]interface{} {
	stats := cm.db.Stats()

	return map[string]interface{}{
		"open_connections":    stats.OpenConnections,
		"in_use":              stats.InUse,
		"idle":                stats.Idle,
		"wait_count":          stats.WaitCount,
		"wait_duration_ms":    stats.WaitDuration.Milliseconds(),
		"max_idle_closed":     stats.MaxIdleClosed,
		"max_lifetime_closed": stats.MaxLifetimeClosed,
		"is_healthy":          cm.IsHealthy(),
		"last_health_check":   cm.lastCheck.Format(time.RFC3339),
	}
}

// GetDB returns the underlying database connection
func (cm *ConnectionManager) GetDB() *sql.DB {
	return cm.db
}

// Close closes the connection manager
func (cm *ConnectionManager) Close() error {
	// Add check to prevent multiple closes
	select {
	case <-cm.done:
		log.Printf("⚠️ Connection manager already closed")
		return nil
	default:
		close(cm.done)
		log.Printf("🔌 Closing connection manager")
		return cm.db.Close()
	}
}

// performHealthCheck performs a single health check
func (cm *ConnectionManager) performHealthCheck(ctx context.Context) error {
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Simple ping test
	if err := cm.db.PingContext(checkCtx); err != nil {
		cm.updateHealthStatus(false)
		return fmt.Errorf("ping failed: %w", err)
	}

	// Test basic query
	var result int
	if err := cm.db.QueryRowContext(checkCtx, "SELECT 1").Scan(&result); err != nil {
		cm.updateHealthStatus(false)
		return fmt.Errorf("test query failed: %w", err)
	}

	cm.updateHealthStatus(true)
	return nil
}

// updateHealthStatus updates the health status thread-safely
func (cm *ConnectionManager) updateHealthStatus(healthy bool) {
	cm.healthMux.Lock()
	defer cm.healthMux.Unlock()

	if cm.isHealthy != healthy {
		if healthy {
			log.Printf("✅ Database health restored")
		} else {
			log.Printf("❌ Database health degraded")
		}
	}

	cm.isHealthy = healthy
	cm.lastCheck = time.Now()
}

// startHealthMonitoring starts background health monitoring
func (cm *ConnectionManager) startHealthMonitoring() {
	ticker := time.NewTicker(cm.checkFreq)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

			if err := cm.performHealthCheck(ctx); err != nil {
				log.Printf("⚠️ Health check failed: %v", err)
			}

			cancel()
		case <-cm.done:
			log.Printf("🔌 Health monitoring stopped")
			return
		}
	}
}

// getPoolConfig returns environment-specific pool configuration
func getPoolConfig() *PoolConfig {
	return buildPoolConfig()
}
