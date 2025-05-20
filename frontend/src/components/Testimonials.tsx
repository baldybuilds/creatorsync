import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { Star, Quote } from 'lucide-react';
import * as Tabs from '@radix-ui/react-tabs';
import { cn } from '@/lib/utils';

interface Testimonial {
  id: string;
  name: string;
  role: string;
  platform: string;
  followers: string;
  quote: string;
  avatar: string;
  rating: number;
}

export function Testimonials() {
  const [isVisible, setIsVisible] = useState(false);
  const [activeTab, setActiveTab] = useState('tab1');

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsVisible(true);
        }
      },
      { threshold: 0.1 }
    );

    const section = document.getElementById('testimonials-section');
    if (section) observer.observe(section);

    return () => {
      if (section) observer.unobserve(section);
    };
  }, []);

  const testimonials: Testimonial[] = [
    {
      id: 'tab1',
      name: 'Alex Rivera',
      role: 'Twitch Streamer',
      platform: 'Twitch',
      followers: '250K',
      quote: 'CreatorSync has completely transformed my content workflow. I used to spend hours editing clips for different platforms, but now I can do it all in minutes. The time I save goes straight back into creating better streams.',
      avatar: 'https://i.pravatar.cc/150?img=1',
      rating: 5
    },
    {
      id: 'tab2',
      name: 'Samantha Chen',
      role: 'Gaming Content Creator',
      platform: 'YouTube',
      followers: '500K',
      quote: 'The multi-platform export feature is a game-changer. I can optimize my videos for YouTube Shorts, TikTok, and Instagram Reels all at once, with perfect aspect ratios and formatting for each platform. My engagement has increased by 40% since using CreatorSync.',
      avatar: 'https://i.pravatar.cc/150?img=5',
      rating: 5
    },
    {
      id: 'tab3',
      name: 'Marcus Johnson',
      role: 'Esports Commentator',
      platform: 'Multiple',
      followers: '180K',
      quote: 'As someone who covers live events, I need to get highlights out quickly. CreatorSync\'s browser-based editor and fast rendering have been invaluable. I can clip, edit, and publish highlights while the tournament is still ongoing.',
      avatar: 'https://i.pravatar.cc/150?img=8',
      rating: 4
    }
  ];

  return (
    <section id="testimonials-section" className="py-24 px-4 sm:px-6 lg:px-8 relative overflow-hidden">

      <div className="max-w-7xl mx-auto relative z-10">
        <motion.div
          className="text-center mb-16"
          initial={{ opacity: 0, y: 20 }}
          animate={isVisible ? { opacity: 1, y: 0 } : { opacity: 0, y: 20 }}
          transition={{ duration: 0.6 }}
        >
          <div className="flex items-center justify-center gap-2 px-4 py-2 rounded-full bg-brand-500/20 border border-brand-500/30 text-brand-300 text-sm mb-6 backdrop-blur-xl w-fit mx-auto">
            <Star className="w-4 h-4" />
            <span>Creator Testimonials</span>
          </div>
          
          <h2 className="text-4xl md:text-5xl font-bold mb-6">
            <span className="text-light-surface-900 dark:text-dark-surface-100">Loved by</span>{' '}
            <span className="text-gradient">content creators</span>
          </h2>

          <p className="text-xl text-light-surface-700 dark:text-dark-surface-300 max-w-3xl mx-auto leading-relaxed">
            See what creators are saying about how CreatorSync has transformed their content workflow.
          </p>
        </motion.div>

        <Tabs.Root 
          value={activeTab} 
          onValueChange={setActiveTab}
          className="w-full"
        >
          <Tabs.List 
            className="flex flex-wrap justify-center gap-4 mb-12"
            aria-label="Creator testimonials"
          >
            {testimonials.map((testimonial) => (
              <Tabs.Trigger
                key={testimonial.id}
                value={testimonial.id}
                className={cn(
                  "group relative px-6 py-3 rounded-full transition-all duration-300 outline-none",
                  activeTab === testimonial.id 
                    ? "bg-brand-500/20 text-brand-500 border border-brand-500/30" 
                    : "bg-light-surface-100/80 dark:bg-dark-surface-800/80 text-light-surface-700 dark:text-dark-surface-300 hover:bg-light-surface-200/80 dark:hover:bg-dark-surface-700/80 border border-light-surface-200/50 dark:border-dark-surface-700/50"
                )}
              >
                <div className="flex items-center gap-3">
                  <img 
                    src={testimonial.avatar} 
                    alt={testimonial.name} 
                    className={cn(
                      "w-8 h-8 rounded-full object-cover transition-all duration-300",
                      activeTab === testimonial.id 
                        ? "ring-2 ring-brand-500 ring-offset-2 ring-offset-light-surface-50 dark:ring-offset-dark-surface-950" 
                        : ""
                    )}
                  />
                  <span className="font-medium">{testimonial.name}</span>
                </div>
              </Tabs.Trigger>
            ))}
          </Tabs.List>

          {testimonials.map((testimonial) => (
            <Tabs.Content 
              key={testimonial.id}
              value={testimonial.id}
              className="focus:outline-none"
            >
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                transition={{ duration: 0.4 }}
                className="max-w-4xl mx-auto bg-light-surface-100/95 dark:bg-dark-surface-900/95 backdrop-blur-sm rounded-2xl p-8 border border-light-surface-200/50 dark:border-dark-surface-800/50"
              >
                <div className="flex flex-col md:flex-row gap-8 items-center md:items-start">
                  <div className="flex-shrink-0">
                    <img 
                      src={testimonial.avatar} 
                      alt={testimonial.name} 
                      className="w-24 h-24 rounded-full object-cover ring-4 ring-brand-500/20"
                    />
                  </div>
                  
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-1">
                      {[...Array(5)].map((_, i) => (
                        <Star 
                          key={i} 
                          className={cn(
                            "w-5 h-5", 
                            i < testimonial.rating ? "text-yellow-400 fill-yellow-400" : "text-light-surface-300 dark:text-dark-surface-700"
                          )} 
                        />
                      ))}
                    </div>
                    
                    <div className="relative">
                      <Quote className="absolute -top-2 -left-2 w-8 h-8 text-brand-500/20" />
                      <p className="text-xl text-light-surface-900 dark:text-dark-surface-100 italic mb-6 pl-6">
                        {testimonial.quote}
                      </p>
                    </div>
                    
                    <div className="flex flex-wrap items-center gap-x-6 gap-y-2 text-light-surface-700 dark:text-dark-surface-300">
                      <div className="font-semibold text-light-surface-900 dark:text-dark-surface-100">
                        {testimonial.name}
                      </div>
                      <div>{testimonial.role}</div>
                      <div className="flex items-center gap-2">
                        <span>{testimonial.platform}</span>
                        <span className="w-1 h-1 rounded-full bg-light-surface-400 dark:bg-dark-surface-600"></span>
                        <span>{testimonial.followers} followers</span>
                      </div>
                    </div>
                  </div>
                </div>
              </motion.div>
            </Tabs.Content>
          ))}
        </Tabs.Root>
      </div>
    </section>
  );
}
