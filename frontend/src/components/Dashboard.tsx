import React, { useState } from 'react';
import { useUser } from '../context/UserContext';
import { authService } from '../services/authService';
import { useNavigate } from 'react-router-dom';
import { TokenModal } from './TokenModal';
import { UserMenu } from './UserMenu';

export function Dashboard() {
  const { user, setUser } = useUser();
  const navigate = useNavigate();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [tokenInfo, setTokenInfo] = useState<{
    token: string | null;
    decodedToken: any | null;
    expiresIn: number | null;
    isValid: boolean;
  }>({
    token: null,
    decodedToken: null,
    expiresIn: null,
    isValid: false,
  });

  const handleLogout = () => {
    authService.logout();
    setUser(null);
    navigate('/login');
  };

  const handleTokenDebug = async () => {
    const token = authService.getToken();
    if (token) {
      const info = await authService.getTokenInfo();
      const [header, payload] = token.split('.').slice(0, 2);
      const decodedToken = JSON.parse(atob(payload));
      setTokenInfo({
        token: token,
        decodedToken: decodedToken,
        expiresIn: info.expires_in !== undefined ? info.expires_in : null,
        isValid: info.expires_in > 0,
      });
      setIsModalOpen(true);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <nav className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <h1 className="text-xl font-semibold">Dashboard</h1>
            </div>
            <div className="flex items-center">
              <UserMenu username={user?.name} onTokenDebug={handleTokenDebug} onLogout={handleLogout} />
            </div>
          </div>
        </div>
      </nav>
      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          {/* Other dashboard content can go here */}
        </div>
      </main>
      <TokenModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        tokenInfo={tokenInfo}
      />
    </div>
  );
} 