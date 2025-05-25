'use client';

import { useState, useEffect, useCallback } from 'react';
import { motion } from 'framer-motion';
import { useAuth } from '@clerk/nextjs';
import { 
    Users, 
    Video, 
    Sparkles, 
    ArrowRight, 
    TrendingUp,
    Eye,
    Heart,
    RefreshCw,
    Play,
    Clock,
    Target
} from 'lucide-react';

// Types for real data
interface DashboardOverview {
    totalViews: number;
    videoCount: number;
    averageViewsPerVideo: number;
    totalWatchTimeHours: number;
    currentFollowers: number;
    currentSubscribers: number;
    followerChange: number;
    subscriberChange: number;
}

interface RecentVideo {
    id: number;
    title: string;
    view_count: number;
    published_at: string;
    video_type: string;
}

interface OverviewData {
    overview: DashboardOverview;
    recentVideos: RecentVideo[];
    connectionStatus: {
        twitch_connected: boolean;
    };
}

// Utility functions
const formatNumber = (num: number | undefined | null): string => {
    if (num === undefined || num === null || isNaN(num)) return '0';
    if (num >= 1000000) return `${(num / 1000000).toFixed(1)}M`;
    if (num >= 1000) return `${(num / 1000).toFixed(1)}K`;
    return Math.round(num).toString();
};

const formatDuration = (hours: number | undefined | null): string => {
    if (hours === undefined || hours === null || isNaN(hours)) return '0 min';
    if (hours >= 1) return `${Math.round(hours)} hrs`;
    return `${Math.round(hours * 60)} min`;
};

const getTimeAgo = (dateString: string): string => {
    const date = new Date(dateString);
    const now = new Date();
    const diffInHours = Math.floor((now.getTime() - date.getTime()) / (1000 * 60 * 60));
    
    if (diffInHours < 1) return 'Just now';
    if (diffInHours < 24) return `${diffInHours}h ago`;
    const diffInDays = Math.floor(diffInHours / 24);
    if (diffInDays < 7) return `${diffInDays}d ago`;
    return date.toLocaleDateString();
};

// Enhanced Metric Card Component
const EnhancedMetricCard = ({
    title,
    value,
    subtitle,
    change,
    icon: Icon,
    gradient = "from-emerald-500/10 to-teal-500/5",
    iconColor = "text-emerald-400"
}: {
    title: string;
    value: string | number;
    subtitle: string;
    change?: { value: number; period: string };
    icon: React.ComponentType<{ className?: string }>;
    gradient?: string;
    iconColor?: string;
}) => (
    <motion.div
        initial={{ y: 20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        className="bg-gray-900/50 backdrop-blur-xl border border-gray-800/50 rounded-2xl p-6 relative overflow-hidden group hover:border-emerald-500/30 transition-all duration-300"
    >
        <div className={`absolute inset-0 bg-gradient-to-br ${gradient} opacity-0 group-hover:opacity-100 transition-opacity duration-500`}></div>
        <div className="relative z-10">
            <div className={`flex items-center ${iconColor} mb-4`}>
                <Icon className="w-5 h-5 mr-2" />
                <h3 className="text-sm font-medium text-gray-400">{title}</h3>
            </div>
            <p className="text-3xl font-bold text-white mb-2">
                {typeof value === 'number' ? formatNumber(value) : value}
            </p>
            <div className="flex items-center justify-between">
                <p className="text-sm text-gray-500">{subtitle}</p>
                {change && (
                    <span className={`text-xs px-2 py-1 rounded-full ${
                        change.value >= 0 
                            ? 'bg-emerald-500/20 text-emerald-300' 
                            : 'bg-red-500/20 text-red-300'
                    }`}>
                        {change.value >= 0 ? '+' : ''}{change.value}% {change.period}
                    </span>
                )}
            </div>
        </div>
    </motion.div>
);

// Zero Data State Component
const ZeroDataState = ({ isConnected }: { isConnected: boolean }) => (
    <motion.div
        initial={{ y: 20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        className="max-w-2xl mx-auto text-center py-16"
    >
        <div className="bg-gradient-to-br from-emerald-500/10 to-blue-500/10 rounded-3xl p-8 border border-emerald-500/20">
            <div className="w-20 h-20 bg-gradient-to-br from-emerald-500 to-blue-500 rounded-full mx-auto mb-6 flex items-center justify-center">
                <Sparkles className="w-10 h-10 text-white" />
            </div>
            
            <h2 className="text-3xl font-bold text-white mb-4">
                {isConnected ? "Welcome to Your Creator Journey!" : "Ready to Get Started?"}
            </h2>
            
            <p className="text-gray-300 text-lg mb-6 leading-relaxed">
                {isConnected 
                    ? "We can see your Twitch account is connected, but it looks like you're just getting started. That's exciting! Every successful creator started exactly where you are now."
                    : "Connect your Twitch account to start tracking your content performance and grow your audience with data-driven insights."
                }
            </p>
            
            {isConnected ? (
                <div className="space-y-4">
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
                        <div className="bg-gray-800/50 rounded-xl p-4">
                            <Video className="w-6 h-6 text-emerald-400 mx-auto mb-2" />
                            <p className="text-gray-300">Create your first stream or upload content</p>
                        </div>
                        <div className="bg-gray-800/50 rounded-xl p-4">
                            <Users className="w-6 h-6 text-blue-400 mx-auto mb-2" />
                            <p className="text-gray-300">Build your community</p>
                        </div>
                        <div className="bg-gray-800/50 rounded-xl p-4">
                            <TrendingUp className="w-6 h-6 text-purple-400 mx-auto mb-2" />
                            <p className="text-gray-300">Track your growth here</p>
                        </div>
                    </div>
                    
                    <div className="bg-emerald-500/10 border border-emerald-500/20 rounded-xl p-4 mt-6">
                        <p className="text-emerald-300 text-sm">
                            💡 <strong>Pro Tip:</strong> Once you start creating content, come back here to see detailed analytics about your performance, audience growth, and content insights!
                        </p>
                    </div>
                </div>
            ) : (
                <button 
                    onClick={() => window.location.href = '/settings'}
                    className="bg-gradient-to-r from-emerald-500 to-blue-500 text-white px-8 py-3 rounded-xl font-semibold hover:shadow-lg hover:shadow-emerald-500/25 transition-all duration-300 flex items-center mx-auto"
                >
                    Connect Twitch Account
                    <ArrowRight className="w-5 h-5 ml-2" />
                </button>
            )}
        </div>
    </motion.div>
);

export function OverviewSection() {
    const { isLoaded, isSignedIn, getToken } = useAuth();
    const [data, setData] = useState<OverviewData | null>(null);
    const [loading, setLoading] = useState(true);
    const [refreshing, setRefreshing] = useState(false);

    // Helper function to get API base URL
    const getApiBaseUrl = () => {
        const hostname = typeof window !== 'undefined' ? window.location.hostname : '';
        const nodeEnv = process.env.NODE_ENV;
        const appEnv = process.env.NEXT_PUBLIC_APP_ENV;

        if (hostname === 'dev.creatorsync.app') {
            return 'https://api-dev.creatorsync.app';
        } else if (appEnv === 'staging') {
            return 'https://api-dev.creatorsync.app';
        } else if (nodeEnv === 'production') {
            return 'https://api.creatorsync.app';
        } else {
            return 'http://localhost:8080';
        }
    };

    // Fetch overview data
    const fetchOverviewData = useCallback(async () => {
        if (!isLoaded || !isSignedIn) return;

        try {
            const token = await getToken();
            const apiBaseUrl = getApiBaseUrl();

            const response = await fetch(`${apiBaseUrl}/api/analytics/overview`, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
            });

            if (response.ok) {
                const data = await response.json();
                console.log('Overview API Response:', data);
                
                // Handle new backend response format with connection_status check
                if (data.connection_status && !data.connection_status.twitch_connected) {
                    // User is disconnected - set empty state regardless of cached data
                    setData({
                        overview: {
                            totalViews: 0,
                            videoCount: 0,
                            averageViewsPerVideo: 0,
                            totalWatchTimeHours: 0,
                            currentFollowers: 0,
                            currentSubscribers: 0,
                            followerChange: 0,
                            subscriberChange: 0,
                        },
                        recentVideos: [],
                        connectionStatus: data.connection_status
                    });
                    return;
                }
                
                // User is connected - process the data normally
                let analyticsData;
                let connectionStatus;
                
                if (data.overview && data.connection_status) {
                    // New format with connection status
                    analyticsData = data;
                    connectionStatus = data.connection_status;
                } else if (data.analytics && data.connection_status) {
                    // Alternative new format
                    analyticsData = data.analytics;
                    connectionStatus = data.connection_status;
                } else {
                    // Legacy format - assume connected if we got data
                    analyticsData = data;
                    connectionStatus = { twitch_connected: true, settings_url: '/settings' };
                }
                
                // Transform the analytics data for overview
                const overviewData: OverviewData = {
                    overview: {
                        totalViews: analyticsData.overview?.totalViews || 0,
                        videoCount: analyticsData.overview?.videoCount || 0,
                        averageViewsPerVideo: analyticsData.overview?.averageViewsPerVideo || 0,
                        totalWatchTimeHours: analyticsData.overview?.totalWatchTimeHours || 0,
                        currentFollowers: analyticsData.overview?.currentFollowers || analyticsData.overview?.CurrentFollowers || 0,
                        currentSubscribers: analyticsData.overview?.currentSubscribers || analyticsData.overview?.CurrentSubscribers || 0,
                        followerChange: analyticsData.overview?.followerChange || analyticsData.overview?.FollowerChange || 0,
                        subscriberChange: analyticsData.overview?.subscriberChange || analyticsData.overview?.SubscriberChange || 0,
                    },
                    recentVideos: (analyticsData.recentVideos || []).slice(0, 3),
                    connectionStatus: connectionStatus
                };

                setData(overviewData);
            } else {
                console.error('Failed to fetch overview data:', response.status);
                // If fetch fails, assume disconnected state
                setData({
                    overview: {
                        totalViews: 0,
                        videoCount: 0,
                        averageViewsPerVideo: 0,
                        totalWatchTimeHours: 0,
                        currentFollowers: 0,
                        currentSubscribers: 0,
                        followerChange: 0,
                        subscriberChange: 0,
                    },
                    recentVideos: [],
                    connectionStatus: { twitch_connected: false }
                });
            }
        } catch (error) {
            console.error('Error fetching overview data:', error);
            // On error, assume disconnected state
            setData({
                overview: {
                    totalViews: 0,
                    videoCount: 0,
                    averageViewsPerVideo: 0,
                    totalWatchTimeHours: 0,
                    currentFollowers: 0,
                    currentSubscribers: 0,
                    followerChange: 0,
                    subscriberChange: 0,
                },
                recentVideos: [],
                connectionStatus: { twitch_connected: false }
            });
        } finally {
            setLoading(false);
        }
    }, [isLoaded, isSignedIn, getToken]);

    const handleRefresh = async () => {
        if (refreshing) return;
        setRefreshing(true);
        await fetchOverviewData();
        setRefreshing(false);
    };

    // Initial load
    useEffect(() => {
        if (isLoaded && isSignedIn) {
            fetchOverviewData();
        }
    }, [isLoaded, isSignedIn, fetchOverviewData]);

    // Loading state
    if (!isLoaded || !isSignedIn || loading) {
        return (
            <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 flex items-center justify-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-500"></div>
            </div>
        );
    }

    const safeOverview = data?.overview || {
        totalViews: 0,
        videoCount: 0,
        averageViewsPerVideo: 0,
        totalWatchTimeHours: 0,
        currentFollowers: 0,
        currentSubscribers: 0,
        followerChange: 0,
        subscriberChange: 0,
    };

    // Calculate engagement metrics
    const engagementRate = safeOverview.videoCount > 0 ? safeOverview.totalViews / safeOverview.videoCount : 0;
    const viewsPerFollower = safeOverview.currentFollowers > 0 ? safeOverview.totalViews / safeOverview.currentFollowers : 0;

    // Check if we should show zero data state
    const isDisconnected = data?.connectionStatus?.twitch_connected === false;
    const hasZeroData = data?.connectionStatus?.twitch_connected === true && 
                        safeOverview.totalViews === 0 && 
                        safeOverview.videoCount === 0 && 
                        safeOverview.currentFollowers === 0;

    // Show zero data state for disconnected users or connected users with no data
    if (isDisconnected || hasZeroData) {
        return (
            <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 p-6">
                <ZeroDataState isConnected={!isDisconnected} />
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 p-6">
            <div className="max-w-7xl mx-auto">
                {/* Header */}
                <div className="flex items-center justify-between mb-8">
                    <div>
                        <div className="flex items-center mb-2">
                            <Sparkles className="w-6 h-6 text-emerald-500 mr-2" />
                            <span className="text-sm text-emerald-400 font-medium">Dashboard Overview</span>
                        </div>
                        <h1 className="text-4xl font-bold mb-2">
                            <span className="text-white">Welcome back, </span>
                            <span className="bg-gradient-to-r from-emerald-400 to-blue-400 bg-clip-text text-transparent">Creator</span>
                        </h1>
                        <p className="text-gray-400 text-lg">Track your channel's performance and recent activity at a glance.</p>
                    </div>
                    <button
                        onClick={handleRefresh}
                        disabled={refreshing}
                        className="flex items-center px-4 py-2 bg-gray-700 text-white rounded-lg hover:bg-gray-600 transition-colors disabled:opacity-50"
                    >
                        <RefreshCw className={`w-4 h-4 mr-2 ${refreshing ? 'animate-spin' : ''}`} />
                        Refresh
                    </button>
                </div>

                {/* Key Metrics Grid */}
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                    <EnhancedMetricCard
                        title="Followers"
                        value={safeOverview.currentFollowers}
                        subtitle="Growing your community"
                        change={{ value: safeOverview.followerChange, period: "this week" }}
                        icon={Heart}
                        gradient="from-rose-500/10 to-pink-500/5"
                        iconColor="text-rose-400"
                    />
                    <EnhancedMetricCard
                        title="Subscribers"
                        value={safeOverview.currentSubscribers}
                        subtitle="Supporting your channel"
                        change={{ value: safeOverview.subscriberChange, period: "this month" }}
                        icon={Users}
                        gradient="from-violet-500/10 to-purple-500/5"
                        iconColor="text-violet-400"
                    />
                    <EnhancedMetricCard
                        title="Total Views"
                        value={safeOverview.totalViews}
                        subtitle={`${formatNumber(engagementRate)} avg per video`}
                        icon={Eye}
                        gradient="from-emerald-500/10 to-teal-500/5"
                        iconColor="text-emerald-400"
                    />
                    <EnhancedMetricCard
                        title="Content"
                        value={safeOverview.videoCount}
                        subtitle={`${formatDuration(safeOverview.totalWatchTimeHours)} watch time`}
                        icon={Video}
                        gradient="from-blue-500/10 to-cyan-500/5"
                        iconColor="text-blue-400"
                    />
                </div>

                {/* Main Content Grid */}
                <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                    {/* Recent Content Performance */}
                    <div className="lg:col-span-2">
                        <motion.div
                            initial={{ y: 20, opacity: 0 }}
                            animate={{ y: 0, opacity: 1 }}
                            className="bg-gray-900/50 backdrop-blur-xl border border-gray-800/50 rounded-2xl p-6"
                        >
                            <div className="flex items-center justify-between mb-6">
                                <div>
                                    <h3 className="text-xl font-bold text-white flex items-center">
                                        <TrendingUp className="w-5 h-5 mr-2 text-emerald-400" />
                                        Recent Content Performance
                                    </h3>
                                    <p className="text-gray-400 text-sm mt-1">Your latest videos and their performance</p>
                                </div>
                                <span className="text-xs text-gray-500 bg-gray-800/50 px-2 py-1 rounded-full">
                                    Last 7 days
                                </span>
                            </div>

                            {data?.recentVideos && data.recentVideos.length > 0 ? (
                                <div className="space-y-4">
                                    {data.recentVideos.map((video, index) => (
                                        <motion.div
                                            key={`recent-video-${video.id}-${index}`}
                                            initial={{ x: -20, opacity: 0 }}
                                            animate={{ x: 0, opacity: 1 }}
                                            transition={{ delay: index * 0.1 }}
                                            className="flex items-center gap-4 p-4 bg-gray-800/30 rounded-xl hover:bg-gray-800/50 transition-all duration-300 group"
                                        >
                                            <div className="w-12 h-12 bg-gradient-to-br from-emerald-500/20 to-blue-500/20 rounded-lg flex items-center justify-center">
                                                <Play className="w-6 h-6 text-emerald-400" />
                                            </div>
                                            <div className="flex-1">
                                                <h4 className="text-white font-medium group-hover:text-emerald-300 transition-colors line-clamp-1">
                                                    {video.title}
                                                </h4>
                                                <div className="flex items-center gap-4 mt-1 text-sm text-gray-400">
                                                    <span className="flex items-center">
                                                        <Eye className="w-4 h-4 mr-1" />
                                                        {formatNumber(video.view_count)} views
                                                    </span>
                                                    <span className="flex items-center">
                                                        <Clock className="w-4 h-4 mr-1" />
                                                        {getTimeAgo(video.published_at)}
                                                    </span>
                                                    <span className={`px-2 py-1 rounded-full text-xs ${
                                                        video.video_type?.toLowerCase().includes('clip') 
                                                            ? 'bg-emerald-500/20 text-emerald-300' 
                                                            : 'bg-blue-500/20 text-blue-300'
                                                    }`}>
                                                        {video.video_type?.toLowerCase().includes('clip') ? 'Clip' : 'Broadcast'}
                                                    </span>
                                                </div>
                                            </div>
                                            <div className="text-right">
                                                <div className="text-lg font-bold text-white">
                                                    {formatNumber(video.view_count)}
                                                </div>
                                                <div className="text-xs text-gray-500">views</div>
                                            </div>
                                        </motion.div>
                                    ))}
                                </div>
                            ) : (
                                <div className="text-center py-8">
                                    <Video className="w-12 h-12 text-gray-600 mx-auto mb-3" />
                                    <p className="text-gray-400">No recent content found</p>
                                    <p className="text-gray-500 text-sm">Start creating to see your performance here!</p>
                                </div>
                            )}
                        </motion.div>
                    </div>

                    {/* Quick Insights & Tips */}
                    <div className="space-y-6">
                        {/* Performance Insights */}
                        <motion.div
                            initial={{ y: 20, opacity: 0 }}
                            animate={{ y: 0, opacity: 1 }}
                            transition={{ delay: 0.3 }}
                            className="bg-gray-900/50 backdrop-blur-xl border border-gray-800/50 rounded-2xl p-6"
                        >
                            <h3 className="text-lg font-semibold text-white mb-4 flex items-center">
                                <Target className="w-5 h-5 mr-2 text-emerald-400" />
                                Quick Insights
                            </h3>
                            <div className="space-y-4">
                                {safeOverview.currentFollowers > 0 && (
                                    <div className="p-3 bg-emerald-500/10 border border-emerald-500/20 rounded-lg">
                                        <p className="text-sm text-emerald-300">
                                            <strong>Great!</strong> You have {formatNumber(safeOverview.currentFollowers)} followers.
                                        </p>
                                    </div>
                                )}
                                {safeOverview.videoCount > 5 && (
                                    <div className="p-3 bg-blue-500/10 border border-blue-500/20 rounded-lg">
                                        <p className="text-sm text-blue-300">
                                            <strong>Active Creator:</strong> {safeOverview.videoCount} pieces of content created!
                                        </p>
                                    </div>
                                )}
                                {viewsPerFollower > 0 && (
                                    <div className="p-3 bg-purple-500/10 border border-purple-500/20 rounded-lg">
                                        <p className="text-sm text-purple-300">
                                            <strong>Engagement:</strong> {formatNumber(viewsPerFollower)} views per follower.
                                        </p>
                                    </div>
                                )}
                                {safeOverview.videoCount === 0 && (
                                    <div className="p-3 bg-yellow-500/10 border border-yellow-500/20 rounded-lg">
                                        <p className="text-sm text-yellow-300">
                                            <strong>Get Started:</strong> Create your first content to see insights here!
                                        </p>
                                    </div>
                                )}
                            </div>
                        </motion.div>

                        {/* Creator Tip */}
                        <motion.div
                            initial={{ y: 20, opacity: 0 }}
                            animate={{ y: 0, opacity: 1 }}
                            transition={{ delay: 0.4 }}
                            className="bg-gradient-to-br from-emerald-500/10 via-blue-500/5 to-purple-500/10 border border-emerald-500/20 rounded-2xl p-6"
                        >
                            <div className="flex items-center gap-2 text-emerald-400 mb-3">
                                <Sparkles className="w-5 h-5" />
                                <h3 className="font-semibold">Creator Tip</h3>
                            </div>
                            <p className="text-gray-300 mb-4">
                                Auto-schedule your clips across multiple platforms to maximize reach and engagement. Consistency is key to growing your audience!
                            </p>
                            <button className="flex items-center text-emerald-400 hover:text-emerald-300 transition-colors font-medium">
                                Try Scheduling <ArrowRight className="w-4 h-4 ml-1" />
                            </button>
                        </motion.div>
                    </div>
                </div>
            </div>
        </div>
    );
} 