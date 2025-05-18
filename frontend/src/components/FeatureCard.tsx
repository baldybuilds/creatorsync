import { useState, useEffect } from 'react';
import type { LucideIcon } from 'lucide-react';
import { ArrowRight, ChevronRight } from 'lucide-react';

export interface FeatureCardProps {
    icon: LucideIcon;
    title: string;
    description: string;
    delay?: number;
    isActive?: boolean;
    onClick?: () => void;
}

export function FeatureCard({
    icon: Icon,
    title,
    description,
    delay = 0,
    isActive = false,
    onClick
}: FeatureCardProps) {
    const [isVisible, setIsVisible] = useState(false);
    const [isFlipped, setIsFlipped] = useState(false);

    useEffect(() => {
        const timer = setTimeout(() => setIsVisible(true), delay);
        return () => clearTimeout(timer);
    }, [delay]);

    const handleClick = () => {
        if (onClick) {
            onClick();
        } else {
            setIsFlipped(!isFlipped);
        }
    };

    return (
        <div
            className={`relative h-[300px] perspective-dramatic ${isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}`}
            onClick={handleClick}
        >
            <div className={`w-full h-full transition-all duration-700 ${isFlipped ? 'rotate-y-180' : ''}`} style={{ transformStyle: 'preserve-3d' }}>
                {/* Front of card */}
                <div className={`card p-8 absolute w-full h-full cursor-pointer transition-all duration-500 hover:scale-[1.02] ${isActive ? 'ring-2 ring-brand-500 ring-offset-2 ring-offset-surface-900' : ''}`} style={{ backfaceVisibility: 'hidden' }}>
                    <div className="flex flex-col items-center text-center h-full justify-between">
                        <div className="p-6 rounded-2xl bg-brand-500/10 border border-brand-500/20 transition-all duration-300 group-hover:scale-110 group-hover:bg-brand-500/20 transform hover:rotate-12">
                            <Icon className="w-10 h-10 text-brand-500 group-hover:text-accent-500 transition-colors duration-300" />
                        </div>

                        <div className="flex-1 flex flex-col justify-center">
                            <h3 className="text-2xl font-bold text-surface-100 mt-6 mb-4 group-hover:text-brand-300 transition-colors duration-300">
                                {title}
                            </h3>

                            <p className="text-surface-400 text-lg leading-relaxed group-hover:text-surface-300 transition-colors duration-300">
                                {description}
                            </p>
                        </div>

                        <button className="mt-4 text-brand-400 flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity duration-300 hover:text-brand-300">
                            <span>Learn more</span>
                            <ChevronRight className="w-4 h-4" />
                        </button>
                    </div>

                    {/* Hover glow effect */}
                    <div className={`absolute inset-0 rounded-2xl transition-opacity duration-500 shadow-glow ${isActive ? 'opacity-100' : 'opacity-0 group-hover:opacity-100'}`} />
                </div>

                {/* Back of card */}
                <div className="card p-8 absolute w-full h-full bg-surface-800 cursor-pointer transition-all duration-500 hover:scale-[1.02] rotate-y-180" style={{ backfaceVisibility: 'hidden' }}>
                    <div className="flex flex-col h-full justify-between">
                        <div className="flex justify-between items-center">
                            <h3 className="text-2xl font-bold text-brand-300">{title}</h3>
                            <div className="p-2 rounded-full bg-brand-500/10 border border-brand-500/20">
                                <Icon className="w-6 h-6 text-brand-500" />
                            </div>
                        </div>

                        <div className="flex-1 flex flex-col justify-center py-6">
                            <p className="text-surface-200 text-lg leading-relaxed mb-4">
                                {description}
                            </p>
                            <ul className="space-y-2">
                                <li className="flex items-center gap-2 text-surface-300">
                                    <div className="w-1.5 h-1.5 rounded-full bg-brand-500"></div>
                                    <span>Advanced customization options</span>
                                </li>
                                <li className="flex items-center gap-2 text-surface-300">
                                    <div className="w-1.5 h-1.5 rounded-full bg-brand-500"></div>
                                    <span>Seamless integration with your workflow</span>
                                </li>
                                <li className="flex items-center gap-2 text-surface-300">
                                    <div className="w-1.5 h-1.5 rounded-full bg-brand-500"></div>
                                    <span>Continuous updates and improvements</span>
                                </li>
                            </ul>
                        </div>

                        <button className="self-end flex items-center gap-1 text-brand-400 hover:text-brand-300 transition-colors">
                            <span>Try it now</span>
                            <ArrowRight className="w-4 h-4" />
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}