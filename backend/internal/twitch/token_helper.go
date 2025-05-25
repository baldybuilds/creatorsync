package twitch

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/baldybuilds/creatorsync/internal/clerk"
	"github.com/baldybuilds/creatorsync/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/nicklaw5/helix/v2"
	"golang.org/x/oauth2"
)

type TwitchTokenHelper struct {
	oauthConfig *OAuthConfig
	db          database.Service
}

func NewTwitchTokenHelper(db database.Service) (*TwitchTokenHelper, error) {
	oauthConfig, err := NewOAuthConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth config: %w", err)
	}

	return &TwitchTokenHelper{
		oauthConfig: oauthConfig,
		db:          db,
	}, nil
}

func (h *TwitchTokenHelper) GetValidTokenForUser(ctx context.Context, clerkUserID string) (*oauth2.Token, error) {
	return h.oauthConfig.GetValidToken(ctx, h.db, clerkUserID)
}

func (h *TwitchTokenHelper) GetTwitchUserID(ctx context.Context, clerkUserID string) (string, error) {
	return h.oauthConfig.GetTwitchUserID(ctx, h.db, clerkUserID)
}

func (h *TwitchTokenHelper) GetHelixClientForUser(ctx context.Context, clerkUserID string) (*helix.Client, error) {
	token, err := h.GetValidTokenForUser(ctx, clerkUserID)
	if err != nil {
		return nil, err
	}

	client, err := helix.NewClient(&helix.Options{
		ClientID:        os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret:    os.Getenv("TWITCH_CLIENT_SECRET"),
		UserAccessToken: token.AccessToken,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Helix client: %w", err)
	}

	return client, nil
}

type TwitchRequestContext struct {
	UserID       string
	TwitchUserID string
	Token        *oauth2.Token
	HelixClient  *helix.Client
	ClerkUser    *clerk.User
}

func (h *TwitchTokenHelper) GetTwitchRequestContext(c *fiber.Ctx) (*TwitchRequestContext, error) {
	user, err := clerk.GetUserFromContext(c)
	if err != nil {
		return nil, fmt.Errorf("user not authenticated")
	}

	token, err := h.GetValidTokenForUser(c.Context(), user.ID)
	if err != nil {
		if err.Error() == "twitch not connected" {
			return nil, fmt.Errorf("twitch account not connected")
		}
		if err.Error() == "twitch re-authentication required" {
			return nil, fmt.Errorf("twitch re-authentication required")
		}
		return nil, fmt.Errorf("failed to get valid token: %w", err)
	}

	twitchUserID, err := h.GetTwitchUserID(c.Context(), user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get Twitch user ID: %w", err)
	}

	helixClient, err := h.GetHelixClientForUser(c.Context(), user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create Helix client: %w", err)
	}

	return &TwitchRequestContext{
		UserID:       user.ID,
		TwitchUserID: twitchUserID,
		Token:        token,
		HelixClient:  helixClient,
		ClerkUser:    user,
	}, nil
}

func HandleTwitchError(c *fiber.Ctx, err error) error {
	statusCode := fiber.StatusInternalServerError
	errorMessage := err.Error()

	switch errorMessage {
	case "user not authenticated":
		statusCode = fiber.StatusUnauthorized
	case "twitch account not connected":
		statusCode = fiber.StatusBadRequest
		return c.Status(statusCode).JSON(fiber.Map{
			"error":              errorMessage,
			"reconnect_required": false,
			"connect_required":   true,
		})
	case "twitch re-authentication required":
		statusCode = fiber.StatusUnauthorized
		return c.Status(statusCode).JSON(fiber.Map{
			"error":              errorMessage,
			"reconnect_required": true,
			"connect_required":   false,
		})
	}

	log.Printf("❌ Twitch API error: %v", err)
	return c.Status(statusCode).JSON(fiber.Map{
		"error": errorMessage,
	})
}
