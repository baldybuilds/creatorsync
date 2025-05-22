"use client";

import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { SignOutButton, useAuth } from '@clerk/nextjs';
import { Button } from '../../components/ui/button';
import { Users, BarChart2, Video, Calendar, LogOut, Settings as SettingsIcon, Sparkles, ArrowRight, Plus, Play } from 'lucide-react';
import { cn } from '@/lib/utils';
import Image from 'next/image';

// Dashboard components with modern styling
const OverviewContent = () => (
    <div>
        <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-brand-500/20 border border-brand-500/30 text-brand-500 text-sm mb-6 w-fit">
            <Sparkles className="w-4 h-4" />
            <span>Dashboard Overview</span>
        </div>

        <h2 className="text-4xl font-bold mb-6">
            <span className="text-light-surface-900 dark:text-dark-surface-100">Welcome back,</span>{' '}
            <span className="text-gradient">Creator</span>
        </h2>

        <p className="text-xl text-light-surface-700 dark:text-dark-surface-300 max-w-3xl mb-10">
            Here's an overview of your content performance and recent activity.
        </p>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-10">
            <motion.div
                initial={{ y: 20, opacity: 0 }}
                animate={{ y: 0, opacity: 1 }}
                transition={{ duration: 0.4, delay: 0.1 }}
                className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm p-6 rounded-xl shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50 relative overflow-hidden group hover:border-brand-500/30 transition-all duration-300"
            >
                <div className="absolute top-0 right-0 w-32 h-32 bg-gradient-to-br from-blue-500/10 to-indigo-500/5 rounded-full -mr-10 -mt-10 opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
                <div className="flex items-center text-brand-500 mb-3">
                    <Users className="w-7 h-7 mr-3" />
                    <h3 className="text-sm font-medium">Followers</h3>
                </div>
                <p className="text-3xl font-bold text-light-surface-900 dark:text-dark-surface-100">1,234</p>
                <p className="text-sm text-light-surface-600 dark:text-dark-surface-400 mt-1 flex items-center">
                    <span className="text-emerald-500 mr-1">↑</span> +23 since last week
                </p>
            </motion.div>

            <motion.div
                initial={{ y: 20, opacity: 0 }}
                animate={{ y: 0, opacity: 1 }}
                transition={{ duration: 0.4, delay: 0.2 }}
                className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm p-6 rounded-xl shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50 relative overflow-hidden group hover:border-brand-500/30 transition-all duration-300"
            >
                <div className="absolute top-0 right-0 w-32 h-32 bg-gradient-to-br from-purple-500/10 to-pink-500/5 rounded-full -mr-10 -mt-10 opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
                <div className="flex items-center text-brand-500 mb-3">
                    <Users className="w-7 h-7 mr-3" />
                    <h3 className="text-sm font-medium">Subscribers</h3>
                </div>
                <p className="text-3xl font-bold text-light-surface-900 dark:text-dark-surface-100">567</p>
                <p className="text-sm text-light-surface-600 dark:text-dark-surface-400 mt-1 flex items-center">
                    <span className="text-emerald-500 mr-1">↑</span> +5 since last month
                </p>
            </motion.div>

            <motion.div
                initial={{ y: 20, opacity: 0 }}
                animate={{ y: 0, opacity: 1 }}
                transition={{ duration: 0.4, delay: 0.3 }}
                className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm p-6 rounded-xl shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50 relative overflow-hidden group hover:border-brand-500/30 transition-all duration-300"
            >
                <div className="absolute top-0 right-0 w-32 h-32 bg-gradient-to-br from-amber-500/10 to-orange-500/5 rounded-full -mr-10 -mt-10 opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
                <div className="flex items-center text-brand-500 mb-3">
                    <BarChart2 className="w-7 h-7 mr-3" />
                    <h3 className="text-sm font-medium">Avg. Views</h3>
                </div>
                <p className="text-3xl font-bold text-light-surface-900 dark:text-dark-surface-100">8,765</p>
                <p className="text-sm text-light-surface-600 dark:text-dark-surface-400 mt-1 flex items-center">
                    <span className="text-rose-500 mr-1">↓</span> -2% vs previous 7 streams
                </p>
            </motion.div>
        </div>

        <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ duration: 0.4, delay: 0.4 }}
            className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm p-6 rounded-xl shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50 mb-10"
        >
            <h3 className="text-xl font-bold text-light-surface-900 dark:text-dark-surface-100 mb-4">Recent Activity</h3>
            <ul className="space-y-4">
                <li className="flex items-center gap-3 p-3 hover:bg-light-surface-200/30 dark:hover:bg-dark-surface-800/30 rounded-lg transition-colors">
                    <div className="w-10 h-10 rounded-full bg-brand-500/20 flex items-center justify-center text-brand-500">
                        <Video className="w-5 h-5" />
                    </div>
                    <div>
                        <p className="text-light-surface-900 dark:text-dark-surface-100 font-medium">New Clip: "Epic Win!"</p>
                        <p className="text-sm text-light-surface-600 dark:text-dark-surface-400">5 mins ago</p>
                    </div>
                    <Button variant="ghost" size="sm" className="ml-auto">
                        <ArrowRight className="w-4 h-4" />
                    </Button>
                </li>
                <li className="flex items-center gap-3 p-3 hover:bg-light-surface-200/30 dark:hover:bg-dark-surface-800/30 rounded-lg transition-colors">
                    <div className="w-10 h-10 rounded-full bg-brand-500/20 flex items-center justify-center text-brand-500">
                        <Calendar className="w-5 h-5" />
                    </div>
                    <div>
                        <p className="text-light-surface-900 dark:text-dark-surface-100 font-medium">Scheduled Post: "Stream starting soon!"</p>
                        <p className="text-sm text-light-surface-600 dark:text-dark-surface-400">1 hour ago</p>
                    </div>
                    <Button variant="ghost" size="sm" className="ml-auto">
                        <ArrowRight className="w-4 h-4" />
                    </Button>
                </li>
                <li className="flex items-center gap-3 p-3 hover:bg-light-surface-200/30 dark:hover:bg-dark-surface-800/30 rounded-lg transition-colors">
                    <div className="w-10 h-10 rounded-full bg-brand-500/20 flex items-center justify-center text-brand-500">
                        <Users className="w-5 h-5" />
                    </div>
                    <div>
                        <p className="text-light-surface-900 dark:text-dark-surface-100 font-medium">New Follower: "CoolDude123"</p>
                        <p className="text-sm text-light-surface-600 dark:text-dark-surface-400">3 hours ago</p>
                    </div>
                    <Button variant="ghost" size="sm" className="ml-auto">
                        <ArrowRight className="w-4 h-4" />
                    </Button>
                </li>
            </ul>
        </motion.div>

        <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ duration: 0.4, delay: 0.5 }}
            className="bg-brand-500/10 border border-brand-500/20 p-6 rounded-xl"
        >
            <div className="flex items-center gap-2 text-brand-500 mb-3">
                <Sparkles className="w-5 h-5" />
                <h3 className="font-semibold">Creator Tip</h3>
            </div>
            <p className="text-light-surface-700 dark:text-dark-surface-300 mb-4">
                Did you know you can auto-schedule your clips to be posted across multiple platforms? Try our new scheduling feature to maximize your reach.
            </p>
            <Button variant="default" size="sm">
                Try Scheduling <ArrowRight className="w-4 h-4 ml-1" />
            </Button>
        </motion.div>
    </div>
);

// Define a proper type for Twitch videos
interface TwitchVideo {
    id: string;
    title: string;
    url: string;
    thumbnail_url: string;
    created_at: string;
    duration: string;
    view_count: number;
    type: string;
}


//TODO: Seperate sections in to their own files with their own components
const ContentSection = () => {
    const { getToken } = useAuth();
    const [videos, setVideos] = useState<TwitchVideo[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchVideos = async () => {
            setIsLoading(true);
            setError(null);
            try {
                // Get the API base URL based on environment
                const apiBaseUrl = process.env.NODE_ENV === 'production'
                    ? 'https://api.creatorsync.app'
                    : 'http://localhost:8080';

                // Use the simplest approach to get the token
                let token;
                try {
                    // Just use the standard getToken method without any templates
                    token = await getToken();
                    
                    if (!token) {
                        throw new Error("User not authenticated or token not available.");
                    }
                } catch (tokenError) {
                    console.error('Failed to get authentication token:', tokenError);
                    throw new Error("Authentication failed: Unable to get valid token");
                }
                
                // Add detailed logging for production debugging
                console.log('Environment:', process.env.NODE_ENV);
                console.log('Token type:', typeof token);
                console.log('Token length:', token.length);
                console.log('Token first 10 chars:', token.substring(0, 10) + '...');
                console.log('Token format check:', token.includes('.') ? 'Contains dots (likely JWT)' : 'No dots (not standard JWT)');
                
                // Create the authorization header with proper format
                // Make sure the Bearer prefix has exactly one space after it
                // In production, we'll log more details to help diagnose the issue
                const authHeader = `Bearer ${token}`;
                
                if (process.env.NODE_ENV === 'production') {
                    console.log('Production auth debugging:');
                    console.log('- Auth header starts with:', authHeader.substring(0, 15) + '...');
                    console.log('- Auth header length:', authHeader.length);
                    console.log('- Token structure check:', token.split('.').length === 3 ? 'Valid JWT structure (3 parts)' : 'Not standard JWT structure');
                    // Don't log the actual token for security reasons
                }
                
                // Make the API request with the same format that works locally
                const requestOptions = {
                    method: 'GET',
                    headers: {
                        'Authorization': authHeader,
                        'Content-Type': 'application/json'
                    }
                };

                // Make the request
                const response = await fetch(`${apiBaseUrl}/api/twitch/videos`, requestOptions);
                if (!response.ok) {
                    let errorMessage = `HTTP error! status: ${response.status}`;
                    try {
                        const errorData = await response.json();
                        // Log the full error response for debugging
                        console.error('Full API Error Response:', JSON.stringify(errorData));
                        
                        // Create a more detailed error message
                        if (errorData) {
                            if (typeof errorData === 'string') {
                                errorMessage = errorData;
                            } else if (errorData.error) {
                                errorMessage = errorData.error;
                            } else if (errorData.errors && errorData.errors.length > 0) {
                                // Handle Clerk-specific error format
                                const clerkErrors = errorData.errors;
                                errorMessage = clerkErrors.map((err: { message?: string; code?: string; long_message?: string }) => 
                                    `${err.message || err.code}: ${err.long_message || ''}`
                                ).join('; ');
                                
                                // If this is an authentication error, log additional details
                                if (response.status === 401) {
                                    console.error('Authentication error detected. Token details:', {
                                        length: token.length,
                                        format: token.includes('.') ? 'JWT format' : 'Not JWT format',
                                        authHeader: authHeader.substring(0, 15) + '...'
                                    });
                                }
                            } else {
                                errorMessage = JSON.stringify(errorData);
                            }
                        }
                    } catch (jsonError) {
                        console.error('Failed to parse error response:', jsonError);
                    }
                    throw new Error(errorMessage);
                }
                const data = await response.json();
                setVideos(data.videos || []);
            } catch (e: unknown) {
                console.error("Failed to fetch videos:", e);
                const errorMessage = e instanceof Error ? e.message : "Failed to load videos.";
                setError(errorMessage);
            }
            setIsLoading(false);
        };

        fetchVideos();
    }, [getToken]);

    if (isLoading) {
        return (
            <div className="flex justify-center items-center h-64">
                <p className="text-xl text-light-surface-700 dark:text-dark-surface-300">Loading content...</p>
            </div>
        );
    }

    if (error) {
        return (
            <div className="flex flex-col justify-center items-center h-64">
                <p className="text-xl text-red-500">Error: {error}</p>
                <Button onClick={() => { /* Consider adding a retry mechanism */ }} className="mt-4">Try Again</Button>
            </div>
        );
    }

    if (videos.length === 0) {
        return (
            <div className="flex flex-col justify-center items-center h-64">
                <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-purple-500/20 border border-purple-500/30 text-purple-500 text-sm mb-6 w-fit">
                    <Video className="w-4 h-4" />
                    <span>Content Library</span>
                </div>
                <h2 className="text-4xl font-bold mb-6">
                    <span className="text-light-surface-900 dark:text-dark-surface-100">Your</span>{' '}
                    <span className="bg-clip-text text-transparent bg-gradient-to-r from-purple-500 to-pink-500">Content</span>
                </h2>
                <p className="text-xl text-light-surface-700 dark:text-dark-surface-300 max-w-3xl mb-10 text-center">
                    No videos found. Once you have some videos on Twitch, they'll appear here!
                </p>
            </div>
        );
    }

    return (
        <div>
            <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-purple-500/20 border border-purple-500/30 text-purple-500 text-sm mb-6 w-fit">
                <Video className="w-4 h-4" />
                <span>Content Library</span>
            </div>

            <h2 className="text-4xl font-bold mb-6">
                <span className="text-light-surface-900 dark:text-dark-surface-100">Your</span>{' '}
                <span className="bg-clip-text text-transparent bg-gradient-to-r from-purple-500 to-pink-500">Content</span>
            </h2>

            <p className="text-xl text-light-surface-700 dark:text-dark-surface-300 max-w-3xl mb-10">
                Manage your VODs, highlights, and clips all in one place.
            </p>

            <motion.div
                initial={{ y: 20, opacity: 0 }}
                animate={{ y: 0, opacity: 1 }}
                transition={{ duration: 0.4 }}
                className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm p-8 rounded-xl shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50 mb-8"
            >
                <div className="flex justify-between items-center mb-6">
                    <h3 className="text-xl font-bold text-light-surface-900 dark:text-dark-surface-100">Recent Clips</h3>
                    <Button variant="outline" size="sm">
                        View All
                    </Button>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {videos.map((video) => (
                        <div key={video.id} className="bg-light-surface-200/50 dark:bg-dark-surface-800/50 rounded-lg overflow-hidden group hover:border-purple-500/30 border border-transparent transition-all duration-300">
                            <a href={video.url} target="_blank" rel="noopener noreferrer" className="block">
                                <div className="aspect-video bg-light-surface-300/50 dark:bg-dark-surface-700/50 relative">
                                    {video.thumbnail_url && (
                                        <Image
                                            src={video.thumbnail_url.replace('%{width}', '480').replace('%{height}', '270')}
                                            alt={video.title || 'Twitch video thumbnail'}
                                            className="object-cover"
                                            fill
                                            sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
                                        />
                                    )}
                                    <div className="absolute inset-0 flex items-center justify-center bg-black/30 opacity-0 group-hover:opacity-100 transition-opacity">
                                        <Play className="w-12 h-12 text-white" />
                                    </div>
                                </div>
                            </a>
                            <div className="p-4">
                                <h4 className="font-medium text-light-surface-900 dark:text-dark-surface-100 mb-1 truncate" title={video.title}>{video.title || 'Untitled Video'}</h4>
                                <p className="text-sm text-light-surface-600 dark:text-dark-surface-400">{new Date(video.created_at).toLocaleDateString()} • {video.duration}</p>
                                <p className="text-xs text-light-surface-500 dark:text-dark-surface-500">Views: {video.view_count}</p>
                            </div>
                        </div>
                    ))}
                </div>
            </motion.div>
        </div>
    );
};

const AnalyticsSection = () => (
    <div>
        <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-emerald-500/20 border border-emerald-500/30 text-emerald-500 text-sm mb-6 w-fit">
            <BarChart2 className="w-4 h-4" />
            <span>Performance Metrics</span>
        </div>

        <h2 className="text-4xl font-bold mb-6">
            <span className="text-light-surface-900 dark:text-dark-surface-100">Detailed</span>{' '}
            <span className="bg-clip-text text-transparent bg-gradient-to-r from-emerald-500 to-teal-500">Analytics</span>
        </h2>

        <p className="text-xl text-light-surface-700 dark:text-dark-surface-300 max-w-3xl mb-10">
            Track your growth and understand your audience better with in-depth analytics.
        </p>

        <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ duration: 0.4 }}
            className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm p-8 rounded-xl shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50"
        >
            <div className="h-64 bg-light-surface-200/50 dark:bg-dark-surface-800/50 rounded-lg flex items-center justify-center mb-6">
                <p className="text-light-surface-600 dark:text-dark-surface-400">Analytics visualization will appear here</p>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {['Views', 'Watch Time', 'Engagement'].map((metric) => (
                    <div key={metric} className="bg-light-surface-200/50 dark:bg-dark-surface-800/50 p-4 rounded-lg">
                        <h4 className="text-sm font-medium text-light-surface-600 dark:text-dark-surface-400 mb-1">{metric}</h4>
                        <p className="text-2xl font-bold text-light-surface-900 dark:text-dark-surface-100">--</p>
                    </div>
                ))}
            </div>
        </motion.div>
    </div>
);

const ScheduledSection = () => (
    <div>
        <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-amber-500/20 border border-amber-500/30 text-amber-500 text-sm mb-6 w-fit">
            <Calendar className="w-4 h-4" />
            <span>Content Calendar</span>
        </div>

        <h2 className="text-4xl font-bold mb-6">
            <span className="text-light-surface-900 dark:text-dark-surface-100">Scheduled</span>{' '}
            <span className="bg-clip-text text-transparent bg-gradient-to-r from-amber-500 to-orange-500">Posts</span>
        </h2>

        <p className="text-xl text-light-surface-700 dark:text-dark-surface-300 max-w-3xl mb-10">
            Plan and schedule your content across all platforms from one dashboard.
        </p>

        <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ duration: 0.4 }}
            className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm p-8 rounded-xl shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50 mb-8"
        >
            <div className="flex justify-between items-center mb-6">
                <h3 className="text-xl font-bold text-light-surface-900 dark:text-dark-surface-100">Upcoming Posts</h3>
                <Button variant="default" size="sm">
                    <Plus className="w-4 h-4 mr-1" /> New Post
                </Button>
            </div>
            <div className="space-y-4">
                {['Tomorrow, 10:00 AM', 'Friday, 6:00 PM', 'Sunday, 12:00 PM'].map((time, index) => (
                    <div key={index} className="flex items-center gap-4 p-4 bg-light-surface-200/50 dark:bg-dark-surface-800/50 rounded-lg">
                        <div className="w-10 h-10 rounded-full bg-amber-500/20 flex items-center justify-center text-amber-500">
                            <Calendar className="w-5 h-5" />
                        </div>
                        <div className="flex-1">
                            <h4 className="font-medium text-light-surface-900 dark:text-dark-surface-100">Stream Announcement</h4>
                            <p className="text-sm text-light-surface-600 dark:text-dark-surface-400">{time}</p>
                        </div>
                        <div className="flex items-center gap-2">
                            <span className="px-2 py-1 text-xs rounded-full bg-amber-500/10 text-amber-500">Twitter</span>
                            <Button variant="ghost" size="sm">
                                <ArrowRight className="w-4 h-4" />
                            </Button>
                        </div>
                    </div>
                ))}
            </div>
        </motion.div>
    </div>
);

const SettingsSection = () => (
    <div>
        <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-blue-500/20 border border-blue-500/30 text-blue-500 text-sm mb-6 w-fit">
            <SettingsIcon className="w-4 h-4" />
            <span>Configuration</span>
        </div>

        <h2 className="text-4xl font-bold mb-6">
            <span className="text-light-surface-900 dark:text-dark-surface-100">Account</span>{' '}
            <span className="bg-clip-text text-transparent bg-gradient-to-r from-blue-500 to-cyan-500">Settings</span>
        </h2>

        <p className="text-xl text-light-surface-700 dark:text-dark-surface-300 max-w-3xl mb-10">
            Configure your CreatorSync account, integrations, and preferences.
        </p>

        <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ duration: 0.4 }}
            className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm p-8 rounded-xl shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50 mb-8"
        >
            <div className="space-y-6">
                <div>
                    <h3 className="text-xl font-bold text-light-surface-900 dark:text-dark-surface-100 mb-4">Connected Accounts</h3>
                    <div className="space-y-3">
                        {['Twitch', 'YouTube', 'Twitter'].map((platform) => (
                            <div key={platform} className="flex items-center justify-between p-4 bg-light-surface-200/50 dark:bg-dark-surface-800/50 rounded-lg">
                                <div className="flex items-center gap-3">
                                    <div className="w-8 h-8 rounded-full bg-blue-500/20 flex items-center justify-center text-blue-500">
                                        {platform.charAt(0)}
                                    </div>
                                    <span className="font-medium text-light-surface-900 dark:text-dark-surface-100">{platform}</span>
                                </div>
                                <Button variant="outline" size="sm">
                                    Connect
                                </Button>
                            </div>
                        ))}
                    </div>
                </div>
                <div>
                    <h3 className="text-xl font-bold text-light-surface-900 dark:text-dark-surface-100 mb-4">Preferences</h3>
                    <div className="space-y-3">
                        {['Notifications', 'Theme', 'Privacy'].map((setting) => (
                            <div key={setting} className="flex items-center justify-between p-4 bg-light-surface-200/50 dark:bg-dark-surface-800/50 rounded-lg">
                                <span className="font-medium text-light-surface-900 dark:text-dark-surface-100">{setting}</span>
                                <Button variant="ghost" size="sm">
                                    <SettingsIcon className="w-4 h-4 mr-1" /> Configure
                                </Button>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </motion.div>
    </div>
);


const DashboardPage = () => {
    const [activeSection, setActiveSection] = useState('Overview');

    const sections = [
        { name: 'Overview', icon: BarChart2, component: <OverviewContent /> },
        { name: 'Content', icon: Video, component: <ContentSection /> },
        { name: 'Analytics', icon: BarChart2, component: <AnalyticsSection /> }, // Can use a different icon for detailed analytics
        { name: 'Scheduled', icon: Calendar, component: <ScheduledSection /> },
        { name: 'Settings', icon: SettingsIcon, component: <SettingsSection /> },
        {
            name: 'Sign Out', icon: LogOut, isSignOut: true as const
        },
    ];

    const renderSection = () => {
        const section = sections.find(s => s.name === activeSection);
        return section ? section.component : <OverviewContent />;
    };

    // Animation is now handled directly by framer-motion

    return (
        <div className="flex min-h-screen bg-[var(--bg-primary)] text-[var(--text-primary)] overflow-x-hidden relative">
            {/* Unified Background Elements */}
            <div className="fixed inset-0 overflow-hidden pointer-events-none">
                {/* Subtle gradient overlay */}
                <div className="absolute inset-0 bg-gradient-to-br from-brand-500/5 via-transparent to-accent-500/10" />

                {/* High-resolution gradient orbs with improved blending */}
                <div
                    className="absolute opacity-30 animate-float"
                    style={{
                        top: '15%',
                        left: '10%',
                        width: '40rem',
                        height: '40rem',
                        background: 'radial-gradient(circle, rgba(99,102,241,0.15) 0%, rgba(99,102,241,0) 70%)',
                        filter: 'blur(60px)',
                        animationDelay: '0s',
                        transform: 'translate3d(0, 0, 0)'
                    }}
                />
                <div
                    className="absolute opacity-30 animate-float"
                    style={{
                        bottom: '10%',
                        right: '5%',
                        width: '35rem',
                        height: '35rem',
                        background: 'radial-gradient(circle, rgba(236,72,153,0.1) 0%, rgba(236,72,153,0) 70%)',
                        filter: 'blur(60px)',
                        animationDelay: '2s',
                        transform: 'translate3d(0, 0, 0)'
                    }}
                />
            </div>

            {/* Sidebar Navigation */}
            <motion.aside
                initial={{ x: -100, opacity: 0 }}
                animate={{ x: 0, opacity: 1 }}
                transition={{ duration: 0.5, ease: "easeOut" }}
                className="w-64 bg-light-surface-50/85 dark:bg-dark-surface-950/85 backdrop-blur-xl text-light-surface-900 dark:text-dark-surface-100 p-6 space-y-4 border-r border-light-surface-200/50 dark:border-dark-surface-800/50 fixed top-0 left-0 h-full shadow-lg z-20">
                <div className="text-2xl font-bold mb-8">
                    <span className="text-light-surface-900 dark:text-dark-surface-100">Creator</span>
                    <span className="text-gradient">Sync</span>
                </div>
                <nav className="space-y-2">
                    {sections.map((section) => {
                        if (section.name === 'Sign Out' && section.isSignOut) {
                            return (
                                <SignOutButton key={section.name}>
                                    <Button
                                        variant="ghost"
                                        size="lg"
                                        className="w-full justify-start text-light-surface-700 dark:text-dark-surface-300 hover:bg-brand-500/10 hover:text-brand-500 focus-visible:ring-brand-500/20"
                                    >
                                        <LogOut className="w-5 h-5 mr-2" />
                                        {section.name}
                                    </Button>
                                </SignOutButton>
                            );
                        }
                        const isActive = activeSection === section.name;
                        return (
                            <button
                                key={section.name}
                                onClick={() => setActiveSection(section.name)}
                                className={cn(
                                    "flex items-center w-full px-4 py-3 rounded-lg group transition-all duration-300",
                                    isActive
                                        ? "bg-brand-500/20 text-brand-500 border border-brand-500/30"
                                        : "text-light-surface-700 dark:text-dark-surface-300 hover:bg-light-surface-100/80 dark:hover:bg-dark-surface-800/80 hover:text-light-surface-900 dark:hover:text-dark-surface-100 border border-transparent"
                                )}
                            >
                                <section.icon className={cn(
                                    "w-5 h-5 mr-3 transition-all duration-300",
                                    isActive ? "text-brand-500" : "text-light-surface-500 dark:text-dark-surface-400 group-hover:text-light-surface-900 dark:group-hover:text-dark-surface-100"
                                )} />
                                <span>{section.name}</span>
                            </button>
                        );
                    })}
                </nav>
                <div className="absolute bottom-8 left-0 right-0 px-6">
                    <div className="p-4 rounded-xl bg-brand-500/10 border border-brand-500/20">
                        <div className="flex items-center gap-2 text-brand-500 text-sm mb-2">
                            <Sparkles className="w-4 h-4" />
                            <span>Pro Tip</span>
                        </div>
                        <p className="text-sm text-light-surface-700 dark:text-dark-surface-300 mb-3">
                            Connect your Twitch account to auto-import your latest clips.
                        </p>
                        <Button variant="default" size="sm" className="w-full">
                            <Plus className="w-4 h-4 mr-1" /> Connect Account
                        </Button>
                    </div>
                </div>
            </motion.aside>

            {/* Main Content Area */}
            <motion.main
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.5, delay: 0.2 }}
                className="flex-1 p-8 ml-64 relative z-10"
            >
                <div className="max-w-7xl mx-auto">
                    <motion.div
                        initial={{ y: 20, opacity: 0 }}
                        animate={{ y: 0, opacity: 1 }}
                        transition={{ duration: 0.6 }}
                    >
                        {renderSection()}
                    </motion.div>
                </div>
            </motion.main>
        </div>
    );
};

export default DashboardPage;

