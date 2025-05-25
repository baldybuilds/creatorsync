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

	// Debug endpoint to check data status
	protected.Get("/debug/data-status", h.GetDataStatus)

}

// GetDashboardOverview returns summary metrics for the dashboard
func (h *Handlers) GetDashboardOverview(c *fiber.Ctx) error {
	user, err := clerk.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}
	userID := user.ID

	// Ensure user exists in database before proceeding
	analyticsRepo := NewRepository(h.service.(*service).db)
	existingUser, err := analyticsRepo.GetUserByClerkID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify user account",
		})
	}

	if existingUser == nil {
		// Create minimal user record for cross-environment compatibility
		user := &User{
			ID:          userID,
			ClerkUserID: userID,
			Username:    fmt.Sprintf("user_%s", userID[len(userID)-10:]),
			DisplayName: "User",
			Email:       "",
		}

		if err := analyticsRepo.CreateOrUpdateUser(c.Context(), user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to initialize user account",
			})
		}
	}

	overview, err := h.service.GetDashboardOverview(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get dashboard overview",
			"details": err.Error(),
		})
	}

	return c.JSON(overview)
}

// GetAnalyticsChartData returns chart data for analytics visualization
func (h *Handlers) GetAnalyticsChartData(c *fiber.Ctx) error {
	user, err := clerk.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}
	userID := user.ID

	// Get days parameter (default to 30)
	daysStr := c.Query("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30
	}

	chartData, err := h.service.GetAnalyticsChartData(c.Context(), userID, days)
	if err != nil {
		log.Printf("Error getting chart data for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get chart data",
		})
	}

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
		log.Printf("❌ GetEnhancedAnalytics: Failed to get user ID: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Get request ID for tracking
	requestID := c.Query("requestId", fmt.Sprintf("server-%d", time.Now().UnixNano()))
	log.Printf("🔍 GetEnhancedAnalytics: Starting for user %s (requestId: %s)", userID, requestID)

	// Check for existing active request for this user
	if existing, loaded := h.activeRequests.LoadOrStore(userID, &RequestInfo{
		UserID:    userID,
		StartTime: time.Now(),
		RequestID: requestID,
	}); loaded {
		existingReq := existing.(*RequestInfo)
		elapsed := time.Since(existingReq.StartTime)

		// If there's a request that's been running for less than 30 seconds, reject this one
		if elapsed < 30*time.Second {
			log.Printf("🚫 GetEnhancedAnalytics: Rejecting duplicate request %s for user %s (existing: %s, elapsed: %v)",
				requestID, userID, existingReq.RequestID, elapsed)
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Request already in progress",
				"details":     fmt.Sprintf("Existing request %s started %v ago", existingReq.RequestID, elapsed),
				"requestId":   requestID,
				"retry_after": "30",
			})
		} else {
			// Old request, replace it
			log.Printf("🔄 GetEnhancedAnalytics: Replacing stale request for user %s (old: %s, new: %s)",
				userID, existingReq.RequestID, requestID)
			h.activeRequests.Store(userID, &RequestInfo{
				UserID:    userID,
				StartTime: time.Now(),
				RequestID: requestID,
			})
		}
	}

	// Ensure cleanup happens regardless of how the function exits
	defer func() {
		h.activeRequests.Delete(userID)
		log.Printf("🧹 GetEnhancedAnalytics: Cleaned up active request for user %s (requestId: %s)", userID, requestID)
	}()

	// Create a context with timeout to prevent hanging requests
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Add early database health check
	if h.service.(*service).db.Health()["status"] != "up" {
		log.Printf("❌ GetEnhancedAnalytics: Database unhealthy for request %s", requestID)
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Database temporarily unavailable",
		})
	}

	// Ensure user exists in database before proceeding
	analyticsRepo := NewRepository(h.service.(*service).db)
	existingUser, err := analyticsRepo.GetUserByClerkID(ctx, userID)
	if err != nil {
		log.Printf("❌ GetEnhancedAnalytics: Failed to check user %s (requestId: %s): %v", userID, requestID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify user account",
		})
	}

	if existingUser == nil {
		log.Printf("🔍 GetEnhancedAnalytics: User %s doesn't exist, creating minimal record (requestId: %s)", userID, requestID)
		// Create minimal user record for cross-environment compatibility
		user := &User{
			ID:          userID,
			ClerkUserID: userID,
			Username:    fmt.Sprintf("user_%s", userID[len(userID)-10:]),
			DisplayName: "User",
			Email:       "",
		}

		if err := analyticsRepo.CreateOrUpdateUser(ctx, user); err != nil {
			log.Printf("❌ GetEnhancedAnalytics: Failed to create user %s (requestId: %s): %v", userID, requestID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to initialize user account",
			})
		}
		log.Printf("✅ GetEnhancedAnalytics: Created minimal user record for %s (requestId: %s)", userID, requestID)
	}

	// Get days parameter (default to 30)
	daysStr := c.Query("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30
	}

	log.Printf("📊 GetEnhancedAnalytics: Fetching analytics for user %s (days: %d, requestId: %s)", userID, days, requestID)

	// Get analytics with error handling and retries
	var analytics *EnhancedAnalytics
	var analyticsErr error

	// Single attempt with comprehensive error handling - no retries to avoid amplifying concurrency issues
	startTime := time.Now()
	analytics, analyticsErr = h.service.GetEnhancedAnalytics(ctx, userID, days)
	duration := time.Since(startTime)

	if analyticsErr != nil {
		log.Printf("❌ GetEnhancedAnalytics: Failed for user %s (requestId: %s, duration: %v): %v", userID, requestID, duration, analyticsErr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":     "Failed to get enhanced analytics",
			"details":   analyticsErr.Error(),
			"requestId": requestID,
		})
	}

	log.Printf("✅ GetEnhancedAnalytics: Successfully retrieved analytics for user %s (requestId: %s, duration: %v)", userID, requestID, duration)

	// Add comprehensive logging before response
	log.Printf("📊 GetEnhancedAnalytics: Response data for %s (requestId: %s) - Videos: %d total, Followers: %d, Views: %d",
		userID, requestID, analytics.Overview.VideoCount, analytics.Overview.CurrentFollowers, analytics.Overview.TotalViews)
	log.Printf("📊 GetEnhancedAnalytics: Response arrays for %s (requestId: %s) - Top videos: %d, Recent videos: %d",
		userID, requestID, len(analytics.TopVideos), len(analytics.RecentVideos))

	// Try to send the JSON response with error handling
	log.Printf("📤 GetEnhancedAnalytics: Sending JSON response for user %s (requestId: %s)", userID, requestID)

	// Set response headers to prevent caching issues
	c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")
	c.Set("X-Request-ID", requestID)

	if err := c.JSON(analytics); err != nil {
		log.Printf("❌ GetEnhancedAnalytics: Failed to send JSON response for user %s (requestId: %s): %v", userID, requestID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":     "Failed to serialize response",
			"requestId": requestID,
		})
	}

	log.Printf("✅ GetEnhancedAnalytics: Successfully sent response for user %s (requestId: %s, total duration: %v)", userID, requestID, time.Since(startTime))
	return nil
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

	// Trigger data collection in background
	h.backgroundCollectionMgr.TriggerUserCollection(userID)

	return c.JSON(fiber.Map{
		"message":   "Data collection triggered successfully",
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

	err = h.service.RefreshChannelData(c.Context(), userID)
	if err != nil {
		log.Printf("Error refreshing channel data for user %s: %v", userID, err)
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

	// Check database health
	if health := h.service.(*service).db.Health(); health["status"] == "up" {
		status["database_health"] = "healthy"
	} else {
		status["database_health"] = "unhealthy"
		if errorMsg, exists := health["error"]; exists {
			status["errors"] = append(status["errors"].([]string), fmt.Sprintf("Database: %v", errorMsg))
		}
	}

	// Check if user exists
	analyticsRepo := NewRepository(h.service.(*service).db)
	existingUser, err := analyticsRepo.GetUserByClerkID(c.Context(), userID)
	if err != nil {
		status["errors"] = append(status["errors"].([]string), fmt.Sprintf("User lookup failed: %v", err))
	} else if existingUser != nil {
		status["user_exists"] = true
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

// HealthCheck returns the health status of the analytics service
func (h *Handlers) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "healthy",
		"service":   "analytics",
		"timestamp": time.Now().Unix(),
	})
}
