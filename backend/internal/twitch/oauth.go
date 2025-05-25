package twitch

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/baldybuilds/creatorsync/internal/database"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

type UserTwitchTokens struct {
	ClerkUserID           string    `db:"clerk_user_id"`
	TwitchUserID          string    `db:"twitch_user_id"`
	EncryptedAccessToken  string    `db:"encrypted_access_token"`
	EncryptedRefreshToken string    `db:"encrypted_refresh_token"`
	Scopes                string    `db:"scopes"`
	ExpiresAt             time.Time `db:"expires_at"`
	CreatedAt             time.Time `db:"created_at"`
	UpdatedAt             time.Time `db:"updated_at"`
}

type OAuthConfig struct {
	config        *oauth2.Config
	encryptionKey string
}

func NewOAuthConfig() (*OAuthConfig, error) {
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	redirectURIBase := os.Getenv("TWITCH_REDIRECT_URI_BASE")
	scopesString := os.Getenv("TWITCH_SCOPES")
	encryptionKey := os.Getenv("TWITCH_TOKEN_ENCRYPTION_KEY")

	if clientID == "" || clientSecret == "" || redirectURIBase == "" || encryptionKey == "" {
		return nil, fmt.Errorf("missing required Twitch OAuth environment variables")
	}

	scopes := []string{}
	if scopesString != "" {
		scopes = strings.Split(scopesString, ",")
		for i, scope := range scopes {
			scopes[i] = strings.TrimSpace(scope)
		}
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURIBase + "/auth/twitch/callback",
		Scopes:       scopes,
		Endpoint:     twitch.Endpoint,
	}

	return &OAuthConfig{
		config:        config,
		encryptionKey: encryptionKey,
	}, nil
}

func (oc *OAuthConfig) GetAuthURL(state string) string {
	return oc.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (oc *OAuthConfig) GetScopes() []string {
	return oc.config.Scopes
}

func (oc *OAuthConfig) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return oc.config.Exchange(ctx, code)
}

func (oc *OAuthConfig) TokenSource(ctx context.Context, token *oauth2.Token) oauth2.TokenSource {
	return oc.config.TokenSource(ctx, token)
}

func (oc *OAuthConfig) StoreTokens(ctx context.Context, db database.Service, clerkUserID, twitchUserID string, token *oauth2.Token) error {
	encryptedAccessToken, err := encryptToken(token.AccessToken, oc.encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt access token: %w", err)
	}

	encryptedRefreshToken, err := encryptToken(token.RefreshToken, oc.encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt refresh token: %w", err)
	}

	scopesString := strings.Join(oc.config.Scopes, ",")

	// Check if this user is switching to a different Twitch account
	var currentTwitchUserID string
	currentQuery := `SELECT twitch_user_id FROM user_twitch_tokens WHERE clerk_user_id = $1`
	err = db.GetDB().QueryRowContext(ctx, currentQuery, clerkUserID).Scan(&currentTwitchUserID)

	isAccountSwitch := false
	if err == nil && currentTwitchUserID != twitchUserID {
		isAccountSwitch = true
		log.Printf("🔄 Account switch detected for user %s: %s -> %s", clerkUserID, currentTwitchUserID, twitchUserID)
	}

	// First, delete any existing tokens for this Twitch user ID to avoid conflicts
	_, err = db.GetDB().ExecContext(ctx, `DELETE FROM user_twitch_tokens WHERE twitch_user_id = $1 AND clerk_user_id != $2`, twitchUserID, clerkUserID)
	if err != nil {
		return fmt.Errorf("failed to clean existing tokens: %w", err)
	}

	// If this is an account switch, clear all old analytics data
	if isAccountSwitch {
		log.Printf("🧹 Clearing old analytics data for account switch (user: %s)", clerkUserID)

		// Clear analytics data in a goroutine to avoid blocking the auth flow
		go func() {
			// Clear cache entries
			if _, err := db.GetDB().Exec(`DELETE FROM cache_entries WHERE user_id = $1`, clerkUserID); err != nil {
				log.Printf("⚠️ Failed to clear cache during account switch for user %s: %v", clerkUserID, err)
			}

			// Clear video analytics
			if _, err := db.GetDB().Exec(`DELETE FROM video_analytics WHERE user_id = $1`, clerkUserID); err != nil {
				log.Printf("⚠️ Failed to clear video analytics during account switch for user %s: %v", clerkUserID, err)
			}

			// Clear channel analytics
			if _, err := db.GetDB().Exec(`DELETE FROM channel_analytics WHERE user_id = $1`, clerkUserID); err != nil {
				log.Printf("⚠️ Failed to clear channel analytics during account switch for user %s: %v", clerkUserID, err)
			}

			// Clear stream sessions
			if _, err := db.GetDB().Exec(`DELETE FROM stream_sessions WHERE user_id = $1`, clerkUserID); err != nil {
				log.Printf("⚠️ Failed to clear stream sessions during account switch for user %s: %v", clerkUserID, err)
			}

			log.Printf("✅ Analytics data cleanup completed for account switch (user: %s)", clerkUserID)
		}()
	}

	query := `
		INSERT INTO user_twitch_tokens (
			clerk_user_id, twitch_user_id, encrypted_access_token, 
			encrypted_refresh_token, scopes, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (clerk_user_id) DO UPDATE SET
			twitch_user_id = EXCLUDED.twitch_user_id,
			encrypted_access_token = EXCLUDED.encrypted_access_token,
			encrypted_refresh_token = EXCLUDED.encrypted_refresh_token,
			scopes = EXCLUDED.scopes,
			expires_at = EXCLUDED.expires_at,
			updated_at = CURRENT_TIMESTAMP`

	_, err = db.GetDB().ExecContext(ctx, query, clerkUserID, twitchUserID, encryptedAccessToken, encryptedRefreshToken, scopesString, token.Expiry)
	if err != nil {
		return fmt.Errorf("failed to store tokens: %w", err)
	}

	if isAccountSwitch {
		log.Printf("✅ Account switch completed: Stored new Twitch tokens for user %s (New Twitch ID: %s)", clerkUserID, twitchUserID)
	} else {
		log.Printf("✅ Stored Twitch tokens for user %s (Twitch ID: %s)", clerkUserID, twitchUserID)
	}

	return nil
}

func (oc *OAuthConfig) GetStoredTokens(ctx context.Context, db database.Service, clerkUserID string) (*oauth2.Token, error) {
	var tokens UserTwitchTokens
	query := `
		SELECT clerk_user_id, twitch_user_id, encrypted_access_token, 
			   encrypted_refresh_token, scopes, expires_at, created_at, updated_at
		FROM user_twitch_tokens
		WHERE clerk_user_id = $1`

	err := db.GetDB().QueryRowContext(ctx, query, clerkUserID).Scan(
		&tokens.ClerkUserID, &tokens.TwitchUserID, &tokens.EncryptedAccessToken,
		&tokens.EncryptedRefreshToken, &tokens.Scopes, &tokens.ExpiresAt,
		&tokens.CreatedAt, &tokens.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("twitch not connected")
		}
		return nil, fmt.Errorf("failed to get stored tokens: %w", err)
	}

	accessToken, err := decryptToken(tokens.EncryptedAccessToken, oc.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt access token: %w", err)
	}

	refreshToken, err := decryptToken(tokens.EncryptedRefreshToken, oc.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt refresh token: %w", err)
	}

	return &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expiry:       tokens.ExpiresAt,
	}, nil
}

func (oc *OAuthConfig) GetValidToken(ctx context.Context, db database.Service, clerkUserID string) (*oauth2.Token, error) {
	storedToken, err := oc.GetStoredTokens(ctx, db, clerkUserID)
	if err != nil {
		return nil, err
	}

	tokenSource := oc.TokenSource(ctx, storedToken)
	freshToken, err := tokenSource.Token()
	if err != nil {
		log.Printf("⚠️ Failed to refresh token for user %s: %v", clerkUserID, err)
		if deleteErr := oc.DeleteStoredTokens(ctx, db, clerkUserID); deleteErr != nil {
			log.Printf("⚠️ Failed to delete invalid tokens for user %s: %v", clerkUserID, deleteErr)
		}
		return nil, fmt.Errorf("twitch re-authentication required")
	}

	if freshToken.AccessToken != storedToken.AccessToken || freshToken.RefreshToken != storedToken.RefreshToken {
		var tokens UserTwitchTokens
		query := `SELECT twitch_user_id FROM user_twitch_tokens WHERE clerk_user_id = $1`
		if err := db.GetDB().QueryRowContext(ctx, query, clerkUserID).Scan(&tokens.TwitchUserID); err != nil {
			return nil, fmt.Errorf("failed to get twitch user ID: %w", err)
		}

		if err := oc.StoreTokens(ctx, db, clerkUserID, tokens.TwitchUserID, freshToken); err != nil {
			log.Printf("⚠️ Failed to update refreshed tokens for user %s: %v", clerkUserID, err)
		} else {
			log.Printf("✅ Updated refreshed tokens for user %s", clerkUserID)
		}
	}

	return freshToken, nil
}

func (oc *OAuthConfig) DeleteStoredTokens(ctx context.Context, db database.Service, clerkUserID string) error {
	query := `DELETE FROM user_twitch_tokens WHERE clerk_user_id = $1`
	_, err := db.GetDB().ExecContext(ctx, query, clerkUserID)
	if err != nil {
		return fmt.Errorf("failed to delete tokens: %w", err)
	}
	log.Printf("🗑️ Deleted stored tokens for user %s", clerkUserID)
	return nil
}

func (oc *OAuthConfig) GetTwitchUserID(ctx context.Context, db database.Service, clerkUserID string) (string, error) {
	var twitchUserID string
	query := `SELECT twitch_user_id FROM user_twitch_tokens WHERE clerk_user_id = $1`
	err := db.GetDB().QueryRowContext(ctx, query, clerkUserID).Scan(&twitchUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("twitch not connected")
		}
		return "", fmt.Errorf("failed to get twitch user ID: %w", err)
	}
	return twitchUserID, nil
}
