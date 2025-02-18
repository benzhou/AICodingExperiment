import axios from 'axios';

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
}

interface TokenInfo {
  token: string;
  expires_in: number;
}

export const authService = {
  tokenCheckInterval: null as NodeJS.Timeout | null,

  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await axios.post(`${API_URL}/auth/login`, data);
    this.setToken(response.data.token, response.data.expires_in);
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
    localStorage.setItem('token', token);
    localStorage.setItem('tokenExpiry', (Date.now() + expiresIn * 1000).toString());
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    
    // Set up token expiry check
    if (this.tokenCheckInterval) {
      clearInterval(this.tokenCheckInterval);
    }
    this.tokenCheckInterval = setInterval(() => this.checkTokenExpiry(), 60000); // Check every minute
  },

  getToken(): string | null {
    return localStorage.getItem('token');
  },

  async checkTokenExpiry() {
    const expiry = localStorage.getItem('tokenExpiry');
    if (expiry && Date.now() > parseInt(expiry)) {
      this.logout();
      window.location.href = '/login';
    }
  },

  async getTokenInfo(): Promise<TokenInfo> {
    const response = await axios.get(`${API_URL}/auth/token-info`);
    return response.data;
  },

  logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('tokenExpiry');
    delete axios.defaults.headers.common['Authorization'];
    if (this.tokenCheckInterval) {
      clearInterval(this.tokenCheckInterval);
    }
  },

  isAuthenticated(): boolean {
    const token = this.getToken();
    const expiry = localStorage.getItem('tokenExpiry');
    if (!token || !expiry) {
      return false;
    }
    return Date.now() < parseInt(expiry);
  }
}; 