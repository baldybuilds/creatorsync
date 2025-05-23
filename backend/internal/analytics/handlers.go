package analytics

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/baldybuilds/creatorsync/internal/clerk"
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	service                 Service
	backgroundCollectionMgr *BackgroundCollectionManager
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
	api := app.Group("/api/analytics")

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

	// Check if we need to trigger automatic data collection
	h.triggerAutoDataCollectionIfNeeded(userID)

	overview, err := h.service.GetDashboardOverview(c.Context(), userID)
	if err != nil {
		log.Printf("Error getting dashboard overview for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get dashboard overview",
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
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	log.Printf("🔍 Enhanced analytics request for user %s", userID)

	// Check if we need to trigger automatic data collection
	h.triggerAutoDataCollectionIfNeeded(userID)

	// Get days parameter (default to 30)
	daysStr := c.Query("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30
	}

	log.Printf("📊 Fetching enhanced analytics for user %s (days: %d)", userID, days)
	analytics, err := h.service.GetEnhancedAnalytics(c.Context(), userID, days)
	if err != nil {
		log.Printf("❌ Error getting enhanced analytics for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get enhanced analytics",
		})
	}

	log.Printf("✅ Enhanced analytics response for user %s: %+v", userID, analytics.Overview)
	return c.JSON(analytics)
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
		log.Printf("Error getting analytics jobs for user %s: %v", userID, err)
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

// triggerAutoDataCollectionIfNeeded checks if we should automatically collect data for a user
func (h *Handlers) triggerAutoDataCollectionIfNeeded(userID string) {
	log.Printf("🔍 Checking if data collection needed for user %s", userID)

	// Check if user has any analytics data
	hasData, lastUpdate, err := h.service.CheckUserAnalyticsData(context.Background(), userID)
	if err != nil {
		log.Printf("❌ Error checking analytics data for user %s: %v", userID, err)
		return
	}

	log.Printf("📊 Data check for user %s: hasData=%v, lastUpdate=%v", userID, hasData, lastUpdate)

	shouldCollect := false
	reason := ""

	if !hasData {
		// No data exists - trigger collection for new users
		shouldCollect = true
		reason = "no existing data"
	} else if lastUpdate != nil {
		// Check if data is stale (older than 6 hours)
		staleThreshold := time.Now().Add(-6 * time.Hour)
		if lastUpdate.Before(staleThreshold) {
			shouldCollect = true
			reason = "data is stale (older than 6 hours)"
		} else {
			log.Printf("✅ Data is fresh for user %s (last update: %v)", userID, lastUpdate)
		}
	}

	if shouldCollect {
		log.Printf("🔄 Auto-triggering data collection for user %s: %s", userID, reason)
		h.backgroundCollectionMgr.TriggerUserCollection(userID)
	} else {
		log.Printf("⏭️ No data collection needed for user %s", userID)
	}
}

// GetDataStatus returns debug information about user's analytics data
func (h *Handlers) GetDataStatus(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	hasData, lastUpdate, err := h.service.CheckUserAnalyticsData(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"user_id":     userID,
		"has_data":    hasData,
		"last_update": lastUpdate,
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
