package analytics

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// CacheEntry represents a cached analytics response
type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
	UserID    string
	Key       string
}

// AnalyticsCache provides caching for analytics data to prevent over-fetching
type AnalyticsCache struct {
	entries map[string]*CacheEntry
	mutex   sync.RWMutex
}

// Cache TTL configurations
const (
	EnhancedAnalyticsTTL = 5 * time.Minute  // Cache enhanced analytics for 5 minutes
	DashboardOverviewTTL = 3 * time.Minute  // Cache dashboard overview for 3 minutes
	ChartDataTTL         = 10 * time.Minute // Cache chart data for 10 minutes
	GrowthAnalysisTTL    = 30 * time.Minute // Cache growth analysis for 30 minutes
)

var globalCache *AnalyticsCache
var cacheOnce sync.Once

// GetAnalyticsCache returns the singleton cache instance
func GetAnalyticsCache() *AnalyticsCache {
	cacheOnce.Do(func() {
		globalCache = &AnalyticsCache{
			entries: make(map[string]*CacheEntry),
		}
		// Start cleanup goroutine
		go globalCache.cleanupExpired()
	})
	return globalCache
}

// generateCacheKey creates a unique cache key
func (ac *AnalyticsCache) generateCacheKey(userID, dataType string, params ...string) string {
	key := fmt.Sprintf("%s:%s", userID, dataType)
	for _, param := range params {
		key += ":" + param
	}
	return key
}

// Get retrieves data from cache if not expired
func (ac *AnalyticsCache) Get(userID, dataType string, params ...string) (interface{}, bool) {
	key := ac.generateCacheKey(userID, dataType, params...)

	ac.mutex.RLock()
	entry, exists := ac.entries[key]
	ac.mutex.RUnlock()

	if !exists {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		// Entry expired, remove it
		ac.mutex.Lock()
		delete(ac.entries, key)
		ac.mutex.Unlock()
		return nil, false
	}

	return entry.Data, true
}

// Set stores data in cache with TTL
func (ac *AnalyticsCache) Set(userID, dataType string, data interface{}, ttl time.Duration, params ...string) {
	key := ac.generateCacheKey(userID, dataType, params...)

	entry := &CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
		UserID:    userID,
		Key:       key,
	}

	ac.mutex.Lock()
	ac.entries[key] = entry
	ac.mutex.Unlock()
}

// InvalidateUser removes all cache entries for a specific user
func (ac *AnalyticsCache) InvalidateUser(userID string) {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	for key, entry := range ac.entries {
		if entry.UserID == userID {
			delete(ac.entries, key)
		}
	}
}

// InvalidateUserDataType removes cache entries for a specific user and data type
func (ac *AnalyticsCache) InvalidateUserDataType(userID, dataType string) {
	prefix := userID + ":" + dataType

	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	for key := range ac.entries {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(ac.entries, key)
		}
	}
}

// cleanupExpired removes expired entries periodically
func (ac *AnalyticsCache) cleanupExpired() {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		ac.mutex.Lock()
		for key, entry := range ac.entries {
			if now.After(entry.ExpiresAt) {
				delete(ac.entries, key)
			}
		}
		ac.mutex.Unlock()
	}
}

// GetCacheStats returns cache statistics for monitoring
func (ac *AnalyticsCache) GetCacheStats() map[string]interface{} {
	ac.mutex.RLock()
	defer ac.mutex.RUnlock()

	now := time.Now()
	totalEntries := len(ac.entries)
	expiredEntries := 0

	userCounts := make(map[string]int)
	typeCounts := make(map[string]int)

	for _, entry := range ac.entries {
		if now.After(entry.ExpiresAt) {
			expiredEntries++
		}

		userCounts[entry.UserID]++

		// Extract data type from key
		parts := strings.Split(entry.Key, ":")
		if len(parts) >= 2 {
			typeCounts[parts[1]]++
		}
	}

	return map[string]interface{}{
		"total_entries":   totalEntries,
		"expired_entries": expiredEntries,
		"active_entries":  totalEntries - expiredEntries,
		"users_cached":    len(userCounts),
		"cache_types":     typeCounts,
	}
}
