'use client';

import { motion } from 'framer-motion';
import { Button } from '@/components/ui/button';
import { Calendar, Clock, Plus, Edit3, Trash2 } from 'lucide-react';

export function ScheduledSection() {
    const scheduledPosts = [
        {
            id: 1,
            title: "Stream starting soon! 🔥",
            platform: "Twitter",
            scheduledTime: "2024-01-15T20:00:00Z",
            status: "pending"
        },
        {
            id: 2,
            title: "New highlight video is live!",
            platform: "Instagram",
            scheduledTime: "2024-01-16T14:30:00Z",
            status: "pending"
        }
    ];

    const formatDateTime = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    return (
        <div className="p-8">
            <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-purple-500/20 border border-purple-500/30 text-purple-500 text-sm mb-6 w-fit">
                <Calendar className="w-4 h-4" />
                <span>Scheduled Content</span>
            </div>

            <div className="flex items-center justify-between mb-6">
                <div>
                    <h2 className="text-4xl font-bold mb-2">
                        <span className="text-light-surface-900 dark:text-dark-surface-100">Scheduled</span>{' '}
                        <span className="text-gradient">Posts</span>
                    </h2>
                    <p className="text-xl text-light-surface-700 dark:text-dark-surface-300">
                        Manage your upcoming social media posts and content releases.
                    </p>
                </div>
                <Button variant="default" className="flex items-center gap-2">
                    <Plus className="w-4 h-4" />
                    Schedule Post
                </Button>
            </div>

            {scheduledPosts.length === 0 ? (
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="text-center py-16"
                >
                    <Calendar className="w-16 h-16 text-light-surface-400 dark:text-dark-surface-600 mx-auto mb-4" />
                    <h3 className="text-xl font-semibold text-light-surface-700 dark:text-dark-surface-300 mb-2">
                        No scheduled posts
                    </h3>
                    <p className="text-light-surface-600 dark:text-dark-surface-400 mb-6">
                        Create your first scheduled post to automate your social media presence.
                    </p>
                    <Button variant="default">
                        <Plus className="w-4 h-4 mr-2" />
                        Schedule Your First Post
                    </Button>
                </motion.div>
            ) : (
                <div className="space-y-4">
                    {scheduledPosts.map((post, index) => (
                        <motion.div
                            key={post.id}
                            initial={{ opacity: 0, y: 20 }}
                            animate={{ opacity: 1, y: 0 }}
                            transition={{ delay: index * 0.1 }}
                            className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm p-6 rounded-xl shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50 hover:border-brand-500/30 transition-all duration-300"
                        >
                            <div className="flex items-center justify-between">
                                <div className="flex-1">
                                    <div className="flex items-center gap-3 mb-2">
                                        <div className="px-3 py-1 bg-brand-500/20 text-brand-500 rounded-full text-xs font-medium">
                                            {post.platform}
                                        </div>
                                        <div className="flex items-center gap-1 text-sm text-light-surface-600 dark:text-dark-surface-400">
                                            <Clock className="w-4 h-4" />
                                            {formatDateTime(post.scheduledTime)}
                                        </div>
                                    </div>
                                    <h3 className="text-lg font-semibold text-light-surface-900 dark:text-dark-surface-100 mb-1">
                                        {post.title}
                                    </h3>
                                    <p className="text-sm text-light-surface-600 dark:text-dark-surface-400">
                                        Status: <span className="capitalize text-amber-500">{post.status}</span>
                                    </p>
                                </div>
                                <div className="flex items-center gap-2 ml-4">
                                    <Button variant="ghost" size="sm">
                                        <Edit3 className="w-4 h-4" />
                                    </Button>
                                    <Button variant="ghost" size="sm" className="text-red-500 hover:text-red-600">
                                        <Trash2 className="w-4 h-4" />
                                    </Button>
                                </div>
                            </div>
                        </motion.div>
                    ))}
                </div>
            )}

            {/* Quick Actions */}
            <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.3 }}
                className="mt-10 bg-gradient-to-br from-purple-500/10 to-purple-600/5 border border-purple-500/20 p-6 rounded-xl"
            >
                <h3 className="text-lg font-semibold text-purple-500 mb-4">Quick Schedule Templates</h3>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <Button variant="outline" className="justify-start h-auto p-4">
                        <div className="text-left">
                            <div className="font-medium">Stream Announcement</div>
                            <div className="text-xs text-light-surface-600 dark:text-dark-surface-400 mt-1">
                                Let followers know you're going live
                            </div>
                        </div>
                    </Button>
                    <Button variant="outline" className="justify-start h-auto p-4">
                        <div className="text-left">
                            <div className="font-medium">Highlight Reel</div>
                            <div className="text-xs text-light-surface-600 dark:text-dark-surface-400 mt-1">
                                Share your best moments
                            </div>
                        </div>
                    </Button>
                    <Button variant="outline" className="justify-start h-auto p-4">
                        <div className="text-left">
                            <div className="font-medium">Community Update</div>
                            <div className="text-xs text-light-surface-600 dark:text-dark-surface-400 mt-1">
                                Keep your audience engaged
                            </div>
                        </div>
                    </Button>
                </div>
            </motion.div>
        </div>
    );
} 