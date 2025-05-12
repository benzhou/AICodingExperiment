import React, { useState, useEffect } from 'react';
import { Card, Form, Select, Button, Radio, message, Typography, Upload, Divider, Tabs, Alert } from 'antd';
import { UploadOutlined, UserOutlined, AppstoreOutlined, SettingOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useUser } from '../context/UserContext';
import { useTenant } from '../context/TenantContext';
import { Logo } from './index';
import UsersManagement from './UsersManagement';
import { userService } from '../services/userService';
import AdminStatusFix from './AdminStatusFix';
import { Link } from 'react-router-dom';

const { Title, Text } = Typography;
const { Option } = Select;
const { TabPane } = Tabs;

const UserProfile: React.FC = () => {
  const { t, i18n } = useTranslation();
  const { user } = useUser();
  const { tenant, uploadLogo } = useTenant();
  const [form] = Form.useForm();
  const [uploading, setUploading] = useState(false);
  const [activeTab, setActiveTab] = useState('1');
  const [roles, setRoles] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showAdminFix, setShowAdminFix] = useState(false);
  
  // Check if user is an admin - using user email for now
  // In a real application, you would check user permissions from the backend
  const isAdmin = user?.email?.endsWith('@admin.com') || localStorage.getItem('isAdmin') === 'true' || (roles && roles.includes('admin'));

  // Initialize form with current settings
  useEffect(() => {
    form.setFieldsValue({
      language: i18n.language,
      theme: localStorage.getItem('theme') || 'light'
    });

    // Get user roles from API
    if (user?.id) {
      setLoading(true);
      userService.getUserRoles(user.id)
        .then(fetchedRoles => {
          if (fetchedRoles) {
            setRoles(fetchedRoles);
            console.log('User roles fetched:', fetchedRoles);
            
            // If user has admin role, set it in localStorage for redundancy
            if (fetchedRoles.includes('admin')) {
              localStorage.setItem('isAdmin', 'true');
              
              // Check if we need to show the admin fix tool
              if (!isAdmin) {
                setShowAdminFix(true);
              }
            }
          } else {
            console.log('No roles returned from API');
            setRoles([]);
          }
        })
        .catch(err => {
          console.error('Error fetching user roles:', err);
          setError('Failed to fetch user roles');
        })
        .finally(() => {
          setLoading(false);
        });
    }
  }, [form, i18n.language, user, isAdmin]);

  // Debug function to make user admin
  const makeUserAdmin = () => {
    localStorage.setItem('isAdmin', 'true');
    window.location.reload();
  };

  const handleSubmit = (values: any) => {
    // Save language preference
    i18n.changeLanguage(values.language);
    localStorage.setItem('i18nextLng', values.language);
    
    // Save theme preference
    localStorage.setItem('theme', values.theme);
    
    message.success(t('userProfile.successUpdate'));
  };

  const handleLogoUpload = async (file: any) => {
    try {
      setUploading(true);
      await uploadLogo(file);
      message.success('Logo uploaded successfully');
    } catch (error) {
      console.error('Error uploading logo:', error);
      message.error('Failed to upload logo');
    } finally {
      setUploading(false);
    }
    return false; // Prevent default upload behavior
  };

  const renderUserProfile = () => (
    <>
      <Card className="mb-6">
        <div className="flex items-center mb-4">
          <div className="text-xl font-medium">{user?.name}</div>
          <div className="ml-2 text-gray-500">{user?.email}</div>
        </div>
        
        {/* Debug information */}
        <div style={{ marginTop: '10px', padding: '10px', background: '#f0f0f0', borderRadius: '4px' }}>
          <h4>Debug Information</h4>
          <p>User ID: {user?.id || 'Not found'}</p>
          <p>Email: {user?.email || 'Not found'}</p>
          <p>Admin status: {isAdmin ? 'Yes' : 'No'}</p>
          <p>Roles: {roles && roles.length > 0 ? roles.join(', ') : 'None'}</p>
          <p>localStorage.isAdmin: {localStorage.getItem('isAdmin') || 'Not set'}</p>
          <Button onClick={makeUserAdmin} type="primary" danger style={{ marginRight: '10px' }}>Make Admin (Debug)</Button>
          <Link to="/admin-troubleshooter">
            <Button type="primary">Admin Troubleshooter</Button>
          </Link>
        </div>
        
        {/* Show Admin Status Fix tool if needed */}
        {roles && roles.includes('admin') && !isAdmin && (
          <div style={{ marginTop: '20px' }}>
            <AdminStatusFix />
          </div>
        )}
      </Card>
      
      <Card>
        <Form 
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
        >
          <Form.Item
            name="language"
            label={t('userProfile.language')}
            help={t('userProfile.languageDescription')}
          >
            <Select style={{ width: 200 }}>
              <Option value="en">English</Option>
              <Option value="zh-CN">简体中文</Option>
            </Select>
          </Form.Item>
          
          <Form.Item
            name="theme"
            label={t('userProfile.theme')}
            help={t('userProfile.themeDescription')}
          >
            <Radio.Group>
              <Radio value="light">{t('userProfile.light')}</Radio>
              <Radio value="dark">{t('userProfile.dark')}</Radio>
            </Radio.Group>
          </Form.Item>
          
          <Form.Item>
            <Button type="primary" htmlType="submit">
              {t('userProfile.save')}
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </>
  );

  const renderTenantBranding = () => (
    <Card title="Tenant Branding">
      <div className="mb-4">
        <Title level={4}>Current Logo</Title>
        <div className="border p-4 rounded flex justify-center items-center bg-gray-50" style={{ height: '100px' }}>
          <Logo height={60} width={200} />
        </div>
      </div>
      
      <Divider />
      
      <div>
        <Title level={4}>Upload New Logo</Title>
        <Text type="secondary" className="block mb-4">
          Upload a new logo for your tenant. Recommended size is 200x60 pixels in SVG or PNG format.
        </Text>
        
        <Upload
          beforeUpload={handleLogoUpload}
          showUploadList={false}
          accept=".svg,.png,.jpg,.jpeg"
        >
          <Button 
            icon={<UploadOutlined />} 
            loading={uploading}
          >
            Select Logo File
          </Button>
        </Upload>
      </div>
    </Card>
  );

  // Show loading or error messages
  if (loading) {
    return (
      <div className="max-w-5xl mx-auto">
        <Title level={2}>{t('userProfile.title')}</Title>
        <Card>
          <div style={{ textAlign: 'center', padding: '20px' }}>Loading user information...</div>
        </Card>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-5xl mx-auto">
        <Title level={2}>{t('userProfile.title')}</Title>
        <Alert message={error} type="error" showIcon />
        {renderUserProfile()}
      </div>
    );
  }

  // If we have a role mismatch, show the admin fix prominently
  if (showAdminFix) {
    return (
      <div className="max-w-5xl mx-auto">
        <Title level={2}>{t('userProfile.title')}</Title>
        <Alert
          message="Admin Role Configuration Issue"
          description="You have admin role in the database but it's not properly configured in the application."
          type="warning"
          showIcon
          style={{ marginBottom: '20px' }}
        />
        <AdminStatusFix />
        {renderUserProfile()}
      </div>
    );
  }

  return (
    <div className="max-w-5xl mx-auto">
      <Title level={2}>{t('userProfile.title')}</Title>
      
      {/* Always show debug info in the header for testing */}
      <Alert
        message="Admin Status Debug"
        description={`isAdmin: ${isAdmin}, email: ${user?.email}, roles: ${roles ? roles.join(', ') : 'None'}`}
        type="info"
        showIcon
        style={{ marginBottom: '20px' }}
      />
      
      {isAdmin ? (
        <Tabs 
          defaultActiveKey="1" 
          activeKey={activeTab}
          onChange={setActiveTab}
          tabPosition="left"
          style={{ minHeight: '500px' }}
        >
          <TabPane 
            tab={
              <span>
                <UserOutlined />
                {t('userProfile.tabs.profile', 'Profile')}
              </span>
            } 
            key="1"
          >
            {renderUserProfile()}
          </TabPane>
          
          <TabPane 
            tab={
              <span>
                <AppstoreOutlined />
                {t('userProfile.tabs.branding', 'Branding')}
              </span>
            } 
            key="2"
          >
            {renderTenantBranding()}
          </TabPane>
          
          <TabPane 
            tab={
              <span>
                <SettingOutlined />
                {t('userProfile.tabs.userManagement', 'User Management')}
              </span>
            } 
            key="3"
          >
            <UsersManagement />
          </TabPane>
        </Tabs>
      ) : (
        renderUserProfile()
      )}
    </div>
  );
};

export default UserProfile; 