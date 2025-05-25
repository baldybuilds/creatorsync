import React, { useState, useEffect } from 'react';
import { createPortal } from 'react-dom';
import { X, Shield, Eye, Video, Users, Settings, ExternalLink, CheckCircle } from 'lucide-react';
import { Button } from './button';
import { initiateOAuth } from '@/utils/oauth';

interface Permission {
    icon: React.ComponentType<{ className?: string }>;
    title: string;
    description: string;
    required: boolean;
}

interface ConnectionModalProps {
    isOpen: boolean;
    onClose: () => void;
    platform: 'twitch' | 'youtube' | 'tiktok';
    getToken?: () => Promise<string | null>;
}

const twitchPermissions: Permission[] = [
    {
        icon: Users,
        title: 'Follower & Subscriber Data',
        description: 'View your follower count and subscriber list to provide accurate analytics',
        required: true,
    },
    {
        icon: Video,
        title: 'Content Management',
        description: 'Access your videos, clips, and stream archives for content analysis',
        required: true,
    },
    {
        icon: Eye,
        title: 'Channel Analytics',
        description: 'View stream performance, viewership data, and engagement metrics',
        required: true,
    },
    {
        icon: Settings,
        title: 'Channel Information',
        description: 'Read basic channel info like display name and profile settings',
        required: true,
    },
];

const platformConfig = {
    twitch: {
        name: 'Twitch',
        color: '#9146FF',
        bgGradient: 'from-purple-600 to-purple-700',
        permissions: twitchPermissions,
    },
    youtube: {
        name: 'YouTube',
        color: '#FF0000',
        bgGradient: 'from-red-600 to-red-700',
        permissions: [],
    },
    tiktok: {
        name: 'TikTok',
        color: '#000000',
        bgGradient: 'from-gray-800 to-black',
        permissions: [],
    },
};

export function ConnectionModal({ isOpen, onClose, platform, getToken }: ConnectionModalProps) {
    const [step, setStep] = useState<'permissions' | 'connecting' | 'success'>('permissions');
    const [mounted, setMounted] = useState(false);
    const config = platformConfig[platform];

    // Ensure this only renders on the client side
    useEffect(() => {
        setMounted(true);
        
        // Create modal root if it doesn't exist
        if (!document.getElementById('modal-root')) {
            const modalRoot = document.createElement('div');
            modalRoot.id = 'modal-root';
            document.body.appendChild(modalRoot);
        }

        return () => {
            // Clean up on unmount
            const modalRoot = document.getElementById('modal-root');
            if (modalRoot && modalRoot.children.length === 0) {
                document.body.removeChild(modalRoot);
            }
        };
    }, []);

    const handleConnect = async () => {
        setStep('connecting');
        
        try {
            await initiateOAuth(platform, getToken);
        } catch (error) {
            console.error('OAuth initiation failed:', error);
            setStep('permissions');
        }
    };

    const handleClose = () => {
        if (step !== 'connecting') {
            onClose();
            setStep('permissions');
        }
    };

    // Don't render anything on the server or before mounting
    if (!mounted || !isOpen) {
        return null;
    }

    const modalRoot = document.getElementById('modal-root');
    if (!modalRoot) {
        console.error('Modal root not found');
        return null;
    }

    return createPortal(
        <div className="fixed inset-0 z-[9999] flex items-center justify-center">
            {/* Backdrop */}
            <div
                className="absolute inset-0 bg-black/60 backdrop-blur-sm"
                onClick={handleClose}
            />
            
            {/* Modal Content */}
            <div
                className="relative bg-white dark:bg-gray-900 rounded-2xl shadow-2xl border border-gray-200 dark:border-gray-800 w-full max-w-md mx-4 overflow-hidden"
                onClick={(e) => e.stopPropagation()}
            >
                {step === 'permissions' && (
                    <div className="p-6">
                        {/* Header */}
                        <div className="flex items-center justify-between mb-6">
                            <div className="flex items-center gap-3">
                                <div className={`w-10 h-10 bg-gradient-to-r ${config.bgGradient} rounded-xl flex items-center justify-center`}>
                                    <Video className="w-5 h-5 text-white" />
                                </div>
                                <div>
                                    <h2 className="text-xl font-bold text-gray-900 dark:text-white">
                                        Connect {config.name}
                                    </h2>
                                    <p className="text-sm text-gray-500 dark:text-gray-400">
                                        Secure OAuth connection
                                    </p>
                                </div>
                            </div>
                            <Button
                                variant="ghost"
                                size="sm"
                                onClick={handleClose}
                                className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200"
                            >
                                <X className="w-4 h-4" />
                            </Button>
                        </div>

                        {/* Security Notice */}
                        <div className="mb-6 p-4 bg-blue-50 dark:bg-blue-950/30 border border-blue-200 dark:border-blue-800 rounded-xl">
                            <div className="flex items-start gap-3">
                                <Shield className="w-5 h-5 text-blue-600 dark:text-blue-400 mt-0.5" />
                                <div>
                                    <h3 className="font-semibold text-blue-900 dark:text-blue-100 text-sm">
                                        Secure Connection
                                    </h3>
                                    <p className="text-blue-700 dark:text-blue-300 text-xs mt-1">
                                        We use OAuth 2.0 for secure authentication. Your login credentials are never stored by CreatorSync.
                                    </p>
                                </div>
                            </div>
                        </div>

                        {/* Permissions */}
                        <div className="mb-6">
                            <h3 className="font-semibold text-gray-900 dark:text-white mb-3 text-sm">
                                Required Permissions
                            </h3>
                            <div className="space-y-3">
                                {config.permissions.map((permission, index) => (
                                    <div key={index} className="flex items-start gap-3">
                                        <div className="w-8 h-8 bg-gray-100 dark:bg-gray-800 rounded-lg flex items-center justify-center flex-shrink-0">
                                            <permission.icon className="w-4 h-4 text-gray-600 dark:text-gray-400" />
                                        </div>
                                        <div className="flex-1">
                                            <p className="font-medium text-gray-900 dark:text-white text-sm">
                                                {permission.title}
                                            </p>
                                            <p className="text-gray-500 dark:text-gray-400 text-xs mt-1">
                                                {permission.description}
                                            </p>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>

                        {/* What happens next */}
                        <div className="mb-6 p-4 bg-gray-50 dark:bg-gray-800/50 rounded-xl">
                            <h3 className="font-semibold text-gray-900 dark:text-white text-sm mb-2">
                                What happens next?
                            </h3>
                            <ol className="text-xs text-gray-600 dark:text-gray-400 space-y-1">
                                <li>1. You'll be redirected to {config.name} to authorize</li>
                                <li>2. We'll securely connect your account</li>
                                <li>3. Your analytics will be available immediately</li>
                            </ol>
                        </div>

                        {/* Actions */}
                        <div className="flex gap-3">
                            <Button
                                variant="outline"
                                onClick={handleClose}
                                className="flex-1"
                            >
                                Cancel
                            </Button>
                            <Button
                                onClick={handleConnect}
                                className={`flex-1 bg-gradient-to-r ${config.bgGradient} hover:opacity-90 text-white border-0`}
                            >
                                <ExternalLink className="w-4 h-4 mr-2" />
                                Connect Account
                            </Button>
                        </div>
                    </div>
                )}

                {step === 'connecting' && (
                    <div className="p-8 text-center">
                        <div className="w-16 h-16 bg-gradient-to-r from-blue-500 to-purple-600 rounded-full flex items-center justify-center mx-auto mb-4">
                            <div className="w-6 h-6 border-2 border-white border-t-transparent rounded-full animate-spin" />
                        </div>
                        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
                            Connecting to {config.name}
                        </h3>
                        <p className="text-gray-500 dark:text-gray-400 text-sm">
                            Please complete the authorization in the {config.name} window
                        </p>
                    </div>
                )}

                {step === 'success' && (
                    <div className="p-8 text-center">
                        <div className="w-16 h-16 bg-green-500 rounded-full flex items-center justify-center mx-auto mb-4">
                            <CheckCircle className="w-8 h-8 text-white" />
                        </div>
                        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
                            Successfully Connected!
                        </h3>
                        <p className="text-gray-500 dark:text-gray-400 text-sm">
                            Your {config.name} account is now connected. Loading your analytics...
                        </p>
                    </div>
                )}
            </div>
        </div>,
        modalRoot
    );
} 