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
	universalCollector      analytics.UniversalAnalyticsCollector
	twitchOAuthHandlers     *twitch.TwitchOAuthHandlers
	twitchTokenHelper       *twitch.TwitchTokenHelper
}

func New() (*FiberServer, error) {
	if err := clerk.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize Clerk client: %w", err)
	}

	dbService, err := database.NewDatabaseService()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database service: %w", err)
	}

	standardDB := dbService.GetStandardDB()

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

	// Initialize analytics components with new universal system
	analyticsRepo := analytics.NewRepository(standardDB)
	universalCollector := analytics.NewUniversalAnalyticsCollector(standardDB, analyticsRepo)

	// Register platform collectors
	twitchCollector, err := analytics.NewTwitchCollector(standardDB, analyticsRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Twitch collector: %w", err)
	}
	universalCollector.RegisterPlatform(twitchCollector)

	// Create backward-compatible collector for existing code
	dataCollector := analytics.NewBackwardCompatibleCollector(universalCollector)

	// Initialize other components
	analyticsService := analytics.NewService(dbService, twitchClient)
	backgroundMgr := analytics.NewBackgroundCollectionManager(dataCollector, standardDB)
	analyticsHandlers := analytics.NewHandlers(analyticsService, backgroundMgr)

	// Initialize Twitch OAuth components
	twitchOAuthHandlers, err := twitch.NewTwitchOAuthHandlers(standardDB)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Twitch OAuth handlers: %w", err)
	}

	twitchTokenHelper, err := twitch.NewTwitchTokenHelper(standardDB)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Twitch token helper: %w", err)
	}

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "creatorsync",
			AppName:      "creatorsync",
		}),
		db:                      standardDB,
		analyticsHandlers:       analyticsHandlers,
		backgroundCollectionMgr: backgroundMgr,
		universalCollector:      universalCollector,
		twitchOAuthHandlers:     twitchOAuthHandlers,
		twitchTokenHelper:       twitchTokenHelper,
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
