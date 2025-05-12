import React from 'react';
import { Typography, Card, Alert } from 'antd';
import { useTranslation } from 'react-i18next';
import { SettingOutlined } from '@ant-design/icons';

const { Title, Text } = Typography;

const MatchsetConfiguration: React.FC = () => {
  const { t } = useTranslation();

  return (
    <div>
      <Title level={2}>{t('matchset.title')}</Title>
      <Text type="secondary" className="block mb-6">{t('matchset.description')}</Text>
      
      <Card>
        <Alert
          message="Feature Under Development"
          description="The matchset configuration feature is currently being developed."
          type="info"
          showIcon
          icon={<SettingOutlined />}
        />
      </Card>
    </div>
  );
};

export default MatchsetConfiguration; 