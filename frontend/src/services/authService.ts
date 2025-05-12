import axios from 'axios';
import { jwtDecode } from 'jwt-decode';

const API_URL = 'http://localhost:8080/api/v1';

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

export interface User {
  id: string;
  email: string;
  name: string;
  authProvider: string;
}

export interface AuthResponse {
  token: string;
  user: User;
  expires_in: number;
  roles?: string[];
  is_admin?: boolean;
}

interface TokenInfo {
  token: string;
  expires_in: number;
}

interface DecodedToken {
  exp: number;
  iat: number;
  user_id: string;
  [key: string]: any;
}

export const authService = {
  tokenCheckInterval: null as NodeJS.Timeout | null,

  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await axios.post(`${API_URL}/auth/login`, data);
    this.setToken(response.data.token, response.data.expires_in);
    
    // Set admin status in localStorage based on backend response
    if (response.data.is_admin === true) {
      localStorage.setItem('isAdmin', 'true');
    } else {
      localStorage.removeItem('isAdmin');
    }
    
    return response.data;
  },

  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await axios.post(`${API_URL}/auth/register`, data);
    this.setToken(response.data.token, response.data.expires_in);
    return response.data;
  },

  async loginWithGoogle() {
    window.location.href = `${API_URL}/auth/google`;
  },

  setToken(token: string, expiresIn: number) {
    // Store token in localStorage
    localStorage.setItem('token', token);
    
    // Try to get expiration from JWT first
    try {
      const decoded = jwtDecode<DecodedToken>(token);
      if (decoded.exp) {
        // exp is in seconds, convert to milliseconds
        const expiryTime = decoded.exp * 1000;
        localStorage.setItem('tokenExpiry', expiryTime.toString());
      } else {
        // Fallback to using the expires_in value
        localStorage.setItem('tokenExpiry', (Date.now() + expiresIn * 1000).toString());
      }
    } catch (error) {
      console.error('Error decoding token:', error);
      // Fallback to using the expires_in value
      localStorage.setItem('tokenExpiry', (Date.now() + expiresIn * 1000).toString());
    }
    
    // Set Authorization header for all axios requests
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    
    // Set up token expiry check
    this.setupTokenExpiryCheck();
  },

  setupTokenExpiryCheck() {
    // Clear any existing interval
    if (this.tokenCheckInterval) {
      clearInterval(this.tokenCheckInterval);
    }
    
    // Check token expiry every minute
    this.tokenCheckInterval = setInterval(() => {
      this.checkTokenExpiry();
    }, 60000);
    
    // Also check immediately
    this.checkTokenExpiry();
  },

  getToken(): string | null {
    return localStorage.getItem('token');
  },

  checkTokenExpiry(): boolean {
    const token = this.getToken();
    const expiry = localStorage.getItem('tokenExpiry');
    
    // If no token or expiry, consider as expired
    if (!token || !expiry) {
      console.log('No token or expiry found, considering as expired');
      this.handleExpiredToken();
      return false;
    }

    const expiryTime = parseInt(expiry);
    // If current time is past expiry time (with 10 second buffer)
    if (Date.now() > expiryTime - 10000) {
      console.log('Token has expired or is about to expire');
      this.handleExpiredToken();
      return false;
    }
    
    // If token is valid but will expire in the next 5 minutes, try to refresh it
    if (Date.now() > expiryTime - 5 * 60 * 1000) {
      console.log('Token will expire soon, attempting refresh');
      this.refreshToken().catch(error => {
        console.error('Error refreshing token:', error);
      });
    }
    
    return true;
  },

  handleExpiredToken() {
    console.log('Handling expired token');
    this.logout();
    
    // Only redirect if we're not already on the login page
    if (window.location.pathname !== '/login') {
      // Store the current location to redirect back after login
      sessionStorage.setItem('redirectPath', window.location.pathname);
      window.location.href = '/login?expired=true';
    }
  },

  async refreshToken(): Promise<boolean> {
    try {
      // This is a placeholder for a token refresh endpoint
      // You would need to implement this on your backend
      const response = await axios.post(`${API_URL}/auth/refresh-token`);
      this.setToken(response.data.token, response.data.expires_in);
      return true;
    } catch (error) {
      // If refresh fails, logout the user
      this.handleExpiredToken();
      return false;
    }
  },

  async getTokenInfo(): Promise<TokenInfo> {
    const response = await axios.get(`${API_URL}/auth/token-info`);
    return response.data;
  },

  logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('tokenExpiry');
    localStorage.removeItem('isAdmin');
    delete axios.defaults.headers.common['Authorization'];
    if (this.tokenCheckInterval) {
      clearInterval(this.tokenCheckInterval);
      this.tokenCheckInterval = null;
    }
  },

  isAuthenticated(): boolean {
    const token = this.getToken();
    if (!token) {
      return false;
    }
    
    try {
      // Try to decode the token to check its validity
      const decoded = jwtDecode<DecodedToken>(token);
      
      // Check if the token is expired
      if (decoded.exp) {
        // exp is in seconds, convert to milliseconds
        return Date.now() < decoded.exp * 1000;
      }
      
      // Fallback to localStorage expiry
      const expiry = localStorage.getItem('tokenExpiry');
      if (!expiry) {
        return false;
      }
      
      return Date.now() < parseInt(expiry);
    } catch (error) {
      console.error('Error decoding token:', error);
      return false;
    }
  }
}; 