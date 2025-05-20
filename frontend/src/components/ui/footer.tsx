'use client';

import * as React from 'react';
import { SocialIcon } from 'react-social-icons'
import * as Separator from '@radix-ui/react-separator';
import { cn } from '@/lib/utils';

interface FooterProps {
  className?: string;
}

export function Footer({ className }: FooterProps) {
  return (
    <footer className={cn(
      "bg-light-surface-100 dark:bg-dark-surface-900 border-t border-light-surface-200/50 dark:border-dark-surface-800/50 py-12",
      className
    )}>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          {/* Company Info */}
          <div className="col-span-1 md:col-span-2">
            <div className="text-2xl font-bold mb-4">
              <span className="text-light-surface-900 dark:text-dark-surface-100">Creator</span>
              <span className="text-gradient">Sync</span>
            </div>
            <p className="text-light-surface-600 dark:text-dark-surface-400 text-lg mb-6 max-w-md">
              Transforming the way creators share their content across platforms.
              Built for streamers, by creators.
            </p>
            <div className="flex space-x-4">
              <div className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">
                <SocialIcon url="https://x.com/creatorsyncapp" style={{ width: 30, height: 30 }} />
              </div>
              <div className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">
                <SocialIcon url="https://instagram.com/creatorsyncapp" style={{ width: 30, height: 30 }} />
              </div>
              <div className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">
                <SocialIcon url="https://twitch.tv/creatorsyncapp" style={{ width: 30, height: 30 }} />
              </div>
              <div className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">
                <SocialIcon url="https://youtube.com/creatorsyncapp" style={{ width: 30, height: 30 }} />
              </div>
            </div>
          </div>

          {/* Quick Links */}
          <div>
            <h4 className="text-light-surface-900 dark:text-dark-surface-100 font-semibold mb-4">Product</h4>
            <ul className="space-y-3">
              <li><a href="#features" className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">Features</a></li>
              <li><a href="#pricing" className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">Pricing</a></li>
              <li><a href="#roadmap" className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">Roadmap</a></li>
              <li><a href="#beta" className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">Beta Program</a></li>
            </ul>
          </div>

          {/* Support */}
          <div>
            <h4 className="text-light-surface-900 dark:text-dark-surface-100 font-semibold mb-4">Support</h4>
            <ul className="space-y-3">
              <li><a href="#docs" className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">Documentation</a></li>
              <li><a href="#community" className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">Community</a></li>
              <li><a href="#contact" className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">Contact Us</a></li>
              <li><a href="#faq" className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">FAQ</a></li>
            </ul>
          </div>
        </div>

        <Separator.Root className="my-8 h-px bg-light-surface-200/50 dark:bg-dark-surface-800/50" />

        {/* Bottom Bar */}
        <div className="flex flex-col md:flex-row justify-between items-center">
          <p className="text-light-surface-600 dark:text-dark-surface-400 text-sm">
            &copy; {new Date().getFullYear()} CreatorSync. All rights reserved. Built for creators, by creators.
          </p>
          <div className="mt-4 md:mt-0 flex space-x-4 text-sm">
            <a href="#privacy" className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">Privacy Policy</a>
            <a href="#terms" className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">Terms of Service</a>
            <a href="#cookies" className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">Cookie Policy</a>
          </div>
        </div>
      </div>
    </footer>
  );
}
