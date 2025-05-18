'use client';

import { useState } from 'react';
import { useAuth } from '@clerk/nextjs';
import { Button } from '@/components/ui/button';

export default function TokenDebugPage() {
  const { getToken } = useAuth();
  const [token, setToken] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState(false);

  const fetchToken = async () => {
    setLoading(true);
    try {
      const sessionToken = await getToken();
      setToken(sessionToken || '');
    } catch (error) {
      console.error('Error fetching token:', error);
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = () => {
    if (token) {
      navigator.clipboard.writeText(token);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  const testApiCall = async () => {
    if (!token) return;

    try {
      const response = await fetch('http://localhost:8080/api/user', {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      const data = await response.json();
      console.log('API Response:', data);
      alert(JSON.stringify(data, null, 2));
    } catch (error) {
      console.error('API Error:', error);
      alert(`Error: ${error}`);
    }
  };

  return (
    <div className="p-8 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">JWT Token Debug Page</h1>

      <div className="mb-6">
        <Button
          onClick={fetchToken}
          disabled={loading}
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
        >
          {loading ? 'Loading...' : 'Get JWT Token'}
        </Button>
      </div>

      {token && (
        <div className="space-y-4">
          <div className="p-4 bg-black rounded-md">
            <h2 className="text-lg font-semibold mb-2">Your JWT Token:</h2>
            <div className="flex items-center space-x-2">
              <textarea
                readOnly
                value={token}
                className="w-full h-32 font-mono text-sm"
              />
              <button
                onClick={copyToClipboard}
                className="px-2 py-1 bg-green-600 text-xs rounded hover:bg-green-700"
              >
                {copied ? 'Copied!' : 'Copy'}
              </button>
            </div>
          </div>

          <div className="p-4 bg-black rounded-md">
            <h2 className="text-lg font-semibold mb-2">Test API Call:</h2>
            <div className="flex items-center space-x-2">
              <code className="bg-black p-2 rounded flex-1 overflow-x-auto">
                curl -H "Authorization: Bearer {token}" http://localhost:8080/api/user
              </code>
              <Button
                onClick={testApiCall}
                className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
              >
                Test
              </Button>
            </div>
          </div>
        </div>
      )}

      <div className="mt-8 text-sm text-gray-500">
        <p>Note: This page is for debugging purposes only. Do not share your JWT token with anyone.</p>
      </div>
    </div>
  );
}
