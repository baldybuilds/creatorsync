import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { Code, Sparkles, Zap, Clock } from 'lucide-react';
import * as Separator from '@radix-ui/react-separator';

export function About() {
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

    const section = document.getElementById('about-section');
    if (section) observer.observe(section);

    return () => {
      if (section) observer.unobserve(section);
    };
  }, []);

  const stats = [
    { icon: Clock, value: '70%', label: 'Time Saved', description: 'Creators save up to 70% of their editing time' },
    // { icon: Users, value: '500+', label: 'Beta Testers', description: 'Growing community of content creators' },
    { icon: Zap, value: '10x', label: 'Faster Workflow', description: 'Publish content 10x faster across platforms' }
  ];

  return (
    <section id="about-section" className="py-24 px-4 sm:px-6 lg:px-8 relative overflow-hidden">

      <div className="max-w-7xl mx-auto relative z-10">
        <motion.div
          className="text-center mb-16"
          initial={{ opacity: 0, y: 20 }}
          animate={isVisible ? { opacity: 1, y: 0 } : { opacity: 0, y: 20 }}
          transition={{ duration: 0.6 }}
        >
          <div className="flex items-center justify-center gap-2 px-4 py-2 rounded-full bg-brand-500/20 border border-brand-500/30 text-brand-300 text-sm mb-6 backdrop-blur-xl w-fit mx-auto">
            <Sparkles className="w-4 h-4" />
            <span>Our Mission</span>
          </div>

          <h2 className="text-4xl md:text-5xl font-bold mb-6">
            <span className="text-light-surface-900 dark:text-dark-surface-100">Built by creators,</span>{' '}
            <span className="text-gradient">for creators</span>
          </h2>

          <p className="text-xl text-light-surface-700 dark:text-dark-surface-300 max-w-3xl mx-auto leading-relaxed">
            We understand the challenges of managing content across multiple platforms because we've been there ourselves.
          </p>
        </motion.div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center mb-16">
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={isVisible ? { opacity: 1, x: 0 } : { opacity: 0, x: -20 }}
            transition={{ duration: 0.6, delay: 0.2 }}
          >
            <div className="relative rounded-2xl overflow-hidden bg-gradient-to-br from-brand-500/20 to-accent-500/20 p-1">
              <div className="absolute inset-0 bg-gradient-to-br from-brand-500/30 to-accent-500/30 opacity-30" />
              <div className="relative bg-light-surface-100/95 dark:bg-dark-surface-900/95 backdrop-blur-sm rounded-xl p-8">
                <h3 className="text-2xl font-bold mb-4 text-light-surface-900 dark:text-dark-surface-100">Our Story</h3>
                <p className="text-light-surface-700 dark:text-dark-surface-300 mb-4">
                  CreatorSync began when our founder, a Twitch streamer, spent countless hours manually editing clips for different social platforms.
                  The process was tedious, repetitive, and took time away from creating new content.
                </p>
                <p className="text-light-surface-700 dark:text-dark-surface-300">
                  We built CreatorSync to solve this problem once and for all — a platform that automates the tedious parts of content repurposing
                  while giving creators complete control over their brand and style.
                </p>
                <Separator.Root className="my-6 h-px bg-light-surface-200/50 dark:bg-dark-surface-800/50" />
                <div className="flex items-center gap-3">
                  <Code className="w-5 h-5 text-brand-500" />
                  <p className="text-sm text-light-surface-700 dark:text-dark-surface-300">
                    <span className="font-semibold">Launched in 2025</span> • Currently in private beta
                  </p>
                </div>
              </div>
            </div>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={isVisible ? { opacity: 1, x: 0 } : { opacity: 0, x: 20 }}
            transition={{ duration: 0.6, delay: 0.4 }}
            className="grid grid-cols-1 gap-6"
          >
            {stats.map((stat, index) => (
              <div
                key={index}
                className="bg-light-surface-100/95 dark:bg-dark-surface-900/95 backdrop-blur-sm rounded-xl p-6 border border-light-surface-200/50 dark:border-dark-surface-800/50 flex items-center gap-6"
              >
                <div className="flex-shrink-0 w-12 h-12 rounded-full bg-brand-500/20 flex items-center justify-center">
                  <stat.icon className="w-6 h-6 text-brand-500" />
                </div>
                <div>
                  <div className="flex items-baseline gap-2">
                    <span className="text-3xl font-bold text-light-surface-900 dark:text-dark-surface-100">{stat.value}</span>
                    <span className="text-lg font-medium text-brand-500">{stat.label}</span>
                  </div>
                  <p className="text-light-surface-600 dark:text-dark-surface-400">{stat.description}</p>
                </div>
              </div>
            ))}
          </motion.div>
        </div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={isVisible ? { opacity: 1, y: 0 } : { opacity: 0, y: 20 }}
          transition={{ duration: 0.6, delay: 0.6 }}
          className="text-center"
        >
          <h3 className="text-2xl font-bold mb-4 text-light-surface-900 dark:text-dark-surface-100">Our Vision</h3>
          <p className="text-xl text-light-surface-700 dark:text-dark-surface-300 max-w-3xl mx-auto leading-relaxed">
            We envision a world where creators can focus on what they do best — creating amazing content —
            while CreatorSync handles the technical aspects of repurposing and distribution.
          </p>
        </motion.div>
      </div>
    </section>
  );
}
