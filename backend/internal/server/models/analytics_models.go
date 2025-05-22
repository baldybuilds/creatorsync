package models

import "time"

// VideoAnalyticsSummary holds aggregated analytics data for Twitch videos.
type VideoAnalyticsSummary struct {
	TotalVideosConsidered int            `json:"total_videos_considered"`
	TotalViews            int            `json:"total_views"`
	AverageViewsPerVideo  float64        `json:"average_views_per_video"`
	ContentDistribution   map[string]int `json:"content_distribution"` // e.g., {"archive": 5, "highlight": 3, "upload": 2}
	RequestedPeriodDays   int            `json:"requested_period_days,omitempty"`
	ActualDateRangeStart  *time.Time     `json:"actual_date_range_start,omitempty"` // Oldest video's PublishedAt in the considered set for the period
	ActualDateRangeEnd    time.Time      `json:"actual_date_range_end"`             // Typically time of request or newest video
}
