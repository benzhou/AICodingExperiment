import React, { useState, useEffect } from 'react';
import { Card, Button, Alert, Space, Typography, Spin, Input, Form } from 'antd';
import { useUser } from '../context/UserContext';
import { userService } from '../services/userService';

const { Title, Text } = Typography;

const AdminStatusDebugger: React.FC = () => {
  const { user } = useUser();
  const [roles, setRoles] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [form] = Form.useForm();

  // Load user roles when component mounts
  useEffect(() => {
    if (user?.id) {
      setLoading(true);
      userService.getUserRoles(user.id)
        .then(fetchedRoles => {
          setRoles(fetchedRoles);
          console.log('User roles fetched:', fetchedRoles);
        })
        .catch(err => {
          console.error('Error fetching user roles:', err);
          setError('Failed to fetch user roles');
        })
        .finally(() => {
          setLoading(false);
        });
    }
  }, [user]);

  // Set admin status in localStorage
  const setAdminStatus = (isAdmin: boolean) => {
    if (isAdmin) {
      localStorage.setItem('isAdmin', 'true');
    } else {
      localStorage.removeItem('isAdmin');
    }
    window.location.reload();
  };

  // Add admin role to user
  const addAdminRole = async () => {
    if (!user?.id) return;
    
    try {
      setLoading(true);
      await userService.updateUserRole(user.id, { role: 'admin', operation: 'add' });
      setRoles([...roles, 'admin']);
      localStorage.setItem('isAdmin', 'true');
      window.location.reload();
    } catch (error) {
      console.error('Error adding admin role:', error);
      setError('Failed to add admin role. You might not have permission to do this.');
    } finally {
      setLoading(false);
    }
  };

  // Create admin user
  const createAdminUser = async () => {
    try {
      setLoading(true);
      const values = await form.validateFields();
      
      await userService.createUser({
        email: values.email,
        password: values.password,
        name: values.name,
        role: 'admin'
      });
      
      setError(null);
      form.resetFields();
    } catch (error) {
      console.error('Error creating admin user:', error);
      setError('Failed to create admin user');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Card title="Admin Status Debugger">
        <div style={{ textAlign: 'center', padding: '20px' }}>
          <Spin size="large" />
          <p>Loading...</p>
        </div>
      </Card>
    );
  }

  const isAdmin = roles.includes('admin') || localStorage.getItem('isAdmin') === 'true';

  return (
    <Card title="Admin Status Debugger">
      <Alert
        message={isAdmin ? "You have admin privileges" : "You don't have admin privileges"}
        type={isAdmin ? "success" : "warning"}
        showIcon
        style={{ marginBottom: '20px' }}
      />

      <div style={{ marginBottom: '20px' }}>
        <Title level={5}>User Information</Title>
        <Text>ID: {user?.id || 'Not logged in'}</Text><br />
        <Text>Email: {user?.email || 'Not logged in'}</Text><br />
        <Text>Roles: {roles.length > 0 ? roles.join(', ') : 'None'}</Text><br />
        <Text>localStorage.isAdmin: {localStorage.getItem('isAdmin') || 'Not set'}</Text>
      </div>

      <Space direction="vertical" style={{ width: '100%', marginBottom: '20px' }}>
        <Title level={5}>Quick Actions</Title>
        <Space>
          <Button 
            type="primary" 
            onClick={() => setAdminStatus(true)}
            disabled={isAdmin}
          >
            Make Admin (localStorage)
          </Button>
          <Button 
            danger 
            onClick={() => setAdminStatus(false)}
            disabled={!isAdmin}
          >
            Remove Admin (localStorage)
          </Button>
          <Button 
            type="primary" 
            onClick={addAdminRole}
            disabled={roles.includes('admin')}
          >
            Add Admin Role (API)
          </Button>
        </Space>
      </Space>

      <div style={{ marginTop: '30px' }}>
        <Title level={5}>Create Admin User</Title>
        <Form
          form={form}
          layout="vertical"
          onFinish={createAdminUser}
        >
          <Form.Item
            name="email"
            label="Email"
            rules={[
              { required: true, message: 'Please enter email' },
              { type: 'email', message: 'Please enter a valid email' }
            ]}
          >
            <Input placeholder="admin@example.com" />
          </Form.Item>
          
          <Form.Item
            name="name"
            label="Name"
            rules={[{ required: true, message: 'Please enter name' }]}
          >
            <Input placeholder="Admin User" />
          </Form.Item>
          
          <Form.Item
            name="password"
            label="Password"
            rules={[{ required: true, message: 'Please enter password' }]}
          >
            <Input.Password placeholder="Password" />
          </Form.Item>
          
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              Create Admin User
            </Button>
          </Form.Item>
        </Form>
      </div>

      {error && (
        <Alert
          message="Error"
          description={error}
          type="error"
          showIcon
          style={{ marginTop: '20px' }}
        />
      )}
    </Card>
  );
};

export default AdminStatusDebugger; 