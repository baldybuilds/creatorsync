// OAuth utility functions
export interface OAuthConfig {
    platform: 'twitch' | 'youtube' | 'tiktok';
    clientId: string;
    redirectUri: string;
    scopes: string[];
}

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

export const oauthConfigs = {
    twitch: {
        platform: 'twitch' as const,
        clientId: process.env.NEXT_PUBLIC_TWITCH_CLIENT_ID || '',
        redirectUri: `${typeof window !== 'undefined' ? window.location.origin : ''}/auth/callback/twitch`,
        scopes: [
            'user:read:email',
            'user:read:subscriptions',
            'moderator:read:followers',
            'clips:edit',
            'channel:read:redemptions',
            'moderation:read'
        ],
    },
};

export const initiateOAuth = async (platform: 'twitch' | 'youtube' | 'tiktok', getToken?: () => Promise<string | null>) => {
    if (platform === 'twitch') {
        try {
            // Get the authentication token first
            if (!getToken) {
                throw new Error('Authentication token provider is required');
            }
            
            const token = await getToken();
            if (!token) {
                throw new Error('Authentication failed: No token available');
            }

            const apiBaseUrl = getApiBaseUrl();
            
            // Make authenticated request to get OAuth URL
            const response = await fetch(`${apiBaseUrl}/auth/twitch/initiate`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
            });

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
                throw new Error(errorData.error || `OAuth request failed: ${response.status}`);
            }

            const data = await response.json();
            
            // Store the current location to redirect back after OAuth
            localStorage.setItem('oauth_redirect_after', window.location.pathname);
            
            // Redirect to the OAuth URL provided by the backend
            if (data.oauth_url) {
                window.location.href = data.oauth_url;
            } else {
                throw new Error('OAuth URL not provided by backend');
            }
        } catch (error) {
            console.error('OAuth initiation error:', error);
            alert(`Failed to connect ${platform}: ${error instanceof Error ? error.message : 'Unknown error'}`);
        }
    } else {
        // For other platforms, show coming soon message
        alert(`${platform} integration coming soon!`);
    }
};

export const handleOAuthCallback = async (platform: string, code: string, state?: string) => {
    try {
        const apiBaseUrl = getApiBaseUrl();
        
        const response = await fetch(`${apiBaseUrl}/auth/${platform}/callback`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ code, state }),
        });

        if (!response.ok) {
            throw new Error(`OAuth callback failed: ${response.status}`);
        }

        const result = await response.json();
        
        // Redirect back to where the user was
        const redirectPath = localStorage.getItem('oauth_redirect_after') || '/dashboard';
        localStorage.removeItem('oauth_redirect_after');
        
        window.location.href = redirectPath;
        
        return result;
    } catch (error) {
        console.error('OAuth callback error:', error);
        throw error;
    }
}; 