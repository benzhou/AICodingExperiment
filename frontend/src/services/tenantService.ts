import axios from 'axios';

// Define the Tenant interface
export interface Tenant {
  id: string;
  name: string;
  logoUrl?: string;
  primaryColor?: string;
  secondaryColor?: string;
  domain?: string;
  settings?: Record<string, any>;
}

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

// TenantService for managing tenant-specific configurations
export const tenantService = {
  // Get current tenant based on JWT token or domain
  async getCurrentTenant(): Promise<Tenant | null> {
    try {
      const response = await apiClient.get('/tenants/current');
      return response.data;
    } catch (error) {
      console.error('Error fetching current tenant:', error);
      return null;
    }
  },

  // Get tenant by ID
  async getTenantById(id: string): Promise<Tenant | null> {
    try {
      const response = await apiClient.get(`/tenants/${id}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching tenant with ID ${id}:`, error);
      return null;
    }
  },

  // Update tenant settings (including logo)
  async updateTenant(id: string, data: Partial<Tenant>): Promise<Tenant | null> {
    try {
      const response = await apiClient.put(`/tenants/${id}`, data);
      return response.data;
    } catch (error) {
      console.error(`Error updating tenant with ID ${id}:`, error);
      return null;
    }
  },

  // Upload a custom logo for a tenant
  async uploadTenantLogo(id: string, logoFile: File): Promise<string | null> {
    try {
      const formData = new FormData();
      formData.append('logo', logoFile);
      
      const response = await apiClient.post(`/tenants/${id}/logo`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      
      return response.data.logoUrl;
    } catch (error) {
      console.error(`Error uploading logo for tenant ${id}:`, error);
      return null;
    }
  },

  // Get tenant information based on domain
  async getTenantByDomain(domain: string): Promise<Tenant | null> {
    try {
      const response = await apiClient.get(`/tenants/domain/${domain}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching tenant for domain ${domain}:`, error);
      return null;
    }
  },
}; 