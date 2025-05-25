"use client";

import { useState, useEffect, useCallback } from 'react';
import { motion } from 'framer-motion';
import { useAuth } from '@clerk/nextjs';
import {
    TrendingUp,
    Eye,
    Video,
    RefreshCw,
    BarChart3,
    Heart,
    Users
} from 'lucide-react';
import {
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
    ScatterChart,
    Scatter,
    Cell
} from 'recharts';
import React from 'react';

// Enhanced types for new dashboard design
interface VideoBasedOverview {
    totalViews: number;
    videoCount: number;
    averageViewsPerVideo: number;
    totalWatchTimeHours: number;
    currentFollowers: number;
    currentSubscribers: number;
    followerChange: number;
    subscriberChange: number;
}

interface ChartDataPoint {
    date: string;
    value: number;
}

interface ContentTypeData {
    date: string;
    broadcasts: number;
    clips: number;
    uploads: number;
}

interface PerformanceData {
    viewsOverTime: ChartDataPoint[];
    contentDistribution: ContentTypeData[];
}

interface VideoAnalytics {
    id: number;
    video_id: string;
    title: string;
    video_type: string;
    view_count: number;
    published_at: string;
}

interface EnhancedAnalytics {
    overview: VideoBasedOverview;
    performance: PerformanceData;
    topVideos: VideoAnalytics[];
    recentVideos: VideoAnalytics[];
}

interface ContentDataPoint {
    id: number;
    title: string;
    views: number;
    date: number;
    displayDate: string;
    daysSince: number;
    type: 'clip' | 'broadcast';
}

// Utility functions
const formatNumber = (num: number | undefined | null): string => {
    // Handle undefined, null, or NaN values
    if (num === undefined || num === null || isNaN(num)) {
        return '0';
    }

    if (num >= 1000000) {
        return `${(num / 1000000).toFixed(1)}M`;
    }
    if (num >= 1000) {
        return `${(num / 1000).toFixed(1)}K`;
    }
    return Math.round(num).toString();
};

const formatDuration = (hours: number | undefined | null): string => {
    // Handle undefined, null, or NaN values
    if (hours === undefined || hours === null || isNaN(hours)) {
        return '0 min';
    }

    if (hours >= 1) {
        return `${Math.round(hours)} hrs`;
    }
    return `${Math.round(hours * 60)} min`;
};

// Enhanced Metric Card for the new design
const EnhancedMetricCard = ({
    title,
    value,
    subtitle,
    icon: Icon,
    gradient = "from-emerald-500/10 to-teal-500/5"
}: {
    title: string;
    value: string | number;
    subtitle: string;
    icon: React.ComponentType<{ className?: string }>;
    gradient?: string;
}) => (
    <motion.div
        initial={{ y: 20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        className="bg-gray-900/50 backdrop-blur-xl border border-gray-800/50 rounded-2xl p-6 relative overflow-hidden group hover:border-emerald-500/30 transition-all duration-300"
    >
        <div className={`absolute inset-0 bg-gradient-to-br ${gradient} opacity-0 group-hover:opacity-100 transition-opacity duration-500`}></div>
        <div className="relative z-10">
            <div className="flex items-center text-emerald-400 mb-4">
                <Icon className="w-5 h-5 mr-2" />
                <h3 className="text-sm font-medium text-gray-400">{title}</h3>
            </div>
            <p className="text-3xl font-bold text-white mb-2">
                {typeof value === 'number' ? formatNumber(value) : value}
            </p>
            <p className="text-sm text-gray-500">{subtitle}</p>
        </div>
    </motion.div>
);

// Time period selector
const TimePeriodSelector = ({
    selected,
    onSelect
}: {
    selected: string;
    onSelect: (period: string) => void;
}) => (
    <div className="flex bg-gray-800/50 rounded-lg p-1">
        {['7', '30', '90'].map((period) => (
            <button
                key={period}
                onClick={() => onSelect(period)}
                className={`px-4 py-2 text-sm rounded-md transition-all duration-200 ${selected === period
                    ? 'bg-emerald-500 text-white font-medium'
                    : 'text-gray-400 hover:text-white hover:bg-gray-700/50'
                    }`}
            >
                {period} Days
            </button>
        ))}
    </div>
);

// Content Performance Timeline Chart
const ContentPerformanceChart = ({
    analytics,
    timeRange,
    onTimeRangeChange
}: {
    analytics: EnhancedAnalytics | null;
    timeRange: string;
    onTimeRangeChange: (period: string) => void;
}) => {
    // Prepare data for scatter plot
    const contentData = React.useMemo(() => {
        if (!analytics) return [];
        
        // Get all available content data from different sources
        const allContent = [
            ...(analytics.topVideos || []),
            ...(analytics.recentVideos || [])
        ];
        
        // If we don't have video arrays, but have video count in overview, show helpful message
        if (allContent.length === 0 && analytics.overview?.videoCount > 0) {
            // Instead of mock data, return empty array to show the empty state
            // The empty state will explain that video details are being processed
            return [];
        }
        
        // Remove duplicates and prepare scatter data
        const uniqueContent = allContent.reduce((acc: ContentDataPoint[], video) => {
            if (!acc.find(v => v.id === video.id)) {
                const publishDate = new Date(video.published_at);
                const daysSincePublish = Math.floor((Date.now() - publishDate.getTime()) / (1000 * 60 * 60 * 24));
                
                // Filter by time range
                const timeRangeNum = parseInt(timeRange);
                if (daysSincePublish <= timeRangeNum) {
                    acc.push({
                        id: video.id,
                        title: video.title,
                        views: video.view_count,
                        date: publishDate.getTime(),
                        displayDate: publishDate.toLocaleDateString(),
                        daysSince: daysSincePublish,
                        type: video.video_type.toLowerCase().includes('clip') ? 'clip' : 'broadcast'
                    });
                }
            }
            return acc;
        }, []);
        
        return uniqueContent.sort((a, b) => a.date - b.date);
    }, [analytics, timeRange]);

    const getContentColor = (type: string) => {
        return type === 'clip' ? '#10b981' : '#3b82f6'; // Green for clips, blue for broadcasts
    };

    const averageViews = contentData.length > 0 
        ? contentData.reduce((sum, item) => sum + item.views, 0) / contentData.length 
        : 0;

    return (
        <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            className="bg-gray-900/50 backdrop-blur-xl border border-gray-800/50 rounded-2xl p-6"
        >
            <div className="flex items-center justify-between mb-6">
                <div>
                    <h3 className="text-xl font-bold text-white">Content Performance Timeline</h3>
                    <p className="text-sm text-gray-400 mt-1">See which content resonates with your audience</p>
                </div>
                <TimePeriodSelector selected={timeRange} onSelect={onTimeRangeChange} />
            </div>
            
            {contentData.length > 0 ? (
                <>
                    <div className="h-80 mb-4">
                        <ResponsiveContainer width="100%" height="100%">
                            <ScatterChart data={contentData} margin={{ top: 20, right: 20, bottom: 20, left: 20 }}>
                                <CartesianGrid strokeDasharray="3 3" stroke="#374151" opacity={0.3} />
                                <XAxis
                                    type="number"
                                    dataKey="date"
                                    scale="time"
                                    domain={['dataMin', 'dataMax']}
                                    stroke="#9ca3af"
                                    fontSize={12}
                                    tickFormatter={(timestamp) => {
                                        const date = new Date(timestamp);
                                        return `${date.getMonth() + 1}/${date.getDate()}`;
                                    }}
                                />
                                <YAxis
                                    stroke="#9ca3af"
                                    fontSize={12}
                                    tickFormatter={(value) => formatNumber(value)}
                                />
                                <Tooltip
                                    contentStyle={{
                                        backgroundColor: '#1f2937',
                                        border: '1px solid #374151',
                                        borderRadius: '8px',
                                        color: '#fff'
                                    }}
                                    formatter={(value: number) => [
                                        `${formatNumber(value)} views`,
                                        'Content'
                                    ]}
                                    labelFormatter={(label: string, payload: any) => { // eslint-disable-line @typescript-eslint/no-explicit-any
                                        if (payload && payload[0] && payload[0].payload) {
                                            const data = payload[0].payload;
                                            return `${data.title} (${data.displayDate})`;
                                        }
                                        return label;
                                    }}
                                />
                                <Scatter dataKey="views" fill="#8884d8">
                                    {contentData.map((entry, index) => (
                                        <Cell key={`cell-${index}`} fill={getContentColor(entry.type)} />
                                    ))}
                                </Scatter>
                                {/* Average line */}
                                {averageViews > 0 && (
                                    <line
                                        x1="0"
                                        y1={`${100 - (averageViews / Math.max(...contentData.map(d => d.views))) * 100}%`}
                                        x2="100%"
                                        y2={`${100 - (averageViews / Math.max(...contentData.map(d => d.views))) * 100}%`}
                                        stroke="#f59e0b"
                                        strokeDasharray="5,5"
                                        strokeWidth={2}
                                        opacity={0.7}
                                    />
                                )}
                            </ScatterChart>
                        </ResponsiveContainer>
                    </div>
                    
                    {/* Legend and Insights */}
                    <div className="flex justify-between items-center">
                        <div className="flex items-center space-x-6">
                            <div className="flex items-center">
                                <div className="w-3 h-3 bg-blue-500 rounded-full mr-2"></div>
                                <span className="text-sm text-gray-400">Broadcasts</span>
                            </div>
                            <div className="flex items-center">
                                <div className="w-3 h-3 bg-emerald-500 rounded-full mr-2"></div>
                                <span className="text-sm text-gray-400">Clips</span>
                            </div>
                            <div className="flex items-center">
                                <div className="w-3 h-3 border-2 border-yellow-500 border-dashed rounded mr-2"></div>
                                <span className="text-sm text-gray-400">Average ({formatNumber(averageViews)} views)</span>
                            </div>
                        </div>
                        <div className="text-right">
                            <p className="text-sm text-gray-400">
                                {contentData.length} pieces of content in last {timeRange} days
                            </p>
                        </div>
                    </div>
                </>
            ) : (
                <div className="text-center py-12">
                    <Video className="w-12 h-12 text-gray-400 mx-auto mb-4" />
                    {analytics?.overview?.videoCount && analytics.overview.videoCount > 0 ? (
                        <>
                            <h4 className="text-lg font-medium text-gray-300 mb-2">Content Timeline Loading</h4>
                            <p className="text-gray-400">
                                We found {analytics.overview.videoCount} videos with {formatNumber(analytics.overview.totalViews)} total views, 
                                but detailed video data is still being processed.
                            </p>
                            <p className="text-gray-500 text-sm mt-2">
                                Try clicking "Update Data" to refresh your content details.
                            </p>
                        </>
                    ) : (
                        <>
                            <h4 className="text-lg font-medium text-gray-300 mb-2">No content in this period</h4>
                            <p className="text-gray-400">Try a longer time range or create more content!</p>
                        </>
                    )}
                </div>
            )}
        </motion.div>
    );
};

export default function AnalyticsPage() {
    const { isLoaded, isSignedIn, getToken } = useAuth();
    const [analytics, setAnalytics] = useState<EnhancedAnalytics | null>(null);
    const [loading, setLoading] = useState(true);
    const [refreshing, setRefreshing] = useState(false);
    const [timeRange, setTimeRange] = useState('30');

    // Helper function to get API base URL based on environment
    const getApiBaseUrl = () => {
        if (typeof window !== 'undefined' && window.location.hostname === 'dev.creatorsync.app') {
            return 'https://api-dev.creatorsync.app';
        } else if (process.env.NEXT_PUBLIC_APP_ENV === 'staging') {
            return 'https://api-dev.creatorsync.app';
        } else if (process.env.NODE_ENV === 'production') {
            return 'https://api.creatorsync.app';
        } else {
            return 'http://localhost:8080';
        }
    };

    // Fetch enhanced analytics data
    const fetchAnalyticsData = useCallback(async () => {
        if (!isLoaded || !isSignedIn) return;

        try {
            const token = await getToken();
            const apiBaseUrl = getApiBaseUrl();

            // Sync user to ensure they exist in the database
            try {
                await fetch(`${apiBaseUrl}/api/user/sync`, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json',
                    },
                });
            } catch {
                // Continue silently if sync fails
            }

            const response = await fetch(`${apiBaseUrl}/api/analytics/enhanced?days=${timeRange}`, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
            });

            if (response.ok) {
                const data = await response.json();
                setAnalytics(data);
            }
            
            setLoading(false);
        } catch {
            setLoading(false);
        }
    }, [isLoaded, isSignedIn, getToken, timeRange]);

    const handleRefresh = async () => {
        setRefreshing(true);
        await fetchAnalyticsData();
        setRefreshing(false);
    };

    const handleManualCollection = async () => {
        if (!isLoaded || !isSignedIn) return;

        try {
            const token = await getToken();
            const apiBaseUrl = getApiBaseUrl();
            
            const response = await fetch(`${apiBaseUrl}/api/analytics/collect`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
            });

            if (response.ok) {
                // Refresh data after a short delay to allow collection to start
                setTimeout(() => {
                    fetchAnalyticsData();
                }, 2000);
            }
        } catch {
            // Handle error silently
        }
    };

    const handleTimeRangeChange = (newTimeRange: string) => {
        setTimeRange(newTimeRange);
    };

    useEffect(() => {
        fetchAnalyticsData();
    }, [fetchAnalyticsData]);

    // Fetch data when time range changes
    useEffect(() => {
        if (isLoaded && isSignedIn) {
            fetchAnalyticsData();
        }
    }, [timeRange, isLoaded, isSignedIn, fetchAnalyticsData]);

    if (!isLoaded || !isSignedIn) {
        return (
            <div className="flex items-center justify-center min-h-screen">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-500"></div>
            </div>
        );
    }

    if (loading) {
        return (
            <div className="flex items-center justify-center min-h-screen">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-500"></div>
            </div>
        );
    }

    const safeOverview = {
        totalViews: analytics?.overview?.totalViews ?? 0,
        videoCount: analytics?.overview?.videoCount ?? 0,
        averageViewsPerVideo: analytics?.overview?.averageViewsPerVideo ?? 0,
        totalWatchTimeHours: analytics?.overview?.totalWatchTimeHours ?? 0,
        currentFollowers: analytics?.overview?.currentFollowers ?? 0,
        currentSubscribers: analytics?.overview?.currentSubscribers ?? 0,
        followerChange: analytics?.overview?.followerChange ?? 0,
        subscriberChange: analytics?.overview?.subscriberChange ?? 0
    };

    // Calculate some streamer-focused metrics
    const viewsPerFollower = safeOverview.currentFollowers > 0 ? safeOverview.totalViews / safeOverview.currentFollowers : 0;
    const engagementRate = safeOverview.totalViews / Math.max(safeOverview.videoCount, 1);

    return (
        <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 p-6">
            <div className="max-w-7xl mx-auto">
                {/* Header */}
                <div className="flex items-center justify-between mb-8">
                    <div>
                        <div className="flex items-center mb-2">
                            <BarChart3 className="w-6 h-6 text-emerald-500 mr-2" />
                            <span className="text-sm text-emerald-400 font-medium">Channel Analytics</span>
                        </div>
                        <h1 className="text-3xl font-bold text-white">Your Channel Performance</h1>
                        <p className="text-gray-400 mt-2">Track your content performance and understand what resonates with your audience.</p>
                    </div>
                    <div className="flex space-x-3">
                        <button
                            onClick={handleManualCollection}
                            className="flex items-center px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition-colors"
                        >
                            Update Data
                        </button>
                        <button
                            onClick={handleRefresh}
                            disabled={refreshing}
                            className="flex items-center px-4 py-2 bg-gray-700 text-white rounded-lg hover:bg-gray-600 transition-colors disabled:opacity-50"
                        >
                            <RefreshCw className={`w-4 h-4 mr-2 ${refreshing ? 'animate-spin' : ''}`} />
                            Refresh
                        </button>
                    </div>
                </div>

                {/* Key Metrics Grid */}
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                    <EnhancedMetricCard
                        title="Total Content"
                        value={safeOverview.videoCount}
                        subtitle="Videos & clips on your channel"
                        icon={Video}
                        gradient="from-purple-500/10 to-pink-500/5"
                    />
                    <EnhancedMetricCard
                        title="Total Views"
                        value={safeOverview.totalViews}
                        subtitle={`${formatNumber(engagementRate)} avg per video`}
                        icon={Eye}
                        gradient="from-emerald-500/10 to-teal-500/5"
                    />
                    <EnhancedMetricCard
                        title="Followers"
                        value={safeOverview.currentFollowers}
                        subtitle={viewsPerFollower > 0 ? `${formatNumber(viewsPerFollower)} views per follower` : "Growing your community"}
                        icon={Heart}
                        gradient="from-rose-500/10 to-pink-500/5"
                    />
                    <EnhancedMetricCard
                        title="Subscribers"
                        value={safeOverview.currentSubscribers}
                        subtitle="Supporting your channel"
                        icon={Users}
                        gradient="from-violet-500/10 to-purple-500/5"
                    />
                </div>

                {/* Content Performance Chart */}
                <div className="mb-8">
                    <ContentPerformanceChart analytics={analytics} timeRange={timeRange} onTimeRangeChange={handleTimeRangeChange} />
                </div>

                {/* Bottom Grid - Actionable Insights */}
                <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                    {/* Content Breakdown */}
                    <motion.div
                        initial={{ y: 20, opacity: 0 }}
                        animate={{ y: 0, opacity: 1 }}
                        className="bg-gray-900/50 backdrop-blur-xl border border-gray-800/50 rounded-2xl p-6"
                    >
                        <h3 className="text-lg font-semibold text-white mb-4 flex items-center">
                            <Video className="w-5 h-5 mr-2 text-emerald-400" />
                            Content Breakdown
                        </h3>
                        <div className="space-y-4">
                            <div>
                                <div className="flex justify-between items-center mb-2">
                                    <span className="text-gray-300">Average Views</span>
                                    <span className="text-white font-semibold">{formatNumber(safeOverview.averageViewsPerVideo)}</span>
                                </div>
                                <div className="w-full bg-gray-700 rounded-full h-2">
                                    <div 
                                        className="bg-emerald-500 h-2 rounded-full" 
                                        style={{ width: `${Math.min((safeOverview.averageViewsPerVideo / Math.max(safeOverview.totalViews / Math.max(safeOverview.videoCount, 1), 1)) * 100, 100)}%` }}
                                    ></div>
                                </div>
                            </div>
                            <div>
                                <div className="flex justify-between items-center mb-2">
                                    <span className="text-gray-300">Total Watch Time</span>
                                    <span className="text-white font-semibold">{formatDuration(safeOverview.totalWatchTimeHours)}</span>
                                </div>
                            </div>
                            <div className="pt-2 border-t border-gray-700">
                                <p className="text-sm text-gray-400">
                                    {safeOverview.videoCount > 0 ? 
                                        `You have ${safeOverview.videoCount} pieces of content generating ${formatNumber(safeOverview.totalViews)} total views.` :
                                        "Start creating content to see your analytics here!"
                                    }
                                </p>
                            </div>
                        </div>
                    </motion.div>

                    {/* Channel Growth Insights */}
                    <motion.div
                        initial={{ y: 20, opacity: 0 }}
                        animate={{ y: 0, opacity: 1 }}
                        className="bg-gray-900/50 backdrop-blur-xl border border-gray-800/50 rounded-2xl p-6"
                    >
                        <h3 className="text-lg font-semibold text-white mb-4 flex items-center">
                            <TrendingUp className="w-5 h-5 mr-2 text-emerald-400" />
                            Channel Health
                        </h3>
                        <div className="space-y-4">
                            <div className="flex justify-between items-center">
                                <span className="text-gray-300">Followers</span>
                                <span className="text-white font-semibold">{formatNumber(safeOverview.currentFollowers)}</span>
                            </div>
                            <div className="flex justify-between items-center">
                                <span className="text-gray-300">Subscribers</span>
                                <span className="text-white font-semibold">{formatNumber(safeOverview.currentSubscribers)}</span>
                            </div>
                            <div className="flex justify-between items-center">
                                <span className="text-gray-300">Sub Rate</span>
                                <span className="text-white font-semibold">
                                    {safeOverview.currentFollowers > 0 ? 
                                        `${((safeOverview.currentSubscribers / safeOverview.currentFollowers) * 100).toFixed(1)}%` : 
                                        '0%'
                                    }
                                </span>
                            </div>
                            <div className="pt-2 border-t border-gray-700">
                                <p className="text-sm text-gray-400">
                                    {safeOverview.currentSubscribers > 0 ? 
                                        `${safeOverview.currentSubscribers} of your ${formatNumber(safeOverview.currentFollowers)} followers are subscribers.` :
                                        "Encourage followers to subscribe for steady support!"
                                    }
                                </p>
                            </div>
                        </div>
                    </motion.div>

                    {/* Quick Actions & Tips */}
                    <motion.div
                        initial={{ y: 20, opacity: 0 }}
                        animate={{ y: 0, opacity: 1 }}
                        className="bg-gray-900/50 backdrop-blur-xl border border-gray-800/50 rounded-2xl p-6"
                    >
                        <h3 className="text-lg font-semibold text-white mb-4 flex items-center">
                            <BarChart3 className="w-5 h-5 mr-2 text-emerald-400" />
                            Quick Insights
                        </h3>
                        <div className="space-y-3">
                            {safeOverview.averageViewsPerVideo > 0 && (
                                <div className="p-3 bg-emerald-500/10 border border-emerald-500/20 rounded-lg">
                                    <p className="text-sm text-emerald-300">
                                        <strong>Great!</strong> Your content averages {formatNumber(safeOverview.averageViewsPerVideo)} views per video.
                                    </p>
                                </div>
                            )}
                            {safeOverview.videoCount > 10 && (
                                <div className="p-3 bg-blue-500/10 border border-blue-500/20 rounded-lg">
                                    <p className="text-sm text-blue-300">
                                        <strong>Consistent Creator:</strong> You have {safeOverview.videoCount} pieces of content!
                                    </p>
                                </div>
                            )}
                            {safeOverview.currentSubscribers === 0 && safeOverview.currentFollowers > 10 && (
                                <div className="p-3 bg-yellow-500/10 border border-yellow-500/20 rounded-lg">
                                    <p className="text-sm text-yellow-300">
                                        <strong>Tip:</strong> Consider encouraging subscriptions during streams!
                                    </p>
                                </div>
                            )}
                            {safeOverview.videoCount === 0 && (
                                <div className="p-3 bg-purple-500/10 border border-purple-500/20 rounded-lg">
                                    <p className="text-sm text-purple-300">
                                        <strong>Get Started:</strong> Create your first content to see analytics here!
                                    </p>
                                </div>
                            )}
                        </div>
                    </motion.div>
                </div>
            </div>
        </div>
    );
} 
