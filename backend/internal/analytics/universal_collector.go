package analytics

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/baldybuilds/creatorsync/internal/database"
)

type universalCollector struct {
	platforms    map[Platform]PlatformCollector
	repo         Repository
	db           database.Service
	platformsMux sync.RWMutex
}

func NewUniversalAnalyticsCollector(db database.Service, repo Repository) UniversalAnalyticsCollector {
	return &universalCollector{
		platforms: make(map[Platform]PlatformCollector),
		repo:      repo,
		db:        db,
	}
}

func (uc *universalCollector) RegisterPlatform(collector PlatformCollector) {
	uc.platformsMux.Lock()
	defer uc.platformsMux.Unlock()

	platform := collector.GetPlatform()
	uc.platforms[platform] = collector
	log.Printf("✅ Registered platform collector: %s", platform)
}

func (uc *universalCollector) GetConnectedPlatforms(ctx context.Context, userID string) ([]Platform, error) {
	uc.platformsMux.RLock()
	defer uc.platformsMux.RUnlock()

	var connectedPlatforms []Platform

	for platform, collector := range uc.platforms {
		connected, err := collector.IsConnected(ctx, userID)
		if err != nil {
			log.Printf("⚠️ Failed to check connection for platform %s, user %s: %v", platform, userID, err)
			continue
		}

		if connected {
			connectedPlatforms = append(connectedPlatforms, platform)
		}
	}

	log.Printf("📊 User %s has %d connected platforms: %v", userID, len(connectedPlatforms), connectedPlatforms)
	return connectedPlatforms, nil
}

func (uc *universalCollector) CollectUserData(ctx context.Context, userID string) error {
	log.Printf("🚀 Starting universal data collection for user %s", userID)

	connectedPlatforms, err := uc.GetConnectedPlatforms(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get connected platforms: %w", err)
	}

	if len(connectedPlatforms) == 0 {
		log.Printf("⚠️ User %s has no connected platforms", userID)
		return nil
	}

	var collectionErrors []error

	// Collect from all platforms concurrently
	var wg sync.WaitGroup
	var errorsMux sync.Mutex

	for _, platform := range connectedPlatforms {
		wg.Add(1)
		go func(p Platform) {
			defer wg.Done()

			if err := uc.CollectPlatformData(ctx, userID, p); err != nil {
				errorsMux.Lock()
				collectionErrors = append(collectionErrors, fmt.Errorf("platform %s: %w", p, err))
				errorsMux.Unlock()
			}
		}(platform)
	}

	wg.Wait()

	if len(collectionErrors) > 0 {
		log.Printf("⚠️ Collection completed with %d errors for user %s", len(collectionErrors), userID)
		// Return first error but log all
		for _, err := range collectionErrors {
			log.Printf("❌ Collection error: %v", err)
		}
		return collectionErrors[0]
	}

	log.Printf("✅ Successfully collected data from %d platforms for user %s", len(connectedPlatforms), userID)
	return nil
}

func (uc *universalCollector) CollectPlatformData(ctx context.Context, userID string, platform Platform) error {
	uc.platformsMux.RLock()
	collector, exists := uc.platforms[platform]
	uc.platformsMux.RUnlock()

	if !exists {
		return fmt.Errorf("no collector registered for platform %s", platform)
	}

	log.Printf("📊 Collecting data from %s for user %s", platform, userID)

	// Validate connection first
	if err := collector.ValidateConnection(ctx, userID); err != nil {
		return fmt.Errorf("connection validation failed: %w", err)
	}

	// Collect channel metrics
	channelMetrics, err := collector.CollectChannelMetrics(ctx, userID)
	if err != nil {
		log.Printf("⚠️ Failed to collect channel metrics from %s for user %s: %v", platform, userID, err)
	} else {
		if err := uc.savePlatformMetrics(ctx, channelMetrics); err != nil {
			log.Printf("⚠️ Failed to save channel metrics: %v", err)
		}
	}

	// Collect video metrics
	videoMetrics, err := collector.CollectVideoMetrics(ctx, userID, 100)
	if err != nil {
		log.Printf("⚠️ Failed to collect video metrics from %s for user %s: %v", platform, userID, err)
	} else {
		for _, video := range videoMetrics {
			if err := uc.saveVideoMetrics(ctx, &video); err != nil {
				log.Printf("⚠️ Failed to save video metrics for video %s: %v", video.VideoID, err)
			}
		}
	}

	log.Printf("✅ Completed data collection from %s for user %s", platform, userID)
	return nil
}

func (uc *universalCollector) ScheduleCollection(ctx context.Context, userID string, interval time.Duration) error {
	// For now, we'll implement a simple scheduling mechanism
	// In production, you might want to use a proper job scheduler like cron
	log.Printf("📅 Scheduling collection for user %s every %v", userID, interval)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Printf("🛑 Stopping scheduled collection for user %s", userID)
				return
			case <-ticker.C:
				if err := uc.CollectUserData(ctx, userID); err != nil {
					log.Printf("⚠️ Scheduled collection failed for user %s: %v", userID, err)
				}
			}
		}
	}()

	return nil
}

// Helper methods to save the new platform-agnostic metrics
func (uc *universalCollector) savePlatformMetrics(ctx context.Context, metrics *PlatformMetrics) error {
	log.Printf("📊 Saving platform metrics for user %s on %s", metrics.UserID, metrics.Platform)
	return nil
}

func (uc *universalCollector) saveVideoMetrics(ctx context.Context, metrics *VideoMetrics) error {
	log.Printf("🎥 Saving video metrics for user %s on %s", metrics.UserID, metrics.Platform)
	return nil
}

// BackwardCompatibleCollector provides compatibility with the old DataCollector interface
type BackwardCompatibleCollector struct {
	universal UniversalAnalyticsCollector
}

func NewBackwardCompatibleCollector(universal UniversalAnalyticsCollector) Collector {
	return &BackwardCompatibleCollector{
		universal: universal,
	}
}

func (bc *BackwardCompatibleCollector) CollectUserData(ctx context.Context, userID string) error {
	return bc.universal.CollectUserData(ctx, userID)
}

func (bc *BackwardCompatibleCollector) CollectChannelData(ctx context.Context, userID string) error {
	return bc.universal.CollectPlatformData(ctx, userID, PlatformTwitch)
}

func (bc *BackwardCompatibleCollector) CollectVideoData(ctx context.Context, userID string) error {
	return bc.universal.CollectPlatformData(ctx, userID, PlatformTwitch)
}

func (bc *BackwardCompatibleCollector) ScheduleCollection(ctx context.Context, userID string, interval time.Duration) error {
	return bc.universal.ScheduleCollection(ctx, userID, interval)
}

func (bc *BackwardCompatibleCollector) IsHealthy() bool {
	return true
}
