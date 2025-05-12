import React, { createContext, useContext, useState, useEffect } from 'react';
import { User, authService } from '../services/authService';
import { userService } from '../services/userService';
import axios from 'axios';

interface UserContextType {
  user: User | null;
  setUser: (user: User | null) => void;
  loading: boolean;
  refreshUser: () => Promise<void>;
}

const UserContext = createContext<UserContextType | undefined>(undefined);

export function UserProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  // Function to refresh user data
  const refreshUser = async () => {
    try {
      const currentUser = user;
      if (!currentUser || !currentUser.id) return;
      
      // Get updated user data including roles
      const userData = await userService.getUserById(currentUser.id);
      
      // Convert between user types if necessary
      const updatedUser: User = {
        id: userData.user.id,
        email: userData.user.email,
        name: userData.user.name,
        authProvider: currentUser.authProvider || 'local' // Preserve auth provider from current user
      };
      
      setUser(updatedUser);
    } catch (error) {
      console.error('Error refreshing user data:', error);
    }
  };

  // Function to load user data from token
  const loadUserFromToken = async () => {
    try {
      setLoading(true);
      const token = authService.getToken();
      
      if (!token) {
        setLoading(false);
        return;
      }
      
      // Set the token in axios for subsequent requests
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      
      try {
        // Try to get token info to validate it
        await authService.getTokenInfo();
        
        // Decode token to get user ID (assuming JWT)
        const tokenParts = token.split('.');
        if (tokenParts.length === 3) {
          try {
            const payload = JSON.parse(atob(tokenParts[1]));
            // Backend uses "user_id" for the user ID in the JWT
            const userId = payload.user_id;
            
            if (userId) {
              // Get user data
              const userData = await userService.getUserById(userId);
              
              if (userData && userData.user) {
                const loadedUser: User = {
                  id: userData.user.id,
                  email: userData.user.email,
                  name: userData.user.name,
                  authProvider: 'local' // Default to local auth provider
                };
                
                setUser(loadedUser);
                
                // Check admin status from roles
                if (userData.roles && userData.roles.includes('admin')) {
                  localStorage.setItem('isAdmin', 'true');
                }
              }
            }
          } catch (error) {
            console.error('Error parsing token payload:', error);
            // Token is invalid or expired
            authService.logout();
          }
        }
      } catch (error) {
        console.error('Error validating token:', error);
        // Token is invalid or expired
        authService.logout();
      }
    } catch (error) {
      console.error('Error loading user from token:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadUserFromToken();
  }, []);

  return (
    <UserContext.Provider value={{ user, setUser, loading, refreshUser }}>
      {children}
    </UserContext.Provider>
  );
}

export function useUser() {
  const context = useContext(UserContext);
  if (context === undefined) {
    throw new Error('useUser must be used within a UserProvider');
  }
  return context;
} 