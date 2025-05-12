import React, { useState } from 'react';
import { Layout, Menu, Dropdown, Button, Select, Avatar, Space, Tooltip } from 'antd';
import { useNavigate, useLocation, Link, Outlet } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import {
  HomeOutlined,
  DatabaseOutlined,
  SyncOutlined,
  SettingOutlined,
  UserOutlined,
  UploadOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  LogoutOutlined,
  DoubleLeftOutlined,
  DoubleRightOutlined,
  BulbOutlined,
} from '@ant-design/icons';
import { useUser } from '../context/UserContext';
import { authService } from '../services/authService';
import { Logo } from './index';
import { useTheme } from '../context/ThemeContext';

const { Header, Sider, Content } = Layout;
const { SubMenu } = Menu;
const { Option } = Select;

interface MainLayoutProps {
  children?: React.ReactNode;
}

const MainLayout: React.FC<MainLayoutProps> = ({ children }) => {
  const [collapsed, setCollapsed] = useState(false);
  const { t, i18n } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();
  const { user, setUser } = useUser();
  const { theme, toggleTheme } = useTheme();

  const handleLanguageChange = (lang: string) => {
    i18n.changeLanguage(lang);
    localStorage.setItem('i18nextLng', lang);
  };

  const handleLogout = () => {
    authService.logout();
    setUser(null);
    navigate('/login');
  };

  // Function to determine which menu keys should be open by default
  const getOpenKeys = () => {
    const path = location.pathname;
    if (path.includes('/datasources')) return ['dataSources'];
    if (path.includes('/transaction-match')) return ['transactionMatch'];
    return [];
  };

  // Function to determine which menu key should be selected
  const getSelectedKey = () => {
    const path = location.pathname;
    
    if (path === '/dashboard') return ['dashboard'];
    if (path === '/datasources') return ['dataSourceDefinition'];
    if (path === '/datasources/upload') return ['uploadDataSources'];
    if (path.match(/^\/datasources\/.*\/imports$/)) return ['dataSourceImports'];
    if (path.match(/^\/imports\/.*\/transactions$/)) return ['dataSourceImports'];
    if (path === '/transaction-match/matchset') return ['matchset'];
    if (path === '/transaction-match/matched') return ['matchedTransactions'];
    if (path === '/transaction-match/unmatched') return ['unmatchedTransactions'];
    
    return ['dashboard']; // Default
  };

  const userMenu = (
    <Menu>
      <Menu.Item key="profile" icon={<UserOutlined />}>
        <Link to="/profile">{t('userProfile.title')}</Link>
      </Menu.Item>
      <Menu.Item key="adminTrouble" icon={<SettingOutlined />}>
        <Link to="/admin-troubleshooter">Admin Troubleshooter</Link>
      </Menu.Item>
      <Menu.Divider />
      <Menu.Item key="logout" icon={<LogoutOutlined />} onClick={handleLogout}>
        {t('auth.labels.logout')}
      </Menu.Item>
    </Menu>
  );

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider 
        trigger={null} 
        collapsible 
        collapsed={collapsed}
        width={250}
        style={{
          overflow: 'auto',
          height: '100vh',
          position: 'fixed',
          left: 0,
          top: 0,
          bottom: 0,
          display: 'flex',
          flexDirection: 'column',
        }}
      >
        <div className="logo" style={{ 
          height: '64px', 
          margin: '16px', 
          display: 'flex', 
          alignItems: 'center',
          justifyContent: collapsed ? 'center' : 'flex-start',
          overflow: 'hidden',
        }}>
          {collapsed ? (
            <Logo height={40} width={40} />
          ) : (
            <Logo height={40} width={160} />
          )}
        </div>
        
        <Menu 
          mode="inline" 
          defaultOpenKeys={getOpenKeys()}
          selectedKeys={getSelectedKey()}
          style={{ flex: 1, paddingBottom: '50px' }}
        >
          <Menu.Item key="dashboard" icon={<HomeOutlined />} onClick={() => navigate('/dashboard')}>
            {t('navigation.dashboard')}
          </Menu.Item>
          
          <SubMenu 
            key="dataSources" 
            icon={<DatabaseOutlined />} 
            title={t('navigation.dataSources')}
          >
            <Menu.Item 
              key="dataSourceDefinition" 
              onClick={() => navigate('/datasources')}
            >
              {t('navigation.dataSourceDefinition')}
            </Menu.Item>
            <Menu.Item 
              key="uploadDataSources" 
              icon={<UploadOutlined />}
              onClick={() => navigate('/datasources/upload')}
            >
              {t('navigation.uploadDataSources')}
            </Menu.Item>
          </SubMenu>
          
          <SubMenu 
            key="transactionMatch" 
            icon={<SyncOutlined />} 
            title={t('navigation.transactionMatch')}
          >
            <Menu.Item 
              key="matchset" 
              icon={<SettingOutlined />}
              onClick={() => navigate('/transaction-match/matchset')}
            >
              {t('navigation.matchset')}
            </Menu.Item>
            <Menu.Item 
              key="matchedTransactions" 
              icon={<CheckCircleOutlined />}
              onClick={() => navigate('/transaction-match/matched')}
            >
              {t('navigation.matchedTransactions')}
            </Menu.Item>
            <Menu.Item 
              key="unmatchedTransactions" 
              icon={<CloseCircleOutlined />}
              onClick={() => navigate('/transaction-match/unmatched')}
            >
              {t('navigation.unmatchedTransactions')}
            </Menu.Item>
          </SubMenu>
        </Menu>
        
        <div style={{ 
          padding: '16px', 
          textAlign: 'center', 
          borderTop: '1px solid var(--border-color)',
          position: 'absolute',
          bottom: 0,
          width: '100%',
          background: 'var(--sidebar-bg)',
        }}>
          <Button
            type="text"
            icon={collapsed ? <DoubleRightOutlined /> : <DoubleLeftOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            style={{ fontSize: '16px' }}
          />
        </div>
      </Sider>
      
      <Layout style={{ marginLeft: collapsed ? 80 : 250, transition: 'all 0.2s' }}>
        <Header style={{ 
          padding: '0 16px', 
          display: 'flex', 
          alignItems: 'center',
          justifyContent: 'flex-end',
          boxShadow: '0 1px 4px var(--shadow-color)',
          position: 'sticky',
          top: 0,
          zIndex: 1,
        }}>
          <Space size="large" align="center">
            <Select 
              defaultValue={i18n.language} 
              style={{ width: 120 }} 
              onChange={handleLanguageChange}
              dropdownMatchSelectWidth={false}
            >
              <Option value="en">English</Option>
              <Option value="zh-CN">简体中文</Option>
            </Select>
            
            <Tooltip title={theme === 'light' ? t('navigation.darkMode') : t('navigation.lightMode')}>
              <Button
                type="text"
                icon={<BulbOutlined />}
                onClick={toggleTheme}
                aria-label={theme === 'light' ? "Switch to dark mode" : "Switch to light mode"}
              />
            </Tooltip>
            
            <Dropdown overlay={userMenu} trigger={['click']}>
              <div style={{ 
                display: 'flex', 
                alignItems: 'center', 
                cursor: 'pointer',
                color: 'var(--text-primary)'
              }}>
                <Avatar icon={<UserOutlined />} style={{ marginRight: 8 }} />
                <span style={{ 
                  maxWidth: '120px', 
                  overflow: 'hidden', 
                  textOverflow: 'ellipsis', 
                  whiteSpace: 'nowrap',
                }}>
                  {user?.name}
                </span>
              </div>
            </Dropdown>
          </Space>
        </Header>
        
        <Content style={{ margin: '24px 16px', padding: 24, minHeight: 280, overflow: 'initial' }}>
          {children || <Outlet />}
        </Content>
      </Layout>
    </Layout>
  );
};

export default MainLayout; 