package twitch

import (
	"time"
)

// ChannelInfo represents a Twitch channel
type ChannelInfo struct {
	ID              string   `json:"id"`
	BroadcasterID   string   `json:"broadcaster_id"`
	BroadcasterName string   `json:"broadcaster_name"`
	GameName        string   `json:"game_name"`
	GameID          string   `json:"game_id"`
	Title           string   `json:"title"`
	Language        string   `json:"language"`
	ThumbnailURL    string   `json:"thumbnail_url"`
	StartedAt       string   `json:"started_at"`
	TagIDs          []string `json:"tag_ids"`
	IsMature        bool     `json:"is_mature"`
}

// ChannelResponse represents the response from the Twitch Get Channel endpoint
type ChannelResponse struct {
	Data []ChannelInfo `json:"data"`
}

// VideoInfo represents a Twitch video
type VideoInfo struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	UserName      string    `json:"user_name"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	PublishedAt   time.Time `json:"published_at"`
	URL           string    `json:"url"`
	ThumbnailURL  string    `json:"thumbnail_url"`
	ViewCount     int       `json:"view_count"`
	Language      string    `json:"language"`
	Type          string    `json:"type"`
	Duration      string    `json:"duration"`
	MutedSegments []struct {
		Duration int `json:"duration"`
		Offset   int `json:"offset"`
	} `json:"muted_segments"`
}

// VideosResponse represents the response from the Twitch Get Videos endpoint
type VideosResponse struct {
	Data       []VideoInfo `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

// TokenValidationResponse represents the response from Twitch's token validation endpoint
type TokenValidationResponse struct {
	ClientID  string   `json:"client_id"`
	Login     string   `json:"login"`
	UserID    string   `json:"user_id"`
	Scopes    []string `json:"scopes"`
	ExpiresIn int      `json:"expires_in"`
}
