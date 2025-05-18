"use client";

import { useState, useEffect } from 'react';
import { LogIn, Menu, X, LayoutDashboard } from 'lucide-react';
import { Button } from '../components/ui/button';
import { Hero } from '../components/Hero';
import { SignedIn, SignedOut, SignInButton, UserButton } from '@clerk/clerk-react';
import { useFeatureFlagEnabled } from 'posthog-js/react';
import { usePostHogUser } from '../hooks/usePostHogUser';
import Link from 'next/link'

const LandingPage = () => {
  const [isScrolled, setIsScrolled] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  usePostHogUser();

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 50);
    };

    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  const navItems = [
    { name: 'Features', href: '#features' },
    { name: 'About', href: '#about' },
    { name: 'Contact', href: '#contact' }
  ];

  const clerkAuthButtonEnabled = useFeatureFlagEnabled('clerk-auth-button');

  return (
    <div className="min-h-screen bg-[var(--bg-primary)] text-[var(--text-primary)] overflow-x-hidden">
      {/* Navigation */}
      <nav className={`fixed top-0 left-0 right-0 w-full z-50 transition-all duration-500 ${isScrolled
        ? 'bg-light-surface-50/80 dark:bg-dark-surface-950/80 backdrop-blur-xl border-b border-light-surface-200/50 dark:border-dark-surface-800/50'
        : 'bg-light-surface-50/20 dark:bg-dark-surface-950/20 backdrop-blur-sm'
        }`}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            {/* Logo */}
            <div className="flex items-center">
              <div className="text-2xl font-bold">
                <span className="text-light-surface-900 dark:text-dark-surface-100">Creator</span>
                <span className="text-gradient">Sync</span>
              </div>
            </div>

            {/* Auth Buttons */}
            {clerkAuthButtonEnabled && (
              <div className="hidden md:block">
                <SignedOut>
                  <SignInButton>
                    <Button variant="default" size="sm"><LogIn className="mr-2 h-4 w-4" />Sign In</Button>
                  </SignInButton>
                </SignedOut>
                <SignedIn>
                  <div className="flex items-center gap-4">
                    <Link href="/dashboard">
                      <Button variant="ghost" size="sm"><LayoutDashboard className="mr-2 h-4 w-4" />Dashboard</Button>
                    </Link>
                    <UserButton />
                  </div>
                </SignedIn>
              </div>
            )}
            {/* Mobile menu button */}
            <div className="md:hidden">
              <button
                onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-300 dark:hover:text-dark-surface-100 transition-colors"
              >
                {isMobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
              </button>
            </div>
          </div>
        </div>

        {/* Mobile Navigation Menu */}
        {isMobileMenuOpen && (
          <div className="md:hidden bg-light-surface-100/95 dark:bg-dark-surface-900/95 backdrop-blur-xl border-b border-light-surface-200/50 dark:border-dark-surface-800/50">
            <div className="px-2 pt-2 pb-3 space-y-1 sm:px-3">
              {navItems.map((item) => (
                <a
                  key={item.name}
                  href={item.href}
                  className="block px-3 py-2 rounded-lg text-light-surface-700 hover:text-light-surface-900 hover:bg-light-surface-200/50 dark:text-dark-surface-300 dark:hover:text-dark-surface-100 dark:hover:bg-dark-surface-800/50 transition-all duration-300"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  {item.name}
                </a>
              ))}
              <div className="pt-4">
                <SignedOut>
                  <SignInButton>
                    <Button variant="default" size="sm" className="w-full"><LogIn className="mr-2 h-4 w-4" />Sign In</Button>
                  </SignInButton>
                </SignedOut>
                <SignedIn>
                  <div className="space-y-3 py-2">
                    <Link
                      href="/dashboard"
                      className="flex items-center gap-3 px-3 py-2 rounded-lg text-light-surface-700 hover:text-light-surface-900 hover:bg-light-surface-200/50 dark:text-dark-surface-300 dark:hover:text-dark-surface-100 dark:hover:bg-dark-surface-800/50 transition-all duration-300"
                      onClick={() => setIsMobileMenuOpen(false)}
                    >
                      <LayoutDashboard className="w-5 h-5" />
                      <span>Dashboard</span>
                    </Link>
                    <div className="pt-3 mt-2 border-t border-light-surface-200/50 dark:border-dark-surface-800/50 flex justify-center">
                      <UserButton />
                    </div>
                  </div>
                </SignedIn>
              </div>
            </div>
          </div>
        )}
      </nav>

      {/* Main Content */}
      <Hero />

      {/* Footer */}
      <footer className="bg-light-surface-100 dark:bg-dark-surface-900 border-t border-light-surface-200/50 dark:border-dark-surface-800/50 py-12">
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
            </div>

            {/* Quick Links */}
            <div>
              <h4 className="text-light-surface-900 dark:text-dark-surface-100 font-semibold mb-4">Product (coming soon)</h4>
              <ul className="space-y-3">
                <li><a href="#" className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">Features</a></li>
                <li><a href="#" className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">Pricing</a></li>
              </ul>
            </div>

            {/* Support */}
            <div>
              <h4 className="text-light-surface-900 dark:text-dark-surface-100 font-semibold mb-4">Support (coming soon)</h4>
              <ul className="space-y-3">
                <li><a href="#" className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">Documentation</a></li>
                <li><a href="#" className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">Community</a></li>
              </ul>
            </div>
          </div>

          {/* Bottom Bar */}
          <div className="mt-12 pt-8 border-t border-light-surface-200/50 dark:border-dark-surface-800/50 flex flex-col md:flex-row justify-between items-center">
            <p className="text-light-surface-600 dark:text-dark-surface-400 text-sm">
              &copy; 2025 CreatorSync. All rights reserved. Built for creators, by creators.
            </p>
            <div className="mt-4 md:mt-0 flex space-x-4 text-sm">
              <a href="#" className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">Privacy Policy</a>
              <a href="#" className="text-light-surface-600 hover:text-light-surface-900 dark:text-dark-surface-400 dark:hover:text-dark-surface-100 transition-colors">Terms of Service</a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}

export default LandingPage;
