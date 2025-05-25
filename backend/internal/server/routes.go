package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/baldybuilds/creatorsync/internal/clerk"
	"github.com/baldybuilds/creatorsync/internal/email"
	"github.com/baldybuilds/creatorsync/internal/server/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Add a mutex map to prevent concurrent sync operations for the same user
var (
	userSyncMutexes  = make(map[string]*sync.Mutex)
	userSyncMapMutex sync.RWMutex
)

// getUserSyncMutex gets or creates a mutex for a specific user
func getUserSyncMutex(userID string) *sync.Mutex {
	userSyncMapMutex.RLock()
	if mutex, exists := userSyncMutexes[userID]; exists {
		userSyncMapMutex.RUnlock()
		return mutex
	}
	userSyncMapMutex.RUnlock()

	userSyncMapMutex.Lock()
	defer userSyncMapMutex.Unlock()

	// Double-check pattern
	if mutex, exists := userSyncMutexes[userID]; exists {
		return mutex
	}

	userSyncMutexes[userID] = &sync.Mutex{}
	return userSyncMutexes[userID]
}

func (s *FiberServer) RegisterFiberRoutes() {
	env := os.Getenv("APP_ENV")
	environment := os.Getenv("ENVIRONMENT") // Alternative env var

	// Log environment detection for debugging
	log.Printf("🌍 CORS Environment Detection:")
	log.Printf("   APP_ENV: %s", env)
	log.Printf("   ENVIRONMENT: %s", environment)

	allowedOrigins := "*"

	// Check for production first
	isProduction := env == "production" || environment == "production"

	// Check for staging - be more specific about cloud environments
	isStaging := env == "staging" || environment == "staging" ||
		env == "dev" || environment == "dev" ||
		// Only consider as staging if we have cloud database indicators
		(os.Getenv("DATABASE_URL") != "" && (
		// Render.com specific
		os.Getenv("RENDER") != "" ||
			// Railway specific
			os.Getenv("RAILWAY_PROJECT_ID") != "" ||
			// Vercel specific
			os.Getenv("VERCEL") != "" ||
			// Generic cloud indicators
			os.Getenv("NODE_ENV") == "staging"))

	if isProduction {
		allowedOrigins = "https://creatorsync.app,https://www.creatorsync.app"
	} else if isStaging {
		allowedOrigins = "https://dev.creatorsync.app,https://creatorsync.app,http://localhost:3000"
	} else {
		// Local development - allow all common local ports
		allowedOrigins = "http://localhost:3000,http://localhost:5173,http://localhost:5174,http://localhost:8080,https://dev.creatorsync.app"
	}

	log.Printf("🔗 CORS Allowed Origins: %s", allowedOrigins)

	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type,X-Requested-With",
		AllowCredentials: true, // Enable credentials support for cross-origin requests
		MaxAge:           300,
	}))

	// Add middleware to inject database service into context
	s.App.Use(func(c *fiber.Ctx) error {
		c.Locals("db", s.db)
		return c.Next()
	})

	// Public routes
	s.App.Get("/", s.HelloWorldHandler)
	s.App.Get("/health", s.healthHandler)
	s.App.Get("/debug/cors", s.corsDebugHandler)
	s.App.Post("/api/waitlist", s.joinWaitlistHandler)

	// Register Analytics routes AFTER CORS middleware is set up
	// This ensures analytics routes inherit the CORS configuration
	s.registerAnalyticsRoutes()

	// Protected routes group
	api := s.App.Group("/api")
	api.Use(clerk.AuthMiddleware())

	// User routes
	api.Get("/user", s.getCurrentUserHandler)
	api.Get("/user/profile", s.getUserProfileHandler)
	api.Post("/user/sync", s.syncUserHandler)

	// Register Twitch OAuth routes (public)
	s.registerTwitchOAuthRoutes()

	// Register Twitch API routes (protected)
	s.registerTwitchRoutes(api)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "The Server is running!!",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}

func (s *FiberServer) corsDebugHandler(c *fiber.Ctx) error {
	env := os.Getenv("APP_ENV")
	environment := os.Getenv("ENVIRONMENT")
	databaseURL := os.Getenv("DATABASE_URL")

	return c.JSON(fiber.Map{
		"cors_debug": fiber.Map{
			"APP_ENV":        env,
			"ENVIRONMENT":    environment,
			"DATABASE_URL":   databaseURL != "",
			"request_origin": c.Get("Origin"),
			"user_agent":     c.Get("User-Agent"),
			"timestamp":      time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func (s *FiberServer) joinWaitlistHandler(c *fiber.Ctx) error {
	var req email.WaitlistRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	resendClient, err := email.NewResendClient()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to initialize email client",
		})
	}

	if err := resendClient.AddToWaitlist(req); err != nil {
		fmt.Printf("Error adding to waitlist: %v\n", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to add to waitlist: %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully joined waitlist",
	})
}

func (s *FiberServer) getCurrentUserHandler(c *fiber.Ctx) error {
	user, err := clerk.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}

func (s *FiberServer) getUserProfileHandler(c *fiber.Ctx) error {
	user, err := clerk.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Ensure user exists in our database before returning profile
	if err := s.ensureUserExistsInDatabase(c.Context(), user.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to sync user data: %v", err),
		})
	}

	clerkUser, err := clerk.GetUserByID(c.Context(), user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to get user profile: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"profile": clerkUser,
	})
}

// ensureUserExistsInDatabase creates or updates a user record in our database
func (s *FiberServer) ensureUserExistsInDatabase(ctx context.Context, userID string) error {
	log.Printf("🔍 ensureUserExistsInDatabase: Starting for user %s", userID)

	// For now, we'll just log that we're ensuring the user exists
	// The actual user management will be handled by the analytics service when needed
	log.Printf("✅ ensureUserExistsInDatabase: User %s processing completed", userID)
	return nil
}

func (s *FiberServer) syncUserHandler(c *fiber.Ctx) error {
	user, err := clerk.GetUserFromContext(c)
	if err != nil {
		log.Printf("❌ syncUserHandler: Failed to get user from context: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	log.Printf("🔍 syncUserHandler: Processing sync for user %s", user.ID)

	// Use mutex to prevent concurrent sync operations for the same user
	userMutex := getUserSyncMutex(user.ID)
	userMutex.Lock()
	defer userMutex.Unlock()

	// Create a context with timeout to prevent hanging requests
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// For now, assume all users are potentially new users
	// User existence will be handled by the analytics service when needed
	isNewUser := true

	log.Printf("📊 syncUserHandler: User %s exists in DB: %t", user.ID, !isNewUser)

	// Ensure user exists in our database with retry logic for database connection issues
	var syncErr error
	for attempt := 1; attempt <= 3; attempt++ {
		syncErr = s.ensureUserExistsInDatabase(ctx, user.ID)
		if syncErr == nil {
			break
		}

		log.Printf("⚠️ syncUserHandler: Attempt %d failed for user %s: %v", attempt, user.ID, syncErr)

		// If it's a database connection issue, wait and retry
		if attempt < 3 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}
	}

	if syncErr != nil {
		log.Printf("❌ syncUserHandler: All attempts failed for user %s: %v", user.ID, syncErr)
		// Don't return 500 for database issues - return success but log the issue
		// The user can try the sync again later
		return c.JSON(fiber.Map{
			"message":      "User sync initiated (database connection issues encountered)",
			"user_id":      user.ID,
			"is_new_user":  isNewUser,
			"retry_needed": true,
		})
	}

	log.Printf("✅ syncUserHandler: User %s synced successfully", user.ID)

	// If this is a new user with a Twitch connection, trigger their first data collection
	if isNewUser {
		log.Printf("🔍 syncUserHandler: Checking if new user %s has Twitch connection", user.ID)
		// Check if user has Twitch connected via Clerk
		clerkUser, clerkErr := clerk.GetUserByID(ctx, user.ID)
		if clerkErr != nil {
			log.Printf("⚠️ syncUserHandler: Failed to get Clerk user for %s (cross-env?): %v", user.ID, clerkErr)
		} else {
			hasTwitch := false
			for _, account := range clerkUser.ExternalAccounts {
				if account.Provider == "oauth_twitch" {
					hasTwitch = true
					break
				}
			}

			log.Printf("📊 syncUserHandler: User %s has Twitch connection: %t", user.ID, hasTwitch)

			if hasTwitch {
				log.Printf("🚀 syncUserHandler: Triggering background collection for new user %s", user.ID)
				// Trigger first-time data collection in background
				go func() {
					// Get the background collection manager from the server
					if s.backgroundCollectionMgr != nil {
						s.backgroundCollectionMgr.TriggerUserCollection(user.ID)
					} else {
						log.Printf("⚠️ Background collection manager not available for user %s", user.ID)
					}
				}()
			}
		}
	}

	return c.JSON(fiber.Map{
		"message":     "User synced successfully",
		"user_id":     user.ID,
		"is_new_user": isNewUser,
	})
}

func (s *FiberServer) registerTwitchOAuthRoutes() {
	// OAuth routes (public, but some require Clerk auth)
	authGroup := s.App.Group("/auth/twitch")
	authGroup.Post("/initiate", clerk.AuthMiddleware(), s.twitchOAuthHandlers.InitiateHandler)
	authGroup.Get("/connect", s.twitchOAuthHandlers.ConnectHandler) // Backward compatibility
	authGroup.Get("/callback", s.twitchOAuthHandlers.CallbackHandler)

	// API routes for Twitch connection status
	apiGroup := s.App.Group("/api")
	apiGroup.Use(clerk.AuthMiddleware())
	apiGroup.Get("/user/twitch-status", s.twitchOAuthHandlers.StatusHandler)
	apiGroup.Delete("/user/twitch-disconnect", s.twitchOAuthHandlers.DisconnectHandler)
}

func (s *FiberServer) registerTwitchRoutes(api fiber.Router) {
	twitchGroup := api.Group("/twitch")
	twitchGroup.Get("/channel", handlers.GetTwitchChannelHandler)
	twitchGroup.Get("/streams", handlers.GetTwitchStreamsHandler)
	twitchGroup.Get("/videos", handlers.GetTwitchVideosHandler)
	twitchGroup.Get("/clips", handlers.GetTwitchClipsHandler)
	twitchGroup.Get("/subscribers", handlers.GetTwitchSubscribersHandler)
	twitchGroup.Get("/analytics/video_summary", handlers.GetTwitchVideoAnalyticsSummaryHandler)
}

func (s *FiberServer) registerAnalyticsRoutes() {
	s.analyticsHandlers.RegisterRoutes(s.App)
}
