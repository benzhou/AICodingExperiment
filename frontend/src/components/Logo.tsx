import React from 'react';
import { useTenant } from '../context/TenantContext';

interface LogoProps {
  height?: number;
  width?: number;
  className?: string;
}

const Logo: React.FC<LogoProps> = ({ height = 40, width = 120, className = '' }) => {
  const { tenant } = useTenant();
  
  // If tenant has a custom logo, use it
  if (tenant?.logoUrl) {
    return (
      <img 
        src={tenant.logoUrl} 
        alt={`${tenant.name} Logo`}
        height={height}
        width={width}
        className={className}
      />
    );
  }
  
  // Otherwise use the default Matcha logo
  return (
    <img 
      src="/logos/matcha.svg" 
      alt="Matcha Logo"
      height={height}
      width={width}
      className={className}
    />
  );
};

export default Logo; 