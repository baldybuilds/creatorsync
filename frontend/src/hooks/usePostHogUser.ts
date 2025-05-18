import { useUser } from '@clerk/clerk-react';
import { useEffect } from 'react';
import posthog from 'posthog-js';

export function usePostHogUser() {
    const { user, isSignedIn } = useUser();

    useEffect(() => {
        if (isSignedIn && user) {
            posthog.identify(user.id, {
                email: user.primaryEmailAddress?.emailAddress,
                name: user.fullName,
                clerk_id: user.id
            });
        } else {
            posthog.reset();
        }
    }, [isSignedIn, user]);
}
