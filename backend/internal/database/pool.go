package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PoolManager struct {
	db             *sql.DB
	config         *PoolConfig
	connStr        string
	healthChecker  *HealthChecker
	metrics        *PoolMetrics
	circuitBreaker *CircuitBreaker
	mu             sync.RWMutex
	closed         int32
}

type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	Environment     string

	HealthCheckInterval     time.Duration
	CircuitBreakerThreshold int
	CircuitBreakerTimeout   time.Duration

	AcquireTimeout time.Duration
	QueryTimeout   time.Duration
}

type PoolMetrics struct {
	TotalConnections    int64
	ActiveConnections   int64
	IdleConnections     int64
	WaitCount           int64
	WaitDuration        int64
	FailedConnections   int64
	CircuitBreakerTrips int64

	LastHealthCheck time.Time
	HealthStatus    bool

	mu sync.RWMutex
}

type CircuitBreaker struct {
	failures    int64
	lastFailure time.Time
	state       int32 // 0: closed, 1: open, 2: half-open
	threshold   int
	timeout     time.Duration
	mu          sync.RWMutex
}

const (
	CircuitClosed = iota
	CircuitOpen
	CircuitHalfOpen
)

func NewPoolManager() (*PoolManager, error) {
	config := buildPoolConfig()
	connStr := buildConnectionString()

	log.Printf("🔗 Initializing enhanced pool manager for %s environment", config.Environment)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	applyPoolConfig(db, config)

	pm := &PoolManager{
		db:      db,
		config:  config,
		connStr: connStr,
		metrics: &PoolMetrics{
			HealthStatus: false,
		},
		circuitBreaker: &CircuitBreaker{
			threshold: config.CircuitBreakerThreshold,
			timeout:   config.CircuitBreakerTimeout,
			state:     CircuitClosed,
		},
	}

	pm.healthChecker = NewHealthChecker(pm)

	if err := pm.performInitialHealthCheck(); err != nil {
		db.Close()
		return nil, fmt.Errorf("initial health check failed: %w", err)
	}

	go pm.healthChecker.Start()
	go pm.startMetricsCollection()

	log.Printf("✅ Enhanced pool manager initialized successfully")
	return pm, nil
}

func (pm *PoolManager) GetConnection(ctx context.Context) (*RequestConnection, error) {
	if atomic.LoadInt32(&pm.closed) == 1 {
		return nil, fmt.Errorf("pool manager is closed")
	}

	if !pm.circuitBreaker.CanExecute() {
		atomic.AddInt64(&pm.metrics.FailedConnections, 1)
		return nil, fmt.Errorf("circuit breaker is open")
	}

	timeout := pm.config.AcquireTimeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	connCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	conn, err := pm.acquireConnection(connCtx)
	if err != nil {
		pm.circuitBreaker.RecordFailure()
		atomic.AddInt64(&pm.metrics.FailedConnections, 1)
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}

	pm.circuitBreaker.RecordSuccess()
	atomic.AddInt64(&pm.metrics.TotalConnections, 1)

	return NewRequestConnection(conn, pm.config.QueryTimeout), nil
}

func (pm *PoolManager) acquireConnection(ctx context.Context) (*sql.DB, error) {
	start := time.Now()
	defer func() {
		atomic.AddInt64(&pm.metrics.WaitDuration, time.Since(start).Nanoseconds())
		atomic.AddInt64(&pm.metrics.WaitCount, 1)
	}()

	if err := pm.db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("connection health check failed: %w", err)
	}

	return pm.db, nil
}

func (pm *PoolManager) GetMetrics() *PoolMetrics {
	pm.metrics.mu.RLock()
	defer pm.metrics.mu.RUnlock()

	stats := pm.db.Stats()

	atomic.StoreInt64(&pm.metrics.ActiveConnections, int64(stats.InUse))
	atomic.StoreInt64(&pm.metrics.IdleConnections, int64(stats.Idle))

	metricsCopy := *pm.metrics
	return &metricsCopy
}

func (pm *PoolManager) IsHealthy() bool {
	pm.metrics.mu.RLock()
	defer pm.metrics.mu.RUnlock()
	return pm.metrics.HealthStatus
}

func (pm *PoolManager) performInitialHealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := pm.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	var result int
	if err := pm.db.QueryRowContext(ctx, "SELECT 1").Scan(&result); err != nil {
		return fmt.Errorf("test query failed: %w", err)
	}

	pm.updateHealthStatus(true)
	return nil
}

func (pm *PoolManager) updateHealthStatus(healthy bool) {
	pm.metrics.mu.Lock()
	defer pm.metrics.mu.Unlock()

	if pm.metrics.HealthStatus != healthy {
		if healthy {
			log.Printf("✅ Database health restored")
		} else {
			log.Printf("❌ Database health degraded")
		}
	}

	pm.metrics.HealthStatus = healthy
	pm.metrics.LastHealthCheck = time.Now()
}

func (pm *PoolManager) startMetricsCollection() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if atomic.LoadInt32(&pm.closed) == 1 {
				return
			}
			pm.collectMetrics()
		}
	}
}

func (pm *PoolManager) collectMetrics() {
	stats := pm.db.Stats()

	atomic.StoreInt64(&pm.metrics.ActiveConnections, int64(stats.InUse))
	atomic.StoreInt64(&pm.metrics.IdleConnections, int64(stats.Idle))

	if stats.OpenConnections > pm.config.MaxOpenConns*8/10 {
		log.Printf("⚠️ High connection usage: %d/%d", stats.OpenConnections, pm.config.MaxOpenConns)
	}

	if stats.WaitCount > 100 {
		log.Printf("⚠️ High wait count detected: %d", stats.WaitCount)
	}
}

func (pm *PoolManager) Close() error {
	if !atomic.CompareAndSwapInt32(&pm.closed, 0, 1) {
		log.Printf("⚠️ Pool manager already closed")
		return nil
	}

	pm.healthChecker.Stop()

	log.Printf("🔌 Closing enhanced pool manager")
	return pm.db.Close()
}

func (cb *CircuitBreaker) CanExecute() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	state := atomic.LoadInt32(&cb.state)

	switch state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		if time.Since(cb.lastFailure) > cb.timeout {
			atomic.StoreInt32(&cb.state, CircuitHalfOpen)
			return true
		}
		return false
	case CircuitHalfOpen:
		return true
	default:
		return false
	}
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	atomic.StoreInt64(&cb.failures, 0)
	atomic.StoreInt32(&cb.state, CircuitClosed)
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	failures := atomic.AddInt64(&cb.failures, 1)
	cb.lastFailure = time.Now()

	if int(failures) >= cb.threshold {
		atomic.StoreInt32(&cb.state, CircuitOpen)
		log.Printf("🔴 Circuit breaker opened after %d failures", failures)
	}
}

func buildPoolConfig() *PoolConfig {
	env := getEnvironment()

	baseConfig := &PoolConfig{
		Environment:             env,
		HealthCheckInterval:     30 * time.Second,
		CircuitBreakerThreshold: 5,
		CircuitBreakerTimeout:   60 * time.Second,
		AcquireTimeout:          30 * time.Second,
		QueryTimeout:            30 * time.Second,
	}

	switch env {
	case "production":
		baseConfig.MaxOpenConns = 15
		baseConfig.MaxIdleConns = 5
		baseConfig.ConnMaxLifetime = 15 * time.Minute
		baseConfig.ConnMaxIdleTime = 5 * time.Minute

	case "staging":
		baseConfig.MaxOpenConns = 20
		baseConfig.MaxIdleConns = 8
		baseConfig.ConnMaxLifetime = 10 * time.Minute
		baseConfig.ConnMaxIdleTime = 2 * time.Minute

	default:
		baseConfig.MaxOpenConns = 25
		baseConfig.MaxIdleConns = 10
		baseConfig.ConnMaxLifetime = 10 * time.Minute
		baseConfig.ConnMaxIdleTime = 1 * time.Minute
	}

	return baseConfig
}

func applyPoolConfig(db *sql.DB, config *PoolConfig) {
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	log.Printf("✅ Applied %s pool config (MaxOpen: %d, MaxIdle: %d)",
		config.Environment, config.MaxOpenConns, config.MaxIdleConns)
}

func getEnvironment() string {
	env := os.Getenv("APP_ENV")
	environment := os.Getenv("ENVIRONMENT")
	databaseURL := os.Getenv("DATABASE_URL")

	if env == "production" || environment == "production" {
		return "production"
	}

	if env == "staging" || environment == "staging" ||
		env == "dev" || environment == "dev" ||
		(databaseURL != "" && env != "production" && environment != "production") {
		return "staging"
	}

	return "development"
}

func buildConnectionString() string {
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		log.Println("🔗 Using DATABASE_URL for connection")
		return databaseURL
	}

	database := os.Getenv("POSTGRES_DB_DATABASE")
	password := os.Getenv("POSTGRES_DB_PASSWORD")
	username := os.Getenv("POSTGRES_DB_USERNAME")
	port := os.Getenv("POSTGRES_DB_PORT")
	host := os.Getenv("POSTGRES_DB_HOST")
	schema := os.Getenv("POSTGRES_DB_SCHEMA")

	if schema == "" {
		schema = "public"
	}

	sslMode := os.Getenv("POSTGRES_SSL_MODE")
	if sslMode == "" {
		sslMode = "require"
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s",
		username, password, host, port, database, sslMode, schema)

	log.Println("🔗 Using individual environment variables for connection")
	return connStr
}
