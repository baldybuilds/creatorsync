package server

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/baldybuilds/creatorsync/internal/analytics"
	"github.com/baldybuilds/creatorsync/internal/clerk"
	"github.com/baldybuilds/creatorsync/internal/database"
	"github.com/baldybuilds/creatorsync/internal/twitch"
)

type FiberServer struct {
	*fiber.App

	db                      database.Service
	analyticsHandlers       *analytics.Handlers
	backgroundCollectionMgr *analytics.BackgroundCollectionManager
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
		db:                      db,
		analyticsHandlers:       analyticsHandlers,
		backgroundCollectionMgr: backgroundMgr,
	}

	return server, nil
}

// StartBackgroundServices starts all background services like analytics collection
func (s *FiberServer) StartBackgroundServices() error {
	log.Printf("🚀 Starting background analytics collection...")

	ctx := context.Background()
	if err := s.backgroundCollectionMgr.Start(ctx); err != nil {
		return fmt.Errorf("failed to start background collection manager: %w", err)
	}

	log.Printf("✅ Background analytics collection started successfully")
	return nil
}

// StopBackgroundServices stops all background services
func (s *FiberServer) StopBackgroundServices() error {
	log.Printf("🛑 Stopping background services...")

	if err := s.backgroundCollectionMgr.Stop(); err != nil {
		return fmt.Errorf("failed to stop background collection manager: %w", err)
	}

	log.Printf("✅ Background services stopped")
	return nil
}
