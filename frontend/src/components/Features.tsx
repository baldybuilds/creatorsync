import { useState, useEffect, useRef, useMemo } from 'react';
import { motion } from 'framer-motion';
import { Play, Edit3, Download, Zap, Users, Clock, Sparkles, ArrowRight, CheckCircle } from 'lucide-react';
import * as Tabs from '@radix-ui/react-tabs';
import * as Accordion from '@radix-ui/react-accordion';
import { cn } from '@/lib/utils';

type FeatureType = {
  id: string;
  icon: React.ElementType;
  title: string;
  description: string;
  benefits: string[];
  color: string;
};

export function Features() {
  const [isVisible, setIsVisible] = useState(false);
  const [activeTab, setActiveTab] = useState('feature1');
  const [isAutoRotating, setIsAutoRotating] = useState(true);
  const autoRotateIntervalRef = useRef<NodeJS.Timeout | null>(null);

  const features = useMemo<FeatureType[]>(() => [
    
    {
      id: 'feature1',
      icon: Play,
      title: "Auto-Import Clips",
      description: "Instantly access your latest Twitch clips and streams from a unified dashboard. Connect your accounts once and we'll handle the rest.",
      benefits: [
        "Automatic clip detection and import",
        "Smart categorization by game, date, and popularity",
        "Bulk import and organization tools"
      ],
      color: "from-blue-500 to-indigo-500"
    },
    {
      id: 'feature2',
      icon: Edit3,
      title: "Browser-Based Editor",
      description: "No downloads needed. Edit your clips directly in the browser with our powerful editor designed specifically for short-form content.",
      benefits: [
        "Professional editing tools in your browser",
        "Custom templates for different platforms",
        "Collaborative editing with team members"
      ],
      color: "from-purple-500 to-pink-500"
    },
    {
      id: 'feature3',
      icon: Zap,
      title: "Lightning Fast Rendering",
      description: "High-quality video rendering in seconds, not minutes. Our cloud-based rendering engine delivers professional results instantly.",
      benefits: [
        "Cloud-based rendering powered by AWS",
        "4K support with no quality loss",
        "Background processing while you work on other clips"
      ],
      color: "from-amber-500 to-orange-500"
    },
    {
      id: 'feature4',
      icon: Users,
      title: "Multi-Platform Export",
      description: "Optimize for TikTok, YouTube Shorts, Instagram Reels, and more with one click. Each export is perfectly formatted for its destination.",
      benefits: [
        "Platform-specific aspect ratios and formats",
        "Automatic caption generation",
        "Hashtag recommendations for each platform"
      ],
      color: "from-emerald-500 to-teal-500"
    },
    {
      id: 'feature5',
      icon: Clock,
      title: "Time-Saving Automation",
      description: "Batch process multiple clips and schedule posts. Set it and forget it with our powerful automation tools that save you hours.",
      benefits: [
        "Batch processing for multiple clips",
        "Smart scheduling for optimal posting times",
        "Automated cross-posting across platforms"
      ],
      color: "from-rose-500 to-red-500"
    },
    {
      id: 'feature6',
      icon: Download,
      title: "Instant Downloads",
      description: "Get your edited clips immediately. No waiting, no hassle, no complicated exports. Just click and download in any format you need.",
      benefits: [
        "Multiple format options (MP4, MOV, GIF)",
        "Custom quality settings",
        "Direct downloads to your device"
      ],
      color: "from-cyan-500 to-blue-500"
    }
  ], []);

  // Handle visibility detection
  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsVisible(true);
        }
      },
      { threshold: 0.1 }
    );

    const section = document.getElementById('features-section');
    if (section) observer.observe(section);

    return () => {
      if (section) observer.unobserve(section);
    };
  }, []);

  // Auto-rotate tabs
  useEffect(() => {
    if (isAutoRotating) {
      // Initial delay before starting rotation
      const initialDelay = setTimeout(() => {
        autoRotateIntervalRef.current = setInterval(() => {
          // Find current index and get next index
          const currentIndex = features.findIndex(f => f.id === activeTab);
          const nextIndex = (currentIndex + 1) % features.length;
          setActiveTab(features[nextIndex].id);
        }, 5000); // Rotate every 5 seconds
      }, 3000); // 3 second initial delay

      return () => {
        clearTimeout(initialDelay);
        if (autoRotateIntervalRef.current) {
          clearInterval(autoRotateIntervalRef.current);
        }
      };
    }
  }, [isAutoRotating, activeTab, features]);

  // Pause auto-rotation when user manually selects a tab
  const handleTabChange = (tabId: string) => {
    setIsAutoRotating(false);
    setActiveTab(tabId);
  };

  return (
    <section id="features-section" className="py-24 px-4 sm:px-6 lg:px-8 relative overflow-hidden">

      <div className="max-w-7xl mx-auto relative z-10">
        <motion.div
          className="text-center mb-16"
          initial={{ opacity: 0, y: 20 }}
          animate={isVisible ? { opacity: 1, y: 0 } : { opacity: 0, y: 20 }}
          transition={{ duration: 0.6 }}
        >
          <div className="flex items-center justify-center gap-2 px-4 py-2 rounded-full bg-brand-500/20 border border-brand-500/30 text-brand-300 text-sm mb-6 backdrop-blur-xl w-fit mx-auto">
            <Sparkles className="w-4 h-4" />
            <span>Powerful Features</span>
          </div>

          <h2 className="text-4xl md:text-5xl font-bold mb-6">
            <span className="text-light-surface-900 dark:text-dark-surface-100">Everything you need</span>{' '}
            <span className="text-gradient">to go viral</span>
          </h2>

          <p className="text-xl text-light-surface-700 dark:text-dark-surface-300 max-w-3xl mx-auto leading-relaxed">
            Stop wasting hours on manual editing. CreatorSync automates your content workflow
            so you can focus on what matters—creating amazing content.
          </p>
        </motion.div>

        {/* Desktop View - Tabs */}
        <div className="hidden md:block">
          <Tabs.Root
            value={activeTab}
            onValueChange={handleTabChange}
            className="w-full"
          >
            <div className="flex flex-col items-center space-y-6 mb-12">
              <Tabs.List
                className="flex flex-wrap justify-center gap-4"
                aria-label="Features"
              >
                {features.map((feature) => (
                  <Tabs.Trigger
                    key={feature.id}
                    value={feature.id}
                    className={cn(
                      "group relative px-6 py-3 rounded-full transition-all duration-300 outline-none",
                      activeTab === feature.id
                        ? "bg-brand-500/20 text-brand-500 border border-brand-500/30"
                        : "bg-light-surface-100/80 dark:bg-dark-surface-800/80 text-light-surface-700 dark:text-dark-surface-300 hover:bg-light-surface-200/80 dark:hover:bg-dark-surface-700/80 border border-light-surface-200/50 dark:border-dark-surface-700/50"
                    )}
                  >
                    <div className="flex items-center gap-3">
                      <feature.icon className="w-5 h-5" />
                      <span className="font-medium">{feature.title}</span>
                    </div>
                  </Tabs.Trigger>
                ))}
              </Tabs.List>

              <button
                onClick={() => setIsAutoRotating(!isAutoRotating)}
                className={`px-4 py-2 rounded-full text-sm transition-colors duration-300 flex items-center gap-2 ${isAutoRotating ? 'bg-brand-500/20 text-brand-500 border border-brand-500/30' : 'bg-light-surface-100/80 dark:bg-dark-surface-800/80 text-light-surface-700 dark:text-dark-surface-300 border border-light-surface-200/50 dark:border-dark-surface-700/50'}`}
              >
                <div className={`w-2 h-2 rounded-full ${isAutoRotating ? 'bg-brand-500 animate-pulse' : 'bg-light-surface-400 dark:bg-dark-surface-600'}`}></div>
                {isAutoRotating ? 'Auto-rotation enabled' : 'Auto-rotation disabled'}
              </button>
            </div>

            {features.map((feature) => (
              <Tabs.Content
                key={feature.id}
                value={feature.id}
                className="focus:outline-none"
              >
                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -20 }}
                  transition={{ duration: 0.4 }}
                  className="grid grid-cols-1 lg:grid-cols-2 gap-8 items-center"
                >
                  {/* Feature Content */}
                  <div className="order-2 lg:order-1">
                    <h3 className="text-3xl font-bold mb-4 text-light-surface-900 dark:text-dark-surface-100">
                      {feature.title}
                    </h3>
                    <p className="text-xl text-light-surface-700 dark:text-dark-surface-300 mb-8">
                      {feature.description}
                    </p>

                    <div className="space-y-4">
                      {feature.benefits.map((benefit, idx) => (
                        <div key={idx} className="flex items-start gap-3">
                          <div className="flex-shrink-0 mt-1">
                            <CheckCircle className="w-5 h-5 text-brand-500" />
                          </div>
                          <p className="text-light-surface-700 dark:text-dark-surface-300">
                            {benefit}
                          </p>
                        </div>
                      ))}
                    </div>

                    <div className="mt-8">
                      <button className="inline-flex items-center gap-2 px-6 py-3 rounded-lg bg-brand-500 hover:bg-brand-600 text-white font-medium transition-colors">
                        Try this feature
                        <ArrowRight className="w-4 h-4" />
                      </button>
                    </div>
                  </div>

                  {/* Feature Visual */}
                  <div className="order-1 lg:order-2 flex justify-center">
                    <div className="relative">
                      <div className={`absolute inset-0 bg-gradient-to-br ${feature.color} opacity-20 blur-2xl rounded-full`}></div>
                      <div className="relative bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm rounded-2xl p-8 border border-light-surface-200/50 dark:border-dark-surface-800/50 shadow-xl">
                        <div className="w-full aspect-video rounded-lg bg-light-surface-200/50 dark:bg-dark-surface-800/50 flex items-center justify-center overflow-hidden">
                          <div className={`p-6 rounded-full bg-gradient-to-br ${feature.color}`}>
                            <feature.icon className="w-16 h-16 text-white" />
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </motion.div>
              </Tabs.Content>
            ))}
          </Tabs.Root>
        </div>

        {/* Mobile View - Accordion */}
        <div className="md:hidden">
          <div className="flex items-center justify-center mb-4">
            <button
              onClick={() => setIsAutoRotating(!isAutoRotating)}
              className={`px-4 py-2 rounded-full text-sm transition-colors duration-300 flex items-center gap-2 ${isAutoRotating ? 'bg-brand-500/20 text-brand-500 border border-brand-500/30' : 'bg-light-surface-100/80 dark:bg-dark-surface-800/80 text-light-surface-700 dark:text-dark-surface-300 border border-light-surface-200/50 dark:border-dark-surface-700/50'}`}
            >
              <div className={`w-2 h-2 rounded-full ${isAutoRotating ? 'bg-brand-500' : 'bg-light-surface-400 dark:bg-dark-surface-600'}`}></div>
              {isAutoRotating ? 'Auto-rotate on' : 'Auto-rotate off'}
            </button>
          </div>
          <Accordion.Root type="single" collapsible className="space-y-4">
            {features.map((feature) => (
              <Accordion.Item
                key={feature.id}
                value={feature.id}
                className="overflow-hidden rounded-xl border border-light-surface-200/50 dark:border-dark-surface-800/50 bg-light-surface-100/90 dark:bg-dark-surface-900/90 backdrop-blur-sm"
              >
                <Accordion.Trigger className="group flex w-full items-center justify-between px-5 py-4 text-left">
                  <div className="flex items-center gap-3">
                    <div className={`p-2 rounded-full bg-gradient-to-br ${feature.color}`}>
                      <feature.icon className="w-5 h-5 text-white" />
                    </div>
                    <h3 className="text-lg font-semibold text-light-surface-900 dark:text-dark-surface-100">
                      {feature.title}
                    </h3>
                  </div>
                  <div className="rounded-full border border-light-surface-300/50 dark:border-dark-surface-700/50 p-1">
                    <ArrowRight className="w-4 h-4 text-light-surface-700 dark:text-dark-surface-300 transition-transform duration-300 group-data-[state=open]:rotate-90" />
                  </div>
                </Accordion.Trigger>
                <Accordion.Content className="overflow-hidden data-[state=closed]:animate-accordion-up data-[state=open]:animate-accordion-down">
                  <div className="px-5 pb-5 pt-2">
                    <p className="text-light-surface-700 dark:text-dark-surface-300 mb-4">
                      {feature.description}
                    </p>

                    <div className="space-y-2">
                      {feature.benefits.map((benefit, idx) => (
                        <div key={idx} className="flex items-start gap-2">
                          <div className="flex-shrink-0 mt-1">
                            <CheckCircle className="w-4 h-4 text-brand-500" />
                          </div>
                          <p className="text-sm text-light-surface-700 dark:text-dark-surface-300">
                            {benefit}
                          </p>
                        </div>
                      ))}
                    </div>
                  </div>
                </Accordion.Content>
              </Accordion.Item>
            ))}
          </Accordion.Root>
        </div>
      </div>
    </section>
  );
}