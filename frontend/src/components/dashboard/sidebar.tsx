'use client';

import { motion } from 'framer-motion';
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { SignOutButton } from "@clerk/nextjs";
import { 
    BarChart2, 
    Video, 
    Calendar, 
    Settings as SettingsIcon, 
    LogOut,
    ChevronLeft,
    ChevronRight
} from 'lucide-react';

import { useSidebar } from './sidebar-context';

interface SidebarProps {
    className?: string;
}

export function Sidebar({ className }: SidebarProps) {
    const { isCollapsed, setIsCollapsed, activeSection, setActiveSection } = useSidebar();

    const sections = [
        { name: 'Overview', icon: BarChart2, id: 'overview' },
        { name: 'Content', icon: Video, id: 'content' },
        { name: 'Analytics', icon: BarChart2, id: 'analytics' },
        { name: 'Scheduled', icon: Calendar, id: 'scheduled' },
        { name: 'Settings', icon: SettingsIcon, id: 'settings' },
    ];

    return (
        <motion.aside
            initial={{ x: -100, opacity: 0 }}
            animate={{ 
                x: 0, 
                opacity: 1,
                width: isCollapsed ? '80px' : '256px'
            }}
            transition={{ duration: 0.3, ease: "easeOut" }}
            className={cn(
                'fixed top-0 left-0 h-full z-30 transition-all duration-300',
                'bg-light-surface-50/30 dark:bg-dark-surface-950/30 backdrop-blur-xl',
                'border-r border-light-surface-200/20 dark:border-dark-surface-800/20',
                'shadow-2xl',
                className
            )}
        >
            {/* Collapse/Expand Button */}
            <button
                onClick={() => setIsCollapsed(!isCollapsed)}
                className={cn(
                    'absolute -right-3 top-8 z-40',
                    'w-6 h-6 rounded-full',
                    'bg-light-surface-100 dark:bg-dark-surface-900',
                    'border border-light-surface-200 dark:border-dark-surface-800',
                    'flex items-center justify-center',
                    'hover:bg-light-surface-200 dark:hover:bg-dark-surface-800',
                    'transition-all duration-200',
                    'text-light-surface-600 dark:text-dark-surface-400'
                )}
            >
                {isCollapsed ? (
                    <ChevronRight className="w-3 h-3" />
                ) : (
                    <ChevronLeft className="w-3 h-3" />
                )}
            </button>

            <div className={cn(
                "space-y-4 h-full flex flex-col transition-all duration-300",
                isCollapsed ? "p-3" : "p-6"
            )}>
                {/* Logo */}
                <motion.div 
                    className="mb-8"
                    animate={{
                        opacity: isCollapsed ? 0 : 1,
                        x: isCollapsed ? -20 : 0
                    }}
                    transition={{ duration: 0.2 }}
                >
                    {!isCollapsed && (
                        <div className="text-2xl font-bold cursor-pointer">
                            <span className="text-light-surface-900 dark:text-dark-surface-100">Creator</span>
                            <span className="text-gradient">Sync</span>
                        </div>
                    )}
                </motion.div>

                {/* Navigation */}
                <nav className={cn(
                    "flex-1 transition-all duration-300",
                    isCollapsed ? "space-y-3" : "space-y-2"
                )}>
                    {sections.map((section, index) => {
                        const isActive = activeSection === section.id;
                        
                        return (
                            <motion.div
                                key={section.name}
                                initial={{ opacity: 0, x: -20 }}
                                animate={{ opacity: 1, x: 0 }}
                                transition={{ delay: index * 0.1 }}
                            >
                                <button
                                    onClick={() => setActiveSection(section.id)}
                                    className={cn(
                                        "flex items-center rounded-xl group transition-all duration-300 relative overflow-hidden w-full",
                                        isActive
                                            ? "bg-gradient-to-r from-brand-500/20 to-brand-600/10 text-brand-500 border border-brand-500/30 shadow-lg shadow-brand-500/10"
                                            : "text-light-surface-700 dark:text-dark-surface-300 hover:bg-light-surface-100/80 dark:hover:bg-dark-surface-800/80 hover:text-light-surface-900 dark:hover:text-dark-surface-100 border border-transparent hover:border-light-surface-300/50 dark:hover:border-dark-surface-700/50",
                                        isCollapsed ? "justify-center px-3 py-3 w-12 h-12" : "justify-start px-4 py-3"
                                    )}
                                >
                                    {/* Animated background gradient */}
                                    {isActive && (
                                        <motion.div
                                            layoutId="activeTab"
                                            className="absolute inset-0 bg-gradient-to-r from-brand-500/10 to-transparent rounded-xl"
                                            transition={{ type: "spring", bounce: 0.2, duration: 0.6 }}
                                        />
                                    )}
                                    
                                    <section.icon className={cn(
                                        "transition-all duration-300 relative z-10",
                                        isActive ? "text-brand-500" : "text-light-surface-500 dark:text-dark-surface-400 group-hover:text-light-surface-900 dark:group-hover:text-dark-surface-100",
                                        isCollapsed ? "w-5 h-5" : "w-5 h-5 mr-3"
                                    )} />
                                    
                                    {!isCollapsed && (
                                        <motion.span
                                            className="relative z-10 font-medium"
                                            initial={{ opacity: 0 }}
                                            animate={{ opacity: 1 }}
                                            exit={{ opacity: 0 }}
                                        >
                                            {section.name}
                                        </motion.span>
                                    )}

                                    {/* Tooltip for collapsed state */}
                                    {isCollapsed && (
                                        <div className="absolute left-full ml-3 px-3 py-2 bg-light-surface-100/95 dark:bg-dark-surface-800/95 backdrop-blur-xl border border-light-surface-200/50 dark:border-dark-surface-700/50 text-light-surface-900 dark:text-dark-surface-100 text-xs rounded-lg opacity-0 group-hover:opacity-100 transition-all duration-200 pointer-events-none whitespace-nowrap z-50 shadow-lg">
                                            {section.name}
                                        </div>
                                    )}
                                </button>
                            </motion.div>
                        );
                    })}
                </nav>

                {/* Sign Out Button */}
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.5 }}
                >
                    <SignOutButton>
                        <Button
                            variant="ghost"
                            size={isCollapsed ? "icon" : "lg"}
                            className={cn(
                                "transition-all duration-300 text-light-surface-700 dark:text-dark-surface-300 hover:bg-rose-500/10 hover:text-rose-500 focus-visible:ring-rose-500/20 group relative",
                                isCollapsed ? "w-12 h-12 justify-center" : "w-full justify-start"
                            )}
                        >
                            <LogOut className={cn(
                                "transition-all duration-300",
                                isCollapsed ? "w-5 h-5" : "w-5 h-5 mr-2"
                            )} />
                            {!isCollapsed && "Sign Out"}
                            
                            {/* Tooltip for collapsed sign out button */}
                            {isCollapsed && (
                                <div className="absolute left-full ml-3 px-3 py-2 bg-light-surface-100/95 dark:bg-dark-surface-800/95 backdrop-blur-xl border border-light-surface-200/50 dark:border-dark-surface-700/50 text-light-surface-900 dark:text-dark-surface-100 text-xs rounded-lg opacity-0 group-hover:opacity-100 transition-all duration-200 pointer-events-none whitespace-nowrap z-50 shadow-lg">
                                    Sign Out
                                </div>
                            )}
                        </Button>
                    </SignOutButton>
                </motion.div>
            </div>
        </motion.aside>
    );
} 