/**
 * Application configuration
 * 
 * This file contains environment-specific configuration settings.
 * Values are determined based on the current environment (development, production, etc.)
 */

// Determine the base API URL based on the environment
const getApiBaseUrl = (): string => {
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