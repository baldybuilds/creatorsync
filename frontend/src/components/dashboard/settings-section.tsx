'use client';

import { motion } from 'framer-motion';
import { Button } from '@/components/ui/button';
import { Settings, Bell, Shield, Palette, User, ExternalLink } from 'lucide-react';

export function SettingsSection() {
    return (
        <div className="p-8">
            <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-gray-500/20 border border-gray-500/30 text-gray-500 text-sm mb-6 w-fit">
                <Settings className="w-4 h-4" />
                <span>Settings</span>
            </div>

            <h2 className="text-4xl font-bold mb-6">
                <span className="text-light-surface-900 dark:text-dark-surface-100">Account</span>{' '}
                <span className="text-gradient">Settings</span>
            </h2>

            <p className="text-xl text-light-surface-700 dark:text-dark-surface-300 max-w-3xl mb-10">
                Manage your account preferences, notifications, and platform integrations.
            </p>

            <div className="space-y-6 max-w-4xl">
                {/* Profile Settings */}
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.1 }}
                    className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm p-6 rounded-xl shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50"
                >
                    <div className="flex items-center gap-3 mb-4">
                        <User className="w-6 h-6 text-brand-500" />
                        <h3 className="text-xl font-semibold text-light-surface-900 dark:text-dark-surface-100">Profile Settings</h3>
                    </div>
                    <p className="text-light-surface-600 dark:text-dark-surface-400 mb-4">
                        Update your profile information and display preferences.
                    </p>
                    <Button variant="outline" size="sm">
                        Edit Profile
                    </Button>
                </motion.div>

                {/* Notifications */}
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.2 }}
                    className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm p-6 rounded-xl shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50"
                >
                    <div className="flex items-center gap-3 mb-4">
                        <Bell className="w-6 h-6 text-brand-500" />
                        <h3 className="text-xl font-semibold text-light-surface-900 dark:text-dark-surface-100">Notifications</h3>
                    </div>
                    <p className="text-light-surface-600 dark:text-dark-surface-400 mb-4">
                        Configure email and push notification preferences.
                    </p>
                    <Button variant="outline" size="sm">
                        Manage Notifications
                    </Button>
                </motion.div>

                {/* Privacy & Security */}
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.3 }}
                    className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm p-6 rounded-xl shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50"
                >
                    <div className="flex items-center gap-3 mb-4">
                        <Shield className="w-6 h-6 text-brand-500" />
                        <h3 className="text-xl font-semibold text-light-surface-900 dark:text-dark-surface-100">Privacy & Security</h3>
                    </div>
                    <p className="text-light-surface-600 dark:text-dark-surface-400 mb-4">
                        Control your privacy settings and security preferences.
                    </p>
                    <Button variant="outline" size="sm">
                        Security Settings
                    </Button>
                </motion.div>

                {/* Appearance */}
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.4 }}
                    className="bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm p-6 rounded-xl shadow-md border border-light-surface-200/50 dark:border-dark-surface-800/50"
                >
                    <div className="flex items-center gap-3 mb-4">
                        <Palette className="w-6 h-6 text-brand-500" />
                        <h3 className="text-xl font-semibold text-light-surface-900 dark:text-dark-surface-100">Appearance</h3>
                    </div>
                    <p className="text-light-surface-600 dark:text-dark-surface-400 mb-4">
                        Customize the look and feel of your dashboard.
                    </p>
                    <Button variant="outline" size="sm">
                        Theme Settings
                    </Button>
                </motion.div>

                {/* Connected Accounts */}
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.5 }}
                    className="bg-gradient-to-br from-brand-500/10 to-brand-600/5 border border-brand-500/20 p-6 rounded-xl"
                >
                    <div className="flex items-center gap-3 mb-4">
                        <ExternalLink className="w-6 h-6 text-brand-500" />
                        <h3 className="text-xl font-semibold text-brand-500">Connected Accounts</h3>
                    </div>
                    <p className="text-light-surface-700 dark:text-dark-surface-300 mb-6">
                        Manage your connected social media accounts and streaming platforms.
                    </p>
                    <div className="space-y-3">
                        <div className="flex items-center justify-between p-3 bg-white/50 dark:bg-black/20 rounded-lg">
                            <div className="flex items-center gap-3">
                                <div className="w-8 h-8 rounded bg-purple-500 flex items-center justify-center">
                                    <span className="text-white text-xs font-bold">T</span>
                                </div>
                                <span className="font-medium text-light-surface-900 dark:text-dark-surface-100">Twitch</span>
                            </div>
                            <Button variant="outline" size="sm">
                                Connected
                            </Button>
                        </div>
                        <div className="flex items-center justify-between p-3 bg-white/30 dark:bg-black/10 rounded-lg opacity-60">
                            <div className="flex items-center gap-3">
                                <div className="w-8 h-8 rounded bg-red-500 flex items-center justify-center">
                                    <span className="text-white text-xs font-bold">Y</span>
                                </div>
                                <span className="font-medium text-light-surface-900 dark:text-dark-surface-100">YouTube</span>
                            </div>
                            <Button variant="ghost" size="sm">
                                Connect
                            </Button>
                        </div>
                    </div>
                </motion.div>
            </div>
        </div>
    );
} 