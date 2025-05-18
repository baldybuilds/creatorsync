import { useState, useEffect } from 'react';

export interface StatCardProps {
    number: string;
    label: string;
    delay?: number;
}

export function StatCard({ number, label, delay = 0 }: StatCardProps) {
    const [count, setCount] = useState(0);
    const [isVisible, setIsVisible] = useState(false);

    useEffect(() => {
        const timer = setTimeout(() => {
            setIsVisible(true);

            // Parse the number and animate to it
            const target = parseInt(number.replace(/\D/g, ''));
            if (target && !isNaN(target)) {
                let current = 0;
                const increment = target / 60; // Animate over ~1 second

                const counting = setInterval(() => {
                    current += increment;
                    if (current >= target) {
                        setCount(target);
                        clearInterval(counting);
                    } else {
                        setCount(Math.floor(current));
                    }
                }, 16); // ~60fps

                return () => clearInterval(counting);
            }
        }, delay);

        return () => clearTimeout(timer);
    }, [number, delay]);

    const formatNumber = (num: number) => {
        if (number.includes('+')) return `${num.toLocaleString()}+`;
        return num.toLocaleString();
    };

    return (
        <div className={`text-center transition-all duration-500 ${isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'
            }`}>
            <div className="text-4xl md:text-5xl font-bold text-brand-500 mb-2">
                {formatNumber(count)}
            </div>
            <div className="text-surface-400 text-sm uppercase tracking-wide font-medium">
                {label}
            </div>
        </div>
    );
}