import { useState } from 'react';
import { Sparkles, Send, CheckCircle, XCircle } from 'lucide-react';
import { Button } from './ui/button';
import { config } from '../config';

interface WaitlistFormProps {
    className?: string;
}

export function WaitlistForm({ className = '' }: WaitlistFormProps) {
    const [email, setEmail] = useState('');
    const [name, setName] = useState('');
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [submitStatus, setSubmitStatus] = useState<'idle' | 'success' | 'error'>('idle');
    const [errorMessage, setErrorMessage] = useState('');

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!email) {
            setErrorMessage('Please enter your email address');
            setSubmitStatus('error');
            return;
        }

        setIsSubmitting(true);
        setSubmitStatus('idle');
        setErrorMessage('');

        try {
            const response = await fetch(`${config.apiBaseUrl}/api/waitlist`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email, name }),
            });

            if (!response.ok) {
                let errorMessage = 'Failed to join waitlist';
                try {
                    const contentType = response.headers.get('content-type');
                    if (contentType && contentType.includes('application/json')) {
                        const text = await response.text();
                        if (text) {
                            const errorData = JSON.parse(text);
                            errorMessage = errorData.error || errorMessage;
                        }
                    }
                } catch (parseError) {
                    console.error('Error parsing error response:', parseError);
                }
                throw new Error(errorMessage);
            }

            setSubmitStatus('success');
            setEmail('');
            setName('');
        } catch (error) {
            setSubmitStatus('error');
            setErrorMessage(error instanceof Error ? error.message : 'An unexpected error occurred');
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className={`${className}`}>
            <div className="max-w-md mx-auto">
                {submitStatus === 'success' ? (
                    <div className="text-center p-6 bg-brand-500/10 border border-brand-500/20 rounded-xl">
                        <CheckCircle className="w-12 h-12 text-brand-500 mx-auto mb-4" />
                        <h3 className="text-xl font-bold text-light-surface-900 dark:text-dark-surface-100 mb-2">You're on the list!</h3>
                        <p className="text-light-surface-700 dark:text-dark-surface-300">
                            Thank you for joining our waitlist. We'll notify you when we're ready to welcome you to our beta program.
                        </p>
                    </div>
                ) : (
                    <form onSubmit={handleSubmit} className="space-y-4">
                        <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-brand-500/20 border border-brand-500/30 text-brand-300 text-sm mb-4 backdrop-blur-xl w-fit mx-auto">
                            <Sparkles className="w-4 h-4" />
                            <span>Join our waitlist</span>
                        </div>

                        <div>
                            <input
                                type="text"
                                placeholder="Your name (optional)"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                                className="w-full px-4 py-3 rounded-lg bg-light-surface-100/80 dark:bg-dark-surface-800/80 border border-light-surface-300/50 dark:border-dark-surface-700/50 text-light-surface-900 dark:text-dark-surface-100 placeholder:text-light-surface-500 dark:placeholder:text-dark-surface-500 focus:outline-none focus:ring-2 focus:ring-brand-500/50 transition-all"
                            />
                        </div>

                        <div>
                            <input
                                type="email"
                                placeholder="Your email address"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                required
                                className="w-full px-4 py-3 rounded-lg bg-light-surface-100/80 dark:bg-dark-surface-800/80 border border-light-surface-300/50 dark:border-dark-surface-700/50 text-light-surface-900 dark:text-dark-surface-100 placeholder:text-light-surface-500 dark:placeholder:text-dark-surface-500 focus:outline-none focus:ring-2 focus:ring-brand-500/50 transition-all"
                            />
                        </div>

                        {submitStatus === 'error' && (
                            <div className="flex items-center gap-2 text-red-400 text-sm">
                                <XCircle className="w-4 h-4" />
                                <span>{errorMessage || 'Failed to join waitlist. Please try again.'}</span>
                            </div>
                        )}

                        <Button
                            type="submit"
                            variant="default"
                            size="lg"
                            disabled={isSubmitting}
                            className="w-full shadow-glow hover:shadow-brand-500/50"
                        >
                            {isSubmitting ? 'Submitting...' : 'Join Waitlist'}
                            <Send className="ml-2 h-4 w-4" />
                        </Button>

                        <p className="text-center text-light-surface-600 dark:text-dark-surface-400 text-sm">
                            We'll never share your email with anyone else.
                        </p>
                    </form>
                )}
            </div>
        </div>
    );
}
