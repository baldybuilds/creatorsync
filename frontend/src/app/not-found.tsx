'use client';

import Link from 'next/link';
import { Button } from '../components/ui/button';
import { LogIn, Home, AlertTriangle } from 'lucide-react';
import { SignedOut, SignInButton } from '@clerk/clerk-react';
import { useFeatureFlagEnabled } from 'posthog-js/react';
import { usePostHogUser } from '../hooks/usePostHogUser';

export default function NotFound() {
    usePostHogUser(); // Initialize PostHog for user tracking consistency
    const clerkAuthButtonEnabled = useFeatureFlagEnabled('clerk-auth-button');

    return (
        <div className="min-h-screen flex flex-col items-center justify-center bg-[var(--bg-primary)] text-[var(--text-primary)] p-6 text-center">
            <AlertTriangle className="w-24 h-24 text-yellow-500 mb-6" />
            <h1 className="text-5xl md:text-6xl font-bold mb-4">
                <span className="text-light-surface-900 dark:text-dark-surface-100">404</span> - <span className="text-gradient">Page Not Found</span>
            </h1>
            <p className="text-lg text-[var(--text-secondary)] mb-10 max-w-md">
                Sorry, the page you are looking for could not be found. It might have been removed, had its name changed, or is temporarily unavailable.
            </p>

            <div className="flex flex-col sm:flex-row items-center gap-4">
                <Link href="/">
                    <Button variant="outline" size="lg" className="bg-transparent hover:bg-light-surface-100 dark:hover:bg-dark-surface-800">
                        <Home className="mr-2 h-5 w-5" />
                        Go to Homepage
                    </Button>
                </Link>

                {clerkAuthButtonEnabled && (
                    <SignedOut>
                        {/* Default redirect after sign-in can be configured in Clerk dashboard or here via redirectUrl */}
                        <SignInButton mode="modal">
                            <Button variant="default" size="lg">
                                <LogIn className="mr-2 h-5 w-5" />
                                Sign In
                            </Button>
                        </SignInButton>
                    </SignedOut>
                )}
            </div>

            {!clerkAuthButtonEnabled && (
                <p className="mt-8 text-sm text-[var(--text-secondary)]">
                    Access is currently limited during our MVP phase. Please check back later!
                </p>
            )}

            {/* Minimal style for text-gradient if not globally defined */}
            <style jsx>{`
        .text-gradient {
          background: linear-gradient(to right, var(--color-primary-500, #3b82f6), var(--color-secondary-500, #8b5cf6));
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          background-clip: text;
          text-fill-color: transparent;
        }
      `}</style>
        </div>
    );
}
