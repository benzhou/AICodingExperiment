import React, { useEffect, useState } from 'react';
import { authService } from '../services/authService';

interface TokenModalProps {
  isOpen: boolean;
  onClose: () => void;
  tokenInfo: {
    token: string | null;
    decodedToken: any | null;
    expiresIn: number | null;
    isValid: boolean;
  };
}

export const TokenModal: React.FC<TokenModalProps> = ({ isOpen, onClose, tokenInfo }) => {
  const [currentTokenInfo, setCurrentTokenInfo] = useState(tokenInfo);

  useEffect(() => {
    let interval: NodeJS.Timeout;

    const fetchTokenInfo = async () => {
      if (isOpen) {
        const info = await authService.getTokenInfo();
        const token = authService.getToken();

        if (token) {
          const [header, payload] = token.split('.').slice(0, 2);
          const decodedToken = JSON.parse(atob(payload));

          setCurrentTokenInfo({
            token: token,
            decodedToken: decodedToken,
            expiresIn: info.expires_in,
            isValid: info.expires_in > 0,
          });
        } else {
          setCurrentTokenInfo({
            token: null,
            decodedToken: null,
            expiresIn: null,
            isValid: false,
          });
        }
      }
    };

    if (isOpen) {
      fetchTokenInfo(); // Initial fetch
      interval = setInterval(fetchTokenInfo, 10000); // Fetch every 10 seconds
    }

    return () => clearInterval(interval); // Cleanup on unmount
  }, [isOpen]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 flex items-center justify-center z-50">
      <div className="absolute inset-0 bg-black opacity-50" onClick={onClose}></div>
      <div className="bg-white rounded-lg shadow-lg z-10 p-6 max-w-md w-full">
        <h2 className="text-lg font-semibold mb-2">Token Debug</h2>
        <div className="space-y-4">
          <div>
            <h3 className="font-medium">Status</h3>
            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
              currentTokenInfo.isValid ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
            }`}>
              {currentTokenInfo.isValid ? 'Valid' : 'Invalid'}
            </span>
          </div>

          {currentTokenInfo.expiresIn !== null && (
            <div>
              <h3 className="font-medium">Expires In</h3>
              <p>{Math.max(0, currentTokenInfo.expiresIn)} seconds</p>
            </div>
          )}

          {currentTokenInfo.decodedToken && (
            <div>
              <h3 className="font-medium">Token Payload</h3>
              <pre className="mt-1 bg-gray-800 text-gray-100 p-2 rounded text-sm overflow-x-auto">
                {JSON.stringify(currentTokenInfo.decodedToken, null, 2)}
              </pre>
            </div>
          )}

          <div>
            <h3 className="font-medium">Raw Token</h3>
            <div className="mt-1 bg-gray-800 text-gray-100 p-2 rounded text-sm overflow-x-auto">
              <code className="break-all">{currentTokenInfo.token}</code>
            </div>
          </div>
        </div>
        <button onClick={onClose} className="mt-4 bg-blue-500 text-white px-4 py-2 rounded">
          Close
        </button>
      </div>
    </div>
  );
}; 