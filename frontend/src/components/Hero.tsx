import { useState, useEffect } from 'react';
import { WaitlistForm } from './WaitlistForm';

export function Hero() {
    const [isVisible, setIsVisible] = useState(false);

    useEffect(() => {
        setIsVisible(true);
    }, []);

    return (
        <section className="min-h-screen flex flex-col justify-center items-center px-4 sm:px-6 lg:px-8 relative overflow-hidden">
            {/* Background Elements */}
            <div className="absolute inset-0">
                <div className="absolute inset-0 bg-gradient-to-br from-brand-500/5 via-transparent to-accent-500/5" />
                <div className="absolute top-1/4 left-1/4 w-96 h-96 rounded-full bg-brand-500/10 blur-3xl animate-float" />
                <div className="absolute bottom-1/4 right-1/4 w-80 h-80 rounded-full bg-accent-500/10 blur-3xl animate-float" style={{ animationDelay: '2s' }} />
            </div>

            <div className="max-w-6xl text-center relative z-10">
                {/* Badge */}
                {/* <div className={`inline-flex items-center gap-2 px-4 py-2 rounded-full bg-surface-800/80 backdrop-blur-xl border border-surface-700/50 text-surface-300 text-sm mb-8 transition-all duration-1000 ${isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'
                    }`}>
                    <Sparkles className="w-4 h-4 text-brand-500 animate-pulse" />
                    <span>Coming Soon • Join our waitlist</span>
                </div> */}

                {/* Main Heading */}
                <div className={`space-y-6 transition-all duration-1000 ${isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'
                    }`}>
                    <h1 className="text-6xl md:text-8xl lg:text-9xl font-bold leading-tight">
                        <span className="text-light-surface-900 dark:text-dark-surface-100">Creator</span>
                        <span className="text-gradient">Sync</span>
                    </h1>

                    <p className="text-xl md:text-3xl text-light-surface-700 dark:text-dark-surface-300 max-w-5xl mx-auto leading-relaxed">
                        Transform your Twitch clips & vods into viral moments across{' '}
                        <span className="text-brand-500">all platforms</span>.{' '}
                        <br className="hidden md:block" />
                        Edit, render, and share—all from one{' '}
                        <span className="text-gradient">powerful dashboard</span>.
                    </p>
                </div>

                {/* Waitlist Form */}
                <WaitlistForm className="mt-6" />

                {/* Stats Section */}
                {/* <div className={`mt-20 transition-all duration-1000 ${isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'
                    }`} style={{ transitionDelay: '400ms' }}>
                    <p className="text-surface-400 text-lg mb-8">Trusted by creators worldwide</p>
                    <div className="grid grid-cols-1 sm:grid-cols-3 gap-8 max-w-3xl mx-auto">
                        <StatCard number="50000" label="Clips Processed" delay={600} />
                        <StatCard number="12000" label="Active Creators" delay={800} />
                        <StatCard number="500" label="Hours Saved Daily" delay={1000} />
                    </div>
                </div> */}
            </div>

            {/* Scroll indicator */}
            {/* <div className="absolute bottom-8 left-1/2 transform -translate-x-1/2 animate-bounce">
                <div className="w-6 h-10 border-2 border-surface-400 rounded-full flex justify-center p-1">
                    <div className="w-1 h-3 bg-surface-400 rounded-full animate-pulse"></div>
                </div>
            </div> */}
        </section>
    );
}