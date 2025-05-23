package server

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/baldybuilds/creatorsync/internal/analytics"
	"github.com/baldybuilds/creatorsync/internal/clerk"
	"github.com/baldybuilds/creatorsync/internal/email"
	"github.com/baldybuilds/creatorsync/internal/server/handlers"
	"github.com/baldybuilds/creatorsync/internal/twitch"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {
	env := os.Getenv("APP_ENV")
	allowedOrigins := "*"

	if env == "production" {
		allowedOrigins = "https://creatorsync.app,https://www.creatorsync.app"
	} else if env == "staging" {
		allowedOrigins = "https://dev.creatorsync.app"
	} else {
		allowedOrigins = "http://localhost:3000,http://localhost:5173,http://localhost:5174"
	}

	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
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
	// Check if user already exists in our database
	analyticsRepo := analytics.NewRepository(s.db.GetDB())
	existingUser, err := analyticsRepo.GetUserByClerkID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		return nil // User already exists
	}

	// User doesn't exist, let's create them
	// Get user's Clerk profile
	clerkUser, err := clerk.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user from Clerk: %w", err)
	}

	// Initialize user with basic info from Clerk
	user := &analytics.User{
		ID:          userID,
		ClerkUserID: userID,
	}

	// Safely set email if available
	if len(clerkUser.EmailAddresses) > 0 {
		user.Email = clerkUser.EmailAddresses[0].EmailAddress
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

	// Try to get Twitch info if available
	for _, account := range clerkUser.ExternalAccounts {
		if account.Provider == "oauth_twitch" {
			user.TwitchUserID = account.ProviderUserID
			if account.Username != nil {
				user.Username = *account.Username
			}

			// Try to get additional Twitch info if we have OAuth token
			if token, tokenErr := clerk.GetOAuthToken(ctx, userID, "oauth_twitch"); tokenErr == nil {
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
						}
					}
				}
			}
			break
		}
	}

	// Create user record in database
	if err := analyticsRepo.CreateOrUpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to create user record: %w", err)
	}

	log.Printf("✅ Created user record for %s (%s)", user.DisplayName, userID)
	return nil
}

func (s *FiberServer) syncUserHandler(c *fiber.Ctx) error {
	user, err := clerk.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Ensure user exists in our database
	if err := s.ensureUserExistsInDatabase(c.Context(), user.ID); err != nil {
		log.Printf("Failed to sync user %s: %v", user.ID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to sync user data: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User synced successfully",
		"user_id": user.ID,
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
