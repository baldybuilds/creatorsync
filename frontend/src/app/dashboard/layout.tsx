"use client";

import { motion } from 'framer-motion';
import { Sidebar } from '@/components/dashboard/sidebar';
import { SidebarProvider, useSidebar } from '@/components/dashboard/sidebar-context';

function DashboardContent({ children }: { children: React.ReactNode }) {
    const { isCollapsed } = useSidebar();
    
    return (
        <div className="flex min-h-screen bg-[var(--bg-primary)] text-[var(--text-primary)] overflow-x-hidden relative">
            {/* Unified Background Elements */}
            <div className="fixed inset-0 overflow-hidden pointer-events-none">
                {/* Subtle gradient overlay */}
                <div className="absolute inset-0 bg-gradient-to-br from-brand-500/5 via-transparent to-accent-500/10" />

                {/* High-resolution gradient orbs with improved blending */}
                <div
                    className="absolute opacity-30 animate-float"
                    style={{
                        top: '15%',
                        left: '10%',
                        width: '40rem',
                        height: '40rem',
                        background: 'radial-gradient(circle, rgba(99,102,241,0.15) 0%, rgba(99,102,241,0) 70%)',
                        filter: 'blur(60px)',
                        animationDelay: '0s',
                        transform: 'translate3d(0, 0, 0)'
                    }}
                />
                <div
                    className="absolute opacity-30 animate-float"
                    style={{
                        bottom: '10%',
                        right: '5%',
                        width: '35rem',
                        height: '35rem',
                        background: 'radial-gradient(circle, rgba(236,72,153,0.1) 0%, rgba(236,72,153,0) 70%)',
                        filter: 'blur(60px)',
                        animationDelay: '2s',
                        transform: 'translate3d(0, 0, 0)'
                    }}
                />
            </div>

            {/* Modern Sidebar */}
            <Sidebar />

            {/* Main Content Area with dynamic margin */}
            <motion.main
                initial={{ opacity: 0 }}
                animate={{ 
                    opacity: 1,
                    marginLeft: isCollapsed ? '80px' : '256px'
                }}
                transition={{ duration: 0.3 }}
                className="flex-1 relative z-10"
            >
                {children}
            </motion.main>
        </div>
    );
}

export default function DashboardLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return (
        <SidebarProvider>
            <DashboardContent>{children}</DashboardContent>
        </SidebarProvider>
    );
} 