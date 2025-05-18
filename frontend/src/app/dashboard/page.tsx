"use client";

import { useState } from 'react';
import { SignOutButton } from '@clerk/clerk-react';
import { Button } from '../../components/ui/button'; // Adjusted import path
import { Users, BarChart2, Video, Calendar, LogOut, Settings as SettingsIcon } from 'lucide-react'; // Icons for sidebar

// Placeholder components for tab content
const OverviewContent = () => (
    <div>
        <h2 className="text-2xl font-semibold text-foreground mb-6">Platform Overview</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <div className="bg-card text-card-foreground p-6 rounded-lg shadow-md border border-border">
                <div className="flex items-center text-primary mb-3">
                    <Users className="w-7 h-7 mr-3" />
                    <h3 className="text-sm font-medium text-muted-foreground">Followers</h3>
                </div>
                <p className="text-2xl font-bold text-card-foreground">1,234</p>
                <p className="text-xs text-muted-foreground mt-1">+23 since last week</p>
            </div>
            <div className="bg-card text-card-foreground p-6 rounded-lg shadow-md border border-border">
                <div className="flex items-center text-primary mb-3">
                    <Users className="w-7 h-7 mr-3" /> {/* Using Users icon for subs too, can change */}
                    <h3 className="text-sm font-medium text-muted-foreground">Subscribers</h3>
                </div>
                <p className="text-2xl font-bold text-card-foreground">567</p>
                <p className="text-xs text-muted-foreground mt-1">+5 since last month</p>
            </div>
            <div className="bg-card text-card-foreground p-6 rounded-lg shadow-md border border-border">
                <div className="flex items-center text-primary mb-3">
                    <BarChart2 className="w-7 h-7 mr-3" />
                    <h3 className="text-sm font-medium text-muted-foreground">Avg. Views (Last 7 Streams)</h3>
                </div>
                <p className="text-2xl font-bold text-card-foreground">8,765</p>
                <p className="text-xs text-muted-foreground mt-1">-2% vs previous 7 streams</p>
            </div>
        </div>
        {/* Placeholder for more overview charts/data */}
        <div className="bg-card text-card-foreground p-6 rounded-lg shadow-md border border-border">
            <h3 className="text-xl font-semibold text-card-foreground mb-4">Activity Feed (Placeholder)</h3>
            <ul className="space-y-3">
                <li className="text-sm text-muted-foreground">New Clip: "Epic Win!" - 5 mins ago</li>
                <li className="text-sm text-muted-foreground">Scheduled Post: "Stream starting soon!" - 1 hour ago</li>
                <li className="text-sm text-muted-foreground">New Follower: "CoolDude123" - 3 hours ago</li>
            </ul>
        </div>
    </div>
);

const ContentSection = () => <div className="p-4 bg-card text-card-foreground rounded-lg shadow-md border border-border"><h2 className="text-2xl font-semibold text-card-foreground">Content Management</h2><p className="text-muted-foreground mt-2">Recent broadcasts and clips will be shown here. Manage your VODs, highlights, and clips.</p></div>;
const AnalyticsSection = () => <div className="p-4 bg-card text-card-foreground rounded-lg shadow-md border border-border"><h2 className="text-2xl font-semibold text-card-foreground">Detailed Analytics</h2><p className="text-muted-foreground mt-2">A deeper view of your stream and content analytics will be available here.</p></div>;
const ScheduledSection = () => <div className="p-4 bg-card text-card-foreground rounded-lg shadow-md border border-border"><h2 className="text-2xl font-semibold text-card-foreground">Scheduled Posts</h2><p className="text-muted-foreground mt-2">Manage your scheduled social media posts and stream announcements.</p></div>;
const SettingsSection = () => <div className="p-4 bg-card text-card-foreground rounded-lg shadow-md border border-border"><h2 className="text-2xl font-semibold text-card-foreground">Platform Settings</h2><p className="text-muted-foreground mt-2">Configure your CreatorSync account, integrations, and preferences.</p></div>;


const DashboardPage = () => {
    const [activeSection, setActiveSection] = useState('Overview');

    const sections = [
        { name: 'Overview', icon: BarChart2, component: <OverviewContent /> },
        { name: 'Content', icon: Video, component: <ContentSection /> },
        { name: 'Analytics', icon: BarChart2, component: <AnalyticsSection /> }, // Can use a different icon for detailed analytics
        { name: 'Scheduled', icon: Calendar, component: <ScheduledSection /> },
        { name: 'Settings', icon: SettingsIcon, component: <SettingsSection /> },
        {
            name: 'Sign Out', icon: LogOut, isSignOut: true as const
        },
    ];

    const renderSection = () => {
        const section = sections.find(s => s.name === activeSection);
        return section ? section.component : <OverviewContent />;
    };

    return (
        <div className="flex min-h-screen bg-[var(--bg-primary)] text-[var(--text-primary)]">
            {/* Sidebar Navigation */}
            <aside className="w-64 bg-sidebar text-sidebar-foreground p-6 space-y-4 border-r border-sidebar-border fixed top-0 left-0 h-full shadow-lg">
                <div className="text-2xl font-bold mb-8">
                    <span className="text-sidebar-foreground">Creator</span>
                    <span className="text-primary">Sync</span>
                </div>
                <nav className="space-y-2">
                    {sections.map((section) => {
                        if (section.name === 'Sign Out' && section.isSignOut) {
                            return (
                                <SignOutButton key={section.name}>
                                    <Button
                                        variant="ghost"
                                        size="lg" // Matches h-9, px-3. Other nav items are px-3 py-2.5.
                                        className="w-full justify-start text-muted-foreground focus-visible:ring-ring hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
                                    >
                                        <LogOut className="w-5 h-5 mr-0.5" />
                                        {section.name}
                                    </Button>
                                </SignOutButton>
                            );
                        }
                        const isActive = activeSection === section.name;
                        return (
                            <button
                                key={section.name}
                                onClick={() => setActiveSection(section.name)}
                                className={`flex items-center w-full px-3 py-2.5 rounded-lg group transition-all duration-300 ${isActive
                                    ? 'bg-sidebar-accent text-sidebar-accent-foreground shadow-sm'
                                    : 'text-muted-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground'
                                    }`}
                            >
                                <section.icon className={`w-5 h-5 mr-3 transition-all duration-300 ${isActive ? 'text-sidebar-accent-foreground' : 'text-muted-foreground group-hover:text-accent-foreground'}`} />
                                <span className={`transition-all duration-300 ${isActive ? 'text-sidebar-accent-foreground' : 'text-sidebar-foreground group-hover:text-accent-foreground'}`}>{section.name}</span>
                            </button>
                        );
                    })}
                </nav>
                {/* You can add user profile/logout at the bottom of sidebar later */}
            </aside>

            {/* Main Content Area */}
            <main className="flex-1 p-8 ml-64"> {/* ml-64 to offset for sidebar width */}
                {renderSection()}
            </main>
        </div>
    );
};

export default DashboardPage;

