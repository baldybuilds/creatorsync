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

// Subscription represents a Twitch subscriber.
// See https://dev.twitch.tv/docs/api/reference/#get-broadcaster-subscriptions
type Subscription struct {
	BroadcasterID    string `json:"broadcaster_id"`
	BroadcasterLogin string `json:"broadcaster_login"`
	BroadcasterName  string `json:"broadcaster_name"`
	GifterID         string `json:"gifter_id,omitempty"`
	GifterLogin      string `json:"gifter_login,omitempty"`
	GifterName       string `json:"gifter_name,omitempty"`
	IsGift           bool   `json:"is_gift"`
	Tier             string `json:"tier"`
	PlanName         string `json:"plan_name"`
	UserID           string `json:"user_id"`
	UserLogin        string `json:"user_login"`
	UserName         string `json:"user_name"`
}

// SubscriptionsResponse represents the response from the Twitch Get Broadcaster Subscriptions endpoint.
type SubscriptionsResponse struct {
	Data       []Subscription `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
	Total  int `json:"total"`
	Points int `json:"points"` // Subscription points, if requested by scope
}

// FollowersResponse represents the response from the Twitch Get Channel Followers endpoint
type FollowersResponse struct {
	Data       []Follower `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
	Total int `json:"total"`
}

// Follower represents a Twitch follower
type Follower struct {
	UserID     string    `json:"user_id"`
	UserName   string    `json:"user_name"`
	UserLogin  string    `json:"user_login"`
	FollowedAt time.Time `json:"followed_at"`
}

// UsersResponse represents the response from the Twitch Get Users endpoint
type UsersResponse struct {
	Data []User `json:"data"`
}

// User represents a Twitch user
type User struct {
	ID              string    `json:"id"`
	Login           string    `json:"login"`
	DisplayName     string    `json:"display_name"`
	Type            string    `json:"type"`
	BroadcasterType string    `json:"broadcaster_type"`
	Description     string    `json:"description"`
	ProfileImageURL string    `json:"profile_image_url"`
	OfflineImageURL string    `json:"offline_image_url"`
	ViewCount       int       `json:"view_count"`
	Email           string    `json:"email,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// StreamResponse represents the response from the Twitch Get Streams endpoint
type StreamResponse struct {
	Data       []StreamInfo `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

// StreamInfo represents a live Twitch stream
type StreamInfo struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	UserLogin    string    `json:"user_login"`
	UserName     string    `json:"user_name"`
	GameID       string    `json:"game_id"`
	GameName     string    `json:"game_name"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	ViewerCount  int       `json:"viewer_count"`
	StartedAt    time.Time `json:"started_at"`
	Language     string    `json:"language"`
	ThumbnailURL string    `json:"thumbnail_url"`
	TagIDs       []string  `json:"tag_ids"`
	Tags         []string  `json:"tags"`
	IsMature     bool      `json:"is_mature"`
}

// Video represents a Twitch video (for the enhanced client)
type Video struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	UserName     string     `json:"user_name"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	CreatedAt    time.Time  `json:"created_at"`
	PublishedAt  *time.Time `json:"published_at"`
	URL          string     `json:"url"`
	ThumbnailURL string     `json:"thumbnail_url"`
	ViewCount    int        `json:"view_count"`
	Language     string     `json:"language"`
	Type         string     `json:"type"`
	Duration     string     `json:"duration"`
}

// SubscribersResponse represents the response for subscriber data
type SubscribersResponse struct {
	Data       []Subscription `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
	Total int `json:"total"`
}
