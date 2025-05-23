// Analytics API service for frontend
export interface AnalyticsOverview {
    currentFollowers: number;
    followerChange: number;
    followerChangePercent: number;
    currentSubscribers: number;
    subscriberChange: number;
    totalViews: number;
    viewChange: number;
    averageViewers: number;
    viewerChange: number;
    streamsLast30Days: number;
    hoursStreamedLast30: number;
}

export interface ChartDataPoint {
    date: string;
    value: number;
    label?: string;
}

export interface AnalyticsChartData {
    followerGrowth: ChartDataPoint[];
    viewershipTrends: ChartDataPoint[];
    streamFrequency: ChartDataPoint[];
    topGames: ChartDataPoint[];
    videoPerformance: ChartDataPoint[];
}

export interface GrowthMetric {
    current: number;
    previous: number;
    change: number;
    percentChange: number;
    trend: string;
}

export interface GrowthAnalysis {
    period: string;
    metrics: Record<string, GrowthMetric>;
}

export interface VideoAnalytics {
    id: number;
    videoId: string;
    title: string;
    videoType: string;
    duration: number;
    viewCount: number;
    thumbnailUrl: string;
    publishedAt: string;
}

export interface GameAnalytics {
    id: number;
    userId: string;
    gameId: string;
    gameName: string;
    totalStreams: number;
    totalHoursStreamed: number;
    averageViewers: number;
    peakViewers: number;
    totalFollowersGained: number;
    lastStreamedAt: string;
}

export interface ContentPerformance {
    topVideos: VideoAnalytics[];
    topGames: GameAnalytics[];
    insights: string[];
}

export interface AnalyticsJob {
    id: number;
    userId: string;
    jobType: string;
    status: string;
    startedAt?: string;
    completedAt?: string;
    errorMessage?: string;
    dataDate?: string;
    createdAt: string;
}

export interface SystemStats {
    totalUsers: number;
    activeUsers: number;
    totalJobs: number;
    successfulJobs: number;
    failedJobs: number;
    successRate: number;
    averageCollectionTime: string;
    lastCollectionRun: string;
}

class AnalyticsService {
    private baseUrl: string;

    constructor() {
        // Check if we're in staging environment
        if (typeof window !== 'undefined' && window.location.hostname === 'dev.creatorsync.app') {
            this.baseUrl = 'https://api-dev.creatorsync.app';
        } else if (process.env.NEXT_PUBLIC_APP_ENV === 'staging') {
            this.baseUrl = 'https://api-dev.creatorsync.app';
        } else if (process.env.NODE_ENV === 'production') {
            this.baseUrl = 'https://api.creatorsync.app';
        } else {
            this.baseUrl = 'http://localhost:8080';
        }
    }

    private async makeRequest<T>(
        endpoint: string,
        options: RequestInit = {},
        token?: string
    ): Promise<T> {
        const url = `${this.baseUrl}${endpoint}`;

        const headers: Record<string, string> = {
            'Content-Type': 'application/json',
            ...(options.headers as Record<string, string>),
        };

        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        const response = await fetch(url, {
            ...options,
            headers,
        });

        if (!response.ok) {
            throw new Error(`API request failed: ${response.status} ${response.statusText}`);
        }

        return response.json();
    }

    // Get dashboard overview
    async getDashboardOverview(token: string): Promise<AnalyticsOverview> {
        return this.makeRequest<AnalyticsOverview>('/api/analytics/overview', {}, token);
    }

    // Get analytics chart data
    async getAnalyticsChartData(token: string, days: number = 30): Promise<AnalyticsChartData> {
        return this.makeRequest<AnalyticsChartData>(
            `/api/analytics/charts?days=${days}`,
            {},
            token
        );
    }

    // Get detailed analytics
    async getDetailedAnalytics(token: string): Promise<{
        overview: AnalyticsOverview;
        charts: AnalyticsChartData;
        topStreams: unknown[];
        topVideos: VideoAnalytics[];
        topGames: GameAnalytics[];
        recentActivity: unknown[];
    }> {
        return this.makeRequest('/api/analytics/detailed', {}, token);
    }

    // Get growth analysis
    async getGrowthAnalysis(token: string, period: string = 'month'): Promise<GrowthAnalysis> {
        return this.makeRequest<GrowthAnalysis>(
            `/api/analytics/growth?period=${period}`,
            {},
            token
        );
    }

    // Get content performance
    async getContentPerformance(token: string): Promise<ContentPerformance> {
        return this.makeRequest<ContentPerformance>('/api/analytics/content', {}, token);
    }

    // Trigger data collection
    async triggerDataCollection(token: string): Promise<{ message: string; user_id: string; timestamp: number }> {
        return this.makeRequest(
            '/api/analytics/trigger',
            { method: 'POST' },
            token
        );
    }

    // Refresh channel data
    async refreshChannelData(token: string): Promise<{ message: string; user_id: string; timestamp: number }> {
        return this.makeRequest(
            '/api/analytics/refresh',
            { method: 'POST' },
            token
        );
    }

    // Get analytics jobs
    async getAnalyticsJobs(token: string, limit: number = 10): Promise<{ jobs: AnalyticsJob[]; user_id: string; timestamp: number }> {
        return this.makeRequest<{ jobs: AnalyticsJob[]; user_id: string; timestamp: number }>(
            `/api/analytics/jobs?limit=${limit}`,
            {},
            token
        );
    }

    // Get system stats (admin only)
    async getSystemStats(adminKey: string): Promise<SystemStats> {
        return this.makeRequest<SystemStats>(
            '/api/analytics/system',
            { headers: { 'X-Admin-Key': adminKey } }
        );
    }

    // Trigger daily collection (admin only)
    async triggerDailyCollection(adminKey: string): Promise<{ message: string; timestamp: number }> {
        return this.makeRequest(
            '/api/analytics/daily',
            {
                method: 'POST',
                headers: { 'X-Admin-Key': adminKey }
            }
        );
    }

    // Health check
    async healthCheck(): Promise<{ status: string; service: string; timestamp: number }> {
        return this.makeRequest('/api/analytics/health');
    }
}

// Export singleton instance
export const analyticsService = new AnalyticsService();
export default analyticsService; 
