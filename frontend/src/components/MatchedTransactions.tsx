import React from 'react';
import { Typography, Card, Alert } from 'antd';
import { useTranslation } from 'react-i18next';
import { CheckCircleOutlined } from '@ant-design/icons';

const { Title, Text } = Typography;

const MatchedTransactions: React.FC = () => {
  const { t } = useTranslation();

  return (
    <div>
      <Title level={2}>{t('transactions.matched.title')}</Title>
      <Text type="secondary" className="block mb-6">{t('transactions.matched.description')}</Text>
      
      <Card>
        <Alert
          message="Feature Under Development"
          description="The matched transactions view is currently being developed."
          type="info"
          showIcon
          icon={<CheckCircleOutlined />}
        />
      </Card>
    </div>
  );
};

export default MatchedTransactions; 