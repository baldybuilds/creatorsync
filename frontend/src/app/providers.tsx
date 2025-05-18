// app/providers.tsx
'use client'

import posthog from 'posthog-js'
import { PostHogProvider as PHProvider } from 'posthog-js/react'
import { useEffect } from 'react'

export function PostHogProvider({ children }: { children: React.ReactNode }) {
    useEffect(() => {
        const apiKey = process.env.NEXT_PUBLIC_POSTHOG_KEY;
        const apiHost = process.env.NEXT_PUBLIC_POSTHOG_HOST;

        if (apiKey) {
            const config: { api_host?: string; capture_pageview: boolean; } = {
                capture_pageview: false
            };

            if (apiHost) {
                config.api_host = apiHost;
            }

            posthog.init(apiKey, config);
        } else {
            console.warn('NEXT_PUBLIC_POSTHOG_KEY is not set. PostHog analytics will not be initialized.');
        }
    }, [])

    return (
        <PHProvider client={posthog}>
            {children}
        </PHProvider>
    )
}