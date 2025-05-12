import React, { useState, useEffect, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import { 
  Table, 
  Card, 
  Typography, 
  Space, 
  Button, 
  Spin, 
  Empty, 
  Pagination, 
  Badge, 
  Tooltip,
  Modal
} from 'antd';
import { ReloadOutlined, InfoCircleOutlined, EyeOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { dataSourceService, ImportRecord, RawTransaction } from '../services/dataSourceService';
import { JSONTree } from 'react-json-tree';

const { Title, Text } = Typography;

const RawTransactionsList: React.FC = () => {
  const { t } = useTranslation();
  const { importId } = useParams<{ importId: string }>();
  const [loading, setLoading] = useState<boolean>(true);
  const [transactions, setTransactions] = useState<RawTransaction[]>([]);
  const [importRecord, setImportRecord] = useState<ImportRecord | null>(null);
  const [total, setTotal] = useState<number>(0);
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);
  const [detailsModalVisible, setDetailsModalVisible] = useState<boolean>(false);
  const [selectedTransaction, setSelectedTransaction] = useState<RawTransaction | null>(null);

  // Fetch raw transactions with pagination
  const fetchTransactions = useCallback(async () => {
    if (!importId) return;
    
    try {
      setLoading(true);
      
      // Fetch the import record details
      const importDetails = await dataSourceService.getImportById(importId);
      setImportRecord(importDetails);
      
      // Fetch raw transactions with pagination
      const response = await dataSourceService.getRawTransactionsByImport(
        importId,
        pageSize,
        (currentPage - 1) * pageSize
      );
      
      setTransactions(response.data);
      setTotal(response.pagination.total);
    } catch (error) {
      console.error('Error fetching raw transactions:', error);
    } finally {
      setLoading(false);
    }
  }, [importId, currentPage, pageSize]);

  useEffect(() => {
    fetchTransactions();
  }, [fetchTransactions]);

  const handleRefresh = () => {
    fetchTransactions();
  };

  const handlePageChange = (page: number, pageSize?: number) => {
    setCurrentPage(page);
    if (pageSize) setPageSize(pageSize);
  };

  const showTransactionDetails = (transaction: RawTransaction) => {
    setSelectedTransaction(transaction);
    setDetailsModalVisible(true);
  };

  const columns = [
    {
      title: t('rawTransactions.rowNumber'),
      dataIndex: 'rowNumber',
      key: 'rowNumber',
      width: 100
    },
    {
      title: t('rawTransactions.data'),
      key: 'data',
      render: (record: RawTransaction) => {
        // Display a summary of the data (first few key-value pairs)
        const data = record.data;
        const entries = Object.entries(data).slice(0, 3);
        
        return (
          <Space direction="vertical" size="small" style={{ width: '100%' }}>
            {entries.map(([key, value]) => (
              <div key={key}>
                <Text strong>{key}:</Text> {' '}
                <Text>{JSON.stringify(value).substring(0, 50)}</Text>
              </div>
            ))}
            {Object.keys(data).length > 3 && (
              <Text type="secondary">
                {t('rawTransactions.moreFields', { count: Object.keys(data).length - 3 })}
              </Text>
            )}
          </Space>
        );
      }
    },
    {
      title: t('rawTransactions.status'),
      key: 'status',
      width: 100,
      render: (record: RawTransaction) => (
        record.errorMessage ? (
          <Tooltip title={record.errorMessage}>
            <Badge status="error" text={t('rawTransactions.error')} />
          </Tooltip>
        ) : (
          <Badge status="success" text={t('rawTransactions.success')} />
        )
      )
    },
    {
      title: t('rawTransactions.createdAt'),
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (timestamp: number) => new Date(timestamp).toLocaleString(),
      width: 180
    },
    {
      title: t('common.actions'),
      key: 'actions',
      width: 100,
      render: (_: any, record: RawTransaction) => (
        <Button 
          type="text" 
          icon={<EyeOutlined />} 
          onClick={() => showTransactionDetails(record)}
        />
      )
    }
  ];

  // JSON tree theme for the details modal
  const jsonTheme = {
    scheme: 'monokai',
    base00: '#272822',
    base01: '#383830',
    base02: '#49483e',
    base03: '#75715e',
    base04: '#a59f85',
    base05: '#f8f8f2',
    base06: '#f5f4f1',
    base07: '#f9f8f5',
    base08: '#f92672',
    base09: '#fd971f',
    base0A: '#f4bf75',
    base0B: '#a6e22e',
    base0C: '#a1efe4',
    base0D: '#66d9ef',
    base0E: '#ae81ff',
    base0F: '#cc6633'
  };

  return (
    <div>
      <Title level={2}>
        {importRecord ? t('rawTransactions.titleWithFileName', { name: importRecord.fileName }) : t('rawTransactions.title')}
      </Title>
      
      <Card>
        {importRecord && (
          <div className="mb-4 bg-gray-50 p-4 rounded">
            <Space size="large">
              <div>
                <Text type="secondary">{t('imports.fileName')}:</Text>{' '}
                <Text strong>{importRecord.fileName}</Text>
              </div>
              <div>
                <Text type="secondary">{t('imports.status')}:</Text>{' '}
                <Text strong>{importRecord.status}</Text>
              </div>
              <div>
                <Text type="secondary">{t('imports.rowCount')}:</Text>{' '}
                <Text strong>{importRecord.rowCount}</Text>
              </div>
              <div>
                <Text type="secondary">{t('imports.successRate')}:</Text>{' '}
                <Text 
                  strong 
                  type={importRecord.errorCount > 0 ? 'danger' : 'success'}
                >
                  {importRecord.rowCount > 0 
                    ? Math.round((importRecord.successCount / importRecord.rowCount) * 100) 
                    : 0}%
                </Text>
              </div>
            </Space>
          </div>
        )}
        
        <div className="flex justify-between items-center mb-4">
          <Text type="secondary">
            {t('rawTransactions.description')}
          </Text>
          <Button 
            icon={<ReloadOutlined />} 
            onClick={handleRefresh}
            loading={loading}
          >
            {t('common.refresh')}
          </Button>
        </div>
        
        {loading ? (
          <div className="flex justify-center my-8">
            <Spin size="large" />
          </div>
        ) : transactions.length > 0 ? (
          <>
            <Table 
              columns={columns} 
              dataSource={transactions} 
              rowKey="id"
              pagination={false}
            />
            
            <div className="flex justify-end mt-4">
              <Pagination
                current={currentPage}
                pageSize={pageSize}
                total={total}
                onChange={handlePageChange}
                showSizeChanger
                showTotal={(total) => t('pagination.showTotal', { total })}
              />
            </div>
          </>
        ) : (
          <Empty 
            description={t('rawTransactions.noTransactions')} 
            image={Empty.PRESENTED_IMAGE_SIMPLE}
          />
        )}
      </Card>
      
      {/* Transaction Details Modal */}
      <Modal
        title={t('rawTransactions.detailsTitle')}
        open={detailsModalVisible}
        onCancel={() => setDetailsModalVisible(false)}
        width={800}
        footer={[
          <Button key="close" onClick={() => setDetailsModalVisible(false)}>
            {t('common.close')}
          </Button>
        ]}
      >
        {selectedTransaction && (
          <div>
            <div className="mb-4">
              <Space direction="vertical" style={{ width: '100%' }}>
                <div>
                  <Text type="secondary">{t('rawTransactions.rowNumber')}:</Text>{' '}
                  <Text strong>{selectedTransaction.rowNumber}</Text>
                </div>
                <div>
                  <Text type="secondary">{t('rawTransactions.createdAt')}:</Text>{' '}
                  <Text strong>{new Date(selectedTransaction.createdAt).toLocaleString()}</Text>
                </div>
                {selectedTransaction.errorMessage && (
                  <div>
                    <Text type="danger">{t('rawTransactions.errorMessage')}:</Text>{' '}
                    <Text type="danger">{selectedTransaction.errorMessage}</Text>
                  </div>
                )}
              </Space>
            </div>
            
            <Card title={t('rawTransactions.dataJson')} size="small">
              <JSONTree 
                data={selectedTransaction.data} 
                theme={jsonTheme}
                invertTheme={false}
              />
            </Card>
          </div>
        )}
      </Modal>
    </div>
  );
};

export default RawTransactionsList; 