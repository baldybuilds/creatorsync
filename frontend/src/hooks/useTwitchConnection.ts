import { useState, useEffect, useCallback } from 'react';
import { useAuth } from '@clerk/nextjs';
import { twitchService, TwitchApiError, TwitchService } from '@/services/twitch';

export interface TwitchConnectionState {
    isConnected: boolean;
    isLoading: boolean;
    error: string | null;
    needsReconnect: boolean;
    needsConnect: boolean;
}

export function useTwitchConnection() {
    const { getToken } = useAuth();
    const [state, setState] = useState<TwitchConnectionState>({
        isConnected: false,
        isLoading: true,
        error: null,
        needsReconnect: false,
        needsConnect: false,
    });

    // Check connection status
    const checkStatus = useCallback(async () => {
        try {
            setState(prev => ({ ...prev, isLoading: true, error: null }));
            
            const token = await getToken();
            if (!token) {
                setState(prev => ({ 
                    ...prev, 
                    isLoading: false, 
                    isConnected: false,
                    needsConnect: true 
                }));
                return;
            }

            const status = await twitchService.getConnectionStatus(token);
            setState(prev => ({
                ...prev,
                isConnected: status.connected,
                isLoading: false,
                needsConnect: !status.connected,
                needsReconnect: false,
            }));
        } catch (error) {
            console.error('Failed to check Twitch connection status:', error);
            setState(prev => ({
                ...prev,
                isLoading: false,
                error: 'Failed to check connection status',
                isConnected: false,
                needsConnect: true,
            }));
        }
    }, [getToken]);

    // Connect to Twitch
    const connect = useCallback(async () => {
        try {
            setState(prev => ({ ...prev, isLoading: true, error: null }));
            
            const token = await getToken();
            if (!token) {
                setState(prev => ({ ...prev, error: 'Not authenticated', isLoading: false }));
                return;
            }

            await twitchService.connectToTwitch(token);
            // Note: If successful, the browser will redirect to Twitch OAuth, so we won't reach here
        } catch (error) {
            console.error('Failed to initiate Twitch connection:', error);
            setState(prev => ({ 
                ...prev, 
                error: error instanceof Error ? error.message : 'Failed to connect to Twitch',
                isLoading: false 
            }));
        }
    }, [getToken]);

    // Disconnect from Twitch
    const disconnect = useCallback(async () => {
        try {
            setState(prev => ({ ...prev, isLoading: true, error: null }));
            
            const token = await getToken();
            if (!token) {
                setState(prev => ({ ...prev, error: 'Not authenticated' }));
                return;
            }

            await twitchService.disconnectFromTwitch(token);
            setState(prev => ({
                ...prev,
                isConnected: false,
                isLoading: false,
                needsConnect: true,
                needsReconnect: false,
            }));
        } catch (error) {
            console.error('Failed to disconnect from Twitch:', error);
            setState(prev => ({
                ...prev,
                isLoading: false,
                error: 'Failed to disconnect from Twitch',
            }));
        }
    }, [getToken]);

    // Handle API errors (for use in other components)
    const handleApiError = useCallback((error: Error & TwitchApiError) => {
        if (error.reconnect_required) {
            setState(prev => ({
                ...prev,
                needsReconnect: true,
                needsConnect: false,
                isConnected: false,
                error: 'Twitch re-authentication required',
            }));
        } else if (error.connect_required) {
            setState(prev => ({
                ...prev,
                needsConnect: true,
                needsReconnect: false,
                isConnected: false,
                error: 'Twitch account not connected',
            }));
        } else {
            setState(prev => ({
                ...prev,
                error: error.message || 'An error occurred',
            }));
        }
    }, []);

    // Clear error
    const clearError = useCallback(() => {
        setState(prev => ({ ...prev, error: null }));
    }, []);

    // Handle OAuth callback on component mount
    useEffect(() => {
        const callback = TwitchService.handleOAuthCallback();
        
        if (callback.success) {
            // Connection successful
            setState(prev => ({
                ...prev,
                isConnected: true,
                needsConnect: false,
                needsReconnect: false,
            }));
            
            // Trigger data prefetch for faster loading
            setTimeout(async () => {
                try {
                    const token = await getToken();
                    if (token && typeof window !== 'undefined') {
                        // Pre-fetch common data that users typically view first
                        fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}/api/twitch/videos?limit=10`, {
                            headers: { 'Authorization': `Bearer ${token}` }
                        }).catch(() => {});
                        fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}/api/twitch/channel`, {
                            headers: { 'Authorization': `Bearer ${token}` }
                        }).catch(() => {});
                    }
                } catch {
                    // Silently fail - this is just for prefetching
                }
            }, 1000); // Small delay to ensure connection is fully established
        } else if (callback.error) {
            // Connection failed
            const errorMessage = TwitchService.getErrorMessage(callback.error);
            setState(prev => ({
                ...prev,
                error: errorMessage,
                isConnected: false,
                needsConnect: true,
            }));
        }
    }, [getToken]);

    // Check status on mount
    useEffect(() => {
        checkStatus();
    }, [checkStatus]);

    return {
        ...state,
        connect,
        disconnect,
        checkStatus,
        handleApiError,
        clearError,
    };
} 