"use client";

import { useState, useEffect, useCallback } from 'react';
import { motion } from 'framer-motion';
import { useAuth } from '@clerk/nextjs';
import { ConnectionModal } from '@/components/ui/connection-modal';
import {
    TrendingUp,
    Eye,
    Video,
    RefreshCw,
    BarChart3,
    Heart,
    Users
} from 'lucide-react';
// Recharts no longer needed for the content stream design
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

interface ConnectionStatus {
    twitch_connected: boolean;
    settings_url: string;
    account_switched?: boolean;
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
    duration_seconds?: number;
    duration_formatted?: string;
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

// Content Performance Stream - Modern Interactive Timeline
const ContentPerformanceChart = ({
    analytics,
    timeRange,
    onTimeRangeChange
}: {
    analytics: EnhancedAnalytics | null;
    timeRange: string;
    onTimeRangeChange: (period: string) => void;
}) => {
    // Prepare data for the content stream
    const contentData = React.useMemo(() => {
        if (!analytics) return [];

        const allContent = [
            ...(analytics.topVideos || []),
            ...(analytics.recentVideos || [])
        ];

        if (allContent.length === 0 && analytics.overview?.videoCount > 0) {
            return [];
        }

        // Create a Map to ensure truly unique content by video_id
        // Note: Backend already filters by date range, so no need to filter again here
        const uniqueContentMap = new Map<string, ContentDataPoint>();
        
        allContent.forEach(video => {
            const publishDate = new Date(video.published_at);
            const daysSincePublish = Math.floor((Date.now() - publishDate.getTime()) / (1000 * 60 * 60 * 24));

            const contentPoint: ContentDataPoint = {
                id: video.id,
                title: video.title,
                views: video.view_count,
                date: publishDate.getTime(),
                displayDate: publishDate.toLocaleDateString(),
                daysSince: daysSincePublish,
                type: video.video_type.toLowerCase().includes('clip') ? 'clip' : 'broadcast'
            };
            
            // Use video_id as key to prevent duplicates, keep the one with more views if duplicate
            const key = `${video.video_id}-${video.id}`;
            const existing = uniqueContentMap.get(key);
            if (!existing || contentPoint.views > existing.views) {
                uniqueContentMap.set(key, contentPoint);
            }
        });

        const uniqueContent = Array.from(uniqueContentMap.values());

        if (uniqueContent.length === 0 && allContent.length > 0) {
            const recentContentMap = new Map<string, ContentDataPoint>();
            
            allContent.slice(0, 10).forEach(video => {
                const publishDate = new Date(video.published_at);
                const daysSincePublish = Math.floor((Date.now() - publishDate.getTime()) / (1000 * 60 * 60 * 24));

                const contentPoint: ContentDataPoint = {
                    id: video.id,
                    title: video.title,
                    views: video.view_count,
                    date: publishDate.getTime(),
                    displayDate: publishDate.toLocaleDateString(),
                    daysSince: daysSincePublish,
                    type: video.video_type.toLowerCase().includes('clip') ? 'clip' : 'broadcast'
                };

                const key = `${video.video_id}-${video.id}`;
                const existing = recentContentMap.get(key);
                if (!existing || contentPoint.views > existing.views) {
                    recentContentMap.set(key, contentPoint);
                }
            });

            const recentContent = Array.from(recentContentMap.values());
            return recentContent.sort((a, b) => a.date - b.date);
        }

        return uniqueContent.sort((a, b) => a.date - b.date);
    }, [analytics, timeRange]);



    const maxViews = contentData.length > 0 ? Math.max(...contentData.map(d => d.views)) : 0;
    const averageViews = contentData.length > 0
        ? contentData.reduce((sum, item) => sum + item.views, 0) / contentData.length
        : 0;

    // Performance tiers for visual distinction
    const getPerformanceTier = (views: number) => {
        if (views >= averageViews * 1.5) return 'excellent';
        if (views >= averageViews) return 'good';
        if (views >= averageViews * 0.5) return 'average';
        return 'needs-work';
    };

    const getPerformanceColor = (tier: string) => {
        switch (tier) {
            case 'excellent': return 'from-emerald-400 to-green-600';
            case 'good': return 'from-blue-400 to-cyan-600';
            case 'average': return 'from-yellow-400 to-orange-500';
            case 'needs-work': return 'from-gray-400 to-gray-600';
            default: return 'from-gray-400 to-gray-600';
        }
    };

    return (
        <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            className="bg-gray-900/50 backdrop-blur-xl border border-gray-800/50 rounded-2xl p-6"
        >
            <div className="flex items-center justify-between mb-6">
                <div>
                    <h3 className="text-xl font-bold text-white">Content Performance Stream</h3>
                    <p className="text-sm text-gray-400 mt-1">Navigate your content journey • Track what resonates</p>
                </div>
                <TimePeriodSelector selected={timeRange} onSelect={onTimeRangeChange} />
            </div>

            {contentData.length > 0 ? (
                <>
                    {/* Content Stream Timeline */}
                    <div className="relative">
                        {/* Timeline River */}
                        <div className="absolute left-8 top-0 bottom-0 w-1 bg-gradient-to-b from-emerald-500/30 via-blue-500/30 to-purple-500/30 rounded-full"></div>
                        
                        {/* Content Bubbles */}
                        <div className="space-y-6 relative z-10">
                            {contentData.map((content, index) => {
                                const tier = getPerformanceTier(content.views);
                                const sizeMultiplier = Math.max(0.6, Math.min(1.4, content.views / maxViews * 1.2 + 0.4));
                                
                                return (
                                    <motion.div
                                        key={`content-${content.id}-${index}`}
                                        initial={{ x: -100, opacity: 0 }}
                                        animate={{ x: 0, opacity: 1 }}
                                        transition={{ delay: index * 0.1 }}
                                        className="flex items-center gap-6 group"
                                    >
                                        {/* Timeline Node */}
                                        <div className="relative flex-shrink-0">
                                            <div 
                                                className={`w-6 h-6 rounded-full bg-gradient-to-br ${getPerformanceColor(tier)} shadow-lg border-2 border-gray-800 transition-all duration-300 group-hover:scale-125 group-hover:shadow-xl`}
                                                style={{ transform: `scale(${sizeMultiplier})` }}
                                            >
                                                <div className="absolute inset-0 bg-gradient-to-br from-white/20 to-transparent rounded-full"></div>
                                            </div>
                                            
                                            {/* Performance indicator ring */}
                                            {tier === 'excellent' && (
                                                <div className="absolute -inset-2 border-2 border-emerald-400/30 rounded-full animate-pulse"></div>
                                            )}
                                        </div>

                                        {/* Content Card */}
                                        <motion.div
                                            whileHover={{ scale: 1.02, y: -2 }}
                                            className={`flex-1 bg-gradient-to-r ${
                                                content.type === 'clip' 
                                                    ? 'from-emerald-900/30 to-emerald-800/20 border-emerald-500/20' 
                                                    : 'from-blue-900/30 to-blue-800/20 border-blue-500/20'
                                            } border rounded-xl p-4 backdrop-blur-sm transition-all duration-300 group-hover:shadow-lg cursor-pointer`}
                                        >
                                            <div className="flex items-start justify-between">
                                                <div className="flex-1">
                                                    <div className="flex items-center gap-2 mb-2">
                                                        <span className={`px-2 py-1 text-xs font-medium rounded-full ${
                                                            content.type === 'clip' 
                                                                ? 'bg-emerald-500/20 text-emerald-300' 
                                                                : 'bg-blue-500/20 text-blue-300'
                                                        }`}>
                                                            {content.type === 'clip' ? '🎬 Clip' : '📺 Broadcast'}
                                                        </span>
                                                        <span className={`px-2 py-1 text-xs rounded-full font-medium ${
                                                            tier === 'excellent' ? 'bg-emerald-500/20 text-emerald-300' :
                                                            tier === 'good' ? 'bg-blue-500/20 text-blue-300' :
                                                            tier === 'average' ? 'bg-yellow-500/20 text-yellow-300' :
                                                            'bg-gray-500/20 text-gray-300'
                                                        }`}>
                                                            {tier === 'excellent' ? '🔥 Hot' :
                                                             tier === 'good' ? '✨ Good' :
                                                             tier === 'average' ? '📈 Average' :
                                                             '💭 Potential'}
                                                        </span>
                                                    </div>
                                                    
                                                    <h4 className="text-white font-semibold text-sm group-hover:text-emerald-300 transition-colors line-clamp-2">
                                                        {content.title}
                                                    </h4>
                                                    
                                                                                        <div className="flex items-center gap-4 mt-2 text-xs text-gray-400">
                                        <span>📅 {content.displayDate}</span>
                                        <span>⏰ {content.daysSince} days ago</span>
                                        {analytics && (
                                            (() => {
                                                const video = analytics.recentVideos.find(v => v.id === content.id) || 
                                                             analytics.topVideos.find(v => v.id === content.id);
                                                return video?.duration_formatted && (
                                                    <span>🕒 {video.duration_formatted}</span>
                                                );
                                            })()
                                        )}
                                    </div>
                                                </div>

                                                {/* Performance metrics */}
                                                <div className="text-right">
                                                    <div className="text-2xl font-bold text-white mb-1">
                                                        {formatNumber(content.views)}
                                                    </div>
                                                    <div className="text-xs text-gray-400">views</div>
                                                    
                                                    {/* Performance bar */}
                                                    <div className="w-16 h-1 bg-gray-700 rounded-full mt-2 overflow-hidden">
                                                        <div 
                                                            className={`h-full bg-gradient-to-r ${getPerformanceColor(tier)} transition-all duration-500`}
                                                            style={{ width: `${Math.min(100, (content.views / maxViews) * 100)}%` }}
                                                        ></div>
                                                    </div>
                                                </div>
                                            </div>
                                        </motion.div>
                                    </motion.div>
                                );
                            })}
                        </div>

                        {/* Stream Flow Animation */}
                        <div className="absolute left-8 top-0 bottom-0 w-1 opacity-50">
                            <div className="w-full h-8 bg-gradient-to-b from-emerald-400 to-transparent animate-pulse"></div>
                        </div>
                    </div>

                    {/* Performance Summary */}
                    <div className="mt-8 grid grid-cols-2 md:grid-cols-4 gap-4">
                        {[
                            { label: 'Total Content', value: contentData.length, icon: '📚', color: 'text-purple-400' },
                            { label: 'Avg Performance', value: formatNumber(averageViews), icon: '📊', color: 'text-blue-400' },
                            { label: 'Top Performer', value: formatNumber(maxViews), icon: '🏆', color: 'text-yellow-400' },
                            { 
                                label: 'Content Mix', 
                                value: `${contentData.filter(d => d.type === 'clip').length} clips, ${contentData.filter(d => d.type === 'broadcast').length} broadcasts`, 
                                icon: '🎬', 
                                color: 'text-emerald-400' 
                            }
                        ].map((stat, index) => (
                                                         <motion.div
                                key={`stat-${stat.label}-${index}`}
                                initial={{ y: 20, opacity: 0 }}
                                animate={{ y: 0, opacity: 1 }}
                                transition={{ delay: 0.5 + index * 0.1 }}
                                className="bg-gray-800/30 rounded-xl p-4 text-center hover:bg-gray-800/50 transition-all duration-300"
                            >
                                <div className={`text-2xl mb-2 ${stat.color}`}>{stat.icon}</div>
                                <div className="text-xl font-bold text-white">{stat.value}</div>
                                <div className="text-xs text-gray-400">{stat.label}</div>
                            </motion.div>
                        ))}
                    </div>


                </>
            ) : (
                <div className="text-center py-16">
                    <div className="w-16 h-16 bg-gradient-to-br from-emerald-500/20 to-blue-500/20 rounded-full flex items-center justify-center mx-auto mb-4">
                        <Video className="w-8 h-8 text-emerald-400" />
                    </div>
                    {analytics?.overview?.videoCount && analytics.overview.videoCount > 0 ? (
                        <>
                            <h4 className="text-lg font-medium text-gray-300 mb-2">Content Stream Loading</h4>
                            <p className="text-gray-400 max-w-md mx-auto">
                                We found {analytics.overview.videoCount} videos with {formatNumber(analytics.overview.totalViews)} total views.
                                Your content stream will appear here once processing is complete.
                            </p>
                            <p className="text-gray-500 text-sm mt-2">
                                Try clicking "Update Data" to refresh your content details.
                            </p>
                        </>
                    ) : (
                        <>
                            <h4 className="text-lg font-medium text-gray-300 mb-2">No content in this period</h4>
                            <p className="text-gray-400">Start creating content to see your performance stream!</p>
                        </>
                    )}
                </div>
            )}
        </motion.div>
    );
};

// Zero Data State Component for Analytics Page
const AnalyticsZeroDataState = ({ hasConnection }: { hasConnection: boolean }) => (
    <motion.div
        initial={{ y: 20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 p-6"
    >
        <div className="max-w-4xl mx-auto">
            <div className="flex items-center mb-2">
                <BarChart3 className="w-6 h-6 text-emerald-500 mr-2" />
                <span className="text-sm text-emerald-400 font-medium">Channel Analytics</span>
            </div>
            <h1 className="text-3xl font-bold text-white mb-6">
                {hasConnection ? "Welcome to Your Analytics Journey!" : "Connect Your Twitch Account"}
            </h1>

            <div className="bg-gradient-to-br from-emerald-500/10 to-blue-500/10 rounded-3xl p-12 border border-emerald-500/20 text-center">
                <div className="w-24 h-24 bg-gradient-to-br from-emerald-500 to-blue-500 rounded-full mx-auto mb-8 flex items-center justify-center">
                    {hasConnection ? (
                        <TrendingUp className="w-12 h-12 text-white" />
                    ) : (
                        <BarChart3 className="w-12 h-12 text-white" />
                    )}
                </div>
                
                <h2 className="text-4xl font-bold text-white mb-6">
                    {hasConnection ? "Your Creator Story Begins Here" : "Analytics Awaiting Connection"}
                </h2>
                
                <p className="text-xl text-gray-300 mb-8 max-w-3xl mx-auto leading-relaxed">
                    {hasConnection ? (
                        <>
                            Every great creator starts with their first piece of content. Your Twitch account is connected and ready – 
                            now it's time to <span className="text-emerald-400 font-semibold">create, stream, and grow</span> your community. 
                            As soon as you publish content, you'll see detailed analytics right here.
                        </>
                    ) : (
                        <>
                            Connect your Twitch account to unlock powerful analytics, track your channel's performance, 
                            and get insights to grow your audience.
                        </>
                    )}
                </p>

                {hasConnection ? (
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
                        <div className="bg-gradient-to-br from-purple-500/20 to-pink-500/10 border border-purple-500/30 rounded-2xl p-6">
                            <Video className="w-10 h-10 text-purple-400 mx-auto mb-4" />
                            <h3 className="text-white font-bold mb-2">Stream Content</h3>
                            <p className="text-gray-400 text-sm">Go live and engage with your audience in real-time</p>
                        </div>
                        <div className="bg-gradient-to-br from-emerald-500/20 to-teal-500/10 border border-emerald-500/30 rounded-2xl p-6">
                            <TrendingUp className="w-10 h-10 text-emerald-400 mx-auto mb-4" />
                            <h3 className="text-white font-bold mb-2">Track Growth</h3>
                            <p className="text-gray-400 text-sm">Monitor followers, views, and engagement metrics</p>
                        </div>
                        <div className="bg-gradient-to-br from-blue-500/20 to-cyan-500/10 border border-blue-500/30 rounded-2xl p-6">
                            <Eye className="w-10 h-10 text-blue-400 mx-auto mb-4" />
                            <h3 className="text-white font-bold mb-2">Analyze Performance</h3>
                            <p className="text-gray-400 text-sm">Understand what content resonates with viewers</p>
                        </div>
                        <div className="bg-gradient-to-br from-rose-500/20 to-orange-500/10 border border-rose-500/30 rounded-2xl p-6">
                            <Heart className="w-10 h-10 text-rose-400 mx-auto mb-4" />
                            <h3 className="text-white font-bold mb-2">Build Community</h3>
                            <p className="text-gray-400 text-sm">Grow a loyal following around your content</p>
                        </div>
                    </div>
                ) : (
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-12">
                        <div className="bg-emerald-500/10 border border-emerald-500/20 rounded-xl p-8">
                            <TrendingUp className="w-12 h-12 text-emerald-500 mx-auto mb-4" />
                            <h3 className="text-white font-semibold mb-3">Performance Tracking</h3>
                            <p className="text-gray-400">Monitor views, followers, and subscriber growth over time.</p>
                        </div>
                        <div className="bg-purple-500/10 border border-purple-500/20 rounded-xl p-8">
                            <Video className="w-12 h-12 text-purple-500 mx-auto mb-4" />
                            <h3 className="text-white font-semibold mb-3">Content Analysis</h3>
                            <p className="text-gray-400">Understand which content performs best with your audience.</p>
                        </div>
                        <div className="bg-blue-500/10 border border-blue-500/20 rounded-xl p-8">
                            <Eye className="w-12 h-12 text-blue-500 mx-auto mb-4" />
                            <h3 className="text-white font-semibold mb-3">Audience Insights</h3>
                            <p className="text-gray-400">Discover engagement patterns and optimize your streaming strategy.</p>
                        </div>
                    </div>
                )}

                <div className="flex flex-col sm:flex-row gap-4 justify-center">
                    {hasConnection ? (
                        <>
                            <button
                                onClick={() => window.open('https://www.twitch.tv/broadcast/dashboard/go-live', '_blank')}
                                className="bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white px-8 py-4 rounded-xl font-semibold transition-all duration-200 shadow-lg hover:shadow-purple-500/25"
                            >
                                Start Streaming Now
                            </button>
                            <button
                                onClick={() => window.open('https://creator.twitch.tv/', '_blank')}
                                className="bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-700 hover:to-teal-700 text-white px-8 py-4 rounded-xl font-semibold transition-all duration-200 shadow-lg hover:shadow-emerald-500/25"
                            >
                                Creator Resources
                            </button>
                        </>
                    ) : (
                        <button
                            onClick={() => window.open('/settings', '_self')}
                            className="bg-gradient-to-r from-emerald-600 to-emerald-700 hover:from-emerald-700 hover:to-emerald-800 text-white px-10 py-4 rounded-xl font-semibold transition-all duration-200 shadow-lg hover:shadow-emerald-500/25"
                        >
                            Connect Twitch Account
                        </button>
                    )}
                </div>

                {hasConnection && (
                    <div className="mt-8 p-6 bg-gray-800/50 rounded-2xl border border-gray-700/50">
                        <p className="text-gray-400 text-sm">
                            <strong className="text-emerald-400">Pro Tip:</strong> Your analytics will appear here automatically after you create content. 
                            We'll track your streams, highlights, and clips to give you comprehensive insights into your channel's performance.
                        </p>
                    </div>
                )}
            </div>
        </div>
    </motion.div>
);

export default function AnalyticsPage() {
    const { isLoaded, isSignedIn, getToken } = useAuth();
    const [analytics, setAnalytics] = useState<EnhancedAnalytics | null>(null);
    const [connectionStatus, setConnectionStatus] = useState<ConnectionStatus | null>(null);
    const [loading, setLoading] = useState(true);
    const [refreshing, setRefreshing] = useState(false);
    const [timeRange, setTimeRange] = useState('7');
    const [showConnectionModal, setShowConnectionModal] = useState(false);

    // Add request deduplication
    const [isRequestInProgress, setIsRequestInProgress] = useState(false);
    const [lastRequestId, setLastRequestId] = useState<string | null>(null);



    // Helper function to get API base URL based on environment
    const getApiBaseUrl = () => {
        const hostname = typeof window !== 'undefined' ? window.location.hostname : '';
        const nodeEnv = process.env.NODE_ENV;
        const appEnv = process.env.NEXT_PUBLIC_APP_ENV;

        let apiUrl = '';
        if (hostname === 'dev.creatorsync.app') {
            apiUrl = 'https://api-dev.creatorsync.app';
        } else if (appEnv === 'staging') {
            apiUrl = 'https://api-dev.creatorsync.app';
        } else if (nodeEnv === 'production') {
            apiUrl = 'https://api.creatorsync.app';
        } else {
            apiUrl = 'http://localhost:8080';
        }

        return apiUrl;
    };

    // Test API connectivity
    const testApiConnectivity = async (apiBaseUrl: string) => {
        try {
            const response = await fetch(`${apiBaseUrl}/api/analytics/health`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });

            if (response.ok) {
                return true;
            } else {
                console.error('API health check failed with status:', response.status);
                return false;
            }
        } catch (error) {
            console.error('API connectivity test failed:', error);
            return false;
        }
    };

    // Check connection status separately
    const checkConnectionStatus = useCallback(async () => {
        if (!isLoaded || !isSignedIn) return;

        try {
            const token = await getToken();
            const apiBaseUrl = getApiBaseUrl();

            const response = await fetch(`${apiBaseUrl}/api/analytics/connection-status`, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
            });

            if (response.ok) {
                const data = await response.json();
                const isConnected = data.platforms?.twitch?.connected || false;
                setConnectionStatus({
                    twitch_connected: isConnected,
                    settings_url: data.settings_url || '/settings'
                });
                return isConnected;
            } else {
                console.error('Failed to check connection status:', response.status);
                // Default to disconnected if we can't check
                setConnectionStatus({
                    twitch_connected: false,
                    settings_url: '/settings'
                });
                return false;
            }
        } catch (error) {
            console.error('Connection status check failed:', error);
            // Default to disconnected if we can't check
            setConnectionStatus({
                twitch_connected: false,
                settings_url: '/settings'
            });
            return false;
        }
    }, [isLoaded, isSignedIn, getToken]);

    // Fetch enhanced analytics data with request deduplication
    const fetchAnalyticsData = useCallback(async (forceRefresh = false) => {
        if (!isLoaded || !isSignedIn) return;

        // Always check connection status first
        const isConnected = await checkConnectionStatus();

        // If not connected, stop here - the UI will show the connection prompt
        if (!isConnected) {
            setLoading(false);
            return;
        }

        // Prevent concurrent requests unless forced refresh
        if (isRequestInProgress && !forceRefresh) {
            return;
        }

        const requestId = `${Date.now()}-${Math.random()}`;

        // Check if this is a duplicate request within 1 second
        if (lastRequestId && !forceRefresh) {
            const timeSinceLastRequest = Date.now() - parseInt(lastRequestId.split('-')[0]);
            if (timeSinceLastRequest < 1000) {
                return;
            }
        }

        setIsRequestInProgress(true);
        setLastRequestId(requestId);

        try {
            const token = await getToken();
            const apiBaseUrl = getApiBaseUrl();

            // Test API connectivity first
            const isApiHealthy = await testApiConnectivity(apiBaseUrl);
            if (!isApiHealthy) {
                console.error('API connectivity test failed - server may be down');
                throw new Error('API server is not accessible. Please check if the server is running.');
            }

            // Sync user to ensure they exist in the database with improved error handling
            try {
                const syncResponse = await fetch(`${apiBaseUrl}/api/user/sync`, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json',
                    },
                });

                if (syncResponse.ok) {
                    const syncData = await syncResponse.json();
                    if (syncData.retry_needed) {
                    }
                } else {
                    console.warn('User sync had issues, continuing with analytics fetch...');
                }
            } catch (syncError) {
                console.warn('User sync failed, continuing with analytics fetch...', syncError);
                // Continue anyway - analytics might still work if user already exists
            }

            // Add retry logic for network failures
            let lastError = null;
            const maxRetries = 3;

            for (let attempt = 1; attempt <= maxRetries; attempt++) {
                try {
                    const controller = new AbortController();
                    const timeoutId = setTimeout(() => controller.abort(), 30000); // 30 second timeout

                    const response = await fetch(`${apiBaseUrl}/api/analytics/enhanced?days=${timeRange}&requestId=${requestId}`, {
                        headers: {
                            'Authorization': `Bearer ${token}`,
                            'Content-Type': 'application/json',
                        },
                        signal: controller.signal,
                    });

                    clearTimeout(timeoutId);

                    if (response.ok) {
                        const data = await response.json();

                        // Handle both old and new response formats
                        if (data.analytics && data.connection_status) {
                            // New format with connection status
                            setAnalytics(data.analytics);
                            setConnectionStatus(data.connection_status);
                        } else {
                            // Legacy format - assume connected if we got data
                            setAnalytics(data as EnhancedAnalytics);
                            setConnectionStatus({ twitch_connected: true, settings_url: '/settings' });
                        }

                        setLoading(false);
                        return; // Success - exit retry loop
                    } else if (response.status === 429) {
                        // Handle duplicate request gracefully - don't show error to user
                        // The concurrent request should complete and provide the data
                        return;
                    } else {
                        const errorText = await response.text();
                        const error = new Error(`HTTP ${response.status}: ${errorText}`);
                        console.error(`Analytics request ${requestId} failed with status:`, response.status, errorText);
                        lastError = error;

                        // Don't retry client errors (4xx), only server errors (5xx) and network issues
                        if (response.status >= 400 && response.status < 500) {
                            throw error; // Don't retry client errors
                        }
                    }
                } catch (fetchError) {
                    console.error(`Analytics fetch attempt ${attempt} failed for request ${requestId}:`, fetchError);
                    lastError = fetchError;

                    // If this is an AbortError (timeout), don't retry
                    if (fetchError instanceof Error && fetchError.name === 'AbortError') {
                        throw new Error('Request timed out after 30 seconds');
                    }

                    // Wait before retrying (exponential backoff)
                    if (attempt < maxRetries) {
                        const waitTime = Math.min(1000 * Math.pow(2, attempt - 1), 5000); // Cap at 5 seconds
                        await new Promise(resolve => setTimeout(resolve, waitTime));
                    }
                }
            }

            // If we get here, all retries failed
            throw lastError || new Error('All retry attempts failed');

        } catch (error) {
            console.error(`Analytics request ${requestId} failed with error:`, error);

            // Provide user-friendly error messages
            let userMessage = 'Failed to load analytics data.';
            if (error instanceof Error) {
                if (error.message.includes('Failed to fetch') || error.message.includes('not accessible')) {
                    userMessage = 'Cannot connect to server. Please check your internet connection and try again.';
                } else if (error.message.includes('timed out')) {
                    userMessage = 'Request timed out. The server may be experiencing high load.';
                } else if (error.message.includes('HTTP 5')) {
                    userMessage = 'Server error. Please try again in a few moments.';
                }
            }

            // Show error to user (you can implement a toast/notification system)
            console.error('User-friendly error:', userMessage);

            setLoading(false);
        } finally {
            setIsRequestInProgress(false);
        }
    }, [isLoaded, isSignedIn, getToken, timeRange, isRequestInProgress, lastRequestId, checkConnectionStatus]);

    const handleRefresh = async () => {
        if (refreshing || isRequestInProgress) return;

        setRefreshing(true);
        await fetchAnalyticsData(true); // Force refresh
        setRefreshing(false);
    };

    const handleManualCollection = async () => {
        if (!isLoaded || !isSignedIn || isRequestInProgress) return;

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
                    fetchAnalyticsData(true);
                }, 2000);
            }
        } catch {
            // Handle error silently
        }
    };

    const handleTimeRangeChange = (newTimeRange: string) => {
        if (newTimeRange === timeRange || isRequestInProgress) return;

        setTimeRange(newTimeRange);
    };

    // Initial load effect - only run once when user is loaded
    useEffect(() => {
        if (isLoaded && isSignedIn && !isRequestInProgress) {
            fetchAnalyticsData();
        }
    }, [isLoaded, isSignedIn]); // Only depend on auth state

    // Time range change effect - with debouncing
    useEffect(() => {
        if (isLoaded && isSignedIn && analytics) { // Only if we have initial data
            const timeoutId = setTimeout(() => {
                fetchAnalyticsData();
            }, 300); // 300ms debounce

            return () => clearTimeout(timeoutId);
        }
    }, [timeRange]); // Only depend on timeRange

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

    // Show connection prompt if Twitch is not connected
    // Also show if we have no analytics data and connectionStatus indicates disconnected
    if ((connectionStatus && !connectionStatus.twitch_connected) ||
        (connectionStatus === null && !analytics && !loading)) {
        return <AnalyticsZeroDataState hasConnection={false} />;
    }

    // Check if we should show zero data state for connected users
    const isConnected = connectionStatus?.twitch_connected === true;
    const accountSwitched = connectionStatus?.account_switched === true;
    const hasZeroData = isConnected && analytics && (
        (analytics.overview?.totalViews === 0 && 
         analytics.overview?.videoCount === 0 && 
         analytics.overview?.currentFollowers === 0) ||
        (!analytics.topVideos || analytics.topVideos.length === 0) &&
        (!analytics.recentVideos || analytics.recentVideos.length === 0)
    );

    // Show encouraging zero data state for connected users with no content or account switches
    if (hasZeroData || accountSwitched) {
        return (
            <>
                <AnalyticsZeroDataState hasConnection={true} />
                {/* Still include the connection modal for settings */}
                <ConnectionModal
                    isOpen={showConnectionModal}
                    onClose={() => setShowConnectionModal(false)}
                    platform="twitch"
                    getToken={getToken}
                />
            </>
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
                        {connectionStatus?.twitch_connected ? (
                            <>
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
                            </>
                        ) : (
                            <button
                                onClick={() => setShowConnectionModal(true)}
                                className="flex items-center px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition-colors"
                            >
                                Connect Twitch Account
                            </button>
                        )}
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
                                <div className="w-full bg-gray-700 rounded-full h-2">
                                    <div
                                        className="bg-blue-500 h-2 rounded-full"
                                        style={{ width: `${Math.min((safeOverview.totalWatchTimeHours / Math.max(safeOverview.totalWatchTimeHours * 1.2, 1)) * 100, 100)}%` }}
                                    ></div>
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
