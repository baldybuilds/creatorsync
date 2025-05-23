package server

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/baldybuilds/creatorsync/internal/analytics"
	"github.com/baldybuilds/creatorsync/internal/clerk"
	"github.com/baldybuilds/creatorsync/internal/database"
	"github.com/baldybuilds/creatorsync/internal/twitch"
)

type FiberServer struct {
	*fiber.App

	db                database.Service
	analyticsHandlers *analytics.Handlers
}

func New() (*FiberServer, error) {
	if err := clerk.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize Clerk client: %w", err)
	}

	db := database.New()

	// Initialize Twitch client
	twitchClientID := os.Getenv("TWITCH_CLIENT_ID")
	twitchClientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	if twitchClientID == "" || twitchClientSecret == "" {
		return nil, fmt.Errorf("TWITCH_CLIENT_ID and TWITCH_CLIENT_SECRET must be set")
	}

	twitchClient, err := twitch.NewClient(twitchClientID, twitchClientSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Twitch client: %w", err)
	}

	// Initialize analytics components
	analyticsService := analytics.NewService(db, twitchClient)
	dataCollector := analytics.NewDataCollector(analytics.NewRepository(db), twitchClient)
	backgroundMgr := analytics.NewBackgroundCollectionManager(dataCollector, db)
	analyticsHandlers := analytics.NewHandlers(analyticsService, backgroundMgr)

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "creatorsync",
			AppName:      "creatorsync",
		}),
		db:                db,
		analyticsHandlers: analyticsHandlers,
	}

	return server, nil
}
