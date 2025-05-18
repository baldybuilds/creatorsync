import { useState, useEffect } from 'react';
import { Play, Edit3, Download, Zap, Users, Clock, ChevronLeft, ChevronRight } from 'lucide-react';
import { FeatureCard } from './FeatureCard';
import { motion } from 'framer-motion';

export function Features() {
    const [activeFeatureIndex, setActiveFeatureIndex] = useState<number | null>(null);
    const [isAutoPlaying, setIsAutoPlaying] = useState(true);

    const features = [
        {
            icon: Play,
            title: "Auto-Import Clips",
            description: "Instantly access your latest Twitch clips and streams from a unified dashboard.",
            delay: 0
        },
        {
            icon: Edit3,
            title: "Browser-Based Editor",
            description: "No downloads needed. Edit your clips directly in the browser with our powerful editor.",
            delay: 200
        },
        {
            icon: Zap,
            title: "Lightning Fast Rendering",
            description: "High-quality video rendering in seconds, not minutes. Powered by AWS Lambda.",
            delay: 400
        },
        {
            icon: Users,
            title: "Multi-Platform Export",
            description: "Optimize for TikTok, YouTube Shorts, Instagram Reels, and more with one click.",
            delay: 600
        },
        {
            icon: Clock,
            title: "Time-Saving Automation",
            description: "Batch process multiple clips and schedule posts. Save hours of manual work.",
            delay: 800
        },
        {
            icon: Download,
            title: "Instant Downloads",
            description: "Get your edited clips immediately. No waiting, no hassle, no complicated exports.",
            delay: 1000
        }
    ];

    // Auto-rotate through features
    useEffect(() => {
        if (!isAutoPlaying) return;

        const interval = setInterval(() => {
            setActiveFeatureIndex(prev => {
                if (prev === null) return 0;
                return (prev + 1) % features.length;
            });
        }, 5000);

        return () => clearInterval(interval);
    }, [isAutoPlaying, features.length]);

    // Pause auto-rotation when user interacts with a feature
    const handleFeatureClick = (index: number) => {
        setIsAutoPlaying(false);
        setActiveFeatureIndex(index);
    };

    const navigateFeature = (direction: 'prev' | 'next') => {
        setIsAutoPlaying(false);
        setActiveFeatureIndex(prev => {
            if (prev === null) return direction === 'next' ? 0 : features.length - 1;
            const newIndex = direction === 'next' ? prev + 1 : prev - 1;
            return (newIndex + features.length) % features.length;
        });
    };

    return (
        <section className="py-32 px-4 sm:px-6 lg:px-8 relative overflow-hidden">
            {/* Background with animated gradient */}
            <div className="absolute inset-0 bg-gradient-to-b from-surface-950/50 to-surface-900/50">
                <div className="absolute inset-0 opacity-20">
                    <div className="absolute top-0 -left-4 w-72 h-72 bg-brand-500/30 rounded-full filter blur-3xl animate-blob"></div>
                    <div className="absolute top-0 -right-4 w-72 h-72 bg-accent-500/20 rounded-full filter blur-3xl animate-blob animation-delay-2000"></div>
                    <div className="absolute -bottom-8 left-20 w-72 h-72 bg-purple-500/30 rounded-full filter blur-3xl animate-blob animation-delay-4000"></div>
                </div>
            </div>

            <div className="max-w-7xl mx-auto relative z-10">
                {/* Section Header with animation */}
                <motion.div
                    className="text-center mb-20"
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.6 }}
                >
                    <h2 className="text-5xl md:text-7xl font-bold mb-8">
                        <span className="text-surface-100">Everything you need</span>
                        <br />
                        <motion.span
                            className="text-gradient"
                            initial={{ opacity: 0, scale: 0.9 }}
                            animate={{ opacity: 1, scale: 1 }}
                            transition={{ delay: 0.3, duration: 0.5 }}
                        >
                            to go viral
                        </motion.span>
                    </h2>

                    <p className="text-xl md:text-2xl text-surface-300 max-w-4xl mx-auto leading-relaxed">
                        Stop wasting hours on manual editing. CreatorSync automates your content workflow
                        so you can focus on what matters—<span className="text-brand-500">creating amazing content</span>.
                    </p>
                </motion.div>

                {/* Feature navigation controls */}
                <div className="flex justify-center mb-12 gap-4">
                    <button
                        onClick={() => navigateFeature('prev')}
                        className="p-3 rounded-full bg-surface-800 hover:bg-surface-700 transition-colors duration-300"
                        aria-label="Previous feature"
                    >
                        <ChevronLeft className="w-6 h-6 text-surface-200" />
                    </button>

                    <div className="flex items-center gap-2">
                        {features.map((_, index) => (
                            <button
                                key={index}
                                onClick={() => handleFeatureClick(index)}
                                className={`w-3 h-3 rounded-full transition-all duration-300 ${activeFeatureIndex === index ? 'bg-brand-500 scale-125' : 'bg-surface-600 hover:bg-surface-500'}`}
                                aria-label={`Go to feature ${index + 1}`}
                            />
                        ))}
                    </div>

                    <button
                        onClick={() => navigateFeature('next')}
                        className="p-3 rounded-full bg-surface-800 hover:bg-surface-700 transition-colors duration-300"
                        aria-label="Next feature"
                    >
                        <ChevronRight className="w-6 h-6 text-surface-200" />
                    </button>
                </div>

                {/* Features Grid with improved layout */}
                <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-8 lg:gap-12">
                    {features.map((feature, index) => (
                        <FeatureCard
                            key={index}
                            icon={feature.icon}
                            title={feature.title}
                            description={feature.description}
                            delay={feature.delay}
                            isActive={activeFeatureIndex === index}
                            onClick={() => handleFeatureClick(index)}
                        />
                    ))}
                </div>

                {/* Auto-play toggle */}
                <div className="mt-12 flex justify-center">
                    <button
                        onClick={() => setIsAutoPlaying(!isAutoPlaying)}
                        className={`px-4 py-2 rounded-full text-sm transition-colors duration-300 flex items-center gap-2 ${isAutoPlaying ? 'bg-brand-500/20 text-brand-300' : 'bg-surface-800 text-surface-300 hover:bg-surface-700'}`}
                    >
                        <div className={`w-3 h-3 rounded-full ${isAutoPlaying ? 'bg-brand-500' : 'bg-surface-600'}`}></div>
                        {isAutoPlaying ? 'Auto-play enabled' : 'Auto-play disabled'}
                    </button>
                </div>
            </div>
        </section>
    );
}