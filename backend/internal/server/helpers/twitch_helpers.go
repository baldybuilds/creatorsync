package helpers

import (
	"context"
	"fmt"
	"log"
	"os"

	clerkSDK "github.com/clerk/clerk-sdk-go/v2"
	"github.com/baldybuilds/creatorsync/internal/clerk"
	"github.com/baldybuilds/creatorsync/internal/twitch"
	"github.com/gofiber/fiber/v2"
)

// TwitchRequestContext holds all the necessary Twitch-related information for a handler.
type TwitchRequestContext struct {
	UserID      string
	AccessToken string
	Client      *twitch.Client
	ClerkUser   *clerkSDK.User // Full Clerk user from SDK
	LocalUser   *clerk.User    // Local user representation
}

// GetTwitchRequestContext consolidates the common logic for fetching user details,
// Twitch token, Twitch user ID, and initializing the Twitch client.
// It returns the context or an error that the calling handler should use to respond to the client.
func GetTwitchRequestContext(c *fiber.Ctx) (*TwitchRequestContext, error) {
	user, clerkErr := clerk.GetUserFromContext(c)
	if clerkErr != nil {
		return nil, fmt.Errorf("user not authenticated")
	}

	clerkUser, clerkErr := clerk.GetUserByID(c.Context(), user.ID)
	if clerkErr != nil {
		return nil, fmt.Errorf("failed to get user profile: %v", clerkErr)
	}

	var foundTwitchUserID string
	for _, account := range clerkUser.ExternalAccounts {
		if account.Provider == "oauth_twitch" {
			foundTwitchUserID = account.ProviderUserID
			break
		}
	}

	if foundTwitchUserID == "" {
		return nil, fmt.Errorf("twitch account not connected")
	}

	token, clerkErr := clerk.GetOAuthToken(c.Context(), user.ID, "oauth_twitch")
	if clerkErr != nil {
		return nil, fmt.Errorf("failed to get Twitch token: %v", clerkErr)
	}

	twitchClientID := os.Getenv("TWITCH_CLIENT_ID")
	if twitchClientID == "" {
		log.Println("Error: TWITCH_CLIENT_ID environment variable not set.")
		return nil, fmt.Errorf("twitch client configuration error")
	}

	initializedClient, clientErr := twitch.NewClient(twitchClientID)
	if clientErr != nil {
		return nil, fmt.Errorf("failed to initialize Twitch client: %v", clientErr)
	}

	localUser, _ := clerk.GetUserFromContext(c)

	return &TwitchRequestContext{
		UserID:      foundTwitchUserID,
		AccessToken: token,
		Client:      initializedClient,
		ClerkUser:   clerkUser,
		LocalUser:   localUser,
	}, nil
}

// HandleTwitchError formats a Twitch-related error as a Fiber response
func HandleTwitchError(c *fiber.Ctx, err error) error {
	// Determine appropriate status code based on error message
	statusCode := fiber.StatusInternalServerError
	
	if err.Error() == "user not authenticated" {
		statusCode = fiber.StatusUnauthorized
	} else if err.Error() == "twitch account not connected" {
		statusCode = fiber.StatusBadRequest
	}
	
	return c.Status(statusCode).JSON(fiber.Map{
		"error": err.Error(),
	})
}

// GetTwitchClient is a convenience function that returns just the initialized Twitch client
// for handlers that don't need the full context
func GetTwitchClient(ctx context.Context) (*twitch.Client, error) {
	twitchClientID := os.Getenv("TWITCH_CLIENT_ID")
	if twitchClientID == "" {
		return nil, fmt.Errorf("twitch client configuration error")
	}

	return twitch.NewClient(twitchClientID)
}
