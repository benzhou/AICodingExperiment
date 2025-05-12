import React, { useState, useEffect, useCallback } from 'react';
import { Card, Button, Alert, Space, Typography, Divider, message, List } from 'antd';
import { useUser } from '../context/UserContext';
import { userService } from '../services/userService';
import { CheckCircleOutlined, CloseCircleOutlined, SyncOutlined } from '@ant-design/icons';

const { Title, Text, Paragraph } = Typography;

interface DiagnosticTest {
  name: string;
  description: string;
  status: 'pending' | 'success' | 'error' | 'running';
  result: string | null;
}

const AdminTroubleshooter: React.FC = () => {
  const { user, refreshUser } = useUser();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [databaseRoles, setDatabaseRoles] = useState<string[]>([]);
  const [diagnostics, setDiagnostics] = useState<DiagnosticTest[]>([
    {
      name: 'User Authentication',
      description: 'Verify if user is logged in and has valid authentication',
      status: 'pending',
      result: null
    },
    {
      name: 'Database Role Check',
      description: 'Check if user has admin role in the database',
      status: 'pending',
      result: null
    },
    {
      name: 'Local Storage Check',
      description: 'Check if admin flag is set in local storage',
      status: 'pending',
      result: null
    },
    {
      name: 'Admin UI Access',
      description: 'Verify if admin UI components are accessible',
      status: 'pending',
      result: null
    }
  ]);

  // Define helper functions first with useCallback
  const updateDiagnostic = useCallback((name: string, status: DiagnosticTest['status'], result: string | null = null) => {
    setDiagnostics(prev => prev.map(test => 
      test.name === name ? { ...test, status, result } : test
    ));
  }, []);

  const checkUserAuth = useCallback(async () => {
    updateDiagnostic('User Authentication', 'running');
    
    if (!user || !user.id) {
      updateDiagnostic('User Authentication', 'error', 'User is not logged in');
      throw new Error('User not logged in');
    }
    
    updateDiagnostic('User Authentication', 'success', `Logged in as ${user.email} (ID: ${user.id})`);
  }, [user, updateDiagnostic]);
  
  const checkDatabaseRoles = useCallback(async () => {
    updateDiagnostic('Database Role Check', 'running');
    
    try {
      if (!user?.id) {
        updateDiagnostic('Database Role Check', 'error', 'Cannot check roles - user not logged in');
        return;
      }
      
      const roles = await userService.getUserRoles(user.id);
      setDatabaseRoles(roles);
      
      if (roles.includes('admin')) {
        updateDiagnostic('Database Role Check', 'success', `User has admin role in database (Roles: ${roles.join(', ')})`);
      } else {
        updateDiagnostic('Database Role Check', 'error', `User does not have admin role in database (Roles: ${roles.length ? roles.join(', ') : 'none'})`);
      }
    } catch (err) {
      console.error('Error checking roles:', err);
      updateDiagnostic('Database Role Check', 'error', 'Failed to retrieve roles from database');
    }
  }, [user, updateDiagnostic]);
  
  const checkLocalStorage = useCallback(() => {
    updateDiagnostic('Local Storage Check', 'running');
    
    const isAdmin = localStorage.getItem('isAdmin');
    
    if (isAdmin === 'true') {
      updateDiagnostic('Local Storage Check', 'success', 'Admin flag is set in localStorage');
    } else {
      updateDiagnostic('Local Storage Check', 'error', 'Admin flag is not set in localStorage');
    }
  }, [updateDiagnostic]);
  
  const checkAdminUIAccess = useCallback(() => {
    updateDiagnostic('Admin UI Access', 'running');
    
    const isAdmin = localStorage.getItem('isAdmin') === 'true' || 
                   (user?.email?.endsWith('@admin.com')) || 
                   databaseRoles.includes('admin');
                   
    if (isAdmin) {
      updateDiagnostic('Admin UI Access', 'success', 'Admin UI should be accessible');
    } else {
      updateDiagnostic('Admin UI Access', 'error', 'Admin UI is not accessible');
    }
  }, [databaseRoles, user, updateDiagnostic]);

  // Define runDiagnostics after the helper functions
  const runDiagnostics = useCallback(async () => {
    setLoading(true);
    setError(null);
    
    // Reset diagnostics
    setDiagnostics(prev => prev.map(test => ({
      ...test,
      status: 'pending',
      result: null
    })));
    
    try {
      // 1. Check user authentication
      await checkUserAuth();
      
      // 2. Check database roles (only if user is authenticated)
      if (user?.id) {
        await checkDatabaseRoles();
      }
      
      // 3. Check localStorage
      checkLocalStorage();
      
      // 4. Check admin UI access
      checkAdminUIAccess();
      
    } catch (err) {
      console.error('Error running diagnostics:', err);
      setError('Failed to complete diagnostics. Please try again.');
    } finally {
      setLoading(false);
    }
  }, [user, checkUserAuth, checkDatabaseRoles, checkLocalStorage, checkAdminUIAccess]);
  
  // Run diagnostics when component mounts
  useEffect(() => {
    runDiagnostics();
  }, [runDiagnostics]);

  const fixAdminStatus = async () => {
    setLoading(true);
    
    try {
      // Only proceed if we have database admin role
      if (!databaseRoles.includes('admin')) {
        message.error('Cannot fix - you do not have admin role in the database');
        return;
      }
      
      // Set admin flag in localStorage
      localStorage.setItem('isAdmin', 'true');
      message.success('Admin status has been fixed!');
      
      // Refresh user data
      if (refreshUser) {
        await refreshUser();
      }
      
      // Run diagnostics again
      await runDiagnostics();
      
      // Let user know to try again
      message.info('Please go back to your profile page and reload to see the admin interface');
      
    } catch (err) {
      console.error('Error fixing admin status:', err);
      message.error('Failed to fix admin status');
    } finally {
      setLoading(false);
    }
  };

  const getIconForStatus = (status: DiagnosticTest['status']) => {
    switch (status) {
      case 'success': return <CheckCircleOutlined style={{ color: 'green' }} />;
      case 'error': return <CloseCircleOutlined style={{ color: 'red' }} />;
      case 'running': return <SyncOutlined spin />;
      default: return <SyncOutlined style={{ color: 'gray' }} />;
    }
  };

  return (
    <div className="max-w-5xl mx-auto">
      <Title level={2}>Admin Access Troubleshooter</Title>
      
      <Alert
        message="This tool helps diagnose and fix admin access issues"
        description="If you believe you should have admin access but can't see the admin interface, use this tool to diagnose and fix the problem."
        type="info"
        showIcon
        style={{ marginBottom: '20px' }}
      />
      
      <Card title="Admin Access Diagnostics" loading={loading}>
        <List
          itemLayout="horizontal"
          dataSource={diagnostics}
          renderItem={item => (
            <List.Item>
              <List.Item.Meta
                avatar={getIconForStatus(item.status)}
                title={<Text strong>{item.name}</Text>}
                description={
                  <>
                    <Text type="secondary">{item.description}</Text>
                    {item.result && (
                      <div style={{ marginTop: '5px' }}>
                        <Text mark>{item.result}</Text>
                      </div>
                    )}
                  </>
                }
              />
            </List.Item>
          )}
        />
        
        <Divider />
        
        <Space direction="vertical" style={{ width: '100%' }}>
          <Space>
            <Button 
              type="primary" 
              onClick={runDiagnostics} 
              loading={loading}
              icon={<SyncOutlined />}
            >
              Run Diagnostics Again
            </Button>
            
            <Button 
              type="primary" 
              onClick={fixAdminStatus}
              loading={loading}
              disabled={!databaseRoles.includes('admin')}
              danger
            >
              Fix Admin Access
            </Button>
          </Space>
          
          {error && (
            <Alert message={error} type="error" showIcon style={{ marginTop: '10px' }} />
          )}
          
          {databaseRoles.includes('admin') && localStorage.getItem('isAdmin') !== 'true' && (
            <Alert
              message="Action Required"
              description="You have admin role in the database but not in the application. Click 'Fix Admin Access' above to resolve this issue."
              type="warning"
              showIcon
              style={{ marginTop: '10px' }}
            />
          )}
        </Space>
      </Card>
      
      <Card title="Manual Fix Instructions" style={{ marginTop: '20px' }}>
        <Paragraph>
          If the automatic fix doesn't work, you can try these manual steps:
        </Paragraph>
        
        <ol>
          <li>
            <Paragraph>
              <Text strong>Check your database role:</Text> Confirm your user ID has the admin role in the user_roles table
            </Paragraph>
          </li>
          <li>
            <Paragraph>
              <Text strong>Set localStorage manually:</Text> Open browser developer tools (F12), go to Application tab, 
              find localStorage, and manually add a key "isAdmin" with value "true"
            </Paragraph>
          </li>
          <li>
            <Paragraph>
              <Text strong>Clear browser cache:</Text> Sometimes clearing browser cache can help resolve UI issues
            </Paragraph>
          </li>
          <li>
            <Paragraph>
              <Text strong>Restart the application:</Text> Log out and log back in to refresh your session
            </Paragraph>
          </li>
        </ol>
      </Card>
    </div>
  );
};

export default AdminTroubleshooter; 