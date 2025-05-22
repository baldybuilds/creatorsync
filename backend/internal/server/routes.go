package server

import (
	"fmt"
	"os"

	"github.com/baldybuilds/creatorsync/internal/clerk"
	"github.com/baldybuilds/creatorsync/internal/email"
	"github.com/baldybuilds/creatorsync/internal/server/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {
	env := os.Getenv("APP_ENV")
	allowedOrigins := "*"

	if env == "production" {
		allowedOrigins = "https://creatorsync.app,https://www.creatorsync.app"
	} else {
		allowedOrigins = "http://localhost:3000,http://localhost:5173,http://localhost:5174"
	}

	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Public routes
	s.App.Get("/", s.HelloWorldHandler)
	s.App.Get("/health", s.healthHandler)
	s.App.Post("/api/waitlist", s.joinWaitlistHandler)

	// Protected routes group
	api := s.App.Group("/api")
	api.Use(clerk.AuthMiddleware())

	// User routes
	api.Get("/user", s.getCurrentUserHandler)
	api.Get("/user/profile", s.getUserProfileHandler)

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

func (s *FiberServer) registerTwitchRoutes(api fiber.Router) {
	twitchGroup := api.Group("/twitch")
	twitchGroup.Get("/channel", handlers.GetTwitchChannelHandler)
	twitchGroup.Get("/streams", handlers.GetTwitchStreamsHandler)
	twitchGroup.Get("/videos", handlers.GetTwitchVideosHandler)
	twitchGroup.Get("/callback", handlers.TwitchCallbackHandler)
	twitchGroup.Get("/subscribers", handlers.GetTwitchSubscribersHandler)
	twitchGroup.Get("/analytics/video_summary", handlers.GetTwitchVideoAnalyticsSummaryHandler)
}
