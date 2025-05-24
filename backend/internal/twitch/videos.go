package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetUserVideos retrieves videos for a specific user
func (c *Client) GetUserVideos(ctx context.Context, userAccessToken string, userID string, limit int) ([]VideoInfo, string, error) {
	if limit <= 0 || limit > 100 {
		limit = 20 // Default limit
	}

	url := fmt.Sprintf("%s/videos?user_id=%s&first=%d", twitchAPIBaseURL, userID, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	// When using Clerk OAuth tokens, the client ID should be included in the token's scopes
	// But we'll still set it if available as a fallback
	if c.clientID != "" {
		req.Header.Set("Client-ID", c.clientID)
	}

	// Set the authorization header with the user's OAuth token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userAccessToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("twitch API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var videosResp VideosResponse
	if err := json.NewDecoder(resp.Body).Decode(&videosResp); err != nil {
		return nil, "", fmt.Errorf("failed to decode response: %w", err)
	}

	return videosResp.Data, videosResp.Pagination.Cursor, nil
}

// GetVideosByID retrieves specific videos by their IDs
func (c *Client) GetVideosByID(ctx context.Context, userAccessToken string, videoIDs []string) ([]VideoInfo, error) {
	if len(videoIDs) == 0 {
		return nil, fmt.Errorf("no video IDs provided")
	}

	if len(videoIDs) > 100 {
		return nil, fmt.Errorf("too many video IDs provided, maximum is 100")
	}

	// Build the URL with video IDs as query parameters
	baseURL := fmt.Sprintf("%s/videos", twitchAPIBaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	q := req.URL.Query()
	for _, id := range videoIDs {
		q.Add("id", id)
	}
	req.URL.RawQuery = q.Encode()

	// When using Clerk OAuth tokens, the client ID should be included in the token's scopes
	// But we'll still set it if available as a fallback
	if c.clientID != "" {
		req.Header.Set("Client-ID", c.clientID)
	}

	// Set the authorization header with the user's OAuth token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userAccessToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("twitch API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var videosResp VideosResponse
	if err := json.NewDecoder(resp.Body).Decode(&videosResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return videosResp.Data, nil
}
