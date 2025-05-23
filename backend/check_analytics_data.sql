-- Check analytics data in database
SELECT 'channel_analytics' as table_name, COUNT(*) as count FROM channel_analytics
UNION ALL
SELECT 'video_analytics' as table_name, COUNT(*) as count FROM video_analytics  
UNION ALL
SELECT 'stream_sessions' as table_name, COUNT(*) as count FROM stream_sessions
UNION ALL
SELECT 'analytics_jobs' as table_name, COUNT(*) as count FROM analytics_jobs;

-- Check recent channel analytics
SELECT 'Recent channel analytics' as info;
SELECT user_id, date, followers_count, total_views, created_at 
FROM channel_analytics 
ORDER BY created_at DESC 
LIMIT 5;

-- Check recent video analytics  
SELECT 'Recent video analytics' as info;
SELECT user_id, title, view_count, video_type, published_at 
FROM video_analytics 
ORDER BY created_at DESC 
LIMIT 5;

-- Check recent analytics jobs
SELECT 'Recent analytics jobs' as info;
SELECT user_id, job_type, status, started_at, completed_at, error_message
FROM analytics_jobs 
ORDER BY created_at DESC 
LIMIT 5; 