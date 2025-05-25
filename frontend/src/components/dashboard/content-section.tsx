'use client';

import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { useAuth } from '@clerk/nextjs';
import { ConnectionModal } from '@/components/ui/connection-modal';
import { 
    Video, 
    Download, 
    ExternalLink, 
    PlayCircle, 
    Eye, 
    Clock, 
    RefreshCw,
    Film,
    Monitor,
    TrendingUp
} from 'lucide-react';
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
    const [twitchConnected, setTwitchConnected] = useState<boolean | null>(null);
    const [showConnectionModal, setShowConnectionModal] = useState(false);
    const [refreshing, setRefreshing] = useState(false);

    // Helper function to get API base URL
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

    const fetchContent = async () => {
        if (!isLoaded) {
            return;
        }

        if (!isSignedIn) {
            setError("Please sign in to view your content");
            setIsLoading(false);
            return;
        }

        setError(null);

        try {
            const apiBaseUrl = getApiBaseUrl();

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

            // First, check connection status
            const connectionResponse = await fetch(`${apiBaseUrl}/api/analytics/connection-status`, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
            });

            if (!connectionResponse.ok) {
                throw new Error("Failed to check connection status");
            }

            const connectionData = await connectionResponse.json();
            const isConnected = connectionData.platforms?.twitch?.connected || false;
            setTwitchConnected(isConnected);

            if (!isConnected) {
                setIsLoading(false);
                return;
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
            } catch {
                // Continue anyway if sync fails
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

    useEffect(() => {
        fetchContent();
    }, [isLoaded, isSignedIn, getToken]);

    const handleRefresh = async () => {
        if (refreshing) return;
        setRefreshing(true);
        await fetchContent();
        setRefreshing(false);
    };

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
            <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 flex items-center justify-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-500"></div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 p-6">
                <div className="max-w-4xl mx-auto">
                    <div className="bg-red-500/10 border border-red-500/20 rounded-2xl p-8 text-center">
                        <h3 className="text-red-400 font-semibold mb-4 text-xl">Error Loading Content</h3>
                        <p className="text-red-300/80 mb-6">{error}</p>
                        <button
                            onClick={() => window.location.reload()}
                            className="bg-red-600 hover:bg-red-700 text-white px-6 py-3 rounded-lg font-semibold transition-all duration-200"
                        >
                            Try Again
                        </button>
                    </div>
                </div>
            </div>
        );
    }

    // Show connection prompt if Twitch is not connected
    if (twitchConnected === false || (twitchConnected === null && videos.length === 0 && !isLoading)) {
        return (
            <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 p-6">
                <div className="max-w-4xl mx-auto">
                    <div className="flex items-center mb-2">
                        <Video className="w-6 h-6 text-emerald-500 mr-2" />
                        <span className="text-sm text-emerald-400 font-medium">Your Content</span>
                    </div>
                    <h1 className="text-3xl font-bold text-white mb-6">Connect Your Twitch Account</h1>
                    
                    <div className="bg-gray-900/50 backdrop-blur-xl border border-gray-800/50 rounded-2xl p-8">
                        <div className="text-center">
                            <div className="w-16 h-16 bg-emerald-500/20 rounded-full flex items-center justify-center mx-auto mb-6">
                                <Video className="w-8 h-8 text-emerald-500" />
                            </div>
                            <h2 className="text-2xl font-bold text-white mb-4">Content Library Awaiting Connection</h2>
                            <p className="text-gray-400 mb-8 max-w-2xl mx-auto">
                                Connect your Twitch account to access your complete video library, manage your content, and download clips for social media.
                            </p>
                            
                            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                                <div className="bg-emerald-500/10 border border-emerald-500/20 rounded-xl p-6">
                                    <Video className="w-8 h-8 text-emerald-500 mx-auto mb-3" />
                                    <h3 className="text-white font-semibold mb-2">Video Management</h3>
                                    <p className="text-gray-400 text-sm">Access your broadcasts, highlights, and VODs in one organized library.</p>
                                </div>
                                <div className="bg-purple-500/10 border border-purple-500/20 rounded-xl p-6">
                                    <PlayCircle className="w-8 h-8 text-purple-500 mx-auto mb-3" />
                                    <h3 className="text-white font-semibold mb-2">Clip Collection</h3>
                                    <p className="text-gray-400 text-sm">View, organize, and download your best clips and highlights.</p>
                                </div>
                                <div className="bg-blue-500/10 border border-blue-500/20 rounded-xl p-6">
                                    <Download className="w-8 h-8 text-blue-500 mx-auto mb-3" />
                                    <h3 className="text-white font-semibold mb-2">Content Export</h3>
                                    <p className="text-gray-400 text-sm">Download videos and clips for social media and content repurposing.</p>
                                </div>
                            </div>
                            
                            <button
                                onClick={() => setShowConnectionModal(true)}
                                className="bg-gradient-to-r from-emerald-600 to-emerald-700 hover:from-emerald-700 hover:to-emerald-800 text-white px-8 py-3 rounded-lg font-semibold transition-all duration-200"
                            >
                                Connect Twitch Account
                            </button>
                        </div>
                    </div>
                </div>

                <ConnectionModal
                    isOpen={showConnectionModal}
                    onClose={() => setShowConnectionModal(false)}
                    platform="twitch"
                    getToken={getToken}
                />
            </div>
        );
    }

    // Separate videos by type
    const broadcasts = videos.filter(video => video.type === 'archive' || video.type === 'highlight' || video.type === 'upload');
    const clips = videos.filter(video => video.type === 'clip');

    const renderVideoGrid = (videoList: TwitchVideo[], startIndex: number = 0) => (
        <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
            {videoList.map((video, index) => (
                <motion.div
                    key={`content-${video.id}-${index}`}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: (startIndex + index) * 0.05 }}
                    className="bg-gray-900/60 backdrop-blur-xl border border-gray-800/50 rounded-xl overflow-hidden hover:border-emerald-500/40 transition-all duration-300 group shadow-lg hover:shadow-2xl hover:shadow-emerald-500/10 hover:-translate-y-1"
                >
                    {/* Thumbnail */}
                    <div className="relative aspect-video overflow-hidden">
                        <Image
                            src={video.thumbnail_url.replace('%{width}', '320').replace('%{height}', '180')}
                            alt={video.title}
                            fill
                            sizes="(max-width: 768px) 100vw, (max-width: 1024px) 33vw, 20vw"
                            className="object-cover group-hover:scale-110 transition-transform duration-500"
                        />
                        <div className="absolute inset-0 bg-gradient-to-t from-black/40 via-transparent to-transparent group-hover:from-black/20 transition-all duration-300" />
                        
                        {/* Duration badge */}
                        <div className="absolute top-2 right-2 bg-black/90 text-white text-xs px-2 py-1 rounded-lg flex items-center gap-1 shadow-lg">
                            <Clock className="w-3 h-3" />
                            {formatDuration(video.duration)}
                        </div>
                        
                        {/* Video type badge */}
                        <div className="absolute top-2 left-2">
                            <span className={`text-xs px-2 py-1 rounded-lg font-semibold backdrop-blur-md shadow-lg ${video.type === 'clip'
                                ? 'bg-purple-500/90 text-white border border-purple-400/30' 
                                : 'bg-blue-500/90 text-white border border-blue-400/30'
                                }`}>
                                {video.type === 'clip' ? 'CLIP' : 'VOD'}
                            </span>
                        </div>
                        
                        {/* Play overlay */}
                        <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-all duration-300">
                            <div className="w-14 h-14 bg-emerald-500/30 backdrop-blur-md rounded-full flex items-center justify-center border-2 border-emerald-400/50 shadow-xl group-hover:scale-110 transition-transform duration-300">
                                <PlayCircle className="w-7 h-7 text-white drop-shadow-lg" />
                            </div>
                        </div>

                        {/* View count overlay */}
                        <div className="absolute bottom-2 left-2 bg-black/90 text-white text-xs px-2 py-1 rounded-lg flex items-center gap-1 shadow-lg">
                            <Eye className="w-3 h-3" />
                            <span className="font-medium">{formatViewCount(video.view_count)}</span>
                        </div>
                    </div>

                    {/* Content */}
                    <div className="p-4">
                        <h3 className="font-semibold text-white mb-2 line-clamp-2 leading-tight text-sm group-hover:text-emerald-300 transition-colors">
                            {video.title}
                        </h3>

                        <div className="flex items-center justify-between text-xs text-gray-400 mb-3">
                            <span>{formatDate(video.created_at)}</span>
                        </div>

                        <div className="flex gap-2">
                            <button
                                onClick={() => window.open(video.url, '_blank')}
                                className="flex-1 bg-emerald-600 hover:bg-emerald-700 text-white px-3 py-1.5 rounded-lg transition-all duration-200 font-medium flex items-center justify-center gap-1 text-xs hover:shadow-lg hover:shadow-emerald-500/20"
                            >
                                <ExternalLink className="w-3 h-3" />
                                View
                            </button>
                            <button className="flex-1 bg-gray-700 hover:bg-gray-600 text-white px-3 py-1.5 rounded-lg transition-all duration-200 font-medium flex items-center justify-center gap-1 text-xs hover:shadow-lg">
                                <Download className="w-3 h-3" />
                                Download
                            </button>
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
            className="text-center py-16"
        >
            <div className="w-16 h-16 bg-gray-700/20 rounded-full flex items-center justify-center mx-auto mb-4">
                {type === 'broadcasts' ? (
                    <Monitor className="w-8 h-8 text-gray-600" />
                ) : (
                    <Film className="w-8 h-8 text-gray-600" />
                )}
            </div>
            <h4 className="text-lg font-medium text-gray-300 mb-2">
                No {type} found
            </h4>
            <p className="text-gray-500">
                {type === 'broadcasts'
                    ? 'Your past broadcasts and highlights will appear here.'
                    : 'Your clips will appear here once created.'}
            </p>
        </motion.div>
    );

    if (videos.length === 0) {
        return (
            <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 p-6">
                <div className="max-w-7xl mx-auto">
                    <div className="flex items-center justify-between mb-8">
                        <div>
                            <div className="flex items-center mb-2">
                                <Video className="w-6 h-6 text-emerald-500 mr-2" />
                                <span className="text-sm text-emerald-400 font-medium">Your Content</span>
                            </div>
                            <h1 className="text-4xl font-bold mb-2">
                                <span className="text-white">Your </span>
                                <span className="bg-gradient-to-r from-emerald-400 to-blue-400 bg-clip-text text-transparent">Videos</span>
                            </h1>
                            <p className="text-gray-400 text-lg">Manage and analyze your Twitch content library.</p>
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

                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        className="bg-gray-900/50 backdrop-blur-xl border border-gray-800/50 rounded-2xl p-16 text-center"
                    >
                        <div className="w-20 h-20 bg-emerald-500/20 rounded-full flex items-center justify-center mx-auto mb-6">
                            <Video className="w-10 h-10 text-emerald-400" />
                        </div>
                        <h3 className="text-2xl font-semibold text-white mb-4">
                            No videos found
                        </h3>
                        <p className="text-gray-400 mb-8 max-w-md mx-auto">
                            Your Twitch videos will appear here once they're available. Make sure your Twitch account is connected and you have recent content.
                        </p>
                        <button
                            onClick={() => window.open('https://twitch.tv', '_blank')}
                            className="bg-emerald-600 hover:bg-emerald-700 text-white px-6 py-3 rounded-lg font-semibold transition-colors flex items-center gap-2 mx-auto"
                        >
                            <ExternalLink className="w-4 h-4" />
                            View on Twitch
                        </button>
                    </motion.div>
                </div>
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
                            <Video className="w-6 h-6 text-emerald-500 mr-2" />
                            <span className="text-sm text-emerald-400 font-medium">Your Content</span>
                        </div>
                        <h1 className="text-4xl font-bold mb-2">
                            <span className="text-white">Your </span>
                            <span className="bg-gradient-to-r from-emerald-400 to-blue-400 bg-clip-text text-transparent">Videos</span>
                        </h1>
                        <p className="text-gray-400 text-lg">Manage and analyze your Twitch content. Download clips, view analytics, and optimize your performance.</p>
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

                {/* Content Summary Cards */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                    <motion.div
                        initial={{ y: 20, opacity: 0 }}
                        animate={{ y: 0, opacity: 1 }}
                        transition={{ delay: 0.1 }}
                        className="bg-gray-900/50 backdrop-blur-xl border border-gray-800/50 rounded-2xl p-6 relative overflow-hidden group hover:border-emerald-500/30 transition-all duration-300"
                    >
                        <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/10 to-teal-500/5 opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
                        <div className="relative z-10">
                            <div className="flex items-center text-emerald-400 mb-4">
                                <Video className="w-5 h-5 mr-2" />
                                <h3 className="text-sm font-medium text-gray-400">Total Content</h3>
                            </div>
                            <p className="text-3xl font-bold text-white mb-2">
                                {videos.length}
                            </p>
                            <p className="text-sm text-gray-500">Videos in your library</p>
                        </div>
                    </motion.div>

                    <motion.div
                        initial={{ y: 20, opacity: 0 }}
                        animate={{ y: 0, opacity: 1 }}
                        transition={{ delay: 0.2 }}
                        className="bg-gray-900/50 backdrop-blur-xl border border-gray-800/50 rounded-2xl p-6 relative overflow-hidden group hover:border-purple-500/30 transition-all duration-300"
                    >
                        <div className="absolute inset-0 bg-gradient-to-br from-purple-500/10 to-pink-500/5 opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
                        <div className="relative z-10">
                            <div className="flex items-center text-purple-400 mb-4">
                                <Monitor className="w-5 h-5 mr-2" />
                                <h3 className="text-sm font-medium text-gray-400">Broadcasts</h3>
                            </div>
                            <p className="text-3xl font-bold text-white mb-2">
                                {broadcasts.length}
                            </p>
                            <p className="text-sm text-gray-500">VODs & highlights</p>
                        </div>
                    </motion.div>

                    <motion.div
                        initial={{ y: 20, opacity: 0 }}
                        animate={{ y: 0, opacity: 1 }}
                        transition={{ delay: 0.3 }}
                        className="bg-gray-900/50 backdrop-blur-xl border border-gray-800/50 rounded-2xl p-6 relative overflow-hidden group hover:border-blue-500/30 transition-all duration-300"
                    >
                        <div className="absolute inset-0 bg-gradient-to-br from-blue-500/10 to-cyan-500/5 opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
                        <div className="relative z-10">
                            <div className="flex items-center text-blue-400 mb-4">
                                <Film className="w-5 h-5 mr-2" />
                                <h3 className="text-sm font-medium text-gray-400">Clips</h3>
                            </div>
                            <p className="text-3xl font-bold text-white mb-2">
                                {clips.length}
                            </p>
                            <p className="text-sm text-gray-500">Community clips</p>
                        </div>
                    </motion.div>
                </div>

                {/* Content Sections */}
                <div className="space-y-12">
                    {/* Recent Broadcasts Section */}
                    <div>
                        <motion.div
                            initial={{ y: 20, opacity: 0 }}
                            animate={{ y: 0, opacity: 1 }}
                            transition={{ delay: 0.4 }}
                            className="flex items-center justify-between mb-6"
                        >
                            <div className="flex items-center gap-3">
                                <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-purple-500/20 border border-purple-500/30 text-purple-400 text-sm">
                                    <Monitor className="w-4 h-4" />
                                    <span>Recent Broadcasts</span>
                                </div>
                                {broadcasts.length > 0 && (
                                    <span className="text-sm text-gray-500">
                                        {broadcasts.length} video{broadcasts.length !== 1 ? 's' : ''}
                                    </span>
                                )}
                            </div>
                            {broadcasts.length > 0 && (
                                <div className="flex items-center text-sm text-gray-400">
                                    <TrendingUp className="w-4 h-4 mr-1" />
                                    <span>{broadcasts.reduce((acc, video) => acc + video.view_count, 0).toLocaleString()} total views</span>
                                </div>
                            )}
                        </motion.div>
                        {broadcasts.length > 0 ? renderVideoGrid(broadcasts, 0) : renderEmptyState('broadcasts')}
                    </div>

                    {/* Clips Section */}
                    <div>
                        <motion.div
                            initial={{ y: 20, opacity: 0 }}
                            animate={{ y: 0, opacity: 1 }}
                            transition={{ delay: 0.5 }}
                            className="flex items-center justify-between mb-6"
                        >
                            <div className="flex items-center gap-3">
                                <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-blue-500/20 border border-blue-500/30 text-blue-400 text-sm">
                                    <Film className="w-4 h-4" />
                                    <span>Clips</span>
                                </div>
                                {clips.length > 0 && (
                                    <span className="text-sm text-gray-500">
                                        {clips.length} clip{clips.length !== 1 ? 's' : ''}
                                    </span>
                                )}
                            </div>
                            {clips.length > 0 && (
                                <div className="flex items-center text-sm text-gray-400">
                                    <TrendingUp className="w-4 h-4 mr-1" />
                                    <span>{clips.reduce((acc, video) => acc + video.view_count, 0).toLocaleString()} total views</span>
                                </div>
                            )}
                        </motion.div>
                        {clips.length > 0 ? renderVideoGrid(clips, broadcasts.length) : renderEmptyState('clips')}
                    </div>
                </div>

                {/* Connection Modal */}
                <ConnectionModal
                    isOpen={showConnectionModal}
                    onClose={() => setShowConnectionModal(false)}
                    platform="twitch"
                    getToken={getToken}
                />
            </div>
        </div>
    );
} 
