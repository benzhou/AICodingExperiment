import axios from 'axios';
import { authService } from './authService';

// Create axios instance with base URL pointing to the backend
// Hardcoded for reliability - this is the backend server
const API_URL = 'http://localhost:8080';

// Create the axios instance with improved debugging
const api = axios.create({
  baseURL: API_URL,
  timeout: 15000, // Increased timeout for file uploads
  headers: {
    'Content-Type': 'application/json',
  }
});

// Add request interceptor to include auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    
    // Don't set Content-Type for FormData
    if (config.data instanceof FormData) {
      // Let axios set the Content-Type to multipart/form-data with boundary
      delete config.headers['Content-Type'];
      console.log('Removed Content-Type header for FormData request');
    }
    
    console.log('API Request:', {
      method: config.method,
      url: config.url,
      baseURL: config.baseURL,
      headers: config.headers,
      data: config.data instanceof FormData ? 'FormData object' : config.data
    });
    
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add response interceptor for error handling
api.interceptors.response.use(
  (response) => {
    console.log('API Response:', {
      status: response.status,
      data: response.data
    });
    return response;
  },
  (error) => {
    console.error('API Error:', error);
    
    // Log more detailed error information
    if (error.response) {
      console.error('Response error:', {
        status: error.response.status,
        data: error.response.data,
        headers: error.response.headers
      });
      
      // Handle token expiration
      if (error.response.status === 401) {
        console.log('Unauthorized response detected, logging out user');
        // Logout the user
        authService.logout();
        
        // Redirect to login page
        if (window.location.pathname !== '/login') {
          // Store the current path to redirect back after login
          sessionStorage.setItem('redirectPath', window.location.pathname);
          window.location.href = '/login?expired=true';
        }
      }
    } else if (error.request) {
      console.error('Request error (no response received):', error.request);
    } else {
      console.error('Error setting up request:', error.message);
    }
    
    return Promise.reject(error);
  }
);

export default api;
export { API_URL }; 