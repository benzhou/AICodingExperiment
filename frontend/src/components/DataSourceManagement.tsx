import React, { useState, useEffect, useCallback } from 'react';
import { Button, Card, Form, Input, Table, Modal, message, Typography, Space, Tabs, Pagination } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, SearchOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { dataSourceService, DataSource, SchemaDefinition, PaginatedResponse } from '../services/dataSourceService';
import SchemaDefinitionForm from './SchemaDefinitionForm';
import { Link } from 'react-router-dom';
import { debounce } from 'lodash';

const { Title, Text } = Typography;
const { TextArea } = Input;
const { TabPane } = Tabs;
const { Search } = Input;

const DataSourceManagement: React.FC = () => {
  const { t } = useTranslation();
  const [dataSources, setDataSources] = useState<DataSource[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [modalVisible, setModalVisible] = useState<boolean>(false);
  const [editingDataSource, setEditingDataSource] = useState<DataSource | null>(null);
  const [form] = Form.useForm();
  const [activeTabKey, setActiveTabKey] = useState<string>("basic");
  const [schemaDefinition, setSchemaDefinition] = useState<SchemaDefinition | undefined>();
  const [searchText, setSearchText] = useState<string>('');
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);
  const [total, setTotal] = useState<number>(0);

  const fetchDataSources = useCallback(async () => {
    try {
      setLoading(true);
      if (searchText) {
        // Use the search API if search text is provided
        const response = await dataSourceService.searchDataSources(
          searchText,
          pageSize,
          (currentPage - 1) * pageSize
        );
        setDataSources(response.data);
        setTotal(response.pagination.total);
      } else {
        // Otherwise fetch all data sources
        const sources = await dataSourceService.getDataSources();
        setDataSources(sources);
        setTotal(sources.length);
      }
    } catch (error) {
      console.error('Error fetching data sources:', error);
      message.error(t('dataSources.errorFetch'));
    } finally {
      setLoading(false);
    }
  }, [t, searchText, currentPage, pageSize]);

  // Fetch data sources on component mount and when search/pagination changes
  useEffect(() => {
    fetchDataSources();
  }, [fetchDataSources]);

  // Debounced search handler
  const debouncedSearch = useCallback(
    debounce((value: string) => {
      setSearchText(value);
      setCurrentPage(1); // Reset to first page on new search
    }, 500),
    []
  );

  // Handle search
  const handleSearch = (value: string) => {
    debouncedSearch(value);
  };

  // Handle pagination change
  const handlePageChange = (page: number, pageSize?: number) => {
    setCurrentPage(page);
    if (pageSize) setPageSize(pageSize);
  };

  const handleCreate = () => {
    setEditingDataSource(null);
    form.resetFields();
    setSchemaDefinition(undefined);
    setActiveTabKey("basic");
    setModalVisible(true);
  };

  const handleEdit = (record: DataSource) => {
    setEditingDataSource(record);
    form.setFieldsValue({
      name: record.name,
      description: record.description
    });
    setSchemaDefinition(record.schemaDefinition);
    setActiveTabKey("basic");
    setModalVisible(true);
  };

  const handleDelete = async (id: string) => {
    Modal.confirm({
      title: t('dataSources.deleteConfirm'),
      content: t('dataSources.deleteConfirmNotice'),
      okText: t('dataSources.deleteConfirmYes'),
      okType: 'danger',
      cancelText: t('dataSources.cancel'),
      onOk: async () => {
        try {
          await dataSourceService.deleteDataSource(id);
          message.success(t('dataSources.successDeleted'));
          fetchDataSources();
        } catch (error) {
          console.error('Error deleting data source:', error);
          message.error(t('dataSources.errorDelete'));
        }
      }
    });
  };

  const handleSubmit = async (values: any) => {
    try {
      // Add schema definition to the form values
      const finalValues = { 
        ...values,
        schemaDefinition
      };
      
      if (editingDataSource) {
        // Update existing data source
        await dataSourceService.updateDataSource(editingDataSource.id, finalValues);
        message.success(t('dataSources.successUpdated'));
      } else {
        // Create new data source
        await dataSourceService.createDataSource(finalValues);
        message.success(t('dataSources.successCreated'));
      }
      
      setModalVisible(false);
      fetchDataSources();
    } catch (error) {
      console.error('Error saving data source:', error);
      message.error(t('dataSources.errorSave'));
    }
  };

  const columns = [
    {
      title: t('dataSources.name'),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t('dataSources.description'),
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: t('dataSources.schemaStatus'),
      key: 'schema',
      render: (record: DataSource) => (
        <Text>
          {record.schemaDefinition && record.schemaDefinition.fields.length > 0 
            ? t('dataSources.schemaConfigured', { fields: record.schemaDefinition.fields.length })
            : t('dataSources.schemaNotConfigured')
          }
        </Text>
      )
    },
    {
      title: t('dataSources.createdAt'),
      dataIndex: 'created_at',
      key: 'created_at',
      render: (value: any) => {
        try {
          // Handle both string dates and epoch milliseconds
          const date = typeof value === 'number' 
            ? new Date(value)
            : new Date(value);
          
          return date instanceof Date && !isNaN(date.getTime()) 
            ? date.toLocaleString()
            : 'N/A';
        } catch (error) {
          console.error('Error parsing date:', value, error);
          return 'N/A';
        }
      },
    },
    {
      title: t('dataSources.actions'),
      key: 'actions',
      render: (_: any, record: DataSource) => (
        <Space>
          <Button 
            icon={<EditOutlined />} 
            onClick={() => handleEdit(record)}
            type="primary"
            ghost
          />
          <Link to={`/datasources/${record.id}/imports`}>
            <Button type="primary" size="middle">
              {t('dataSources.viewImports')}
            </Button>
          </Link>
          <Button 
            icon={<DeleteOutlined />}
            onClick={() => handleDelete(record.id)}
            danger
          />
        </Space>
      ),
    },
  ];

  // Change handler for schema definition
  const handleSchemaChange = (newSchema: SchemaDefinition) => {
    setSchemaDefinition(newSchema);
  };

  return (
    <div>
      <Title level={2}>{t('dataSources.title')}</Title>
      
      <Card>
        <div className="flex justify-between items-center mb-4">
          <Search
            placeholder={t('dataSources.search')}
            allowClear
            enterButton={<SearchOutlined />}
            size="middle"
            onSearch={handleSearch}
            onChange={(e) => handleSearch(e.target.value)}
            style={{ width: 300 }}
          />
          <Button 
            type="primary" 
            icon={<PlusOutlined />} 
            onClick={handleCreate}
          >
            {t('dataSources.addSource')}
          </Button>
        </div>
        
        <Table 
          columns={columns} 
          dataSource={dataSources} 
          rowKey="id"
          loading={loading}
          pagination={false}
        />
        
        {total > 0 && (
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
        )}
      </Card>

      {/* Create/Edit Modal with Tabs */}
      <Modal
        title={editingDataSource ? t('dataSources.edit') : t('dataSources.create')} 
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
        width={800}
      >
        <Tabs activeKey={activeTabKey} onChange={setActiveTabKey}>
          <TabPane tab={t('dataSources.basicInfo')} key="basic">
            <Form 
              form={form}
              layout="vertical"
              onFinish={handleSubmit}
            >
              <Form.Item 
                name="name" 
                label={t('dataSources.name')}
                rules={[{ required: true, message: t('dataSources.nameRequired') }]}
              >
                <Input placeholder={t('dataSources.name')} />
              </Form.Item>
              
              <Form.Item 
                name="description" 
                label={t('dataSources.description')}
              >
                <TextArea 
                  rows={4} 
                  placeholder={t('dataSources.description')} 
                />
              </Form.Item>
              
              <div className="flex justify-end">
                <Button className="mr-2" onClick={() => setModalVisible(false)}>
                  {t('dataSources.cancel')}
                </Button>
                <Button type="primary" htmlType="submit">
                  {t('dataSources.save')}
                </Button>
              </div>
            </Form>
          </TabPane>
          
          <TabPane tab={t('dataSources.schema')} key="schema">
            <SchemaDefinitionForm
              value={schemaDefinition}
              onChange={handleSchemaChange}
            />
            
            <div className="flex justify-end mt-4">
              <Button className="mr-2" onClick={() => setModalVisible(false)}>
                {t('dataSources.cancel')}
              </Button>
              <Button type="primary" onClick={() => form.submit()}>
                {t('dataSources.save')}
              </Button>
            </div>
          </TabPane>
        </Tabs>
      </Modal>
    </div>
  );
};

export default DataSourceManagement; 