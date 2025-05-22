-- Migration: 001_create_analytics_tables.sql
-- Description: Create comprehensive analytics tables for creator metrics

-- Users table (if not exists) - stores creator profiles
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    clerk_user_id VARCHAR(255) UNIQUE NOT NULL,
    twitch_user_id VARCHAR(255),
    username VARCHAR(255),
    display_name VARCHAR(255),
    email VARCHAR(255),
    profile_image_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Channel Analytics - daily snapshots of channel metrics
CREATE TABLE IF NOT EXISTS channel_analytics (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    followers_count INTEGER DEFAULT 0,
    following_count INTEGER DEFAULT 0,
    total_views INTEGER DEFAULT 0,
    subscriber_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, date)
);

-- Stream Sessions - individual stream performance
CREATE TABLE IF NOT EXISTS stream_sessions (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE CASCADE,
    stream_id VARCHAR(255) UNIQUE,
    title TEXT,
    game_name VARCHAR(255),
    game_id VARCHAR(255),
    started_at TIMESTAMP WITH TIME ZONE,
    ended_at TIMESTAMP WITH TIME ZONE,
    duration_minutes INTEGER,
    peak_viewers INTEGER DEFAULT 0,
    average_viewers INTEGER DEFAULT 0,
    total_chatters INTEGER DEFAULT 0,
    followers_gained INTEGER DEFAULT 0,
    subscribers_gained INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Video Analytics - VODs, highlights, clips performance
CREATE TABLE IF NOT EXISTS video_analytics (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE CASCADE,
    video_id VARCHAR(255) UNIQUE NOT NULL,
    title TEXT,
    video_type VARCHAR(50), -- 'vod', 'highlight', 'clip', 'upload'
    duration_seconds INTEGER,
    view_count INTEGER DEFAULT 0,
    like_count INTEGER DEFAULT 0,
    comment_count INTEGER DEFAULT 0,
    thumbnail_url TEXT,
    published_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Daily Video Performance - track view changes over time
CREATE TABLE IF NOT EXISTS video_daily_stats (
    id SERIAL PRIMARY KEY,
    video_id VARCHAR(255) REFERENCES video_analytics(video_id) ON DELETE CASCADE,
    date DATE NOT NULL,
    view_count INTEGER DEFAULT 0,
    like_count INTEGER DEFAULT 0,
    comment_count INTEGER DEFAULT 0,
    watch_time_minutes INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(video_id, date)
);

-- Game/Category Performance
CREATE TABLE IF NOT EXISTS game_analytics (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE CASCADE,
    game_id VARCHAR(255),
    game_name VARCHAR(255),
    total_streams INTEGER DEFAULT 0,
    total_hours_streamed DECIMAL(10,2) DEFAULT 0,
    average_viewers DECIMAL(10,2) DEFAULT 0,
    peak_viewers INTEGER DEFAULT 0,
    total_followers_gained INTEGER DEFAULT 0,
    last_streamed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, game_id)
);

-- Social Media Cross-Platform Analytics (for future expansion)
CREATE TABLE IF NOT EXISTS social_analytics (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE CASCADE,
    platform VARCHAR(50), -- 'twitter', 'youtube', 'instagram', 'tiktok'
    date DATE NOT NULL,
    followers_count INTEGER DEFAULT 0,
    posts_count INTEGER DEFAULT 0,
    engagement_rate DECIMAL(5,2) DEFAULT 0,
    impressions INTEGER DEFAULT 0,
    clicks INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, platform, date)
);

-- Analytics Jobs - track data collection status
CREATE TABLE IF NOT EXISTS analytics_jobs (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE CASCADE,
    job_type VARCHAR(100), -- 'daily_channel', 'stream_session', 'video_stats'
    status VARCHAR(50), -- 'pending', 'running', 'completed', 'failed'
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    data_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_channel_analytics_user_date ON channel_analytics(user_id, date DESC);
CREATE INDEX IF NOT EXISTS idx_stream_sessions_user_started ON stream_sessions(user_id, started_at DESC);
CREATE INDEX IF NOT EXISTS idx_video_analytics_user_published ON video_analytics(user_id, published_at DESC);
CREATE INDEX IF NOT EXISTS idx_video_daily_stats_video_date ON video_daily_stats(video_id, date DESC);
CREATE INDEX IF NOT EXISTS idx_game_analytics_user_hours ON game_analytics(user_id, total_hours_streamed DESC);
CREATE INDEX IF NOT EXISTS idx_analytics_jobs_user_type ON analytics_jobs(user_id, job_type, created_at DESC);

-- Create views for common analytics queries
CREATE OR REPLACE VIEW user_analytics_summary AS
SELECT 
    u.id as user_id,
    u.username,
    u.display_name,
    ca.followers_count as current_followers,
    ca.total_views as current_total_views,
    ca.subscriber_count as current_subscribers,
    (
        SELECT COUNT(*)
        FROM stream_sessions ss
        WHERE ss.user_id = u.id
        AND ss.started_at >= CURRENT_DATE - INTERVAL '30 days'
    ) as streams_last_30_days,
    (
        SELECT AVG(average_viewers)
        FROM stream_sessions ss
        WHERE ss.user_id = u.id
        AND ss.started_at >= CURRENT_DATE - INTERVAL '30 days'
    ) as avg_viewers_last_30_days,
    (
        SELECT SUM(duration_minutes)
        FROM stream_sessions ss
        WHERE ss.user_id = u.id
        AND ss.started_at >= CURRENT_DATE - INTERVAL '30 days'
    ) as total_minutes_streamed_last_30_days
FROM users u
LEFT JOIN channel_analytics ca ON u.id = ca.user_id 
    AND ca.date = (
        SELECT MAX(date) 
        FROM channel_analytics ca2 
        WHERE ca2.user_id = u.id
    );
