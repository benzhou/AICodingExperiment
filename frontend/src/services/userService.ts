import axios from 'axios';

// API base URL
const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

// Create axios instance with authorization header
const apiClient = axios.create({
  baseURL: API_URL,
});

// Add request interceptor to include auth token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

export interface User {
  id: string;
  email: string;
  name: string;
}

export interface Role {
  role: string;
}

export interface CreateUserRequest {
  email: string;
  password: string;
  name: string;
  role: string;
}

export interface UpdateRoleRequest {
  role: string;
  operation: 'add' | 'remove';
}

export const userService = {
  // Get user by ID
  async getUserById(userId: string): Promise<{ user: User; roles: string[] }> {
    try {
      const response = await apiClient.get(`/users/${userId}`);
      return response.data;
    } catch (error) {
      console.error('Error fetching user:', error);
      throw error;
    }
  },

  // Get all users (admin only)
  async getAllUsers(): Promise<User[]> {
    try {
      const response = await apiClient.get('/users');
      return response.data;
    } catch (error) {
      console.error('Error fetching users:', error);
      throw error;
    }
  },

  // Get user roles
  async getUserRoles(userId: string): Promise<string[]> {
    try {
      const response = await apiClient.get(`/users/${userId}/roles`);
      return response.data.roles;
    } catch (error) {
      console.error('Error fetching user roles:', error);
      throw error;
    }
  },

  // Update user role (add or remove)
  async updateUserRole(userId: string, roleData: UpdateRoleRequest): Promise<{ roles: string[] }> {
    try {
      const response = await apiClient.put(`/users/${userId}/roles`, roleData);
      return response.data;
    } catch (error) {
      console.error('Error updating user role:', error);
      throw error;
    }
  },

  // Set admin role
  async setAdminRole(userId: string): Promise<void> {
    try {
      await apiClient.put(`/users/${userId}/admin`);
    } catch (error) {
      console.error('Error setting admin role:', error);
      throw error;
    }
  },

  // Create new user with role
  async createUser(userData: CreateUserRequest): Promise<{ user: User; roles: string[] }> {
    try {
      const response = await apiClient.post('/users', userData);
      return response.data;
    } catch (error) {
      console.error('Error creating user:', error);
      throw error;
    }
  },
}; 