package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// GetBroadcasterSubscribers fetches a list of subscriptions for a given broadcaster.
// Required scope: channel:read:subscriptions
// See: https://dev.twitch.tv/docs/api/reference/#get-broadcaster-subscriptions
func (c *Client) GetBroadcasterSubscribers(ctx context.Context, userAccessToken, broadcasterID string, limit int, afterCursor string) (*SubscriptionsResponse, error) {
	if broadcasterID == "" {
		return nil, fmt.Errorf("broadcasterID cannot be empty")
	}

	// Construct the URL
	apiURL := fmt.Sprintf("%s/subscriptions", twitchAPIBaseURL)
	params := url.Values{}
	params.Set("broadcaster_id", broadcasterID)

	if limit <= 0 {
		limit = 20 // Default limit if not specified or invalid
	} else if limit > 100 {
		limit = 100 // Max limit per Twitch API
	}
	params.Set("first", strconv.Itoa(limit))

	if afterCursor != "" {
		params.Set("after", afterCursor)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	if c.clientID != "" {
		req.Header.Set("Client-ID", c.clientID)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userAccessToken))

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("twitch API error getting subscriptions: status %d, body: %s", resp.StatusCode, string(body))
	}

	// Decode response
	var response SubscriptionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode subscriptions response: %w", err)
	}

	return &response, nil
}
