import { useState, useEffect } from 'react';
import { Shield, Zap, Sparkles } from 'lucide-react';

export function CTA() {
    const [isVisible, setIsVisible] = useState(false);

    useEffect(() => {
        const observer = new IntersectionObserver(
            ([entry]) => {
                if (entry.isIntersecting) {
                    setIsVisible(true);
                }
            },
            { threshold: 0.1 }
        );

        const element = document.getElementById('cta-section');
        if (element) observer.observe(element);

        return () => observer.disconnect();
    }, []);

    return (
        <section id="cta-section" className="py-18 px-4 sm:px-6 lg:px-8 relative overflow-hidden">
            {/* Background */}
            <div className="absolute inset-0">
                <div className="absolute inset-0 bg-gradient-to-t from-surface-950 via-surface-900 to-surface-950" />
                <div className="absolute inset-0 bg-[radial-gradient(ellipse_1200px_800px_at_50%_50%,rgba(99,102,241,0.1),transparent_60%)]" />
            </div>

            <div className="max-w-6xl mx-auto relative z-10">
                <div className={`transition-all duration-1000 ${isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'
                    }`}>
                    {/* Main CTA Card */}
                    <div className="card p-12 text-center relative">
                        {/* Badge */}
                        <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-brand-500/20 border border-brand-500/30 text-brand-300 text-sm mb-8 backdrop-blur-xl">
                            <Sparkles className="w-4 h-4" />
                            <span>Limited Time Beta Access</span>
                        </div>

                        {/* Heading */}
                        <h2 className="text-4xl md:text-6xl font-bold mb-6">
                            <span className="text-surface-100">Ready to</span>
                            <br />
                            <span className="text-gradient">transform your workflow</span>
                            <span className="text-surface-100">?</span>
                        </h2>

                        {/* Subtitle */}
                        <p className="text-xl md:text-2xl text-surface-300 mb-12 max-w-3xl mx-auto">
                            Join <span className="text-brand-500 font-semibold">thousands of creators</span> who've already
                            streamlined their content process with CreatorSync.
                        </p>

                        {/* Feature Badges */}
                        <div className="flex flex-wrap justify-center gap-4 mb-12">
                            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-brand-500/20 text-brand-300 border border-brand-500/30 text-sm">
                                <Shield className="w-3 h-3" />
                                <span>No Credit Card Required</span>
                            </div>
                            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-brand-500/20 text-brand-300 border border-brand-500/30 text-sm">
                                <Zap className="w-3 h-3" />
                                <span>Instant Access</span>
                            </div>
                            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-brand-500/20 text-brand-300 border border-brand-500/30 text-sm">
                                <Sparkles className="w-3 h-3" />
                                <span>Premium Features Free</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    );
}