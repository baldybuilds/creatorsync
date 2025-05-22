package analytics

import (
	"context"
	"log"
	"time"

	"github.com/baldybuilds/creatorsync/internal/database"
)

type Scheduler interface {
	Start(ctx context.Context) error
	Stop() error
	ScheduleDailyCollection()
	TriggerUserCollection(userID string)
}

type scheduler struct {
	collector   DataCollector
	db          database.Service
	ticker      *time.Ticker
	stopChannel chan bool
	running     bool
}

func NewScheduler(collector DataCollector, db database.Service) Scheduler {
	return &scheduler{
		collector:   collector,
		db:          db,
		stopChannel: make(chan bool),
		running:     false,
	}
}

func (s *scheduler) Start(ctx context.Context) error {
	if s.running {
		return nil
	}

	log.Println("Starting analytics scheduler...")
	s.running = true

	// Schedule daily collection at 2 AM UTC
	s.ticker = time.NewTicker(1 * time.Hour) // Check every hour

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.checkAndRunDailyCollection(ctx)
			case <-s.stopChannel:
				return
			}
		}
	}()

	log.Println("Analytics scheduler started successfully")
	return nil
}

func (s *scheduler) Stop() error {
	if !s.running {
		return nil
	}

	log.Println("Stopping analytics scheduler...")
	s.running = false

	if s.ticker != nil {
		s.ticker.Stop()
	}

	s.stopChannel <- true
	log.Println("Analytics scheduler stopped")
	return nil
}

func (s *scheduler) ScheduleDailyCollection() {
	ctx := context.Background()
	s.runDailyCollectionForAllUsers(ctx)
}

func (s *scheduler) TriggerUserCollection(userID string) {
	ctx := context.Background()
	go func() {
		if err := s.collector.CollectAllUserData(ctx, userID); err != nil {
			log.Printf("Failed to collect data for user %s: %v", userID, err)
		}
	}()
}

func (s *scheduler) checkAndRunDailyCollection(ctx context.Context) {
	now := time.Now().UTC()

	// Run daily collection at 2 AM UTC
	if now.Hour() == 2 && now.Minute() == 0 {
		log.Println("Starting daily analytics collection...")
		s.runDailyCollectionForAllUsers(ctx)
	}
}

func (s *scheduler) runDailyCollectionForAllUsers(ctx context.Context) {
	// Get all users from database
	users, err := s.getAllUsers(ctx)
	if err != nil {
		log.Printf("Failed to get users for daily collection: %v", err)
		return
	}

	log.Printf("Starting daily collection for %d users", len(users))

	// Process users in batches to avoid overwhelming the API
	batchSize := 10
	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}

		batch := users[i:end]
		s.processBatch(ctx, batch)

		// Wait between batches to respect rate limits
		time.Sleep(30 * time.Second)
	}

	log.Println("Daily collection completed")
}

func (s *scheduler) processBatch(ctx context.Context, users []string) {
	for _, userID := range users {
		go func(uid string) {
			// Add some jitter to avoid hitting rate limits
			time.Sleep(time.Duration(len(uid)%10) * time.Second)

			if err := s.collector.CollectDailyChannelData(ctx, uid); err != nil {
				log.Printf("Failed daily collection for user %s: %v", uid, err)
			} else {
				log.Printf("Completed daily collection for user %s", uid)
			}
		}(userID)
	}
}

func (s *scheduler) getAllUsers(ctx context.Context) ([]string, error) {
	query := `
		SELECT DISTINCT id 
		FROM users 
		WHERE twitch_user_id IS NOT NULL 
		AND twitch_user_id != ''
	`

	rows, err := s.db.GetDB().QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			continue
		}
		users = append(users, userID)
	}

	return users, nil
}

// BackgroundCollectionManager manages all background collection tasks
type BackgroundCollectionManager struct {
	scheduler Scheduler
	collector DataCollector
}

func NewBackgroundCollectionManager(collector DataCollector, db database.Service) *BackgroundCollectionManager {
	scheduler := NewScheduler(collector, db)
	return &BackgroundCollectionManager{
		scheduler: scheduler,
		collector: collector,
	}
}

func (bcm *BackgroundCollectionManager) Start(ctx context.Context) error {
	return bcm.scheduler.Start(ctx)
}

func (bcm *BackgroundCollectionManager) Stop() error {
	return bcm.scheduler.Stop()
}

func (bcm *BackgroundCollectionManager) TriggerUserCollection(userID string) {
	bcm.scheduler.TriggerUserCollection(userID)
}

func (bcm *BackgroundCollectionManager) TriggerDailyCollection() {
	bcm.scheduler.ScheduleDailyCollection()
}
