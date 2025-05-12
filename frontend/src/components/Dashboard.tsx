import React, { useState } from 'react';
import { useUser } from '../context/UserContext';
import { authService } from '../services/authService';
import { Link } from 'react-router-dom';
import { TokenModal } from './TokenModal';
import { useTranslation } from 'react-i18next';
import { Card, Row, Col, Statistic, Typography } from 'antd';
import { 
  DatabaseOutlined, 
  SyncOutlined, 
  CheckCircleOutlined, 
  CloseCircleOutlined,
  PercentageOutlined
} from '@ant-design/icons';

const { Title, Text } = Typography;

export function Dashboard() {
  const { user } = useUser();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [tokenInfo, setTokenInfo] = useState<{
    token: string | null;
    decodedToken: any | null;
    expiresIn: number | null;
    isValid: boolean;
  }>({
    token: null,
    decodedToken: null,
    expiresIn: null,
    isValid: false,
  });
  const { t } = useTranslation();

  const handleTokenDebug = async () => {
    const token = authService.getToken();
    if (token) {
      const info = await authService.getTokenInfo();
      const [header, payload] = token.split('.').slice(0, 2);
      const decodedToken = JSON.parse(atob(payload));
      setTokenInfo({
        token: token,
        decodedToken: decodedToken,
        expiresIn: info.expires_in !== undefined ? info.expires_in : null,
        isValid: info.expires_in > 0,
      });
      setIsModalOpen(true);
    }
  };

  // Placeholder statistics data - would come from backend API in a real application
  const stats = {
    dataSources: 2,
    transactions: 150,
    matched: 120,
    unmatched: 30,
    matchRate: 80
  };

  return (
    <div>
      <Title level={2}>{t('dashboard.welcome')}</Title>
      <Text type="secondary" className="block mb-6">{t('dashboard.description')}</Text>
      
      {/* Statistics Cards */}
      <Row gutter={[16, 16]} className="mb-6">
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic 
              title={t('dashboard.stats.dataSources')}
              value={stats.dataSources}
              prefix={<DatabaseOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic 
              title={t('dashboard.stats.transactions')}
              value={stats.transactions}
              prefix={<SyncOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic 
              title={t('dashboard.stats.matched')}
              value={stats.matched}
              valueStyle={{ color: '#3f8600' }}
              prefix={<CheckCircleOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic 
              title={t('dashboard.stats.unmatched')}
              value={stats.unmatched}
              valueStyle={{ color: '#cf1322' }}
              prefix={<CloseCircleOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic 
              title={t('dashboard.stats.matchRate')}
              value={stats.matchRate}
              suffix="%"
              prefix={<PercentageOutlined />}
            />
          </Card>
        </Col>
      </Row>
      
      {/* Action Cards */}
      <Row gutter={[16, 16]}>
        <Col xs={24} md={12}>
          <Card className="h-full">
            <Title level={4}>{t('dashboard.manageDataSources.title')}</Title>
            <Text className="block mb-4">{t('dashboard.manageDataSources.description')}</Text>
            <Link 
              to="/datasources" 
              className="inline-block bg-blue-600 text-white px-4 py-2 rounded text-sm hover:bg-blue-700"
            >
              {t('dashboard.manageDataSources.action')}
            </Link>
          </Card>
        </Col>
        
        <Col xs={24} md={12}>
          <Card className="h-full">
            <Title level={4}>{t('dashboard.uploadData.title')}</Title>
            <Text className="block mb-4">{t('dashboard.uploadData.description')}</Text>
            <Link 
              to="/datasources/upload" 
              className="inline-block bg-green-600 text-white px-4 py-2 rounded text-sm hover:bg-green-700"
            >
              {t('dashboard.uploadData.action')}
            </Link>
          </Card>
        </Col>
      </Row>
      
      <TokenModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        tokenInfo={tokenInfo}
      />
    </div>
  );
} 