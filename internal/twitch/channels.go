package twitch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetChannelInfo(ctx context.Context, userAccessToken string, broadcasterID string) (*ChannelInfo, error) {
	url := fmt.Sprintf("%s/channels?broadcaster_id=%s", twitchAPIBaseURL, broadcasterID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.clientID != "" {
		req.Header.Set("Client-ID", c.clientID)
	}

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

	var channelResp ChannelResponse
	if err := json.NewDecoder(resp.Body).Decode(&channelResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(channelResp.Data) == 0 {
		return nil, errors.New("no channel data found")
	}

	return &channelResp.Data[0], nil
}
