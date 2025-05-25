import { useCallback } from 'react';
import { useAuth } from '@clerk/nextjs';
import { twitchService, TwitchApiError } from '@/services/twitch';
import { useTwitchConnection } from './useTwitchConnection';

export function useTwitchApi() {
    const { getToken } = useAuth();
    const { handleApiError } = useTwitchConnection();

    const makeApiCall = useCallback(async <T>(
        apiCall: (token: string) => Promise<T>
    ): Promise<T | null> => {
        try {
            const token = await getToken();
            if (!token) {
                throw new Error('Not authenticated');
            }

            return await apiCall(token);
        } catch (error) {
            const twitchError = error as Error & TwitchApiError;
            handleApiError(twitchError);
            return null;
        }
    }, [getToken, handleApiError]);

    // Specific API methods
    const getChannelInfo = useCallback(() => {
        return makeApiCall(token => twitchService.getChannelInfo(token));
    }, [makeApiCall]);

    const getVideos = useCallback((type: string = 'archive', limit: number = 20) => {
        return makeApiCall(token => twitchService.getVideos(token, type, limit));
    }, [makeApiCall]);

    const getSubscribers = useCallback(() => {
        return makeApiCall(token => twitchService.getSubscribers(token));
    }, [makeApiCall]);

    return {
        makeApiCall,
        getChannelInfo,
        getVideos,
        getSubscribers,
    };
} 