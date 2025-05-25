import { config } from '@/config';

export interface TwitchConnectionStatus {
    connected: boolean;
}

export interface TwitchApiError {
    error: string;
    reconnect_required?: boolean;
    connect_required?: boolean;
}

export interface TwitchUser {
    id: string;
    login: string;
    display_name: string;
    type: string;
    broadcaster_type: string;
    description: string;
    profile_image_url: string;
    offline_image_url: string;
    view_count: number;
    email?: string;
    created_at: string;
}

export interface TwitchChannelInfo {
    broadcaster_id: string;
    broadcaster_login: string;
    broadcaster_name: string;
    broadcaster_language: string;
    game_id: string;
    game_name: string;
    title: string;
    delay: number;
    tags: string[];
    content_classification_labels: string[];
    is_branded_content: boolean;
}

export interface TwitchVideo {
    id: string;
    stream_id?: string;
    user_id: string;
    user_login: string;
    user_name: string;
    title: string;
    description: string;
    created_at: string;
    published_at: string;
    url: string;
    thumbnail_url: string;
    viewable: string;
    view_count: number;
    language: string;
    type: string;
    duration: string;
    muted_segments?: unknown[];
}

export interface TwitchSubscriber {
    broadcaster_id: string;
    broadcaster_login: string;
    broadcaster_name: string;
    gifter_id?: string;
    gifter_login?: string;
    gifter_name?: string;
    is_gift: boolean;
    tier: string;
    plan_name: string;
    user_id: string;
    user_name: string;
    user_login: string;
}

class TwitchService {
    private baseUrl: string;

    constructor() {
        this.baseUrl = config.apiBaseUrl;
    }

    private async makeRequest<T>(
        endpoint: string,
        options: RequestInit = {},
        token?: string
    ): Promise<T> {
        const url = `${this.baseUrl}${endpoint}`;

        const headers: Record<string, string> = {
            'Content-Type': 'application/json',
            ...(options.headers as Record<string, string>),
        };

        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        const response = await fetch(url, {
            ...options,
            headers,
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
            const error = new Error(errorData.error || `API request failed: ${response.status}`) as Error & TwitchApiError;
            error.reconnect_required = errorData.reconnect_required;
            error.connect_required = errorData.connect_required;
            throw error;
        }

        return response.json();
    }

    // Check Twitch connection status
    async getConnectionStatus(token: string): Promise<TwitchConnectionStatus> {
        return this.makeRequest<TwitchConnectionStatus>('/api/user/twitch-status', {}, token);
    }

    // Connect to Twitch (secure OAuth flow)
    async connectToTwitch(token: string): Promise<void> {
        if (typeof window === 'undefined') {
            throw new Error('Not in browser environment');
        }

        // Make authenticated request to initiate OAuth flow
        const response = await fetch(`${this.baseUrl}/auth/twitch/initiate`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
            throw new Error(errorData.error || `Failed to initiate OAuth: ${response.status}`);
        }

        const data = await response.json();
        if (!data.oauth_url) {
            throw new Error('No OAuth URL received from server');
        }

        // Show permissions to user before redirecting
        if (data.permissions && data.permissions.length > 0) {
            const permissionsList = data.permissions
                .map((p: { description: string }) => `• ${p.description}`)
                .join('\n');
            
            const confirmed = window.confirm(
                `CreatorSync is requesting access to the following Twitch permissions:\n\n${permissionsList}\n\nDo you want to continue?`
            );
            
            if (!confirmed) {
                return;
            }
        }

        // Redirect to the secure OAuth URL
        window.location.href = data.oauth_url;
    }

    // Disconnect from Twitch
    async disconnectFromTwitch(token: string): Promise<{ message: string }> {
        return this.makeRequest<{ message: string }>(
            '/api/user/twitch-disconnect',
            { method: 'DELETE' },
            token
        );
    }

    // Get Twitch channel information
    async getChannelInfo(token: string): Promise<{ channel: TwitchChannelInfo }> {
        return this.makeRequest<{ channel: TwitchChannelInfo }>('/api/twitch/channel', {}, token);
    }

    // Get Twitch videos
    async getVideos(token: string, type: string = 'archive', limit: number = 20): Promise<{ 
        videos: TwitchVideo[];
        pagination: { cursor?: string };
    }> {
        return this.makeRequest<{ 
            videos: TwitchVideo[];
            pagination: { cursor?: string };
        }>(`/api/twitch/videos?type=${type}&limit=${limit}`, {}, token);
    }

    // Get Twitch subscribers
    async getSubscribers(token: string): Promise<{ subscribers: TwitchSubscriber[] }> {
        return this.makeRequest<{ subscribers: TwitchSubscriber[] }>('/api/twitch/subscribers', {}, token);
    }

    // Handle callback URL parameters
    static handleOAuthCallback(): { success: boolean; error?: string } {
        if (typeof window === 'undefined') {
            return { success: false, error: 'Not in browser environment' };
        }

        const urlParams = new URLSearchParams(window.location.search);
        
        if (urlParams.get('twitch_connected') === 'true') {
            // Clean up URL parameters
            window.history.replaceState({}, document.title, window.location.pathname);
            return { success: true };
        }
        
        const error = urlParams.get('twitch_error');
        if (error) {
            // Clean up URL parameters
            window.history.replaceState({}, document.title, window.location.pathname);
            return { success: false, error };
        }

        return { success: false };
    }

    // Get user-friendly error message
    static getErrorMessage(error: string): string {
        switch (error) {
            case 'oauth_denied':
                return 'You denied access to your Twitch account. Please try connecting again.';
            case 'csrf_failed':
                return 'Security validation failed. Please try connecting again.';
            case 'invalid_callback':
                return 'Invalid callback from Twitch. Please try connecting again.';
            case 'token_exchange_failed':
                return 'Failed to exchange authorization code. Please try connecting again.';
            case 'user_info_failed':
                return 'Failed to get your Twitch user information. Please try connecting again.';
            case 'token_storage_failed':
                return 'Failed to store your Twitch credentials. Please try connecting again.';
            case 'auth_failed':
                return 'Authentication failed. Please log in and try again.';
            case 'client_init_failed':
                return 'Failed to initialize Twitch client. Please try connecting again.';
            default:
                return 'An unexpected error occurred. Please try connecting again.';
        }
    }
}

// Export singleton instance and class for static methods
export const twitchService = new TwitchService();
export { TwitchService }; 