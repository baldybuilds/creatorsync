package analytics

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/baldybuilds/creatorsync/internal/clerk"
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	service                 Service
	backgroundCollectionMgr *BackgroundCollectionManager
	activeRequests          sync.Map // Track active requests to prevent duplicates
}

// RequestInfo tracks active request information
type RequestInfo struct {
	UserID    string
	StartTime time.Time
	RequestID string
}

func NewHandlers(service Service, backgroundCollectionMgr *BackgroundCollectionManager) *Handlers {
	return &Handlers{
		service:                 service,
		backgroundCollectionMgr: backgroundCollectionMgr,
	}
}

// Helper function to get user ID from context
func (h *Handlers) getUserID(c *fiber.Ctx) (string, error) {
	user, err := clerk.GetUserFromContext(c)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

// RegisterRoutes registers all analytics routes
func (h *Handlers) RegisterRoutes(app *fiber.App) {
	// Create analytics API group that inherits from main app (with CORS)
	api := app.Group("/api/analytics")

	// Add a test endpoint to verify CORS is working
	api.Get("/cors-test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":   "CORS test successful",
			"origin":    c.Get("Origin"),
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Public routes (no authentication required)
	api.Get("/health", h.HealthCheck)

	// Protected routes - require authentication
	protected := api.Group("")
	protected.Use(clerk.AuthMiddleware())

	// Add a simple protected test endpoint
	protected.Get("/auth-test", func(c *fiber.Ctx) error {
		user, err := clerk.GetUserFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
			})
		}
		return c.JSON(fiber.Map{
			"message":   "Auth test successful",
			"user_id":   user.ID,
			"origin":    c.Get("Origin"),
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Dashboard overview - returns summary metrics for main dashboard
	protected.Get("/overview", h.GetDashboardOverview)

	// Analytics page - returns detailed analytics
	protected.Get("/detailed", h.GetDetailedAnalytics)

	// Enhanced analytics - returns video-based analytics for new dashboard design
	protected.Get("/enhanced", h.GetEnhancedAnalytics)

	// Chart data for specific time periods
	protected.Get("/charts", h.GetAnalyticsChartData)

	// Growth analysis
	protected.Get("/growth", h.GetGrowthAnalysis)

	// Content performance
	protected.Get("/content", h.GetContentPerformance)

	// Job status
	protected.Get("/jobs", h.GetAnalyticsJobs)

	// Manual data collection triggers
	protected.Post("/collect", h.TriggerDataCollection)
	protected.Post("/refresh", h.RefreshChannelData)

	// Connection status endpoint
	protected.Get("/connection-status", h.GetConnectionStatus)

	// Debug endpoints
	protected.Get("/debug/data-status", h.GetDataStatus)
	protected.Get("/debug/cache-stats", h.GetCacheStats)

}

// GetDashboardOverview returns summary metrics for the dashboard
func (h *Handlers) GetDashboardOverview(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Check cache first
	cache := GetAnalyticsCache()
	if cachedData, found := cache.Get(userID, "overview"); found {
		if response, ok := cachedData.(map[string]interface{}); ok {
			c.Set("X-Cache-Status", "HIT")
			return c.JSON(response)
		} else if overview, ok := cachedData.(*DashboardOverview); ok {
			// Legacy cache format - add connection status
			twitchConnected, err := h.service.CheckTwitchConnection(c.Context(), userID)
			if err != nil {
				log.Printf("❌ Connection status check failed for user %s: %v", userID, err)
				twitchConnected = false
			}
			response := map[string]interface{}{
				"overview": overview,
				"connection_status": map[string]interface{}{
					"twitch_connected": twitchConnected,
					"settings_url":     "/settings",
				},
			}
			c.Set("X-Cache-Status", "HIT")
			return c.JSON(response)
		}
	}

	// Check connection status
	twitchConnected, err := h.service.CheckTwitchConnection(c.Context(), userID)
	if err != nil {
		log.Printf("❌ Connection status check failed for user %s: %v", userID, err)
		twitchConnected = false
	}

	overview, err := h.service.GetDashboardOverview(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get dashboard overview",
		})
	}

	// Create response with connection status
	response := map[string]interface{}{
		"overview": overview,
		"connection_status": map[string]interface{}{
			"twitch_connected": twitchConnected,
			"settings_url":     "/settings",
		},
	}

	// Cache the result
	cache.Set(userID, "overview", response, DashboardOverviewTTL)

	c.Set("X-Cache-Status", "MISS")
	return c.JSON(response)
}

// GetAnalyticsChartData returns chart data for analytics visualization
func (h *Handlers) GetAnalyticsChartData(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Get days parameter (default to 30)
	daysStr := c.Query("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30
	}

	// Check cache first
	cache := GetAnalyticsCache()
	if cachedData, found := cache.Get(userID, "chartdata", daysStr); found {
		if chartData, ok := cachedData.(*AnalyticsChartData); ok {
			c.Set("X-Cache-Status", "HIT")
			return c.JSON(chartData)
		}
	}

	chartData, err := h.service.GetAnalyticsChartData(c.Context(), userID, days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get chart data",
		})
	}

	// Cache the result
	cache.Set(userID, "chartdata", chartData, ChartDataTTL, daysStr)

	c.Set("X-Cache-Status", "MISS")
	return c.JSON(chartData)
}

// GetDetailedAnalytics returns comprehensive analytics for the analytics page
func (h *Handlers) GetDetailedAnalytics(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	analytics, err := h.service.GetDetailedAnalytics(c.Context(), userID)
	if err != nil {
		log.Printf("Error getting detailed analytics for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get detailed analytics",
		})
	}

	return c.JSON(analytics)
}

// GetEnhancedAnalytics returns video-based analytics for the new dashboard design
func (h *Handlers) GetEnhancedAnalytics(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Get days parameter (default to 30)
	daysStr := c.Query("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30
	}

	// Check cache first
	cache := GetAnalyticsCache()
	if cachedData, found := cache.Get(userID, "enhanced", daysStr); found {
		// Handle both old and new cache formats for compatibility
		if response, ok := cachedData.(map[string]interface{}); ok {
			c.Set("Cache-Control", "private, max-age=300") // 5 minutes
			c.Set("X-Cache-Status", "HIT")
			return c.JSON(response)
		} else if analytics, ok := cachedData.(*EnhancedAnalytics); ok {
			// Legacy cache format - add connection status
			twitchConnected, err := h.service.CheckTwitchConnection(c.Context(), userID)
			if err != nil {
				log.Printf("❌ Connection status check failed for user %s: %v", userID, err)
				twitchConnected = false
			}
			response := map[string]interface{}{
				"analytics": analytics,
				"connection_status": map[string]interface{}{
					"twitch_connected": twitchConnected,
					"settings_url":     "/settings",
				},
			}
			c.Set("Cache-Control", "private, max-age=300") // 5 minutes
			c.Set("X-Cache-Status", "HIT")
			return c.JSON(response)
		}
	}

	// Check for concurrent requests to prevent duplicate API calls
	requestKey := fmt.Sprintf("enhanced_%s_%d", userID, days)
	if existing, loaded := h.activeRequests.LoadOrStore(requestKey, &RequestInfo{
		UserID:    userID,
		StartTime: time.Now(),
		RequestID: fmt.Sprintf("req_%d", time.Now().UnixNano()),
	}); loaded {
		existingReq := existing.(*RequestInfo)
		elapsed := time.Since(existingReq.StartTime)

		// If there's a recent request, return 429 with shorter retry time
		if elapsed < 10*time.Second {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Request in progress",
				"retry_after": "10",
			})
		}
		// Replace stale request
		h.activeRequests.Store(requestKey, &RequestInfo{
			UserID:    userID,
			StartTime: time.Now(),
			RequestID: fmt.Sprintf("req_%d", time.Now().UnixNano()),
		})
	}

	defer func() {
		h.activeRequests.Delete(requestKey)
	}()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	// Check connection status
	twitchConnected, err := h.service.CheckTwitchConnection(ctx, userID)
	if err != nil {
		log.Printf("❌ Connection status check failed for user %s: %v", userID, err)
		twitchConnected = false
	}

	// Clear any cached analytics data before fetching fresh data
	// This ensures we get the latest data after data collection
	cache.InvalidateUserDataType(userID, "enhanced")

	// Fetch data from service
	analytics, err := h.service.GetEnhancedAnalytics(ctx, userID, days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get enhanced analytics",
		})
	}

	if analytics == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No analytics data available",
		})
	}

	// Add connection status to response
	response := map[string]interface{}{
		"analytics": analytics,
		"connection_status": map[string]interface{}{
			"twitch_connected": twitchConnected,
			"settings_url":     "/settings",
		},
	}

	// Cache the result (cache the full response including connection status)
	cache.Set(userID, "enhanced", response, EnhancedAnalyticsTTL, daysStr)

	// Set cache headers
	c.Set("Cache-Control", "private, max-age=300") // 5 minutes
	c.Set("X-Cache-Status", "MISS")

	return c.JSON(response)
}

// GetGrowthAnalysis provides growth trend analysis
func (h *Handlers) GetGrowthAnalysis(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	period := c.Query("period", "month")
	if period != "week" && period != "month" && period != "quarter" && period != "year" {
		period = "month"
	}

	analysis, err := h.service.GetGrowthAnalysis(c.Context(), userID, period)
	if err != nil {
		log.Printf("Error getting growth analysis for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get growth analysis",
		})
	}

	return c.JSON(analysis)
}

// GetContentPerformance analyzes video and stream performance
func (h *Handlers) GetContentPerformance(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	performance, err := h.service.GetContentPerformance(c.Context(), userID)
	if err != nil {
		log.Printf("Error getting content performance for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get content performance",
		})
	}

	return c.JSON(performance)
}

// TriggerDataCollection manually triggers data collection for a user
func (h *Handlers) TriggerDataCollection(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Check if there's already a collection in progress for this user
	requestKey := fmt.Sprintf("collection_%s", userID)
	if _, loaded := h.activeRequests.LoadOrStore(requestKey, &RequestInfo{
		UserID:    userID,
		StartTime: time.Now(),
		RequestID: fmt.Sprintf("collect_%d", time.Now().UnixNano()),
	}); loaded {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Data collection already in progress for this user",
		})
	}

	// Ensure cleanup
	defer h.activeRequests.Delete(requestKey)

	// Invalidate cache since we're refreshing data
	cache := GetAnalyticsCache()
	cache.InvalidateUser(userID)

	log.Printf("🚀 Starting manual data collection for user %s", userID)

	// Create context with timeout for the collection process
	ctx, cancel := context.WithTimeout(c.Context(), 4*time.Minute)
	defer cancel()

	// Use service method directly for better control and error handling
	err = h.service.TriggerDataCollection(ctx, userID)
	if err != nil {
		log.Printf("❌ Data collection failed for user %s: %v", userID, err)

		// Check if it's a partial failure (some data was collected)
		if err.Error() != "failed to save any videos" {
			return c.Status(fiber.StatusPartialContent).JSON(fiber.Map{
				"message":   "Data collection completed with some errors",
				"warning":   err.Error(),
				"user_id":   userID,
				"timestamp": time.Now().Unix(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Data collection failed: %v", err),
		})
	}

	log.Printf("✅ Data collection completed successfully for user %s", userID)
	return c.JSON(fiber.Map{
		"message":   "Data collection completed successfully",
		"user_id":   userID,
		"timestamp": time.Now().Unix(),
	})
}

// RefreshChannelData specifically refreshes channel metrics
func (h *Handlers) RefreshChannelData(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Invalidate cache since we're refreshing data
	cache := GetAnalyticsCache()
	cache.InvalidateUser(userID)

	err = h.service.RefreshChannelData(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to refresh channel data",
		})
	}

	return c.JSON(fiber.Map{
		"message":   "Channel data refreshed successfully",
		"user_id":   userID,
		"timestamp": time.Now().Unix(),
	})
}

// GetAnalyticsJobs returns the status of analytics jobs for a user
func (h *Handlers) GetAnalyticsJobs(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Get limit parameter (default to 10)
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	jobs, err := h.service.GetAnalyticsJobs(c.Context(), userID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get analytics jobs",
		})
	}

	return c.JSON(fiber.Map{
		"jobs":      jobs,
		"user_id":   userID,
		"timestamp": time.Now().Unix(),
	})
}

// GetDataStatus returns debug information about data and system status
func (h *Handlers) GetDataStatus(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	status := fiber.Map{
		"user_id":         userID,
		"timestamp":       time.Now().UTC().Format(time.RFC3339),
		"environment":     os.Getenv("APP_ENV"),
		"database_health": "",
		"user_exists":     false,
		"analytics_data": fiber.Map{
			"has_channel_data": false,
			"has_video_data":   false,
			"last_update":      nil,
		},
		"errors": []string{},
	}

	// Check service health
	if h.service.IsHealthy() {
		status["database_health"] = "healthy"
	} else {
		status["database_health"] = "unhealthy"
		status["errors"] = append(status["errors"].([]string), "Database connection unhealthy")
	}

	// Check if user has analytics data (simplified check)
	hasData, _, err := h.service.CheckUserAnalyticsData(c.Context(), userID)
	if err != nil {
		status["errors"] = append(status["errors"].([]string), fmt.Sprintf("User data check failed: %v", err))
	} else {
		status["user_exists"] = hasData
	}

	// Check analytics data
	hasData, lastUpdate, err := h.service.CheckUserAnalyticsData(c.Context(), userID)
	if err != nil {
		status["errors"] = append(status["errors"].([]string), fmt.Sprintf("Analytics data check failed: %v", err))
	} else {
		analyticsData := status["analytics_data"].(fiber.Map)
		analyticsData["has_data"] = hasData
		if lastUpdate != nil {
			analyticsData["last_update"] = lastUpdate.Format(time.RFC3339)
		}
	}

	return c.JSON(status)
}

// GetConnectionStatus returns the user's platform connection status
func (h *Handlers) GetConnectionStatus(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Check cache first for connection status
	cache := GetAnalyticsCache()
	cacheKey := "connection_status"
	if cachedData, found := cache.Get(userID, cacheKey); found {
		if status, ok := cachedData.(map[string]interface{}); ok {
			c.Set("X-Cache-Status", "HIT")
			return c.JSON(status)
		}
	}

	// Check Twitch connection using the universal collector
	connected, err := h.service.CheckTwitchConnection(c.Context(), userID)
	if err != nil {
		log.Printf("❌ Connection status check failed for user %s: %v", userID, err)
		connected = false // Default to false on error
	}

	status := map[string]interface{}{
		"user_id": userID,
		"platforms": map[string]interface{}{
			"twitch": map[string]interface{}{
				"connected":    connected,
				"display_name": "Twitch",
				"color":        "#9146FF",
			},
			"youtube": map[string]interface{}{
				"connected":    false,
				"display_name": "YouTube",
				"color":        "#FF0000",
				"coming_soon":  true,
			},
			"tiktok": map[string]interface{}{
				"connected":    false,
				"display_name": "TikTok",
				"color":        "#000000",
				"coming_soon":  true,
			},
		},
		"has_any_connection": connected,
		"settings_url":       "/settings",
	}

	// Cache the result for 2 minutes (connection status can change)
	cache.Set(userID, cacheKey, status, 2*time.Minute)

	c.Set("X-Cache-Status", "MISS")
	return c.JSON(status)
}

// GetCacheStats returns cache statistics for monitoring
func (h *Handlers) GetCacheStats(c *fiber.Ctx) error {
	cache := GetAnalyticsCache()
	stats := cache.GetCacheStats()

	return c.JSON(fiber.Map{
		"cache_stats": stats,
		"timestamp":   time.Now().Unix(),
	})
}

// HealthCheck returns the health status of the analytics service
func (h *Handlers) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "healthy",
		"service":   "analytics",
		"timestamp": time.Now().Unix(),
	})
}
