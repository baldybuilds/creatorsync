'use client';

import { motion } from 'framer-motion';
import { Button } from '@/components/ui/button';
import { Users, BarChart2, Video, Calendar, Sparkles, ArrowRight } from 'lucide-react';

export function OverviewSection() {
    return (
        <div className="p-8">
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

            {/* Stats Grid */}
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

            {/* Recent Activity */}
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

            {/* Creator Tip */}
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
} 