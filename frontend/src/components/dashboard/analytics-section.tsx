'use client';

import { useState, useEffect, useCallback } from 'react';
import { motion } from 'framer-motion';
import { useAuth } from '@clerk/nextjs';
import { BarChart3, TrendingUp, Users, Eye, Clock, Calendar } from 'lucide-react';

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

export function AnalyticsSection() {
    const { isLoaded, isSignedIn, getToken } = useAuth();
    const [overview, setOverview] = useState<DashboardOverview | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchAnalyticsOverview = useCallback(async () => {
        if (!isLoaded || !isSignedIn) return;

        try {
            const token = await getToken();
            const apiBaseUrl = process.env.NODE_ENV === 'production'
                ? 'https://api.creatorsync.app'
                : 'http://localhost:8080';

            const response = await fetch(`${apiBaseUrl}/api/analytics/overview`, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
            });

            if (response.ok) {
                const data = await response.json();
                setOverview(data);
            } else {
                console.error('Failed to fetch analytics overview:', response.status);
                setError('Failed to load analytics data');
            }
        } catch (error) {
            console.error('Error fetching analytics overview:', error);
            setError('Failed to load analytics data');
        } finally {
            setIsLoading(false);
        }
    }, [isLoaded, isSignedIn, getToken]);

    useEffect(() => {
        fetchAnalyticsOverview();
    }, [fetchAnalyticsOverview]);

    const formatNumber = (num: number | undefined | null): string => {
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
        if (hours === undefined || hours === null || isNaN(hours)) {
            return '0 min';
        }
        
        if (hours >= 1) {
            return `${Math.round(hours)} hrs`;
        }
        return `${Math.round(hours * 60)} min`;
    };

    if (!isLoaded || !isSignedIn) {
        return (
            <div className="p-8">
                <div className="flex items-center justify-center h-64">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-brand-500"></div>
                </div>
            </div>
        );
    }

    if (isLoading) {
        return (
            <div className="p-8">
                <div className="flex items-center justify-center h-64">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-brand-500"></div>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="p-8">
                <div className="bg-red-500/10 border border-red-500/20 rounded-xl p-6">
                    <h3 className="text-red-500 font-semibold mb-2">Error Loading Analytics</h3>
                    <p className="text-red-500/80">{error}</p>
                </div>
            </div>
        );
    }

    const safeOverview = {
        totalViews: overview?.totalViews ?? 0,
        videoCount: overview?.videoCount ?? 0,
        averageViewsPerVideo: overview?.averageViewsPerVideo ?? 0,
        totalWatchTimeHours: overview?.totalWatchTimeHours ?? 0,
        currentFollowers: overview?.currentFollowers ?? 0,
        currentSubscribers: overview?.currentSubscribers ?? 0,
        followerChange: overview?.followerChange ?? 0,
        subscriberChange: overview?.subscriberChange ?? 0
    };
    return (
        <div className="p-8">
            <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-blue-500/20 border border-blue-500/30 text-blue-500 text-sm mb-6 w-fit">
                <BarChart3 className="w-4 h-4" />
                <span>Performance Metrics</span>
            </div>

            <h2 className="text-4xl font-bold mb-6">
                <span className="text-light-surface-900 dark:text-dark-surface-100">Detailed</span>{' '}
                <span className="text-gradient">Analytics</span>
            </h2>

            <p className="text-xl text-light-surface-700 dark:text-dark-surface-300 max-w-3xl mb-10">
                Track your growth and understand your audience better with in-depth analytics.
            </p>

            {/* Main Stats Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-10">
                <motion.div
                    initial={{ y: 20, opacity: 0 }}
                    animate={{ y: 0, opacity: 1 }}
                    transition={{ duration: 0.4, delay: 0.1 }}
                    className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm p-6 rounded-xl shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50"
                >
                    <div className="flex items-center gap-3 mb-4">
                        <Eye className="w-6 h-6 text-brand-500" />
                        <h3 className="text-lg font-semibold text-light-surface-900 dark:text-dark-surface-100">Total Views</h3>
                    </div>
                    <div className="mb-2">
                        <span className="text-4xl font-bold text-light-surface-900 dark:text-dark-surface-100">{formatNumber(safeOverview.totalViews)}</span>
                        <span className="text-light-surface-600 dark:text-dark-surface-400 ml-2">Across {safeOverview.videoCount} videos</span>
                    </div>
                </motion.div>

                <motion.div
                    initial={{ y: 20, opacity: 0 }}
                    animate={{ y: 0, opacity: 1 }}
                    transition={{ duration: 0.4, delay: 0.2 }}
                    className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm p-6 rounded-xl shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50"
                >
                    <div className="flex items-center gap-3 mb-4">
                        <TrendingUp className="w-6 h-6 text-brand-500" />
                        <h3 className="text-lg font-semibold text-light-surface-900 dark:text-dark-surface-100">Average Views</h3>
                    </div>
                    <div className="mb-2">
                        <span className="text-4xl font-bold text-light-surface-900 dark:text-dark-surface-100">{formatNumber(safeOverview.averageViewsPerVideo)}</span>
                        <span className="text-light-surface-600 dark:text-dark-surface-400 ml-2">Per video</span>
                    </div>
                </motion.div>
            </div>

            {/* Performance Over Time Chart */}
            <motion.div
                initial={{ y: 20, opacity: 0 }}
                animate={{ y: 0, opacity: 1 }}
                transition={{ duration: 0.4, delay: 0.3 }}
                className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm p-6 rounded-xl shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50 mb-10"
            >
                <h3 className="text-xl font-semibold text-light-surface-900 dark:text-dark-surface-100 mb-6">Performance Over Time</h3>
                
                <div className="relative h-64 bg-gradient-to-br from-brand-500/5 to-brand-600/10 rounded-lg p-4">
                    <div className="absolute inset-4 border-l border-b border-light-surface-300/50 dark:border-dark-surface-700/50">
                        {/* Y-axis labels */}
                        <div className="absolute -left-8 top-0 text-xs text-light-surface-600 dark:text-dark-surface-400">100</div>
                        <div className="absolute -left-8 top-1/4 text-xs text-light-surface-600 dark:text-dark-surface-400">75</div>
                        <div className="absolute -left-8 top-2/4 text-xs text-light-surface-600 dark:text-dark-surface-400">50</div>
                        <div className="absolute -left-8 top-3/4 text-xs text-light-surface-600 dark:text-dark-surface-400">25</div>
                        <div className="absolute -left-8 bottom-0 text-xs text-light-surface-600 dark:text-dark-surface-400">0</div>
                        
                        {/* X-axis labels */}
                        <div className="absolute -bottom-6 left-0 text-xs text-light-surface-600 dark:text-dark-surface-400">4/25</div>
                        <div className="absolute -bottom-6 left-1/2 -translate-x-1/2 text-xs text-light-surface-600 dark:text-dark-surface-400">4/27</div>
                        <div className="absolute -bottom-6 right-0 text-xs text-light-surface-600 dark:text-dark-surface-400">5/2</div>
                        
                        {/* Simple line chart visualization */}
                        <svg className="w-full h-full">
                            <defs>
                                <linearGradient id="chartGradient" x1="0%" y1="0%" x2="0%" y2="100%">
                                    <stop offset="0%" stopColor="rgb(99, 102, 241)" stopOpacity="0.3"/>
                                    <stop offset="100%" stopColor="rgb(99, 102, 241)" stopOpacity="0"/>
                                </linearGradient>
                            </defs>
                            <path
                                d="M 0,120 Q 50,80 100,60 T 200,40 T 300,45 L 300,160 L 0,160 Z"
                                fill="url(#chartGradient)"
                                className="animate-pulse"
                            />
                            <path
                                d="M 0,120 Q 50,80 100,60 T 200,40 T 300,45"
                                stroke="rgb(99, 102, 241)"
                                strokeWidth="2"
                                fill="none"
                                className="drop-shadow-sm"
                            />
                        </svg>
                    </div>
                </div>
            </motion.div>

            {/* Detailed Metrics */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-10">
                <motion.div
                    initial={{ y: 20, opacity: 0 }}
                    animate={{ y: 0, opacity: 1 }}
                    transition={{ duration: 0.4, delay: 0.4 }}
                    className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm p-6 rounded-xl shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50 text-center"
                >
                    <Eye className="w-8 h-8 text-brand-500 mx-auto mb-3" />
                    <h4 className="text-sm font-medium text-light-surface-700 dark:text-dark-surface-300 mb-1">Views</h4>
                    <p className="text-2xl font-bold text-light-surface-900 dark:text-dark-surface-100">{formatNumber(safeOverview.totalViews)}</p>
                </motion.div>

                <motion.div
                    initial={{ y: 20, opacity: 0 }}
                    animate={{ y: 0, opacity: 1 }}
                    transition={{ duration: 0.4, delay: 0.5 }}
                    className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm p-6 rounded-xl shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50 text-center"
                >
                    <Clock className="w-8 h-8 text-brand-500 mx-auto mb-3" />
                    <h4 className="text-sm font-medium text-light-surface-700 dark:text-dark-surface-300 mb-1">Watch Time (est.)</h4>
                    <p className="text-2xl font-bold text-light-surface-900 dark:text-dark-surface-100">{formatDuration(safeOverview.totalWatchTimeHours)}</p>
                </motion.div>

                <motion.div
                    initial={{ y: 20, opacity: 0 }}
                    animate={{ y: 0, opacity: 1 }}
                    transition={{ duration: 0.4, delay: 0.6 }}
                    className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm p-6 rounded-xl shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50 text-center"
                >
                    <Users className="w-8 h-8 text-brand-500 mx-auto mb-3" />
                    <h4 className="text-sm font-medium text-light-surface-700 dark:text-dark-surface-300 mb-1">Avg. Views/Video</h4>
                    <p className="text-2xl font-bold text-light-surface-900 dark:text-dark-surface-100">{formatNumber(safeOverview.averageViewsPerVideo)}</p>
                </motion.div>
            </div>

            {/* Coming Soon Features */}
            <motion.div
                initial={{ y: 20, opacity: 0 }}
                animate={{ y: 0, opacity: 1 }}
                transition={{ duration: 0.4, delay: 0.7 }}
                className="bg-gradient-to-br from-brand-500/10 to-brand-600/5 border border-brand-500/20 p-6 rounded-xl"
            >
                <div className="flex items-center gap-2 text-brand-500 mb-4">
                    <Calendar className="w-5 h-5" />
                    <h3 className="font-semibold">Coming Soon</h3>
                </div>
                <p className="text-light-surface-700 dark:text-dark-surface-300 mb-4">
                    We're working on advanced analytics features including audience demographics, engagement rates, and revenue tracking.
                </p>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="flex items-center gap-2 text-sm text-light-surface-600 dark:text-dark-surface-400">
                        <div className="w-2 h-2 bg-brand-500 rounded-full"></div>
                        Audience Demographics
                    </div>
                    <div className="flex items-center gap-2 text-sm text-light-surface-600 dark:text-dark-surface-400">
                        <div className="w-2 h-2 bg-brand-500 rounded-full"></div>
                        Engagement Rates
                    </div>
                    <div className="flex items-center gap-2 text-sm text-light-surface-600 dark:text-dark-surface-400">
                        <div className="w-2 h-2 bg-brand-500 rounded-full"></div>
                        Revenue Tracking
                    </div>
                    <div className="flex items-center gap-2 text-sm text-light-surface-600 dark:text-dark-surface-400">
                        <div className="w-2 h-2 bg-brand-500 rounded-full"></div>
                        Growth Predictions
                    </div>
                </div>
            </motion.div>
        </div>
    );
} 