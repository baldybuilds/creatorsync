'use client';

import { useState } from 'react';
import { motion } from 'framer-motion';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { AlertCircle, CheckCircle, ExternalLink, Loader2, RefreshCw, Trash2 } from 'lucide-react';
import { useTwitchConnection } from '@/hooks/useTwitchConnection';

interface TwitchConnectionManagerProps {
    showTitle?: boolean;
    compact?: boolean;
    className?: string;
}

export function TwitchConnectionManager({ 
    showTitle = false, 
    compact = false, 
    className = '' 
}: TwitchConnectionManagerProps) {
    const {
        isConnected,
        isLoading,
        error,
        needsReconnect,

        connect,
        disconnect,
        checkStatus,
        clearError
    } = useTwitchConnection();

    const [isDisconnecting, setIsDisconnecting] = useState(false);

    const handleConnect = async () => {
        clearError();
        await connect();
    };

    const handleDisconnect = async () => {
        if (window.confirm('Are you sure you want to disconnect your Twitch account? This will stop all analytics collection.')) {
            setIsDisconnecting(true);
            try {
                await disconnect();
            } finally {
                setIsDisconnecting(false);
            }
        }
    };

    const handleRefresh = async () => {
        clearError();
        await checkStatus();
    };

    const getStatusBadge = () => {
        if (isLoading) {
            return (
                <Badge variant="outline" className="gap-2">
                    <Loader2 className="w-3 h-3 animate-spin" />
                    Checking...
                </Badge>
            );
        }

        if (isConnected) {
            return (
                <Badge variant="default" className="gap-2 bg-green-500 hover:bg-green-600">
                    <CheckCircle className="w-3 h-3" />
                    Connected
                </Badge>
            );
        }

        if (needsReconnect) {
            return (
                <Badge variant="destructive" className="gap-2">
                    <AlertCircle className="w-3 h-3" />
                    Reconnect Required
                </Badge>
            );
        }

        return (
            <Badge variant="outline" className="gap-2">
                <ExternalLink className="w-3 h-3" />
                Not Connected
            </Badge>
        );
    };

    const getActionButton = () => {
        if (isConnected) {
            return (
                <div className="flex gap-2">
                    <Button
                        variant="outline"
                        size={compact ? "sm" : "default"}
                        onClick={handleRefresh}
                        disabled={isLoading}
                        className="gap-2"
                    >
                        <RefreshCw className={`w-4 h-4 ${isLoading ? 'animate-spin' : ''}`} />
                        {!compact && 'Refresh'}
                    </Button>
                    <Button
                        variant="destructive"
                        size={compact ? "sm" : "default"}
                        onClick={handleDisconnect}
                        disabled={isDisconnecting}
                        className="gap-2"
                    >
                        {isDisconnecting ? (
                            <Loader2 className="w-4 h-4 animate-spin" />
                        ) : (
                            <Trash2 className="w-4 h-4" />
                        )}
                        {!compact && 'Disconnect'}
                    </Button>
                </div>
            );
        }

        return (
            <Button
                onClick={handleConnect}
                disabled={isLoading}
                size={compact ? "sm" : "default"}
                className="gap-2"
                variant={needsReconnect ? "destructive" : "default"}
            >
                {isLoading ? (
                    <Loader2 className="w-4 h-4 animate-spin" />
                ) : (
                    <ExternalLink className="w-4 h-4" />
                )}
                {needsReconnect ? 'Reconnect' : 'Connect'} to Twitch
            </Button>
        );
    };

    const getErrorMessage = () => {
        if (!error) return null;

        return (
            <motion.div
                initial={{ opacity: 0, y: -10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -10 }}
                className="mt-3 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg"
            >
                <div className="flex items-start gap-2">
                    <AlertCircle className="w-4 h-4 text-red-500 mt-0.5 flex-shrink-0" />
                    <div className="flex-1">
                        <p className="text-sm text-red-700 dark:text-red-300">{error}</p>
                        <Button
                            variant="link"
                            size="sm"
                            onClick={clearError}
                            className="h-auto p-0 text-red-600 dark:text-red-400"
                        >
                            Dismiss
                        </Button>
                    </div>
                </div>
            </motion.div>
        );
    };

    if (compact) {
        return (
            <div className={`flex items-center gap-3 ${className}`}>
                {getStatusBadge()}
                {getActionButton()}
                {getErrorMessage()}
            </div>
        );
    }

    return (
        <div className={className}>
            {showTitle && (
                <div className="flex items-center gap-3 mb-4">
                    <div className="w-8 h-8 rounded bg-purple-500 flex items-center justify-center">
                        <span className="text-white text-xs font-bold">T</span>
                    </div>
                    <h3 className="text-xl font-semibold text-light-surface-900 dark:text-dark-surface-100">
                        Twitch Integration
                    </h3>
                </div>
            )}

            <div className="flex items-start justify-between gap-4">
                <div className="flex-1">
                    <div className="flex items-center gap-3 mb-2">
                        <span className="font-medium text-light-surface-900 dark:text-dark-surface-100">
                            Twitch
                        </span>
                        {getStatusBadge()}
                    </div>
                    
                    <p className="text-sm text-light-surface-600 dark:text-dark-surface-400">
                        {isConnected 
                            ? 'Your Twitch account is connected and analytics are being collected.'
                            : needsReconnect 
                                ? 'Your Twitch connection has expired. Please reconnect to continue analytics collection.'
                                : 'Connect your Twitch account to start collecting analytics data.'
                        }
                    </p>
                </div>

                <div className="flex-shrink-0">
                    {getActionButton()}
                </div>
            </div>

            {getErrorMessage()}
        </div>
    );
} 