'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { useTwitchConnection } from '@/hooks/useTwitchConnection';
import { useTwitchApi } from '@/hooks/useTwitchApi';

export function TwitchConnectionTest() {
    const { isConnected, isLoading, error } = useTwitchConnection();
    const { getChannelInfo, getVideos } = useTwitchApi();
    const [testResults, setTestResults] = useState<Record<string, unknown> | null>(null);
    const [testing, setTesting] = useState(false);

    const runTests = async () => {
        if (!isConnected) {
            setTestResults({ error: 'Not connected to Twitch' });
            return;
        }

        setTesting(true);
        try {
            const channelInfo = await getChannelInfo();
            const videos = await getVideos('archive', 5);

            setTestResults({
                channelInfo,
                videos,
                success: true
            });
        } catch (error) {
            setTestResults({ error: error instanceof Error ? error.message : 'Unknown error' });
        } finally {
            setTesting(false);
        }
    };

    return (
        <div className="p-6 border rounded-lg">
            <h3 className="text-lg font-semibold mb-4">Twitch Integration Test</h3>
            
            <div className="space-y-4">
                <div>
                    <strong>Connection Status:</strong> {isConnected ? 'Connected' : 'Not Connected'}
                </div>
                
                {error && (
                    <div className="text-red-500">
                        <strong>Error:</strong> {error}
                    </div>
                )}

                <Button 
                    onClick={runTests} 
                    disabled={!isConnected || isLoading || testing}
                >
                    {testing ? 'Testing...' : 'Test API Calls'}
                </Button>

                {testResults && (
                    <div className="mt-4 p-4 bg-gray-100 rounded">
                        <h4 className="font-semibold mb-2">Test Results:</h4>
                        <pre className="text-sm overflow-auto">
                            {JSON.stringify(testResults, null, 2)}
                        </pre>
                    </div>
                )}
            </div>
        </div>
    );
} 