import React, { useState } from 'react';
import { Card, Button, Alert, message, Typography } from 'antd';
import { userService } from '../services/userService';
import { useUser } from '../context/UserContext';

const { Title, Text } = Typography;

const AdminStatusFix: React.FC = () => {
  const { user, refreshUser } = useUser();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // Fix admin token issue
  const fixAdminStatus = async () => {
    if (!user?.id) {
      setError('User not logged in');
      return;
    }
    
    try {
      setLoading(true);
      setError(null);
      
      // 1. Verify in database if user has admin role
      const roles = await userService.getUserRoles(user.id);
      
      // 2. If user has admin role in DB but not in localStorage, fix it
      if (roles && roles.includes('admin')) {
        localStorage.setItem('isAdmin', 'true');
        message.success('Admin status fixed! Reloading the page...');
        
        // Refresh user context data
        if (refreshUser) {
          await refreshUser();
        }
        
        // Reload the page to apply changes
        setTimeout(() => {
          window.location.reload();
        }, 1500);
      } else {
        setError('You do not have admin role in the database');
      }
    } catch (err) {
      console.error('Error fixing admin status:', err);
      setError('Failed to fix admin status. Please try again.');
    } finally {
      setLoading(false);
    }
  };
  
  return (
    <Card title="Admin Status Fix Tool">
      <Alert
        message="Admin Status Issue Detected"
        description="Your admin role is set in the database but not recognized by the application. Use this tool to fix the issue."
        type="warning"
        showIcon
        style={{ marginBottom: '20px' }}
      />
      
      <div style={{ marginBottom: '20px' }}>
        <Title level={5}>User Information</Title>
        <Text>ID: {user?.id || 'Not logged in'}</Text><br />
        <Text>Email: {user?.email || 'Not logged in'}</Text><br />
        <Text>localStorage.isAdmin: {localStorage.getItem('isAdmin') || 'Not set'}</Text>
      </div>
      
      <Button 
        type="primary" 
        onClick={fixAdminStatus}
        loading={loading}
        style={{ marginBottom: '20px' }}
      >
        Fix Admin Status
      </Button>
      
      {error && (
        <Alert
          message="Error"
          description={error}
          type="error"
          showIcon
        />
      )}
    </Card>
  );
};

export default AdminStatusFix; 