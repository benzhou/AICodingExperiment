import React, { createContext, useState, useContext, useEffect, ReactNode } from 'react';
import { Tenant, tenantService } from '../services/tenantService';

interface TenantContextType {
  tenant: Tenant | null;
  isLoading: boolean;
  error: string | null;
  updateTenant: (data: Partial<Tenant>) => Promise<void>;
  uploadLogo: (file: File) => Promise<void>;
}

const TenantContext = createContext<TenantContextType>({
  tenant: null,
  isLoading: false,
  error: null,
  updateTenant: async () => {},
  uploadLogo: async () => {},
});

export const useTenant = () => useContext(TenantContext);

interface TenantProviderProps {
  children: ReactNode;
}

export const TenantProvider: React.FC<TenantProviderProps> = ({ children }) => {
  const [tenant, setTenant] = useState<Tenant | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  // Fetch current tenant on component mount
  useEffect(() => {
    const fetchTenant = async () => {
      try {
        setIsLoading(true);
        setError(null);
        
        // Try to get tenant from domain
        const hostname = window.location.hostname;
        let currentTenant = null;
        
        if (hostname !== 'localhost' && hostname !== '127.0.0.1') {
          currentTenant = await tenantService.getTenantByDomain(hostname);
        }
        
        // If no tenant found by domain, get from JWT token
        if (!currentTenant) {
          currentTenant = await tenantService.getCurrentTenant();
        }
        
        setTenant(currentTenant);
      } catch (err) {
        setError('Failed to load tenant information');
        console.error('Error loading tenant:', err);
      } finally {
        setIsLoading(false);
      }
    };
    
    fetchTenant();
  }, []);

  // Update tenant information
  const updateTenant = async (data: Partial<Tenant>) => {
    if (!tenant) return;
    
    try {
      setIsLoading(true);
      setError(null);
      
      const updatedTenant = await tenantService.updateTenant(tenant.id, data);
      
      if (updatedTenant) {
        setTenant(updatedTenant);
      }
    } catch (err) {
      setError('Failed to update tenant information');
      console.error('Error updating tenant:', err);
    } finally {
      setIsLoading(false);
    }
  };

  // Upload tenant logo
  const uploadLogo = async (file: File) => {
    if (!tenant) return;
    
    try {
      setIsLoading(true);
      setError(null);
      
      const logoUrl = await tenantService.uploadTenantLogo(tenant.id, file);
      
      if (logoUrl) {
        setTenant(prev => prev ? { ...prev, logoUrl } : null);
      }
    } catch (err) {
      setError('Failed to upload logo');
      console.error('Error uploading logo:', err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <TenantContext.Provider value={{ tenant, isLoading, error, updateTenant, uploadLogo }}>
      {children}
    </TenantContext.Provider>
  );
}; 