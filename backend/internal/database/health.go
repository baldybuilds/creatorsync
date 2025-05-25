package database

import (
	"context"
	"log"
	"sync/atomic"
	"time"
)

type HealthChecker struct {
	poolManager *PoolManager
	stopChan    chan struct{}
	running     int32
}

type HealthCheckResult struct {
	IsHealthy    bool
	Timestamp    time.Time
	ResponseTime time.Duration
	Error        error
}

func NewHealthChecker(pm *PoolManager) *HealthChecker {
	return &HealthChecker{
		poolManager: pm,
		stopChan:    make(chan struct{}),
	}
}

func (hc *HealthChecker) Start() {
	if !atomic.CompareAndSwapInt32(&hc.running, 0, 1) {
		log.Printf("⚠️ Health checker already running")
		return
	}

	log.Printf("🩺 Starting database health checker")
	ticker := time.NewTicker(hc.poolManager.config.HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			hc.performHealthCheck()
		case <-hc.stopChan:
			log.Printf("🩺 Health checker stopped")
			return
		}
	}
}

func (hc *HealthChecker) Stop() {
	if atomic.CompareAndSwapInt32(&hc.running, 1, 0) {
		close(hc.stopChan)
	}
}

func (hc *HealthChecker) performHealthCheck() {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := &HealthCheckResult{
		Timestamp: start,
	}

	if err := hc.basicConnectivityCheck(ctx); err != nil {
		result.IsHealthy = false
		result.Error = err
		hc.handleUnhealthyState(result)
		return
	}

	if err := hc.performanceCheck(ctx); err != nil {
		result.IsHealthy = false
		result.Error = err
		hc.handleUnhealthyState(result)
		return
	}

	result.IsHealthy = true
	result.ResponseTime = time.Since(start)
	hc.handleHealthyState(result)
}

func (hc *HealthChecker) basicConnectivityCheck(ctx context.Context) error {
	return hc.poolManager.db.PingContext(ctx)
}

func (hc *HealthChecker) performanceCheck(ctx context.Context) error {
	var result int
	return hc.poolManager.db.QueryRowContext(ctx, "SELECT 1").Scan(&result)
}

func (hc *HealthChecker) handleHealthyState(result *HealthCheckResult) {
	hc.poolManager.updateHealthStatus(true)

	if result.ResponseTime > 5*time.Second {
		log.Printf("⚠️ Slow database response: %v", result.ResponseTime)
	}
}

func (hc *HealthChecker) handleUnhealthyState(result *HealthCheckResult) {
	hc.poolManager.updateHealthStatus(false)
	log.Printf("❌ Health check failed: %v", result.Error)

	hc.poolManager.circuitBreaker.RecordFailure()

	if hc.shouldAttemptRecovery() {
		go hc.attemptRecovery()
	}
}

func (hc *HealthChecker) shouldAttemptRecovery() bool {
	state := atomic.LoadInt32(&hc.poolManager.circuitBreaker.state)
	return state == CircuitOpen
}

func (hc *HealthChecker) attemptRecovery() {
	log.Printf("🔄 Attempting database recovery...")

	time.Sleep(30 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := hc.basicConnectivityCheck(ctx); err == nil {
		log.Printf("✅ Database recovery successful")
		hc.poolManager.circuitBreaker.RecordSuccess()
	} else {
		log.Printf("❌ Database recovery failed: %v", err)
	}
}
