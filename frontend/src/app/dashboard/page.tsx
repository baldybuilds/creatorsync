"use client";

import { motion } from 'framer-motion';
import { OverviewSection } from '@/components/dashboard/overview-section';
import { ContentSection } from '@/components/dashboard/content-section';
import { AnalyticsSection } from '@/components/dashboard/analytics-section';
import { ScheduledSection } from '@/components/dashboard/scheduled-section';
import { SettingsSection } from '@/components/dashboard/settings-section';
import { useSidebar } from '@/components/dashboard/sidebar-context';

export default function DashboardPage() {
    const { activeSection } = useSidebar();

    const renderSection = () => {
        switch (activeSection) {
            case 'overview':
                return <OverviewSection />;
            case 'content':
                return <ContentSection />;
            case 'analytics':
                return <AnalyticsSection />;
            case 'scheduled':
                return <ScheduledSection />;
            case 'settings':
                return <SettingsSection />;
            default:
                return <OverviewSection />;
        }
    };

            return (
        <div className="min-h-screen">
            {/* Content Area */}
            <motion.div
                key={activeSection}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.3 }}
                className="relative"
            >
                {renderSection()}
            </motion.div>
        </div>
    );
}

