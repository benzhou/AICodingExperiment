import React, { useState, useEffect, useRef, useCallback } from 'react';
import { Button, Card, Form, Select, message, Typography, Steps, Table, Input, Space, Alert } from 'antd';
import { UploadOutlined, InboxOutlined, SaveOutlined } from '@ant-design/icons';
import type { UploadFile, UploadProps } from 'antd/es/upload/interface';
import { useTranslation } from 'react-i18next';
import { dataSourceService, DataSource } from '../services/dataSourceService';
import ServerStatus from './ServerStatus';
import api from '../services/api';
import { useNavigate } from 'react-router-dom';

const { Title, Text } = Typography;
const { Step } = Steps;
const { Option } = Select;

interface ColumnMapping {
  [key: string]: number;
}

// Add a new interface for schema mappings that uses string values
interface SchemaMapping {
  [key: string]: string;
}

const DataSourceUpload: React.FC = () => {
  const { t } = useTranslation();
  const [dataSources, setDataSources] = useState<DataSource[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [currentStep, setCurrentStep] = useState<number>(0);
  const [selectedFile, setSelectedFile] = useState<UploadFile | null>(null);
  const [rawFile, setRawFile] = useState<File | null>(null); // Store the actual File object
  const [uploading, setUploading] = useState<boolean>(false);
  const [selectedDataSource, setSelectedDataSource] = useState<string>('');
  const [dateFormat, setDateFormat] = useState<string>('2006-01-02');
  const [fileUploaded, setFileUploaded] = useState<boolean>(false);
  const [fileUploadSuccess, setFileUploadSuccess] = useState<boolean>(false);
  const [previewUrl, setPreviewUrl] = useState<string>('');
  const [columnMapping, setColumnMapping] = useState<ColumnMapping>({});
  const [previewData, setPreviewData] = useState<any[]>([]);
  const [columnOptions, setColumnOptions] = useState<string[]>([]);
  const [uploadError, setUploadError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const dropAreaRef = useRef<HTMLDivElement>(null);
  const navigate = useNavigate();

  // Required fields for transaction data
  const requiredFields = ['date', 'description', 'amount', 'reference'];
  const optionalFields = ['postDate', 'currency'];
  const allFields = [...requiredFields, ...optionalFields];

  const fetchDataSources = useCallback(async () => {
    try {
      setLoading(true);
      const sources = await dataSourceService.getDataSources();
      setDataSources(sources || []);
    } catch (error) {
      console.error('Error fetching data sources:', error);
      message.error(t('dataSources.errorFetch'));
      setDataSources([]); // Ensure dataSources is always at least an empty array
    } finally {
      setLoading(false);
    }
  }, [t]);

  useEffect(() => {
    fetchDataSources();
  }, [fetchDataSources]);

  // Handle file selection from the Antd Upload component
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const handleFileUpload: UploadProps['onChange'] = (info: any) => {
    console.log('Upload change event:', info.file.status, info.file.name);
    
    if (info.file.originFileObj instanceof File) {
      console.log('Origin file object available:', info.file.originFileObj.name);
      
      // Save the raw File object for later use
      setRawFile(info.file.originFileObj);
      
      // Set the UploadFile object for the UI display
      setSelectedFile(info.file);
      
      // Reset upload status flags when a new file is selected
      setFileUploaded(false);
      setFileUploadSuccess(false);
      setUploadError(null);
      
      message.info('File selected. Click the "Upload File to Server" button to start the upload.');
    } else {
      console.warn('No valid originFileObj in the selected file:', info.file);
      message.error(`${info.file.name} ${t('upload.fileUpload.failed')}: Missing file data`);
    }
  };

  // Handle direct file input selection
  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const file = e.target.files[0];
      console.log('File selected via direct input:', file.name);
      
      // Save the raw File object
      setRawFile(file);
      
      // Create a faux UploadFile object for UI consistency
      const fauxUploadFile = {
        uid: '-1',
        name: file.name,
        size: file.size,
        type: file.type,
        originFileObj: file,
      } as UploadFile;
      
      setSelectedFile(fauxUploadFile);
      
      // Reset upload status flags
      setFileUploaded(false);
      setFileUploadSuccess(false);
      setUploadError(null);
      
      message.info('File selected. Click the "Upload File to Server" button to start the upload.');
    }
  };

  // Drag and drop handlers
  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    if (dropAreaRef.current) {
      dropAreaRef.current.classList.add('drag-active');
    }
  };

  const handleDragLeave = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    if (dropAreaRef.current) {
      dropAreaRef.current.classList.remove('drag-active');
    }
  };

  // Handle file drop
  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    
    if (dropAreaRef.current) {
      dropAreaRef.current.classList.remove('drag-active');
    }

    if (!selectedDataSource) {
      message.error(t('upload.validation.dataSourceRequired'));
      return;
    }

    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      const file = e.dataTransfer.files[0];
      
      // Check if the file is a CSV
      if (!file.name.toLowerCase().endsWith('.csv')) {
        message.error(t('upload.validation.csvOnly'));
        return;
      }

      handleFileUploadToServer(file);
    }
  };

  // Process the uploaded file
  const handleFileUploadToServer = async (file: File) => {
    try {
      setUploading(true);
      const formData = new FormData();
      formData.append('file', file);
      formData.append('dataSourceId', selectedDataSource);

      // Upload the file to get a preview
      const uploadResponse = await api.post('/api/v1/uploads/preview', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });

      if (uploadResponse.data && uploadResponse.data.previewUrl) {
        setPreviewUrl(uploadResponse.data.previewUrl);
        
        // Fetch the data source to get its schema
        const dataSource = await dataSourceService.getDataSourceById(selectedDataSource);
        
        // If the data source has a schema definition with default mappings, use them
        if (dataSource.schemaDefinition?.defaultMappings) {
          const schemaMapping: SchemaMapping = dataSource.schemaDefinition.defaultMappings;
          const newMapping: ColumnMapping = {};
          
          // Find the column indices that match the field names in the schema mappings
          Object.entries(schemaMapping).forEach(([standardField, customField]) => {
            // Find the index of the custom field in the columns array
            const columnIndex = uploadResponse.data.columns.indexOf(customField);
            if (columnIndex !== -1) {
              newMapping[standardField] = columnIndex;
            }
          });
          
          if (Object.keys(newMapping).length > 0) {
            setColumnMapping(newMapping);
          }
        }
        
        // If date format is defined in the schema, use it
        if (dataSource.schemaDefinition?.dateFormat) {
          setDateFormat(dataSource.schemaDefinition.dateFormat);
        }
        
        // Get preview data
        const previewResponse = await api.get(`/api/v1/uploads/preview-data?url=${uploadResponse.data.previewUrl}`);
        if (previewResponse.data && previewResponse.data.data) {
          setPreviewData(previewResponse.data.data);
          
          // Make sure columns is an array of strings
          const columns: string[] = uploadResponse.data.columns || [];
          setColumnOptions(columns);
          
          setCurrentStep(1); // Move to mapping step since we have preview data
        } else {
          message.error(t('upload.validation.noDataParsed'));
        }
      }
    } catch (error) {
      console.error('File upload error:', error);
      message.error(t('upload.fileUpload.failed'));
    } finally {
      setUploading(false);
    }
  };

  // Trigger manual file selection
  const triggerFileSelection = () => {
    if (fileInputRef.current) {
      fileInputRef.current.click();
    }
  };

  const uploadFileToServer = async () => {
    // Make sure we have a file to upload
    if (!rawFile) {
      setUploadError(t('upload.validation.fileRequired'));
      message.error(t('upload.validation.selectFile'));
      return;
    }
    
    // Make sure dataSourceId is set properly
    if (!selectedDataSource) {
      setUploadError(t('upload.validation.dataSourceRequired'));
      message.error(t('upload.validation.selectDataSource'));
      return;
    }
    
    try {
      console.log('Starting upload process with file:', rawFile.name, 'size:', rawFile.size, 'bytes');
      setUploading(true);
      setUploadError(null);
      message.loading({ content: 'Uploading file to server...', key: 'fileUpload' });
      
      // Create a FormData object
      const formData = new FormData();
      formData.append('file', rawFile);
      formData.append('dataSourceId', selectedDataSource);
      
      // Log what we're sending
      console.log('FormData created with keys:', Array.from(formData.keys()));
      console.log('Uploading to data source ID:', selectedDataSource);
      
      // Use a direct XMLHttpRequest for more reliable FormData handling
      const xhr = new XMLHttpRequest();
      const baseUrl = process.env.REACT_APP_API_URL || 'http://localhost:8080';
      const url = `${baseUrl}/api/v1/uploads/preview`;
      
      console.log('Sending upload request to:', url);
      
      xhr.open('POST', url);
      
      // Add authorization header
      const token = localStorage.getItem('token');
      if (token) {
        xhr.setRequestHeader('Authorization', `Bearer ${token}`);
      }
      
      // Handle successful response
      xhr.onload = function() {
        if (xhr.status >= 200 && xhr.status < 300) {
          console.log('Upload successful, status:', xhr.status);
          try {
            const responseText = xhr.responseText;
            console.log('Raw response:', responseText.substring(0, 200) + '...');
            
            const data = JSON.parse(responseText);
            console.log('Parsed response data:', data);
            
            if (data && data.previewUrl) {
              setPreviewUrl(data.previewUrl);
              setFileUploaded(true);
              setFileUploadSuccess(true);
              setUploadError(null);
              
              if (data.preview && data.preview.length > 0) {
                setPreviewData(data.preview);
                console.log('Preview data received, rows:', data.preview.length);
              }
              
              if (data.suggestedMappings) {
                setColumnMapping(data.suggestedMappings);
              }
              
              message.success({ content: 'File uploaded successfully!', key: 'fileUpload' });
              
              // Move to next step automatically
              setCurrentStep(1);
            } else {
              const errorMsg = 'Invalid response data - missing previewUrl';
              console.error(errorMsg, data);
              setUploadError(errorMsg);
              message.error({ content: errorMsg, key: 'fileUpload' });
            }
          } catch (parseError: any) {
            const errorMsg = `Error parsing response: ${parseError.message}`;
            console.error(errorMsg, 'Raw response:', xhr.responseText);
            setUploadError(errorMsg);
            message.error({ content: errorMsg, key: 'fileUpload' });
          }
        } else {
          const errorMsg = `Upload failed with status: ${xhr.status} - ${xhr.statusText}`;
          console.error(errorMsg);
          setUploadError(errorMsg);
          message.error({ content: errorMsg, key: 'fileUpload' });
        }
        setUploading(false);
      };
      
      // Handle network errors
      xhr.onerror = function(e) {
        const errorMsg = 'Network error during upload. Please check your connection and the server status.';
        console.error(errorMsg, e);
        setUploadError(errorMsg);
        message.error({ content: errorMsg, key: 'fileUpload' });
        setUploading(false);
      };
      
      // Show upload progress
      xhr.upload.onprogress = function(event) {
        if (event.lengthComputable) {
          const percentComplete = Math.round((event.loaded / event.total) * 100);
          console.log(`Upload progress: ${percentComplete}%`);
        }
      };
      
      // Send the request
      xhr.send(formData);
      
    } catch (error: any) {
      const errorMsg = `Error in upload process: ${error.message || 'Unknown error'}`;
      console.error(errorMsg);
      setUploadError(errorMsg);
      message.error({ content: errorMsg, key: 'fileUpload' });
      setUploading(false);
    }
  };

  const uploadToDataSource = async () => {
    if (!selectedFile || !selectedDataSource) {
      message.error(t('upload.validation.selectDataSource'));
      return;
    }

    // Validate required mappings
    const missingFields = requiredFields.filter(field => columnMapping[field] === undefined);
    if (missingFields.length > 0) {
      message.error(`${t('upload.validation.requiredFields')}: ${missingFields.join(', ')}`);
      return;
    }

    try {
      setUploading(true);
      
      // Create an import record first
      const dataSource = await dataSourceService.getDataSourceById(selectedDataSource);
      const filename = previewUrl.split('/').pop() || 'unknown.csv';
      
      // Send the final confirmation to process the uploaded file with import tracking
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      const response = await api.post('/api/v1/uploads/process', {
        previewUrl: previewUrl,
        dataSourceId: selectedDataSource,
        dateFormat: dateFormat,
        columnMappings: columnMapping,
        createImportRecord: true,
        filename: filename,
        schemaDefinition: dataSource.schemaDefinition
      });
      
      message.success(t('upload.success'));
      
      // Navigate to the imports page for this data source
      navigate(`/datasources/${selectedDataSource}/imports`);
      
    } catch (error) {
      console.error('Error processing file:', error);
      message.error(t('upload.error'));
    } finally {
      setUploading(false);
    }
  };

  // Handle column mapping changes  
  const handleColumnMappingChange = (field: string, columnIndex: any) => {
    setColumnMapping({
      ...columnMapping,
      [field]: columnIndex
    });
  };
  
  // Render the upload step with improved UI
  const renderUploadStep = () => {
    return (
      <div>
        <Card title={t('upload.selectDataSource')}>
          <Form layout="vertical">
            <Form.Item label={t('upload.dataSource')} required>
              <Select
                placeholder={t('upload.selectDataSource')}
                value={selectedDataSource}
                onChange={(value) => setSelectedDataSource(value)}
                loading={loading}
              >
                {dataSources.map((source) => (
                  <Option key={source.id} value={source.id}>
                    {source.name}
                  </Option>
                ))}
              </Select>
            </Form.Item>
            
            <Form.Item label={t('upload.dateFormat')}>
              <Input
                placeholder="2006-01-02"
                value={dateFormat}
                onChange={(e) => setDateFormat(e.target.value)}
              />
              <Text type="secondary">
                {t('upload.dateFormatHint')}
              </Text>
            </Form.Item>
          </Form>
        </Card>
        
        <Card 
          title={t('upload.fileUpload.title')} 
          style={{ marginTop: '16px' }}
        >
          {/* Hidden direct file input */}
          <input
            type="file"
            ref={fileInputRef}
            style={{ display: 'none' }}
            accept=".csv"
            onChange={handleFileInputChange}
          />
          
          {/* Custom drag and drop area */}
          <div style={{ marginBottom: '20px' }}>
            <div style={{ marginBottom: '8px' }}>
              <Text strong>{t('upload.fileUpload.dragDropTitle')}</Text>
            </div>
            <div
              ref={dropAreaRef}
              className="custom-drop-area"
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onDrop={handleDrop}
              onClick={triggerFileSelection}
            >
              <InboxOutlined style={{ fontSize: '48px', color: '#40a9ff', marginBottom: '8px' }} />
              <p>{t('upload.fileUpload.dragText')}</p>
            </div>
            <div style={{ marginTop: '4px' }}>
              <Text type="secondary">{t('upload.fileUpload.hint')}</Text>
            </div>
          </div>
          
          {/* Show selected file info */}
          {selectedFile && (
            <div style={{ marginBottom: '16px' }}>
              <Alert
                message={t('upload.fileUpload.title')}
                description={`${t('upload.file')}: ${selectedFile.name}, ${t('common.size')}: ${(selectedFile.size || 0) / 1024} KB`}
                type="info"
                showIcon
              />
            </div>
          )}
          
          {/* Upload button */}
          <Button
            type="primary"
            onClick={uploadFileToServer}
            loading={uploading}
            disabled={!selectedFile || !selectedDataSource}
            icon={<UploadOutlined />}
            style={{ marginTop: '16px' }}
          >
            {uploading ? t('common.uploading') : t('upload.fileUpload.title')}
          </Button>
          
          {/* Display error if any */}
          {uploadError && (
            <Alert
              message={t('common.error')}
              description={uploadError}
              type="error"
              showIcon
              style={{ marginTop: '16px' }}
            />
          )}
          
          {/* Display success message if upload successful */}
          {fileUploadSuccess && (
            <Alert
              message={t('common.success')}
              description={t('upload.success')}
              type="success"
              showIcon
              style={{ marginTop: '16px' }}
            />
          )}
          
          {/* Debug information */}
          <div style={{ marginTop: '16px' }}>
            <details>
              <summary>{t('common.debugInfo')}</summary>
              <div style={{ marginTop: '8px', padding: '8px', backgroundColor: '#f5f5f5', borderRadius: '4px' }}>
                <p><strong>API URL:</strong> {process.env.REACT_APP_API_URL || 'http://localhost:8080'}</p>
                <p><strong>{t('upload.dataSource')}:</strong> {selectedDataSource || t('common.none')}</p>
                <p><strong>{t('upload.file')}:</strong> {selectedFile?.name || t('common.none')}</p>
                <p><strong>{t('common.rawFile')}:</strong> {rawFile ? t('common.available') : t('common.notAvailable')}</p>
                <p><strong>{t('common.fileUploaded')}:</strong> {fileUploaded ? t('common.yes') : t('common.no')}</p>
                <p><strong>{t('common.uploadSuccess')}:</strong> {fileUploadSuccess ? t('common.yes') : t('common.no')}</p>
                <p><strong>{t('common.uploadError')}:</strong> {uploadError || t('common.none')}</p>
              </div>
            </details>
          </div>
        </Card>
        
        {/* Server status info */}
        <div style={{ marginTop: '16px' }}>
          <ServerStatus />
        </div>
      </div>
    );
  };
  
  const steps = [
    {
      title: t('upload.fileUpload.title'),
      content: renderUploadStep(),
    },
    {
      title: t('upload.columnMapping.title'),
      content: (
        <Card title={t('upload.columnMapping.title')}>
          {previewData.length > 0 ? (
            <>
              <div className="mb-4">
                <Title level={5}>{t('upload.columnMapping.instructions')}</Title>
                <Text>{t('upload.columnMapping.description')}</Text>
              </div>
              
              <Form layout="vertical">
                {allFields.map((field) => (
                  <Form.Item 
                    key={field}
                    label={(
                      <Space>
                        {requiredFields.includes(field) && (
                          <Text type="danger">*</Text>
                        )}
                        {field}
                      </Space>
                    )}
                    required={requiredFields.includes(field)}
                  >
                    <Select
                      placeholder={t('upload.columnMapping.selectColumn')}
                      value={columnMapping[field]}
                      onChange={(value) => handleColumnMappingChange(field, value)}
                      style={{ width: '100%' }}
                    >
                      {previewData[0]?.map((header: any, index: number) => (
                        <Option key={index} value={index}>
                          {header} (Column {index + 1})
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                ))}
              </Form>
              
              {previewData.length > 1 && (
                <div className="mt-4">
                  <Title level={5}>{t('upload.preview')}</Title>
                  <Table 
                    dataSource={previewData.slice(1, 6).map((row: any[], index: number) => ({
                      key: index,
                      ...row.reduce((acc: Record<string, any>, cell: any, i: number) => {
                        acc[`col${i}`] = cell;
                        return acc;
                      }, {} as any)
                    }))}
                    columns={previewData[0].map((header: any, index: number) => ({
                      title: header,
                      dataIndex: `col${index}`,
                      key: `col${index}`,
                      ellipsis: true,
                    }))}
                    size="small"
                    pagination={false}
                    scroll={{ x: true }}
                  />
                </div>
              )}
            </>
          ) : (
            <Alert
              message={t('upload.noPreviewData')}
              description={t('upload.pleaseUploadFirst')}
              type="info"
              showIcon
            />
          )}
        </Card>
      ),
    },
    {
      title: t('upload.processData.title'),
      content: (
        <Card title={t('upload.processData.title')}>
          {previewUrl ? (
            <>
              <Alert
                message={t('upload.readyToProcess')}
                description={t('upload.processDescription')}
                type="info"
                showIcon
                className="mb-4"
              />
              
              <div className="mb-4">
                <Title level={5}>{t('upload.summary')}</Title>
                <ul>
                  <li>
                    <Text>
                      {t('upload.dataSource')}: {dataSources.find(ds => ds.id === selectedDataSource)?.name || ''}
                    </Text>
                  </li>
                  <li>
                    <Text>
                      {t('upload.file')}: {selectedFile?.name || ''}
                    </Text>
                  </li>
                  <li>
                    <Text>
                      {t('upload.dateFormat')}: {dateFormat}
                    </Text>
                  </li>
                  <li>
                    <Text>{t('upload.mappings')}:</Text>
                    <ul>
                      {Object.entries(columnMapping).map(([field, colIndex]) => (
                        <li key={field}>
                          <Text>
                            {field} → {previewData[0]?.[colIndex as number] || `Column ${colIndex}`}
                          </Text>
                        </li>
                      ))}
                    </ul>
                  </li>
                </ul>
              </div>
              
              <Button
                type="primary"
                icon={<SaveOutlined />}
                onClick={uploadToDataSource}
                loading={uploading}
              >
                {t('upload.processData.submit')}
              </Button>
            </>
          ) : (
            <Alert
              message={t('upload.noDataToProcess')}
              description={t('upload.pleaseUploadFirst')}
              type="warning"
              showIcon
            />
          )}
        </Card>
      ),
    },
  ];

  const next = () => {
    if (currentStep === 0) {
      if (!selectedDataSource) {
        message.error(t('upload.validation.selectDataSource'));
        return;
      }
    }

    if (currentStep === 1) {
      if (!selectedFile) {
        message.error(t('upload.validation.selectFile'));
        return;
      }
      
      // Check if file upload is still in progress
      if (uploading) {
        message.warning(t('upload.validation.waitForUpload'));
        return;
      }
      
      // Check if file was uploaded successfully
      if (!fileUploaded || !fileUploadSuccess) {
        message.error(t('upload.validation.uploadFirst'));
        return;
      }
      
      // Check if we received preview data from the server
      if (previewData.length === 0) {
        message.error(t('upload.validation.noDataParsed'));
        return;
      }
    }

    setCurrentStep(currentStep + 1);
  };

  const prev = () => {
    setCurrentStep(currentStep - 1);
  };

  return (
    <div>
      <Title level={2}>{t('upload.title')}</Title>
      <Text type="secondary" className="block mb-6">{t('upload.description')}</Text>
      
      <Card>
        <Steps current={currentStep} className="mb-8">
          {steps && steps.map(item => (
            <Step key={item.title} title={item.title} />
          ))}
        </Steps>
        
        <div className="mb-6">
          {currentStep === 0 && steps[0].content}
          {currentStep === 1 && steps[1].content}
          {currentStep === 2 && steps[2].content}
        </div>
        
        {/* Add debug status indicators for development - can be removed in production */}
        {process.env.NODE_ENV === 'development' && currentStep === 1 && (
          <div className="mb-4 p-3 border border-gray-300 rounded text-xs">
            <div><strong>Debug Status:</strong></div>
            <div>Selected File: {selectedFile ? selectedFile.name : 'None'}</div>
            <div>Raw File: {rawFile ? `${rawFile.name} (${rawFile.size} bytes)` : 'None'}</div>
            <div>Data Source: {selectedDataSource || 'None'}</div>
            <div>Uploading: {uploading ? 'Yes' : 'No'}</div>
            <div>File Uploaded: {fileUploaded ? 'Yes' : 'No'}</div>
            <div>Upload Success: {fileUploadSuccess ? 'Yes' : 'No'}</div>
            <div>Preview Data: {previewData.length > 0 ? `${previewData.length} rows` : 'None'}</div>
            <div>Preview URL: {previewUrl || 'None'}</div>
          </div>
        )}
        
        {/* Add a debug test button to check server connectivity directly */}
        {process.env.NODE_ENV === 'development' && currentStep === 1 && (
          <div className="mb-4 mt-2">
            <Button 
              type="default" 
              size="small"
              onClick={async () => {
                message.info(t('common.testingBackend'));
                const token = localStorage.getItem('token');
                
                // Log if we have a token
                if (!token) {
                  console.warn('No token available in localStorage');
                  message.error(t('auth.errors.noToken'));
                  return;
                }
                
                console.log('Token available, length:', token.length);
                console.log('Token prefix:', token.substring(0, 15) + '...');
                
                // First try a GET request to health endpoint to check connectivity
                try {
                  const healthResponse = await api.get('/health');
                  console.log('Health check response:', healthResponse.data);
                  message.success(t('common.connectionSuccess'));
                  
                  // Now try to fetch datasources to check authentication
                  try {
                    const dsResponse = await api.get('/api/v1/datasources');
                    console.log('Datasources response:', dsResponse.data);
                    message.success(t('common.authSuccess'));
                  } catch (authErr) {
                    console.error('Auth check failed:', authErr);
                    message.warning(t('common.authFailedServerReachable'));
                  }
                } catch (healthErr) {
                  console.error('Health check failed:', healthErr);
                  message.error(t('common.connectionFailed'));
                }
              }}
            >
              {t('common.testBackendConnection')}
            </Button>
          </div>
        )}
        
        <div className="flex justify-between">
          {currentStep > 0 && (
            <Button onClick={prev}>
              {t('upload.navigation.previous')}
            </Button>
          )}
          
          <div className="ml-auto">
            {currentStep < steps.length - 1 && (
              <Button 
                type="primary" 
                onClick={next}
                disabled={currentStep === 1 && uploading}
              >
                {t('upload.navigation.next')}
              </Button>
            )}
            
            {currentStep === steps.length - 1 && (
              <Button
                type="primary"
                icon={<SaveOutlined />}
                onClick={uploadToDataSource}
                loading={uploading}
              >
                {t('upload.navigation.uploadData')}
              </Button>
            )}
          </div>
        </div>
      </Card>
    </div>
  );
};

export default DataSourceUpload; 