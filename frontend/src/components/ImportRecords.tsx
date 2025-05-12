import React, { useState, useEffect, useCallback } from 'react';
import { useParams, Link } from 'react-router-dom';
import { Table, Card, Typography, Space, Button, Tag, Tooltip, Empty, Spin, Modal, Pagination } from 'antd';
import { 
  ReloadOutlined, 
  FileTextOutlined, 
  DeleteOutlined, 
  ExclamationCircleOutlined, 
  CheckCircleOutlined,
  CloseCircleOutlined,
  SyncOutlined
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { dataSourceService, ImportRecord, DataSource } from '../services/dataSourceService';

const { Title, Text } = Typography;
const { confirm } = Modal;

const ImportRecords: React.FC = () => {
  const { t } = useTranslation();
  const { dataSourceId } = useParams<{ dataSourceId: string }>();
  const [loading, setLoading] = useState<boolean>(true);
  const [imports, setImports] = useState<ImportRecord[]>([]);
  const [dataSource, setDataSource] = useState<DataSource | null>(null);
  const [total, setTotal] = useState<number>(0);
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);

  // Fetch import records with pagination
  const fetchImports = useCallback(async () => {
    if (!dataSourceId) return;
    
    try {
      setLoading(true);
      
      // Fetch the data source details
      const ds = await dataSourceService.getDataSourceById(dataSourceId);
      setDataSource(ds);
      
      // Fetch import records with pagination
      const response = await dataSourceService.getImportsByDataSource(
        dataSourceId,
        pageSize,
        (currentPage - 1) * pageSize
      );
      
      setImports(response.data);
      setTotal(response.pagination.total);
    } catch (error) {
      console.error('Error fetching import records:', error);
    } finally {
      setLoading(false);
    }
  }, [dataSourceId, currentPage, pageSize]);

  useEffect(() => {
    fetchImports();
  }, [fetchImports]);

  const handleRefresh = () => {
    fetchImports();
  };

  const handlePageChange = (page: number, pageSize?: number) => {
    setCurrentPage(page);
    if (pageSize) setPageSize(pageSize);
  };

  const confirmDelete = (importId: string) => {
    confirm({
      title: t('imports.deleteConfirmTitle'),
      icon: <ExclamationCircleOutlined />,
      content: t('imports.deleteConfirmContent'),
      okText: t('imports.deleteConfirmOk'),
      okType: 'danger',
      cancelText: t('imports.deleteConfirmCancel'),
      onOk: async () => {
        try {
          await dataSourceService.deleteImport(importId);
          fetchImports();
        } catch (error) {
          console.error('Error deleting import:', error);
        }
      }
    });
  };

  const getStatusTag = (status: string) => {
    switch (status) {
      case 'Completed':
        return <Tag icon={<CheckCircleOutlined />} color="success">{t('imports.statusCompleted')}</Tag>;
      case 'Failed':
        return <Tag icon={<CloseCircleOutlined />} color="error">{t('imports.statusFailed')}</Tag>;
      case 'Processing':
        return <Tag icon={<SyncOutlined spin />} color="processing">{t('imports.statusProcessing')}</Tag>;
      default:
        return <Tag>{status}</Tag>;
    }
  };

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
  };

  const columns = [
    {
      title: t('imports.fileName'),
      dataIndex: 'fileName',
      key: 'fileName',
      render: (text: string, record: ImportRecord) => (
        <Link to={`/imports/${record.id}`}>
          <Space>
            <FileTextOutlined />
            {text}
          </Space>
        </Link>
      )
    },
    {
      title: t('imports.status'),
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => getStatusTag(status)
    },
    {
      title: t('imports.fileSize'),
      dataIndex: 'fileSize',
      key: 'fileSize',
      render: (size: number) => formatFileSize(size)
    },
    {
      title: t('imports.rowCount'),
      dataIndex: 'rowCount',
      key: 'rowCount',
    },
    {
      title: t('imports.successCount'),
      dataIndex: 'successCount',
      key: 'successCount',
      render: (successCount: number, record: ImportRecord) => (
        <Text type={record.errorCount > 0 ? 'warning' : 'success'}>
          {successCount} / {record.rowCount}
        </Text>
      )
    },
    {
      title: t('imports.createdAt'),
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (timestamp: number) => new Date(timestamp).toLocaleString()
    },
    {
      title: t('imports.actions'),
      key: 'actions',
      render: (_: any, record: ImportRecord) => (
        <Space>
          <Link to={`/imports/${record.id}/transactions`}>
            <Button type="primary" size="small">
              {t('imports.viewTransactions')}
            </Button>
          </Link>
          <Tooltip title={t('imports.delete')}>
            <Button 
              type="text" 
              danger 
              icon={<DeleteOutlined />} 
              onClick={() => confirmDelete(record.id)}
              disabled={record.status === 'Processing'}
            />
          </Tooltip>
        </Space>
      )
    }
  ];

  return (
    <div>
      <Title level={2}>
        {dataSource ? t('imports.titleWithName', { name: dataSource.name }) : t('imports.title')}
      </Title>
      
      <Card>
        <div className="flex justify-between items-center mb-4">
          <Text type="secondary">
            {t('imports.description')}
          </Text>
          <Button 
            icon={<ReloadOutlined />} 
            onClick={handleRefresh}
            loading={loading}
          >
            {t('imports.refresh')}
          </Button>
        </div>
        
        {loading ? (
          <div className="flex justify-center my-8">
            <Spin size="large" />
          </div>
        ) : imports.length > 0 ? (
          <>
            <Table 
              columns={columns} 
              dataSource={imports} 
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
            description={t('imports.noImports')} 
            image={Empty.PRESENTED_IMAGE_SIMPLE}
          />
        )}
      </Card>
    </div>
  );
};

export default ImportRecords; 