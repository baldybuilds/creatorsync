package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestUser represents a test user for local development
type TestUser struct {
	ID                string  `gorm:"primaryKey"`
	ClerkID           string  `gorm:"column:clerk_user_id;uniqueIndex"`
	TwitchUserID      *string `gorm:"column:twitch_user_id"`
	Username          string
	DisplayName       *string `gorm:"column:display_name"`
	Email             string
	ProfileImageURL   *string `gorm:"column:profile_image_url"`
	TwitchAccessToken *string `gorm:"-"`
	TwitchUsername    *string `gorm:"-"`
	TwitchDisplayName *string `gorm:"-"`
	CreatedAt         string  `gorm:"column:created_at"`
	UpdatedAt         string  `gorm:"column:updated_at"`
}

func main() {
	// Load environment variables
	if err := godotenv.Load(".env.local"); err != nil {
		log.Printf("Warning: .env.local file not found, using environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgresql://postgres:localdev123@localhost:5432/creatorsync_local"
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("🌱 Seeding LOCAL DEVELOPMENT database with TEST DATA...")
	fmt.Println("⚠️  This is TEST DATA ONLY - clearly marked for local development")

	// Create test users with obvious test data markers
	testUsers := []TestUser{
		{
			ID:              "test-user-1",
			ClerkID:         "user_test_dev_alice_123456789",
			TwitchUserID:    stringPtr("test_twitch_alice_123"),
			Username:        "test_alice_dev",
			DisplayName:     stringPtr("🧪 TEST Alice (DEV)"),
			Email:           "test.alice.dev@example.com",
			ProfileImageURL: nil,
			CreatedAt:       "2024-01-01T00:00:00Z",
			UpdatedAt:       "2024-01-01T00:00:00Z",
		},
		{
			ID:              "test-user-2",
			ClerkID:         "user_test_dev_bob_987654321",
			TwitchUserID:    stringPtr("test_twitch_bob_456"),
			Username:        "test_bob_dev",
			DisplayName:     stringPtr("🧪 TEST Bob (DEV)"),
			Email:           "test.bob.dev@example.com",
			ProfileImageURL: nil,
			CreatedAt:       "2024-01-01T00:00:00Z",
			UpdatedAt:       "2024-01-01T00:00:00Z",
		},
		{
			ID:              "test-user-3",
			ClerkID:         "user_test_dev_charlie_111222333",
			TwitchUserID:    nil,
			Username:        "test_charlie_dev",
			DisplayName:     stringPtr("🧪 TEST Charlie (DEV)"),
			Email:           "test.charlie.dev@example.com",
			ProfileImageURL: nil,
			CreatedAt:       "2024-01-01T00:00:00Z",
			UpdatedAt:       "2024-01-01T00:00:00Z",
		},
	}

	// Insert test users
	for _, user := range testUsers {
		result := db.Table("users").Create(&user)
		if result.Error != nil {
			log.Printf("Error creating test user %s: %v", user.Username, result.Error)
		} else {
			fmt.Printf("✅ Created test user: %s (%s)\n", user.Username, user.Email)
		}
	}

	fmt.Println("")
	fmt.Println("🎉 Test data seeding completed!")
	fmt.Println("📝 Note: All users have 'test_' prefixes and '🧪 TEST' markers")
	fmt.Println("🔒 This data is completely isolated from staging/production")
	fmt.Println("")
	fmt.Println("To use these test users:")
	fmt.Println("1. Use your Clerk development environment")
	fmt.Println("2. Create test users in Clerk with matching ClerkIDs")
	fmt.Println("3. Or modify the ClerkIDs to match your Clerk development users")
}

func stringPtr(s string) *string {
	return &s
}
