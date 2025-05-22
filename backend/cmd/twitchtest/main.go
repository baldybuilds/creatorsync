package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/baldybuilds/creatorsync/internal/clerk"
	"github.com/baldybuilds/creatorsync/internal/twitch"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file from project root
	projectRoot, err := findProjectRoot()
	if err != nil {
		log.Printf("Warning: Could not find project root: %v", err)
	} else {
		envPath := filepath.Join(projectRoot, ".env")
		err = godotenv.Load(envPath)
		if err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		} else {
			log.Printf("Successfully loaded environment variables from %s", envPath)
		}
	}

	// Check if user ID was provided
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run cmd/twitchtest/main.go <clerk_user_id>")
	}

	// Get the user ID from command line
	userID := os.Args[1]
	fmt.Printf("Testing with Clerk user ID: %s\n", userID)

	// Initialize Clerk
	if err := clerk.Initialize(); err != nil {
		log.Fatalf("Failed to initialize Clerk: %v", err)
	}
	fmt.Printf("Successfully initialized Clerk\n")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get the OAuth token for this user
	token, err := clerk.GetOAuthToken(ctx, userID, "oauth_twitch")
	if err != nil {
		log.Fatalf("Failed to get OAuth token: %v", err)
	}
	fmt.Printf("Successfully retrieved Twitch OAuth token\n")

	// Get the user to find their Twitch ID
	user, err := clerk.GetUserByID(ctx, userID)
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}
	// Safely print user info with nil checks
	firstName := "<no first name>"
	if user.FirstName != nil {
		firstName = *user.FirstName
	}

	lastName := "<no last name>"
	if user.LastName != nil {
		lastName = *user.LastName
	}

	fmt.Printf("Successfully retrieved user: %s %s (%s)\n", firstName, lastName, user.ID)

	// Find the Twitch user ID
	var twitchUserID string
	for _, account := range user.ExternalAccounts {
		if account.Provider == "oauth_twitch" {
			twitchUserID = account.ProviderUserID

			// Safely print username with nil check
			username := "<no username>"
			if account.Username != nil {
				username = *account.Username
			}

			fmt.Printf("Found Twitch account: %s (ID: %s)\n", username, twitchUserID)
			break
		}
	}

	if twitchUserID == "" {
		log.Fatalf("User does not have a connected Twitch account. Please connect Twitch in your application first.")
	}

	// Initialize the Twitch client
	twitchClientID := os.Getenv("TWITCH_CLIENT_ID")
	twitchClientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	if twitchClientID == "" {
		log.Fatalf("TWITCH_CLIENT_ID environment variable not set. This is required to run the twitchtest tool.")
	}
	if twitchClientSecret == "" {
		log.Fatalf("TWITCH_CLIENT_SECRET environment variable not set. This is required to run the twitchtest tool.")
	}
	twitchClient, err := twitch.NewClient(twitchClientID, twitchClientSecret)
	if err != nil {
		log.Fatalf("Failed to create Twitch client: %v", err)
	}
	fmt.Printf("Successfully initialized Twitch client\n")

	// Validate the token
	valid, err := twitchClient.ValidateToken(ctx, token)
	if err != nil {
		log.Fatalf("Failed to validate token: %v", err)
	}

	if !valid {
		log.Fatalf("Token is not valid. You may need to reconnect the Twitch account.")
	}
	fmt.Printf("Token validation successful\n")

	// Get channel info
	channelInfo, err := twitchClient.GetChannelInfo(ctx, token, twitchUserID)
	if err != nil {
		log.Fatalf("Failed to get channel info: %v", err)
	}

	// Print the channel info
	fmt.Printf("\n=== Twitch Channel Info ===\n")
	fmt.Printf("Broadcaster ID: %s\n", channelInfo.BroadcasterID)
	fmt.Printf("Broadcaster Name: %s\n", channelInfo.BroadcasterName)
	fmt.Printf("Game: %s\n", channelInfo.GameName)
	fmt.Printf("Title: %s\n", channelInfo.Title)
	fmt.Printf("Language: %s\n", channelInfo.Language)
	fmt.Printf("Is Mature: %v\n", channelInfo.IsMature)

	// Print environment variables for testing
	fmt.Printf("\n=== For Testing ===\n")
	fmt.Printf("To run tests with this user, set:\n")
	fmt.Printf("export TEST_CLERK_USER_ID=\"%s\"\n", userID)
	fmt.Printf("\nYou can now run tests with:\n")
	fmt.Printf("go test -v ./internal/twitch\n")
}

// findProjectRoot attempts to find the root directory of the project
// by looking for common project files like go.mod
func findProjectRoot() (string, error) {
	// Start with the current working directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up the directory tree looking for go.mod
	for {
		// Check if go.mod exists in this directory
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// We've reached the root of the filesystem without finding go.mod
			return "", fmt.Errorf("could not find project root (no go.mod found)")
		}
		dir = parent
	}
}
