package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/baldybuilds/creatorsync/internal/analytics"
	"github.com/baldybuilds/creatorsync/internal/clerk"
	"github.com/baldybuilds/creatorsync/internal/email"
	"github.com/baldybuilds/creatorsync/internal/server/handlers"
	"github.com/baldybuilds/creatorsync/internal/twitch"

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

	// Check multiple environment indicators for staging
	isStaging := env == "staging" || environment == "staging" ||
		env == "dev" || environment == "dev" ||
		os.Getenv("DATABASE_URL") != "" // Cloud databases often indicate staging/production

	// Check for production
	isProduction := env == "production" || environment == "production"

	if isProduction {
		allowedOrigins = "https://creatorsync.app,https://www.creatorsync.app"
	} else if isStaging {
		allowedOrigins = "https://dev.creatorsync.app,https://creatorsync.app"
	} else {
		allowedOrigins = "http://localhost:3000,http://localhost:5173,http://localhost:5174,https://dev.creatorsync.app"
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

	// Register Analytics routes (includes both public and protected routes)
	s.registerAnalyticsRoutes()

	// Protected routes group
	api := s.App.Group("/api")
	api.Use(clerk.AuthMiddleware())

	// User routes
	api.Get("/user", s.getCurrentUserHandler)
	api.Get("/user/profile", s.getUserProfileHandler)
	api.Post("/user/sync", s.syncUserHandler)

	// Register Twitch routes
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

	// Check if user already exists in our database
	analyticsRepo := analytics.NewRepository(s.db)
	existingUser, err := analyticsRepo.GetUserByClerkID(ctx, userID)
	if err != nil {
		log.Printf("❌ ensureUserExistsInDatabase: Failed to check existing user %s: %v", userID, err)
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		log.Printf("✅ ensureUserExistsInDatabase: User %s already exists", userID)
		return nil // User already exists
	}

	log.Printf("🔍 ensureUserExistsInDatabase: User %s doesn't exist, creating...", userID)

	// User doesn't exist, let's create them
	// Try to get user's Clerk profile, but handle cross-environment cases gracefully
	clerkUser, err := clerk.GetUserByID(ctx, userID)
	if err != nil {
		// If Clerk user doesn't exist (e.g., cross-environment user ID), create basic record
		log.Printf("⚠️ ensureUserExistsInDatabase: Clerk user %s not found in current environment, creating basic user record: %v", userID, err)

		// Create minimal user record with just the Clerk ID
		user := &analytics.User{
			ID:          userID,
			ClerkUserID: userID,
			Username:    fmt.Sprintf("user_%s", userID[5:15]), // Create a basic username from ID
			DisplayName: "Unknown User",
			Email:       "",
		}

		// Create user record in database
		if err := analyticsRepo.CreateOrUpdateUser(ctx, user); err != nil {
			log.Printf("❌ ensureUserExistsInDatabase: Failed to create basic user record for %s: %v", userID, err)
			return fmt.Errorf("failed to create basic user record: %w", err)
		}

		log.Printf("✅ ensureUserExistsInDatabase: Created basic user record for cross-environment user %s", userID)
		return nil
	}

	log.Printf("✅ ensureUserExistsInDatabase: Got Clerk user for %s, building full record", userID)

	// Initialize user with basic info from Clerk
	user := &analytics.User{
		ID:          userID,
		ClerkUserID: userID,
	}

	// Safely set email if available
	if len(clerkUser.EmailAddresses) > 0 {
		user.Email = clerkUser.EmailAddresses[0].EmailAddress
		log.Printf("📧 ensureUserExistsInDatabase: Set email for user %s", userID)
	}

	// Set name fields safely
	if clerkUser.FirstName != nil {
		user.DisplayName = *clerkUser.FirstName
	}
	if clerkUser.LastName != nil && *clerkUser.LastName != "" {
		if user.DisplayName != "" {
			user.DisplayName += " " + *clerkUser.LastName
		} else {
			user.DisplayName = *clerkUser.LastName
		}
	}

	log.Printf("👤 ensureUserExistsInDatabase: Set display name '%s' for user %s", user.DisplayName, userID)

	// Try to get Twitch info if available
	for _, account := range clerkUser.ExternalAccounts {
		if account.Provider == "oauth_twitch" {
			log.Printf("🎮 ensureUserExistsInDatabase: Found Twitch account for user %s", userID)
			user.TwitchUserID = account.ProviderUserID
			if account.Username != nil {
				user.Username = *account.Username
			}

			// Try to get additional Twitch info if we have OAuth token
			if token, tokenErr := clerk.GetOAuthToken(ctx, userID, "oauth_twitch"); tokenErr == nil {
				log.Printf("🔑 ensureUserExistsInDatabase: Got Twitch token for user %s, fetching profile", userID)
				// Initialize Twitch client
				twitchClientID := os.Getenv("TWITCH_CLIENT_ID")
				twitchClientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
				if twitchClientID != "" && twitchClientSecret != "" {
					if twitchClient, clientErr := twitch.NewClient(twitchClientID, twitchClientSecret); clientErr == nil {
						if userInfo, infoErr := twitchClient.GetUserInfo(token); infoErr == nil {
							user.Username = userInfo.Login
							user.DisplayName = userInfo.DisplayName
							user.ProfileImageURL = userInfo.ProfileImageURL
							if userInfo.Email != "" {
								user.Email = userInfo.Email
							}
							log.Printf("✅ ensureUserExistsInDatabase: Enhanced user %s with Twitch profile data", userID)
						} else {
							log.Printf("⚠️ ensureUserExistsInDatabase: Failed to get Twitch user info for %s: %v", userID, infoErr)
						}
					} else {
						log.Printf("⚠️ ensureUserExistsInDatabase: Failed to create Twitch client for %s: %v", userID, clientErr)
					}
				} else {
					log.Printf("⚠️ ensureUserExistsInDatabase: Missing Twitch client credentials")
				}
			} else {
				log.Printf("⚠️ ensureUserExistsInDatabase: Failed to get Twitch token for %s: %v", userID, tokenErr)
			}
			break
		}
	}

	// Create user record in database
	log.Printf("💾 ensureUserExistsInDatabase: Creating database record for user %s", userID)
	if err := analyticsRepo.CreateOrUpdateUser(ctx, user); err != nil {
		log.Printf("❌ ensureUserExistsInDatabase: Failed to create user record for %s: %v", userID, err)
		return fmt.Errorf("failed to create user record: %w", err)
	}

	log.Printf("✅ ensureUserExistsInDatabase: Created user record for %s (%s)", user.DisplayName, userID)
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

	// Check if this is a new user by seeing if they exist in our database
	analyticsRepo := analytics.NewRepository(s.db)
	existingUser, err := analyticsRepo.GetUserByClerkID(ctx, user.ID)
	isNewUser := (err != nil || existingUser == nil)

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

func (s *FiberServer) registerTwitchRoutes(api fiber.Router) {
	twitchGroup := api.Group("/twitch")
	twitchGroup.Get("/channel", handlers.GetTwitchChannelHandler)
	twitchGroup.Get("/streams", handlers.GetTwitchStreamsHandler)
	twitchGroup.Get("/videos", handlers.GetTwitchVideosHandler)
	twitchGroup.Get("/clips", handlers.GetTwitchClipsHandler)
	twitchGroup.Get("/callback", handlers.TwitchCallbackHandler)
	twitchGroup.Get("/subscribers", handlers.GetTwitchSubscribersHandler)
	twitchGroup.Get("/analytics/video_summary", handlers.GetTwitchVideoAnalyticsSummaryHandler)
}

func (s *FiberServer) registerAnalyticsRoutes() {
	s.analyticsHandlers.RegisterRoutes(s.App)
}
