package twitch

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/baldybuilds/creatorsync/internal/clerk"
	"github.com/baldybuilds/creatorsync/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/nicklaw5/helix/v2"
)

type TwitchOAuthHandlers struct {
	oauthConfig *OAuthConfig
	db          database.Service
}

func NewTwitchOAuthHandlers(db database.Service) (*TwitchOAuthHandlers, error) {
	oauthConfig, err := NewOAuthConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth config: %w", err)
	}

	return &TwitchOAuthHandlers{
		oauthConfig: oauthConfig,
		db:          db,
	}, nil
}

func generateState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// InitiateHandler creates a secure OAuth session and returns the auth URL
func (h *TwitchOAuthHandlers) InitiateHandler(c *fiber.Ctx) error {
	user, err := clerk.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Create secure session
	sessionStore := GetSessionStore()
	state, err := sessionStore.CreateSession(user.ID)
	if err != nil {
		log.Printf("❌ Failed to create OAuth session for user %s: %v", user.ID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create OAuth session",
		})
	}

	authURL := h.oauthConfig.GetAuthURL(state)
	log.Printf("🔗 Created OAuth session for user %s with state %s", user.ID, state)

	// Get scopes to show user what they're authorizing
	scopes := h.oauthConfig.GetScopes()
	scopeDescriptions := map[string]string{
		"user:read:email":            "Access your email address",
		"user:read:broadcast":        "View your stream key and preferences",
		"channel:read:subscriptions": "View your subscriber list and subscriber information",
		"clips:edit":                 "Create and edit clips from your streams",
		"channel:read:redemptions":   "View channel point redemptions",
		"moderation:read":            "View your moderation data",
		"analytics:read:extensions":  "View analytics data for your extensions",
		"analytics:read:games":       "View analytics data for games",
	}

	requestedPermissions := []map[string]string{}
	for _, scope := range scopes {
		description := scopeDescriptions[scope]
		if description == "" {
			description = "Access " + scope + " data"
		}
		requestedPermissions = append(requestedPermissions, map[string]string{
			"scope":       scope,
			"description": description,
		})
	}

	return c.JSON(fiber.Map{
		"oauth_url":   authURL,
		"state":       state,
		"permissions": requestedPermissions,
		"scope_count": len(scopes),
	})
}

// ConnectHandler is kept for backward compatibility but redirects to new flow
func (h *TwitchOAuthHandlers) ConnectHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "Use /auth/twitch/initiate endpoint instead",
	})
}

func (h *TwitchOAuthHandlers) CallbackHandler(c *fiber.Ctx) error {
	frontendURL := os.Getenv("APP_FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	state := c.Query("state")
	code := c.Query("code")
	errorParam := c.Query("error")

	if errorParam != "" {
		log.Printf("❌ Twitch OAuth error: %s", errorParam)
		return c.Redirect(frontendURL+"/dashboard?twitch_error=oauth_denied", fiber.StatusTemporaryRedirect)
	}

	if state == "" || code == "" {
		log.Printf("❌ Missing state or code in callback")
		return c.Redirect(frontendURL+"/dashboard?twitch_error=invalid_callback", fiber.StatusTemporaryRedirect)
	}

	// Get user session from secure session store
	sessionStore := GetSessionStore()
	session, exists := sessionStore.GetSession(state)
	if !exists {
		log.Printf("❌ Invalid or expired OAuth session for state: %s", state)
		return c.Redirect(frontendURL+"/dashboard?twitch_error=csrf_failed", fiber.StatusTemporaryRedirect)
	}

	// Clean up the session
	sessionStore.DeleteSession(state)

	userID := session.UserID
	log.Printf("✅ Retrieved user %s from OAuth session", userID)

	token, err := h.oauthConfig.ExchangeCode(c.Context(), code)
	if err != nil {
		log.Printf("❌ Failed to exchange code for user %s: %v", userID, err)
		return c.Redirect(frontendURL+"/dashboard?twitch_error=token_exchange_failed", fiber.StatusTemporaryRedirect)
	}

	helixClient, err := helix.NewClient(&helix.Options{
		ClientID:        os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret:    os.Getenv("TWITCH_CLIENT_SECRET"),
		UserAccessToken: token.AccessToken,
	})
	if err != nil {
		log.Printf("❌ Failed to create Helix client for user %s: %v", userID, err)
		return c.Redirect(frontendURL+"/dashboard?twitch_error=client_init_failed", fiber.StatusTemporaryRedirect)
	}

	usersResp, err := helixClient.GetUsers(&helix.UsersParams{})
	if err != nil || len(usersResp.Data.Users) == 0 {
		log.Printf("❌ Failed to get Twitch user info for user %s: %v", userID, err)
		return c.Redirect(frontendURL+"/dashboard?twitch_error=user_info_failed", fiber.StatusTemporaryRedirect)
	}

	twitchUser := usersResp.Data.Users[0]
	twitchUserIDFromTwitch := twitchUser.ID

	if err := h.oauthConfig.StoreTokens(c.Context(), h.db, userID, twitchUserIDFromTwitch, token); err != nil {
		log.Printf("❌ Failed to store tokens for user %s: %v", userID, err)
		return c.Redirect(frontendURL+"/dashboard?twitch_error=token_storage_failed", fiber.StatusTemporaryRedirect)
	}

	log.Printf("✅ Successfully connected Twitch for user %s (Twitch ID: %s)", userID, twitchUserIDFromTwitch)

	// Invalidate analytics cache to ensure fresh connection status
	// Note: This is a bit of a hack, but necessary to ensure immediate connection detection
	go func() {
		// Give a small delay to ensure token storage completes
		time.Sleep(100 * time.Millisecond)

		// Clear any cached connection status using a simple database query
		// This avoids importing the analytics package and creating circular dependencies
		if _, err := h.db.GetDB().Exec(`DELETE FROM cache_entries WHERE user_id = $1 AND (cache_key = 'connection_status' OR cache_key LIKE '%enhanced%' OR cache_key LIKE '%overview%')`, userID); err != nil {
			log.Printf("⚠️ Failed to clear analytics cache for user %s: %v", userID, err)
		} else {
			log.Printf("🧹 Cleared analytics cache for user %s", userID)
		}
	}()

	// Trigger initial data collection for the newly connected account
	go func() {
		log.Printf("🚀 Triggering initial data collection for user %s", userID)
		// This will be handled by the background collection manager if available
		// For now, we'll just log it and let the frontend request data when needed
	}()

	return c.Redirect(frontendURL+"/dashboard?twitch_connected=true", fiber.StatusTemporaryRedirect)
}

func (h *TwitchOAuthHandlers) StatusHandler(c *fiber.Ctx) error {
	user, err := clerk.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	_, err = h.oauthConfig.GetStoredTokens(c.Context(), h.db, user.ID)
	if err != nil {
		return c.JSON(fiber.Map{
			"connected": false,
		})
	}

	return c.JSON(fiber.Map{
		"connected": true,
	})
}

// DisconnectHandler removes Twitch OAuth tokens and clears all related data
func (h *TwitchOAuthHandlers) DisconnectHandler(c *fiber.Ctx) error {
	user, err := clerk.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	userID := user.ID
	log.Printf("🔌 Disconnecting Twitch for user %s", userID)

	// Delete OAuth tokens
	err = h.oauthConfig.DeleteStoredTokens(c.Context(), h.db, userID)
	if err != nil {
		log.Printf("❌ Failed to delete tokens for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to disconnect Twitch account",
		})
	}

	// Clear all analytics cache and stored data
	go func() {
		log.Printf("🧹 Clearing all data for disconnected user %s", userID)

		// Clear cache entries
		if _, err := h.db.GetDB().Exec(`DELETE FROM cache_entries WHERE user_id = $1`, userID); err != nil {
			log.Printf("⚠️ Failed to clear cache for user %s: %v", userID, err)
		}

		// Clear video analytics
		if _, err := h.db.GetDB().Exec(`DELETE FROM video_analytics WHERE user_id = $1`, userID); err != nil {
			log.Printf("⚠️ Failed to clear video analytics for user %s: %v", userID, err)
		}

		// Clear channel analytics
		if _, err := h.db.GetDB().Exec(`DELETE FROM channel_analytics WHERE user_id = $1`, userID); err != nil {
			log.Printf("⚠️ Failed to clear channel analytics for user %s: %v", userID, err)
		}

		// Clear stream sessions
		if _, err := h.db.GetDB().Exec(`DELETE FROM stream_sessions WHERE user_id = $1`, userID); err != nil {
			log.Printf("⚠️ Failed to clear stream sessions for user %s: %v", userID, err)
		}

		log.Printf("✅ Data cleanup completed for user %s", userID)
	}()

	log.Printf("✅ Successfully disconnected Twitch for user %s", userID)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Twitch account disconnected successfully",
	})
}
