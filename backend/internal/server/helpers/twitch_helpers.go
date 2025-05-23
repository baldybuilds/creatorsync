package helpers

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/baldybuilds/creatorsync/internal/analytics"
	"github.com/baldybuilds/creatorsync/internal/clerk"
	"github.com/baldybuilds/creatorsync/internal/database"
	"github.com/baldybuilds/creatorsync/internal/twitch"
	clerkSDK "github.com/clerk/clerk-sdk-go/v2"
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

// ensureUserExistsInDatabase creates or updates a user record in our database
func ensureUserExistsInDatabase(ctx context.Context, db database.Service, userID string) error {
	// Check if user already exists in our database
	analyticsRepo := analytics.NewRepository(db)
	existingUser, err := analyticsRepo.GetUserByClerkID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		return nil // User already exists
	}

	// User doesn't exist, let's create them from Clerk
	clerkUser, err := clerk.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user from Clerk: %w", err)
	}

	// Initialize user with basic info from Clerk
	user := &analytics.User{
		ID:          userID,
		ClerkUserID: userID,
	}

	// Safely set email if available
	if len(clerkUser.EmailAddresses) > 0 {
		user.Email = clerkUser.EmailAddresses[0].EmailAddress
	}

	// Set name fields safely
	if clerkUser.FirstName != nil {
		user.DisplayName = *clerkUser.FirstName
	}
	if clerkUser.LastName != nil && *clerkUser.LastName != "" {
		if user.DisplayName != "" {
			user.DisplayName += " " + *clerkUser.LastName
		} else {
			user.DisplayName = *clerkUser.LastName
		}
	}

	// Try to get Twitch info if available
	for _, account := range clerkUser.ExternalAccounts {
		if account.Provider == "oauth_twitch" {
			user.TwitchUserID = account.ProviderUserID
			if account.Username != nil {
				user.Username = *account.Username
			}

			// Try to get additional Twitch info if we have OAuth token
			if token, tokenErr := clerk.GetOAuthToken(ctx, userID, "oauth_twitch"); tokenErr == nil {
				// Initialize Twitch client
				twitchClientID := os.Getenv("TWITCH_CLIENT_ID")
				twitchClientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
				if twitchClientID != "" && twitchClientSecret != "" {
					if twitchClient, clientErr := twitch.NewClient(twitchClientID, twitchClientSecret); clientErr == nil {
						if userInfo, infoErr := twitchClient.GetUserInfo(token); infoErr == nil {
							user.Username = userInfo.Login
							user.DisplayName = userInfo.DisplayName
							user.ProfileImageURL = userInfo.ProfileImageURL
							if userInfo.Email != "" {
								user.Email = userInfo.Email
							}
						}
					}
				}
			}
			break
		}
	}

	// Create user record in database
	if err := analyticsRepo.CreateOrUpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to create user record: %w", err)
	}

	log.Printf("✅ Created user record for %s (%s)", user.DisplayName, userID)
	return nil
}

// GetTwitchRequestContext consolidates the common logic for fetching user details,
// Twitch token, Twitch user ID, and initializing the Twitch client.
// It returns the context or an error that the calling handler should use to respond to the client.
func GetTwitchRequestContext(c *fiber.Ctx) (*TwitchRequestContext, error) {
	user, clerkErr := clerk.GetUserFromContext(c)
	if clerkErr != nil {
		return nil, fmt.Errorf("user not authenticated")
	}

	// Get database service from fiber context
	db, ok := c.UserContext().Value("db").(database.Service)
	if !ok {
		// Try to get it from fiber locals
		if dbLocal := c.Locals("db"); dbLocal != nil {
			if dbService, ok := dbLocal.(database.Service); ok {
				db = dbService
			}
		}
	}

	// If we have database access, ensure user exists before proceeding
	if db != nil {
		if err := ensureUserExistsInDatabase(c.Context(), db, user.ID); err != nil {
			log.Printf("⚠️ Failed to sync user %s to database: %v", user.ID, err)
			// Don't fail the request, just log the warning and continue
		}
	}

	clerkUser, clerkErr := clerk.GetUserByID(c.Context(), user.ID)
	if clerkErr != nil {
		return nil, fmt.Errorf("failed to get user from Clerk: %v", clerkErr)
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
	twitchClientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	if twitchClientID == "" || twitchClientSecret == "" {
		log.Println("Error: TWITCH_CLIENT_ID or TWITCH_CLIENT_SECRET environment variable not set.")
		return nil, fmt.Errorf("twitch client configuration error")
	}

	initializedClient, clientErr := twitch.NewClient(twitchClientID, twitchClientSecret)
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
	twitchClientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	if twitchClientID == "" || twitchClientSecret == "" {
		return nil, fmt.Errorf("twitch client configuration error")
	}

	return twitch.NewClient(twitchClientID, twitchClientSecret)
}
