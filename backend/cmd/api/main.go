package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/baldybuilds/creatorsync/internal/server"

	_ "github.com/joho/godotenv/autoload"
)

func gracefulShutdown(fiberServer *server.FiberServer, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// Stop background services first
	if err := fiberServer.StopBackgroundServices(); err != nil {
		log.Printf("Error stopping background services: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := fiberServer.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
	done <- true
}

func main() {
	// Log environment configuration on startup
	log.Printf("🚀 Starting CreatorSync API Server")
	log.Printf("📍 Environment: %s", os.Getenv("APP_ENV"))
	log.Printf("🔌 Port: %s", os.Getenv("PORT"))

	// Log database configuration (without sensitive info)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		log.Printf("🗄️ Database: Using DATABASE_URL (configured)")
	} else {
		log.Printf("🗄️ Database: Using individual environment variables")
		log.Printf("   Host: %s", os.Getenv("POSTGRES_DB_HOST"))
		log.Printf("   Port: %s", os.Getenv("POSTGRES_DB_PORT"))
		log.Printf("   Database: %s", os.Getenv("POSTGRES_DB_DATABASE"))
		log.Printf("   Username: %s", os.Getenv("POSTGRES_DB_USERNAME"))
	}

	// Log authentication configuration
	if clerkKey := os.Getenv("CLERK_SECRET_KEY"); clerkKey != "" {
		keyPrefix := clerkKey[:10] + "..."
		log.Printf("🔐 Clerk: Configured (%s)", keyPrefix)
	} else {
		log.Printf("⚠️ Clerk: NOT CONFIGURED")
	}

	// Log Twitch configuration
	if twitchClientID := os.Getenv("TWITCH_CLIENT_ID"); twitchClientID != "" {
		log.Printf("📺 Twitch: Configured")
	} else {
		log.Printf("⚠️ Twitch: NOT CONFIGURED")
	}

	server, err := server.New()
	if err != nil {
		log.Fatalf("❌ Failed to initialize server: %v", err)
	}

	// Start background services (analytics collection)
	if err := server.StartBackgroundServices(); err != nil {
		log.Fatalf("❌ Failed to start background services: %v", err)
	}

	server.RegisterFiberRoutes()

	done := make(chan bool, 1)

	go func() {
		port, _ := strconv.Atoi(os.Getenv("PORT"))
		log.Printf("🌐 Server listening on port %d", port)
		err := server.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	go gracefulShutdown(server, done)

	<-done
	log.Println("✅ Graceful shutdown complete.")
}
