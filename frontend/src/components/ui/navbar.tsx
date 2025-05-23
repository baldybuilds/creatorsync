'use client';

import * as React from 'react';
import { useState, useEffect } from 'react';
import Link from 'next/link';
import { LogIn, Menu, X, LayoutDashboard } from 'lucide-react';
import { SignedIn, SignedOut, SignInButton, UserButton } from '@clerk/clerk-react';
import { useFeatureFlagEnabled } from 'posthog-js/react';
import * as NavigationMenu from '@radix-ui/react-navigation-menu';
import { Button } from './button';
import { cn } from '@/lib/utils';

interface NavbarProps {
  className?: string;
}

export function Navbar({ className }: NavbarProps) {
  const [isScrolled, setIsScrolled] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const clerkAuthButtonEnabled = useFeatureFlagEnabled('clerk-auth-button');

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

  return (
    <nav
      className={cn(
        'fixed top-4 left-0 right-0 w-full z-50 transition-all duration-300 px-4',
        isScrolled
          ? 'top-2'
          : 'top-4',
        className
      )}
    >
      <div className={cn(
        'mx-50 rounded-xl transition-all duration-500',
        isScrolled
          ? 'bg-light-surface-50/85 dark:bg-dark-surface-950/85 backdrop-blur-xl border border-light-surface-200/50 dark:border-dark-surface-800/50 shadow-lg'
          : 'bg-light-surface-50/30 dark:bg-dark-surface-950/30 backdrop-blur-md'
      )}>

        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            {/* Logo */}
            <div className="flex items-center">
              <Link href="/" className="text-2xl font-bold">
                <span className="text-light-surface-900 dark:text-dark-surface-100">Creator</span>
                <span className="text-gradient">Sync</span>
              </Link>
            </div>

            {/* Desktop Navigation */}
            <NavigationMenu.Root className="hidden md:flex">
              <NavigationMenu.List className="flex space-x-6">
                {navItems.map((item) => (
                  <NavigationMenu.Item key={item.name}>
                    <NavigationMenu.Link
                      asChild
                      className="text-light-surface-700 hover:text-light-surface-900 dark:text-dark-surface-300 dark:hover:text-dark-surface-100 transition-colors font-medium"
                    >
                      <a href={item.href}>{item.name}</a>
                    </NavigationMenu.Link>
                  </NavigationMenu.Item>
                ))}
              </NavigationMenu.List>
            </NavigationMenu.Root>

            {/* Auth Buttons */}
            {(process.env.NODE_ENV !== 'production' || process.env.NEXT_PUBLIC_APP_ENV === 'staging' || clerkAuthButtonEnabled) && (
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

      </div>

      {/* Mobile Navigation Menu */}
      {isMobileMenuOpen && (
        <div className="md:hidden bg-light-surface-100/95 dark:bg-dark-surface-900/95 backdrop-blur-xl rounded-2xl mt-2 border border-light-surface-200/50 dark:border-dark-surface-800/50 mx-4 shadow-lg">
          <div className="px-4 pt-4 pb-4 space-y-2">
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
  );
}
