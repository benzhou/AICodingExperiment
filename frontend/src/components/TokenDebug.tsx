import React, { useState, useEffect } from 'react';
import { authService } from '../services/authService';

interface DecodedToken {
  user_id: string;
  email: string;
  exp: number;
}

const UPDATE_INTERVAL = 10000; // 10 seconds

export function TokenDebug() {
  const [tokenInfo, setTokenInfo] = useState<{
    token: string | null;
    decodedToken: DecodedToken | null;
    expiresIn: number | null;
    isValid: boolean;
  }>({
    token: null,
    decodedToken: null,
    expiresIn: null,
    isValid: false,
  });

  useEffect(() => {
    const updateTokenInfo = async () => {
      const token = authService.getToken();
      if (!token) {
        setTokenInfo({
          token: null,
          decodedToken: null,
          expiresIn: null,
          isValid: false,
        });
        return;
      }

      try {
        // Decode token (JWT is base64 encoded)
        const [header, payload] = token.split('.').slice(0, 2);
        const decodedToken = JSON.parse(atob(payload)) as DecodedToken;
        
        // Get token info from server
        const info = await authService.getTokenInfo();

        setTokenInfo({
          token,
          decodedToken,
          expiresIn: info.expires_in,
          isValid: info.expires_in > 0,
        });
      } catch (error) {
        setTokenInfo({
          token,
          decodedToken: null,
          expiresIn: null,
          isValid: false,
        });
      }
    };

    updateTokenInfo(); // Initial call
    const interval = setInterval(updateTokenInfo, UPDATE_INTERVAL); // Update every 10 seconds

    return () => clearInterval(interval);
  }, []);

  if (!tokenInfo.token) {
    return (
      <div className="p-4 bg-gray-100 rounded-lg">
        <h2 className="text-lg font-semibold mb-2">Token Debug</h2>
        <p className="text-red-600">No token found</p>
      </div>
    );
  }

  return (
    <div className="p-4 bg-gray-100 rounded-lg">
      <h2 className="text-lg font-semibold mb-2">Token Debug</h2>
      
      <div className="space-y-4">
        <div>
          <h3 className="font-medium">Status</h3>
          <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
            tokenInfo.isValid ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
          }`}>
            {tokenInfo.isValid ? 'Valid' : 'Invalid'}
          </span>
        </div>

        {tokenInfo.expiresIn !== null && (
          <div>
            <h3 className="font-medium">Expires In</h3>
            <p>{Math.max(0, tokenInfo.expiresIn)} seconds</p>
          </div>
        )}

        {tokenInfo.decodedToken && (
          <div>
            <h3 className="font-medium">Token Payload</h3>
            <pre className="mt-1 bg-gray-800 text-gray-100 p-2 rounded text-sm overflow-x-auto">
              {JSON.stringify(tokenInfo.decodedToken, null, 2)}
            </pre>
          </div>
        )}

        <div>
          <h3 className="font-medium">Raw Token</h3>
          <div className="mt-1 bg-gray-800 text-gray-100 p-2 rounded text-sm overflow-x-auto">
            <code className="break-all">{tokenInfo.token}</code>
          </div>
        </div>
      </div>
    </div>
  );
} 