// How to run these tests with real data:
//
// Since we're using Clerk for authentication, we only need to set up a few things:
//
// 1. Make sure your .env file has the necessary Clerk credentials:
//    CLERK_SECRET_KEY=your_clerk_secret_key
//
// 2. To test with a real user's Twitch connection:
//    a. Log in to your application using the development Clerk instance
//    b. Connect your Twitch account through the Clerk OAuth flow
//    c. Get your Clerk user ID from the Clerk dashboard
//    d. Run a simple script to extract the token and broadcaster ID (see below)
//
// 3. Run the tests with:
//    go test -v ./internal/tests/twitch
//
// 4. To run a specific test:
//    go test -v ./internal/tests/twitch -run TestGetChannelInfo

package twitch_test

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/baldybuilds/creatorsync/internal/clerk"
	"github.com/baldybuilds/creatorsync/internal/twitch"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {

	projectRoot := findProjectRoot()
	if projectRoot != "" {
		envFile := filepath.Join(projectRoot, ".env")
		err := godotenv.Load(envFile)
		if err != nil {
			log.Printf("Warning: Could not load .env file from %s: %v", envFile, err)
		} else {
			log.Printf("Loaded environment variables from %s", envFile)
		}
	}
}

func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory: %v", err)
		return ""
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	log.Println("Could not find project root with .env file")
	return ""
}

// TestNewClient tests the creation of a new Twitch client
func TestNewClient(t *testing.T) {
	twitchClientID := os.Getenv("TWITCH_CLIENT_ID")
	if twitchClientID == "" {
		t.Log("TWITCH_CLIENT_ID not set, using placeholder for NewClient test")
		twitchClientID = "test_placeholder_client_id" // Or skip, depending on test needs
	}
	client, err := twitch.NewClient(twitchClientID)
	assert.NoError(t, err, "NewClient should not return an error")
	assert.NotNil(t, client, "Client should not be nil")
	t.Log("Client created successfully, will use Clerk for Twitch Auth")
}

// TestValidateToken tests the token validation functionality
func TestValidateToken(t *testing.T) {
	twitchClientID := os.Getenv("TWITCH_CLIENT_ID")
	if twitchClientID == "" {
		t.Log("TWITCH_CLIENT_ID not set, using placeholder for TestValidateToken")
		twitchClientID = "test_placeholder_client_id"
	}
	client, err := twitch.NewClient(twitchClientID)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	valid, err := client.ValidateToken(ctx, "invalid_token")
	assert.NoError(t, err, "ValidateToken should not return an error for invalid tokens")
	assert.False(t, valid, "Invalid token should not validate")

	clerkUserID := os.Getenv("TEST_CLERK_USER_ID")
	if clerkUserID != "" {
		t.Logf("Found TEST_CLERK_USER_ID, attempting to get Twitch token for user %s", clerkUserID)
		err := clerk.Initialize()
		if err != nil {
			t.Logf("Failed to initialize Clerk: %v", err)
			t.Skip("Skipping real token test due to Clerk initialization failure")
			return
		}

		token, err := clerk.GetOAuthToken(ctx, clerkUserID, "oauth_twitch")
		if err != nil {
			t.Logf("Failed to get Twitch  token: %v", err)
			t.Skip("Skipping real token test due to Twitch  token retrieval failure")
			return
		}

		valid, err = client.ValidateToken(ctx, token)
		assert.NoError(t, err)
		assert.True(t, valid, "Valid token from Clerk should validate successfully")
	}
}

// TestGetChannelInfo tests the channel info retrieval using Clerk
func TestGetChannelInfo(t *testing.T) {
	clerkUserID := os.Getenv("TEST_CLERK_USER_ID")
	if clerkUserID == "" {
		t.Skip("Skipping test: Set TEST_CLERK_USER_ID env var with a real Clerk user ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := clerk.Initialize()
	if err != nil {
		t.Logf("Failed to initialize Clerk: %v", err)
		t.Skip("Skipping test due to Clerk initialization failure")
		return
	}

	token, err := clerk.GetOAuthToken(ctx, clerkUserID, "oauth_twitch")
	if err != nil {
		t.Logf("Failed to get Twitch token: %v", err)
		t.Skip("Skipping test due to Twitch token retrieval failure")
		return
	}

	user, err := clerk.GetUserByID(ctx, clerkUserID)
	if err != nil {
		t.Logf("Failed to get user: %v", err)
		t.Skip("Skipping test due to user retrieval failure")
		return
	}

	var twitchUserID string
	for _, account := range user.ExternalAccounts {
		if account.Provider == "oauth_twitch" {
			twitchUserID = account.ProviderUserID
			break
		}
	}

	if twitchUserID == "" {
		t.Skip("User does not have a connected Twitch account")
	}

	twitchClientID := os.Getenv("TWITCH_CLIENT_ID")
	if twitchClientID == "" {
		t.Log("TWITCH_CLIENT_ID not set, using placeholder") // Generic log message
		twitchClientID = "test_placeholder_client_id"
	}
	client, err := twitch.NewClient(twitchClientID)
	assert.NoError(t, err)

	valid, err := client.ValidateToken(ctx, token)
	assert.NoError(t, err)
	assert.True(t, valid, "Token should be valid")

	channelInfo, err := client.GetChannelInfo(ctx, token, twitchUserID)
	assert.NoError(t, err)
	assert.NotNil(t, channelInfo)
	assert.Equal(t, twitchUserID, channelInfo.BroadcasterID)

	t.Logf("Successfully retrieved channel info for %s", channelInfo.BroadcasterName)
}

// TestGetUserVideos tests the video retrieval using Clerk
func TestGetUserVideos(t *testing.T) {
	clerkUserID := os.Getenv("TEST_CLERK_USER_ID")
	if clerkUserID == "" {
		t.Skip("Skipping test: Set TEST_CLERK_USER_ID env var with a real Clerk user ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := clerk.Initialize()
	if err != nil {
		t.Logf("Failed to initialize Clerk: %v", err)
		t.Skip("Skipping test due to Clerk initialization failure")
		return
	}

	token, err := clerk.GetOAuthToken(ctx, clerkUserID, "oauth_twitch")
	if err != nil {
		t.Logf("Failed to get Twitch  token: %v", err)
		t.Skip("Skipping test due to Twitch  token retrieval failure")
		return
	}

	user, err := clerk.GetUserByID(ctx, clerkUserID)
	if err != nil {
		t.Logf("Failed to get user: %v", err)
		t.Skip("Skipping test due to user retrieval failure")
		return
	}

	var twitchUserID string
	for _, account := range user.ExternalAccounts {
		if account.Provider == "oauth_twitch" {
			twitchUserID = account.ProviderUserID
			break
		}
	}

	if twitchUserID == "" {
		t.Skip("User does not have a connected Twitch account")
	}

	twitchClientID := os.Getenv("TWITCH_CLIENT_ID")
	if twitchClientID == "" {
		t.Log("TWITCH_CLIENT_ID not set, using placeholder") // Generic log message
		twitchClientID = "test_placeholder_client_id"
	}
	client, err := twitch.NewClient(twitchClientID)
	assert.NoError(t, err)

	valid, err := client.ValidateToken(ctx, token)
	assert.NoError(t, err)
	assert.True(t, valid, "Token should be valid")

	videos, cursor, err := client.GetUserVideos(ctx, token, twitchUserID, 5)
	assert.NoError(t, err)

	if len(videos) > 0 {
		t.Logf("Successfully retrieved %d videos for user %s", len(videos), twitchUserID)
		t.Logf("First video title: %s", videos[0].Title)
	} else {
		t.Logf("User %s has no videos", twitchUserID)
	}

	if cursor != "" {
		t.Logf("More videos available, pagination cursor: %s", cursor)
	}
}
