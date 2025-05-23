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
    BarChart as RechartsBarChart,
    Bar,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
    AreaChart,
    Area
} from 'recharts';

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
    title: string;
    viewCount: number;
    publishedAt: string;
}

interface EnhancedAnalytics {
    overview: VideoBasedOverview;
    performance: PerformanceData;
    topVideos: VideoAnalytics[];
    recentVideos: VideoAnalytics[];
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
                className={`px-4 py-2 text-sm rounded-md transition-all duration-200 ${
                    selected === period
                        ? 'bg-emerald-500 text-white font-medium'
                        : 'text-gray-400 hover:text-white hover:bg-gray-700/50'
                }`}
            >
                {period} Days
            </button>
        ))}
    </div>
);

// Performance Over Time Chart
const PerformanceChart = ({
    data,
    timeRange
}: {
    data: ChartDataPoint[];
    timeRange: string;
}) => (
    <motion.div
        initial={{ y: 20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        className="bg-gray-900/50 backdrop-blur-xl border border-gray-800/50 rounded-2xl p-6"
    >
        <div className="flex items-center justify-between mb-6">
            <h3 className="text-xl font-bold text-white">Performance Over Time</h3>
            <TimePeriodSelector selected={timeRange} onSelect={() => {}} />
        </div>
        <div className="h-80">
            <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={data}>
                    <defs>
                        <linearGradient id="viewsGradient" x1="0" y1="0" x2="0" y2="1">
                            <stop offset="5%" stopColor="#10b981" stopOpacity={0.3}/>
                            <stop offset="95%" stopColor="#10b981" stopOpacity={0}/>
                        </linearGradient>
                    </defs>
                    <CartesianGrid strokeDasharray="3 3" stroke="#374151" opacity={0.3} />
                    <XAxis 
                        dataKey="date" 
                        stroke="#9ca3af"
                        fontSize={12}
                        tickFormatter={(value) => {
                            const date = new Date(value);
                            return `${date.getMonth() + 1}/${date.getDate()}`;
                        }}
                    />
                    <YAxis stroke="#9ca3af" fontSize={12} />
                    <Tooltip
                        contentStyle={{
                            backgroundColor: '#1f2937',
                            border: '1px solid #374151',
                            borderRadius: '8px',
                            color: '#fff'
                        }}
                        formatter={(value) => [formatNumber(value as number), 'Views']}
                        labelFormatter={(label) => `Date: ${label}`}
                    />
                    <Area
                        type="monotone"
                        dataKey="value"
                        stroke="#10b981"
                        strokeWidth={2}
                        fill="url(#viewsGradient)"
                    />
                </AreaChart>
            </ResponsiveContainer>
        </div>
        <div className="mt-6 text-center text-emerald-400 text-sm">
            → views
        </div>
    </motion.div>
);

// Bottom stats row
const BottomStatsRow = ({
    totalViews,
    watchTime,
    avgViews
}: {
    totalViews: number | undefined | null;
    watchTime: number | undefined | null;
    avgViews: number | undefined | null;
}) => (
    <div className="grid grid-cols-3 gap-6">
        <div className="text-center">
            <p className="text-gray-400 text-sm mb-2">Views</p>
            <p className="text-2xl font-bold text-white">{formatNumber(totalViews)}</p>
        </div>
        <div className="text-center">
            <p className="text-gray-400 text-sm mb-2">Watch Time (est.)</p>
            <p className="text-2xl font-bold text-white">{formatDuration(watchTime)}</p>
        </div>
        <div className="text-center">
            <p className="text-gray-400 text-sm mb-2">Avg. Views/Video</p>
            <p className="text-2xl font-bold text-white">{formatNumber(avgViews)}</p>
        </div>
    </div>
);

// Growth metrics row
const GrowthStatsRow = ({
    followers,
    subscribers,
    followerChange,
    subscriberChange
}: {
    followers: number | undefined | null;
    subscribers: number | undefined | null;
    followerChange: number | undefined | null;
    subscriberChange: number | undefined | null;
}) => (
    <div className="grid grid-cols-2 gap-6">
        <div className="text-center">
            <p className="text-gray-400 text-sm mb-2">Followers</p>
            <p className="text-2xl font-bold text-white">{formatNumber(followers)}</p>
            {followerChange !== 0 && (
                <p className={`text-sm mt-1 ${followerChange && followerChange > 0 ? 'text-emerald-400' : 'text-red-400'}`}>
                    {followerChange && followerChange > 0 ? '+' : ''}{followerChange}
                </p>
            )}
        </div>
        <div className="text-center">
            <p className="text-gray-400 text-sm mb-2">Subscribers</p>
            <p className="text-2xl font-bold text-white">{formatNumber(subscribers)}</p>
            {subscriberChange !== 0 && (
                <p className={`text-sm mt-1 ${subscriberChange && subscriberChange > 0 ? 'text-emerald-400' : 'text-red-400'}`}>
                    {subscriberChange && subscriberChange > 0 ? '+' : ''}{subscriberChange}
                </p>
            )}
        </div>
    </div>
);

// Content Distribution Chart
const ContentDistributionChart = ({
    data
}: {
    data: ContentTypeData[];
}) => (
    <motion.div
        initial={{ y: 20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        className="bg-gray-900/50 backdrop-blur-xl border border-gray-800/50 rounded-2xl p-6"
    >
        <h3 className="text-xl font-bold text-white mb-6">Content Distribution</h3>
        <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
                <RechartsBarChart data={data}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#374151" opacity={0.3} />
                    <XAxis 
                        dataKey="date" 
                        stroke="#9ca3af"
                        fontSize={12}
                        tickFormatter={(value) => {
                            const date = new Date(value);
                            return `${date.getMonth() + 1}/${date.getDate()}`;
                        }}
                    />
                    <YAxis stroke="#9ca3af" fontSize={12} />
                    <Tooltip
                        contentStyle={{
                            backgroundColor: '#1f2937',
                            border: '1px solid #374151',
                            borderRadius: '8px',
                            color: '#fff'
                        }}
                    />
                    <Bar dataKey="broadcasts" stackId="a" fill="#3b82f6" name="Broadcasts" />
                    <Bar dataKey="clips" stackId="a" fill="#10b981" name="Clips" />
                </RechartsBarChart>
            </ResponsiveContainer>
        </div>
        <div className="flex justify-center mt-4 space-x-6">
            <div className="flex items-center">
                <div className="w-3 h-3 bg-blue-500 rounded mr-2"></div>
                <span className="text-sm text-gray-400">Broadcasts</span>
            </div>
            <div className="flex items-center">
                <div className="w-3 h-3 bg-emerald-500 rounded mr-2"></div>
                <span className="text-sm text-gray-400">Clips</span>
            </div>
        </div>
    </motion.div>
);

export default function AnalyticsPage() {
    const { isLoaded, isSignedIn, getToken } = useAuth();
    const [analytics, setAnalytics] = useState<EnhancedAnalytics | null>(null);
    const [loading, setLoading] = useState(true);
    const [refreshing, setRefreshing] = useState(false);
    const [timeRange] = useState('30');

    // Fetch enhanced analytics data
    const fetchAnalyticsData = useCallback(async () => {
        if (!isLoaded || !isSignedIn) return;

        try {
            const token = await getToken();
            const apiBaseUrl = process.env.NODE_ENV === 'production'
                ? 'https://api.creatorsync.app'
                : 'http://localhost:8080';

            const response = await fetch(`${apiBaseUrl}/api/analytics/enhanced?days=${timeRange}`, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
            });

            if (response.ok) {
                const data = await response.json();
                console.log('📊 Analytics data received:', data);
                console.log('📊 Overview data:', data.overview);
                setAnalytics(data);
            } else {
                console.error('Failed to fetch analytics:', response.status);
            }

        } catch (error) {
            console.error('Error fetching analytics data:', error);
        } finally {
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
            const apiBaseUrl = process.env.NODE_ENV === 'production'
                ? 'https://api.creatorsync.app'
                : 'http://localhost:8080';

            console.log('🔄 Triggering manual data collection...');
            const response = await fetch(`${apiBaseUrl}/api/analytics/collect`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
            });

            if (response.ok) {
                const result = await response.json();
                console.log('✅ Manual data collection triggered:', result);
                // Wait a bit then refresh the analytics
                setTimeout(async () => {
                    await fetchAnalyticsData();
                }, 3000);
            } else {
                console.error('❌ Failed to trigger manual collection:', response.status);
            }
        } catch (error) {
            console.error('❌ Error triggering manual collection:', error);
        }
    };



    useEffect(() => {
        fetchAnalyticsData();
    }, [fetchAnalyticsData]);

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

    const safePerformanceData = analytics?.performance?.viewsOverTime || [];
    const safeContentData = analytics?.performance?.contentDistribution || [];

    return (
        <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 p-6">
            <div className="max-w-7xl mx-auto">
                {/* Header */}
                <div className="flex items-center justify-between mb-8">
                    <div>
                        <div className="flex items-center mb-2">
                            <BarChart3 className="w-6 h-6 text-emerald-500 mr-2" />
                            <span className="text-sm text-emerald-400 font-medium">Performance Metrics</span>
                        </div>
                        <h1 className="text-3xl font-bold text-white">Detailed Analytics</h1>
                        <p className="text-gray-400 mt-2">Track your growth and understand your audience better with in-depth analytics.</p>
                    </div>
                    <div className="flex space-x-3">
                        <button
                            onClick={handleManualCollection}
                            className="flex items-center px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition-colors"
                        >
                            Collect Data
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

                {/* Top Metrics Cards */}
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6 mb-8">
                    <EnhancedMetricCard
                        title="Total Views"
                        value={safeOverview.totalViews}
                        subtitle={`Across ${safeOverview.videoCount} videos`}
                        icon={Eye}
                        gradient="from-emerald-500/10 to-teal-500/5"
                    />
                    <EnhancedMetricCard
                        title="Followers"
                        value={safeOverview.currentFollowers}
                        subtitle="Total followers"
                        icon={Heart}
                        gradient="from-rose-500/10 to-pink-500/5"
                    />
                    <EnhancedMetricCard
                        title="Subscribers"
                        value={safeOverview.currentSubscribers}
                        subtitle="Total subscribers"
                        icon={Users}
                        gradient="from-violet-500/10 to-purple-500/5"
                    />
                    <EnhancedMetricCard
                        title="Average Views"
                        value={Math.round(safeOverview.averageViewsPerVideo)}
                        subtitle="Per video"
                        icon={TrendingUp}
                        gradient="from-blue-500/10 to-indigo-500/5"
                    />
                    <EnhancedMetricCard
                        title="Content Count"
                        value={safeOverview.videoCount}
                        subtitle="Videos on your channel"
                        icon={Video}
                        gradient="from-purple-500/10 to-pink-500/5"
                    />
                </div>

                {/* Performance Over Time Chart */}
                <div className="mb-8">
                    <PerformanceChart data={safePerformanceData} timeRange={timeRange} />
                </div>

                {/* Bottom Stats and Content Distribution */}
                <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                    {/* Performance Stats */}
                    <motion.div
                        initial={{ y: 20, opacity: 0 }}
                        animate={{ y: 0, opacity: 1 }}
                        className="bg-gray-900/50 backdrop-blur-xl border border-gray-800/50 rounded-2xl p-6"
                    >
                        <h3 className="text-lg font-semibold text-white mb-4">Performance</h3>
                        <BottomStatsRow
                            totalViews={safeOverview.totalViews}
                            watchTime={safeOverview.totalWatchTimeHours}
                            avgViews={safeOverview.averageViewsPerVideo}
                        />
                    </motion.div>

                    {/* Growth Metrics */}
                    <motion.div
                        initial={{ y: 20, opacity: 0 }}
                        animate={{ y: 0, opacity: 1 }}
                        className="bg-gray-900/50 backdrop-blur-xl border border-gray-800/50 rounded-2xl p-6"
                    >
                        <h3 className="text-lg font-semibold text-white mb-4">Growth</h3>
                        <GrowthStatsRow
                            followers={safeOverview.currentFollowers}
                            subscribers={safeOverview.currentSubscribers}
                            followerChange={safeOverview.followerChange}
                            subscriberChange={safeOverview.subscriberChange}
                        />
                    </motion.div>

                    {/* Content Distribution */}
                    <ContentDistributionChart data={safeContentData} />
                </div>
            </div>
        </div>
    );
} 