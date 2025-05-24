package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ClipInfo represents a Twitch clip
type ClipInfo struct {
	ID              string    `json:"id"`
	URL             string    `json:"url"`
	EmbedURL        string    `json:"embed_url"`
	BroadcasterID   string    `json:"broadcaster_id"`
	BroadcasterName string    `json:"broadcaster_name"`
	CreatorID       string    `json:"creator_id"`
	CreatorName     string    `json:"creator_name"`
	VideoID         string    `json:"video_id"`
	GameID          string    `json:"game_id"`
	Language        string    `json:"language"`
	Title           string    `json:"title"`
	ViewCount       int       `json:"view_count"`
	CreatedAt       time.Time `json:"created_at"`
	ThumbnailURL    string    `json:"thumbnail_url"`
	Duration        float64   `json:"duration"`
	VodOffset       int       `json:"vod_offset"`
	IsFeatured      bool      `json:"is_featured"`
}

// ClipsResponse represents the response from the Twitch Get Clips endpoint
type ClipsResponse struct {
	Data       []ClipInfo `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

// GetClips fetches clips for a specific broadcaster
func (c *Client) GetClips(ctx context.Context, userAccessToken string, broadcasterID string, limit int) ([]ClipInfo, error) {
	baseURL := "https://api.twitch.tv/helix/clips"
	params := url.Values{}
	params.Add("broadcaster_id", broadcasterID)
	params.Add("first", strconv.Itoa(limit))

	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -731)
	params.Add("started_at", startTime.Format(time.RFC3339))
	params.Add("ended_at", endTime.Format(time.RFC3339))

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Client-ID", c.clientID)
	req.Header.Set("Authorization", "Bearer "+userAccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body := make([]byte, 1024)
		resp.Body.Read(body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var clipsResponse ClipsResponse
	if err := json.NewDecoder(resp.Body).Decode(&clipsResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return clipsResponse.Data, nil
}
