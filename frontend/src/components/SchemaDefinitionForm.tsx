import React, { useState } from 'react';
import { Form, Input, Button, Space, Select, Checkbox, Card, Typography, Divider, Row, Col, Tooltip } from 'antd';
import { PlusOutlined, MinusCircleOutlined, InfoCircleOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { SchemaField, SchemaDefinition } from '../services/dataSourceService';

const { Text, Title } = Typography;
const { Option } = Select;

interface SchemaDefinitionFormProps {
  value?: SchemaDefinition;
  onChange?: (value: SchemaDefinition) => void;
}

const defaultSchemaDefinition: SchemaDefinition = {
  fields: [],
  dateFormat: 'YYYY-MM-DD',
  defaultMappings: {},
  requiredFields: ['date', 'amount', 'description']
};

const SchemaDefinitionForm: React.FC<SchemaDefinitionFormProps> = ({ value, onChange }) => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  
  // Use either the provided value or default schema definition
  const [schemaDefinition, setSchemaDefinition] = useState<SchemaDefinition>(
    value || defaultSchemaDefinition
  );

  // Field type options
  const fieldTypes = [
    { value: 'string', label: t('schema.fieldTypes.string') },
    { value: 'number', label: t('schema.fieldTypes.number') },
    { value: 'date', label: t('schema.fieldTypes.date') },
    { value: 'boolean', label: t('schema.fieldTypes.boolean') }
  ];

  // Standard transaction fields for default mappings
  const standardFields = [
    { value: 'date', label: t('schema.standardFields.date') },
    { value: 'description', label: t('schema.standardFields.description') },
    { value: 'amount', label: t('schema.standardFields.amount') },
    { value: 'reference', label: t('schema.standardFields.reference') },
    { value: 'postDate', label: t('schema.standardFields.postDate') },
    { value: 'currency', label: t('schema.standardFields.currency') }
  ];

  // Handle field changes
  const handleFieldChange = (index: number, field: Partial<SchemaField>) => {
    const updatedFields = [...schemaDefinition.fields];
    updatedFields[index] = { ...updatedFields[index], ...field };
    
    // Update schema definition
    const updatedSchema = {
      ...schemaDefinition,
      fields: updatedFields
    };
    
    setSchemaDefinition(updatedSchema);
    if (onChange) {
      onChange(updatedSchema);
    }
  };

  // Add a new field
  const addField = () => {
    const newField: SchemaField = {
      name: '',
      displayName: '',
      type: 'string',
      required: false,
      description: ''
    };
    
    const updatedFields = [...schemaDefinition.fields, newField];
    const updatedSchema = {
      ...schemaDefinition,
      fields: updatedFields
    };
    
    setSchemaDefinition(updatedSchema);
    if (onChange) {
      onChange(updatedSchema);
    }
  };

  // Remove a field
  const removeField = (index: number) => {
    const updatedFields = schemaDefinition.fields.filter((_, i) => i !== index);
    
    // Also remove any default mappings that reference this field
    const fieldName = schemaDefinition.fields[index].name;
    const updatedMappings = { ...schemaDefinition.defaultMappings };
    
    Object.keys(updatedMappings).forEach(key => {
      if (updatedMappings[key] === fieldName) {
        delete updatedMappings[key];
      }
    });
    
    const updatedSchema = {
      ...schemaDefinition,
      fields: updatedFields,
      defaultMappings: updatedMappings
    };
    
    setSchemaDefinition(updatedSchema);
    if (onChange) {
      onChange(updatedSchema);
    }
  };

  // Update date format
  const handleDateFormatChange = (dateFormat: string) => {
    const updatedSchema = {
      ...schemaDefinition,
      dateFormat
    };
    
    setSchemaDefinition(updatedSchema);
    if (onChange) {
      onChange(updatedSchema);
    }
  };

  // Update default mappings
  const handleMappingChange = (standardField: string, customField: string) => {
    const updatedMappings = {
      ...schemaDefinition.defaultMappings,
      [standardField]: customField
    };
    
    // Remove mapping if empty value is selected
    if (!customField) {
      delete updatedMappings[standardField];
    }
    
    const updatedSchema = {
      ...schemaDefinition,
      defaultMappings: updatedMappings
    };
    
    setSchemaDefinition(updatedSchema);
    if (onChange) {
      onChange(updatedSchema);
    }
  };

  // Update required fields
  const handleRequiredFieldsChange = (field: string, isRequired: boolean) => {
    let updatedRequiredFields = [...schemaDefinition.requiredFields];
    
    if (isRequired && !updatedRequiredFields.includes(field)) {
      updatedRequiredFields.push(field);
    } else if (!isRequired) {
      updatedRequiredFields = updatedRequiredFields.filter(f => f !== field);
    }
    
    const updatedSchema = {
      ...schemaDefinition,
      requiredFields: updatedRequiredFields
    };
    
    setSchemaDefinition(updatedSchema);
    if (onChange) {
      onChange(updatedSchema);
    }
  };

  return (
    <div className="schema-definition-form">
      <Card title={t('schema.title')}>
        <div className="mb-4">
          <Text type="secondary">{t('schema.description')}</Text>
        </div>
        
        {/* Fields section */}
        <Title level={5}>{t('schema.fields.title')}</Title>
        <div className="mb-4">
          {schemaDefinition.fields.map((field, index) => (
            <Card 
              key={index} 
              size="small" 
              className="mb-2"
              title={field.name || t('schema.fields.newField')}
              extra={
                <Button 
                  type="text" 
                  danger 
                  icon={<MinusCircleOutlined />} 
                  onClick={() => removeField(index)}
                />
              }
            >
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item 
                    label={t('schema.fields.name')}
                    required
                    tooltip={t('schema.fields.nameTooltip')}
                  >
                    <Input
                      value={field.name}
                      onChange={e => handleFieldChange(index, { name: e.target.value })}
                      placeholder={t('schema.fields.namePlaceholder')}
                    />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item 
                    label={t('schema.fields.displayName')}
                    tooltip={t('schema.fields.displayNameTooltip')}
                  >
                    <Input
                      value={field.displayName}
                      onChange={e => handleFieldChange(index, { displayName: e.target.value })}
                      placeholder={t('schema.fields.displayNamePlaceholder')}
                    />
                  </Form.Item>
                </Col>
              </Row>
              
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item 
                    label={t('schema.fields.type')}
                    required
                  >
                    <Select
                      value={field.type}
                      onChange={type => handleFieldChange(index, { type })}
                    >
                      {fieldTypes.map(type => (
                        <Option key={type.value} value={type.value}>{type.label}</Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item 
                    label={t('schema.fields.format')}
                    tooltip={t('schema.fields.formatTooltip')}
                  >
                    <Input
                      value={field.format}
                      onChange={e => handleFieldChange(index, { format: e.target.value })}
                      placeholder={field.type === 'date' ? 'YYYY-MM-DD' : ''}
                    />
                  </Form.Item>
                </Col>
              </Row>
              
              <Row>
                <Col span={24}>
                  <Form.Item>
                    <Checkbox
                      checked={field.required}
                      onChange={e => handleFieldChange(index, { required: e.target.checked })}
                    >
                      {t('schema.fields.required')}
                    </Checkbox>
                  </Form.Item>
                </Col>
              </Row>
              
              <Row>
                <Col span={24}>
                  <Form.Item label={t('schema.fields.description')}>
                    <Input.TextArea
                      value={field.description}
                      onChange={e => handleFieldChange(index, { description: e.target.value })}
                      placeholder={t('schema.fields.descriptionPlaceholder')}
                      rows={2}
                    />
                  </Form.Item>
                </Col>
              </Row>
            </Card>
          ))}
          
          <Button 
            type="dashed" 
            onClick={addField} 
            block 
            icon={<PlusOutlined />}
          >
            {t('schema.fields.addField')}
          </Button>
        </div>
        
        <Divider />
        
        {/* Date format section */}
        <Title level={5}>{t('schema.dateFormat.title')}</Title>
        <div className="mb-4">
          <Form.Item 
            label={t('schema.dateFormat.label')}
            tooltip={t('schema.dateFormat.tooltip')}
          >
            <Input
              value={schemaDefinition.dateFormat}
              onChange={e => handleDateFormatChange(e.target.value)}
              placeholder="YYYY-MM-DD"
            />
          </Form.Item>
        </div>
        
        <Divider />
        
        {/* Default mappings section */}
        <Title level={5}>{t('schema.mappings.title')}</Title>
        <div className="mb-4">
          <Text type="secondary">{t('schema.mappings.description')}</Text>
          
          {standardFields.map(standardField => (
            <Form.Item 
              key={standardField.value}
              label={
                <Space>
                  {standardField.label}
                  <Tooltip title={t('schema.mappings.fieldTooltip', { field: standardField.label })}>
                    <InfoCircleOutlined />
                  </Tooltip>
                </Space>
              }
            >
              <Select
                value={schemaDefinition.defaultMappings[standardField.value] || undefined}
                onChange={value => handleMappingChange(standardField.value, value)}
                placeholder={t('schema.mappings.selectField')}
                allowClear
              >
                {schemaDefinition.fields.map(field => (
                  <Option key={field.name} value={field.name}>
                    {field.displayName || field.name}
                  </Option>
                ))}
              </Select>
            </Form.Item>
          ))}
        </div>
        
        <Divider />
        
        {/* Required fields section */}
        <Title level={5}>{t('schema.requiredFields.title')}</Title>
        <div className="mb-4">
          <Text type="secondary">{t('schema.requiredFields.description')}</Text>
          
          <Form.Item>
            <Space direction="vertical" style={{ width: '100%' }}>
              {standardFields.map(field => (
                <Checkbox
                  key={field.value}
                  checked={schemaDefinition.requiredFields.includes(field.value)}
                  onChange={e => handleRequiredFieldsChange(field.value, e.target.checked)}
                >
                  {field.label}
                </Checkbox>
              ))}
            </Space>
          </Form.Item>
        </div>
      </Card>
    </div>
  );
};

export default SchemaDefinitionForm; 