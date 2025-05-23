'use client';

import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { useAuth } from '@clerk/nextjs';
import { Button } from '@/components/ui/button';
import { Video, Download, ExternalLink, Calendar, PlayCircle, Eye, Clock } from 'lucide-react';
import Image from 'next/image';

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

interface TwitchClip {
    id: string;
    title: string;
    url: string;
    thumbnail_url: string;
    created_at: string;
    duration: number; // Clips return duration in seconds, not string format
    view_count: number;
    creator_name: string;
    broadcaster_name: string;
}

interface ClerkError {
    long_message?: string;
    message?: string;
    code?: string;
    [key: string]: unknown;
}

export function ContentSection() {
    const { getToken, isLoaded, isSignedIn } = useAuth();
    const [videos, setVideos] = useState<TwitchVideo[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchVideos = async () => {
            if (!isLoaded) {
                return;
            }

            if (!isSignedIn) {
                setError("Please sign in to view your content");
                setIsLoading(false);
                return;
            }

            setIsLoading(true);
            setError(null);

            try {
                // Determine API base URL based on environment
                let apiBaseUrl: string;
                if (typeof window !== 'undefined' && window.location.hostname === 'dev.creatorsync.app') {
                    apiBaseUrl = 'https://api-dev.creatorsync.app';
                } else if (process.env.NEXT_PUBLIC_APP_ENV === 'staging') {
                    apiBaseUrl = 'https://api-dev.creatorsync.app';
                } else if (process.env.NODE_ENV === 'production') {
                    apiBaseUrl = 'https://api.creatorsync.app';
                } else {
                    apiBaseUrl = 'http://localhost:8080';
                }

                let token;
                try {
                    token = await getToken();
                } catch (tokenError: unknown) {
                    console.error('Failed to get token:', tokenError);
                    const errorMessage = tokenError instanceof Error ? tokenError.message : 'Unable to get token';
                    throw new Error(`Authentication failed: ${errorMessage}`);
                }

                if (!token) {
                    throw new Error("Authentication failed: No token available");
                }

                // Sync user to ensure they exist in the database (especially for staging)
                try {
                    await fetch(`${apiBaseUrl}/api/user/sync`, {
                        method: 'POST',
                        headers: {
                            'Authorization': `Bearer ${token}`,
                            'Content-Type': 'application/json',
                        },
                    });
                } catch (syncError) {
                    console.warn('User sync failed, continuing anyway:', syncError);
                    // Don't throw here, continue with the rest of the logic
                }

                // Fetch videos (broadcasts, highlights, uploads)
                const videosResponse = await fetch(`${apiBaseUrl}/api/twitch/videos`, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json',
                    },
                });

                if (!videosResponse.ok) {
                    if (videosResponse.status === 401) {
                        throw new Error("Session expired. Please sign in again.");
                    }

                    let errorMessage = `Failed to fetch videos: ${videosResponse.status}`;
                    try {
                        const errorData = await videosResponse.json();
                        if (errorData.error) {
                            errorMessage = errorData.error;
                        }
                    } catch {
                        // Ignore JSON parsing errors and use the generic message
                    }
                    throw new Error(errorMessage);
                }

                const videosData = await videosResponse.json();
                const fetchedVideos = videosData.videos || [];

                // Try to fetch clips (might not be implemented yet)
                let fetchedClips: TwitchClip[] = [];
                try {
                    const clipsResponse = await fetch(`${apiBaseUrl}/api/twitch/clips`, {
                        headers: {
                            'Authorization': `Bearer ${token}`,
                            'Content-Type': 'application/json',
                        },
                    });

                    if (clipsResponse.ok) {
                        const clipsData = await clipsResponse.json();
                        fetchedClips = clipsData.clips || [];
                    }
                } catch {
                    // Clips endpoint might not be implemented yet, continue without clips
                    console.info('Clips endpoint not available yet');
                }

                // Convert clips to video format for unified display
                const clipsAsVideos: TwitchVideo[] = fetchedClips.map(clip => ({
                    id: clip.id,
                    title: clip.title,
                    url: clip.url,
                    thumbnail_url: clip.thumbnail_url,
                    created_at: clip.created_at,
                    duration: `${Math.floor(clip.duration)}s`, // Convert seconds to string format
                    view_count: clip.view_count,
                    type: 'clip'
                }));

                // Combine videos and clips
                const allContent = [...fetchedVideos, ...clipsAsVideos];
                setVideos(allContent);

            } catch (err: unknown) {
                console.error('Error fetching content:', err);

                if (err instanceof Error) {
                    setError(err.message);
                } else {
                    const clerkError = err as ClerkError;
                    if (clerkError?.long_message) {
                        setError(clerkError.long_message);
                    } else if (clerkError?.message) {
                        setError(clerkError.message);
                    } else {
                        setError('An unexpected error occurred while fetching content');
                    }
                }
            } finally {
                setIsLoading(false);
            }
        };

        fetchVideos();
    }, [isLoaded, isSignedIn, getToken]);

    const formatViewCount = (views: number) => {
        if (views >= 1000000) {
            return `${(views / 1000000).toFixed(1)}M`;
        }
        if (views >= 1000) {
            return `${(views / 1000).toFixed(1)}K`;
        }
        return views.toString();
    };

    const formatDuration = (duration: string) => {
        const match = duration.match(/(\d+)h?(\d+)m?(\d+)?s?/);
        if (match) {
            const hours = parseInt(match[1] || '0');
            const minutes = parseInt(match[2] || '0');
            const seconds = parseInt(match[3] || '0');

            if (hours > 0) {
                return `${hours}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
            } else {
                return `${minutes}:${seconds.toString().padStart(2, '0')}`;
            }
        }
        return duration;
    };

    const formatDate = (dateString: string | Date) => {
        return new Date(dateString).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric'
        });
    };

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
                    <h3 className="text-red-500 font-semibold mb-2">Error Loading Content</h3>
                    <p className="text-red-500/80">{error}</p>
                    <Button
                        variant="outline"
                        size="sm"
                        className="mt-4"
                        onClick={() => window.location.reload()}
                    >
                        Try Again
                    </Button>
                </div>
            </div>
        );
    }

    return (
        <div className="p-8">
            <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-blue-500/20 border border-blue-500/30 text-blue-500 text-sm mb-6 w-fit">
                <Video className="w-4 h-4" />
                <span>Your Content</span>
            </div>

            <h2 className="text-4xl font-bold mb-6">
                <span className="text-light-surface-900 dark:text-dark-surface-100">Your</span>{' '}
                <span className="text-gradient">Videos</span>
            </h2>

            <p className="text-xl text-light-surface-700 dark:text-dark-surface-300 max-w-3xl mb-10">
                Manage and analyze your Twitch content. Download clips, view analytics, and optimize your performance.
            </p>

            {(() => {
                // Separate videos by type
                const broadcasts = videos.filter(video => video.type === 'archive' || video.type === 'highlight' || video.type === 'upload');
                const clips = videos.filter(video => video.type === 'clip');

                const renderVideoGrid = (videoList: TwitchVideo[], startIndex: number = 0) => (
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        {videoList.map((video, index) => (
                            <motion.div
                                key={video.id}
                                initial={{ opacity: 0, y: 20 }}
                                animate={{ opacity: 1, y: 0 }}
                                transition={{ delay: (startIndex + index) * 0.1 }}
                                className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm rounded-xl overflow-hidden shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50 hover:border-brand-500/30 transition-all duration-300 group"
                            >
                                {/* Thumbnail */}
                                <div className="relative aspect-video overflow-hidden">
                                    <Image
                                        src={video.thumbnail_url.replace('%{width}', '320').replace('%{height}', '180')}
                                        alt={video.title}
                                        fill
                                        className="object-cover group-hover:scale-105 transition-transform duration-300"
                                    />
                                    <div className="absolute inset-0 bg-black/20 group-hover:bg-black/10 transition-colors duration-300" />
                                    <div className="absolute top-2 right-2 bg-black/70 text-white text-xs px-2 py-1 rounded-md flex items-center gap-1">
                                        <Clock className="w-3 h-3" />
                                        {formatDuration(video.duration)}
                                    </div>
                                    {/* Video type badge */}
                                    <div className="absolute top-2 left-2">
                                        <span className={`text-xs px-2 py-1 rounded-md font-medium ${video.type === 'clip'
                                            ? 'bg-purple-500/80 text-white'
                                            : 'bg-blue-500/80 text-white'
                                            }`}>
                                            {video.type === 'clip' ? 'CLIP' : 'VOD'}
                                        </span>
                                    </div>
                                    <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                                        <PlayCircle className="w-12 h-12 text-white drop-shadow-lg" />
                                    </div>
                                </div>

                                {/* Content */}
                                <div className="p-4">
                                    <h3 className="font-semibold text-light-surface-900 dark:text-dark-surface-100 mb-2 line-clamp-2 leading-tight">
                                        {video.title}
                                    </h3>

                                    <div className="flex items-center gap-4 text-sm text-light-surface-600 dark:text-dark-surface-400 mb-4">
                                        <div className="flex items-center gap-1">
                                            <Eye className="w-4 h-4" />
                                            {formatViewCount(video.view_count)}
                                        </div>
                                        <div className="flex items-center gap-1">
                                            <Calendar className="w-4 h-4" />
                                            {formatDate(video.created_at)}
                                        </div>
                                    </div>

                                    <div className="flex gap-2">
                                        <Button
                                            variant="outline"
                                            size="sm"
                                            className="flex-1"
                                            onClick={() => window.open(video.url, '_blank')}
                                        >
                                            <ExternalLink className="w-4 h-4 mr-1" />
                                            View
                                        </Button>
                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            className="flex-1"
                                        >
                                            <Download className="w-4 h-4 mr-1" />
                                            Download
                                        </Button>
                                    </div>
                                </div>
                            </motion.div>
                        ))}
                    </div>
                );

                const renderEmptyState = (type: 'broadcasts' | 'clips') => (
                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        className="text-center py-12"
                    >
                        <Video className="w-12 h-12 text-light-surface-400 dark:text-dark-surface-600 mx-auto mb-4" />
                        <h4 className="text-lg font-medium text-light-surface-700 dark:text-dark-surface-300 mb-2">
                            No {type} found
                        </h4>
                        <p className="text-light-surface-600 dark:text-dark-surface-400">
                            {type === 'broadcasts'
                                ? 'Your past broadcasts and highlights will appear here.'
                                : 'Your clips will appear here once created.'}
                        </p>
                    </motion.div>
                );

                if (videos.length === 0) {
                    return (
                        <motion.div
                            initial={{ opacity: 0, y: 20 }}
                            animate={{ opacity: 1, y: 0 }}
                            className="text-center py-16"
                        >
                            <Video className="w-16 h-16 text-light-surface-400 dark:text-dark-surface-600 mx-auto mb-4" />
                            <h3 className="text-xl font-semibold text-light-surface-700 dark:text-dark-surface-300 mb-2">
                                No videos found
                            </h3>
                            <p className="text-light-surface-600 dark:text-dark-surface-400 mb-6">
                                Your Twitch videos will appear here once they're available. Make sure your Twitch account is connected.
                            </p>
                            <Button variant="default">
                                <ExternalLink className="w-4 h-4 mr-2" />
                                View on Twitch
                            </Button>
                        </motion.div>
                    );
                }

                return (
                    <div className="space-y-12">
                        {/* Recent Broadcasts Section */}
                        <div>
                            <div className="flex items-center gap-3 mb-6">
                                <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-blue-500/20 border border-blue-500/30 text-blue-500 text-sm">
                                    <Video className="w-4 h-4" />
                                    <span>Recent Broadcasts</span>
                                </div>
                                {broadcasts.length > 0 && (
                                    <span className="text-sm text-light-surface-600 dark:text-dark-surface-400">
                                        {broadcasts.length} video{broadcasts.length !== 1 ? 's' : ''}
                                    </span>
                                )}
                            </div>
                            {broadcasts.length > 0 ? renderVideoGrid(broadcasts, 0) : renderEmptyState('broadcasts')}
                        </div>

                        {/* Clips Section */}
                        <div>
                            <div className="flex items-center gap-3 mb-6">
                                <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-purple-500/20 border border-purple-500/30 text-purple-500 text-sm">
                                    <PlayCircle className="w-4 h-4" />
                                    <span>Clips</span>
                                </div>
                                {clips.length > 0 && (
                                    <span className="text-sm text-light-surface-600 dark:text-dark-surface-400">
                                        {clips.length} clip{clips.length !== 1 ? 's' : ''}
                                    </span>
                                )}
                            </div>
                            {clips.length > 0 ? renderVideoGrid(clips, broadcasts.length) : renderEmptyState('clips')}
                        </div>
                    </div>
                );
            })()}
        </div>
    );
} 
