package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/baldybuilds/creatorsync/internal/clerk"
	"github.com/baldybuilds/creatorsync/internal/database"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	log.Printf("🔍 CreatorSync Environment Debug Tool")
	log.Printf("=====================================")

	// Check environment variables
	checkEnvironment()

	// Test database connection
	testDatabase()

	// Test Clerk authentication
	testClerk()

	log.Printf("✅ Debug check complete")
}

func checkEnvironment() {
	log.Printf("\n📍 Environment Configuration:")

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "not set"
	}
	log.Printf("   APP_ENV: %s", env)

	port := os.Getenv("PORT")
	if port == "" {
		port = "not set"
	}
	log.Printf("   PORT: %s", port)

	// Database
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		log.Printf("   DATABASE_URL: configured (%d chars)", len(databaseURL))
	} else {
		log.Printf("   DATABASE_URL: not set")
		log.Printf("   POSTGRES_DB_HOST: %s", getEnvOrDefault("POSTGRES_DB_HOST", "not set"))
		log.Printf("   POSTGRES_DB_PORT: %s", getEnvOrDefault("POSTGRES_DB_PORT", "not set"))
		log.Printf("   POSTGRES_DB_DATABASE: %s", getEnvOrDefault("POSTGRES_DB_DATABASE", "not set"))
		log.Printf("   POSTGRES_DB_USERNAME: %s", getEnvOrDefault("POSTGRES_DB_USERNAME", "not set"))
	}

	// Clerk
	if clerkKey := os.Getenv("CLERK_SECRET_KEY"); clerkKey != "" {
		log.Printf("   CLERK_SECRET_KEY: configured (%d chars)", len(clerkKey))
	} else {
		log.Printf("   CLERK_SECRET_KEY: not set")
	}

	// Twitch
	if twitchClientID := os.Getenv("TWITCH_CLIENT_ID"); twitchClientID != "" {
		log.Printf("   TWITCH_CLIENT_ID: configured (%d chars)", len(twitchClientID))
	} else {
		log.Printf("   TWITCH_CLIENT_ID: not set")
	}

	if twitchSecret := os.Getenv("TWITCH_CLIENT_SECRET"); twitchSecret != "" {
		log.Printf("   TWITCH_CLIENT_SECRET: configured (%d chars)", len(twitchSecret))
	} else {
		log.Printf("   TWITCH_CLIENT_SECRET: not set")
	}
}

func testDatabase() {
	log.Printf("\n🗄️ Database Connection Test:")

	db := database.New()
	defer db.Close()

	health := db.Health()
	if health["status"] == "up" {
		log.Printf("   ✅ Database connection: HEALTHY")
		log.Printf("   📊 Open connections: %s", health["open_connections"])
		log.Printf("   💤 Idle connections: %s", health["idle"])
	} else {
		log.Printf("   ❌ Database connection: FAILED")
		if error, exists := health["error"]; exists {
			log.Printf("   Error: %v", error)
		}
	}

	// Test a simple query
	sqlDB := db.GetDB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result int
	err := sqlDB.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		log.Printf("   ❌ Simple query test: FAILED - %v", err)
	} else {
		log.Printf("   ✅ Simple query test: PASSED")
	}

	// Check if tables exist
	var tableCount int
	err = sqlDB.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name IN ('users', 'channel_analytics', 'video_analytics')
	`).Scan(&tableCount)

	if err != nil {
		log.Printf("   ❌ Table check: FAILED - %v", err)
	} else {
		log.Printf("   📋 Required tables found: %d/3", tableCount)
		if tableCount < 3 {
			log.Printf("   ⚠️ Missing tables - migrations may need to be run")
		}
	}
}

func testClerk() {
	log.Printf("\n🔐 Clerk Authentication Test:")

	err := clerk.Initialize()
	if err != nil {
		log.Printf("   ❌ Clerk initialization: FAILED - %v", err)
		return
	}

	log.Printf("   ✅ Clerk initialization: SUCCESS")

	// We can't easily test user authentication without a token,
	// but we can verify the client is configured
	log.Printf("   ℹ️ Clerk client is properly configured")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
