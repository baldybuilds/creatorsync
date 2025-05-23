/**
 * Application configuration
 * 
 * This file contains environment-specific configuration settings.
 * Values are determined based on the current environment (development, production, etc.)
 */

// Determine the base API URL based on the environment
const getApiBaseUrl = (): string => {
    // Check if we're in staging environment (dev.creatorsync.app)
    if (typeof window !== 'undefined' && window.location.hostname === 'dev.creatorsync.app') {
        return 'https://api-dev.creatorsync.app';
    }

    // Check for explicit staging environment variable
    if (process.env.NEXT_PUBLIC_APP_ENV === 'staging') {
        return 'https://api-dev.creatorsync.app';
    }

    // For local development
    if (process.env.NODE_ENV === 'development') {
        return process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
    }

    // For production
    return process.env.NEXT_PUBLIC_API_URL || 'https://api.creatorsync.app';
};

export const config = {
    apiBaseUrl: getApiBaseUrl(),
};
