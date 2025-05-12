import React from 'react';
import { Typography, Card, Alert } from 'antd';
import { useTranslation } from 'react-i18next';
import { CloseCircleOutlined } from '@ant-design/icons';

const { Title, Text } = Typography;

const UnmatchedTransactions: React.FC = () => {
  const { t } = useTranslation();

  return (
    <div>
      <Title level={2}>{t('transactions.unmatched.title')}</Title>
      <Text type="secondary" className="block mb-6">{t('transactions.unmatched.description')}</Text>
      
      <Card>
        <Alert
          message="Feature Under Development"
          description="The unmatched transactions view is currently being developed."
          type="info"
          showIcon
          icon={<CloseCircleOutlined />}
        />
      </Card>
    </div>
  );
};

export default UnmatchedTransactions; 