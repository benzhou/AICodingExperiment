import React, { useState, useEffect } from 'react';
import { 
  Table, Button, Modal, Form, Input, Select, message, 
  Tag, Space, Typography, Card, Tooltip, Popconfirm, Spin
} from 'antd';
import { 
  UserAddOutlined, EditOutlined, DeleteOutlined, 
  CheckCircleOutlined, CloseCircleOutlined
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { 
  userService, 
  User, 
  CreateUserRequest, 
  UpdateRoleRequest 
} from '../services/userService';

const { Title } = Typography;
const { Option } = Select;

// Extended User type that includes roles
interface UserWithRoles extends User {
  roles: string[];
}

const UsersManagement: React.FC = () => {
  const { t } = useTranslation();
  const [users, setUsers] = useState<UserWithRoles[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [modalVisible, setModalVisible] = useState<boolean>(false);
  const [modalMode, setModalMode] = useState<'create' | 'edit'>('create');
  const [currentUser, setCurrentUser] = useState<UserWithRoles | null>(null);
  const [form] = Form.useForm();

  // Fetch all users on component mount
  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    try {
      setIsLoading(true);
      const usersData = await userService.getAllUsers();
      
      // Fetch roles for each user
      const usersWithRoles = await Promise.all(
        usersData.map(async (user) => {
          try {
            const roles = await userService.getUserRoles(user.id);
            return { ...user, roles } as UserWithRoles;
          } catch (error) {
            console.error(`Error fetching roles for user ${user.id}:`, error);
            return { ...user, roles: [] } as UserWithRoles;
          }
        })
      );
      
      setUsers(usersWithRoles);
    } catch (error) {
      message.error('Failed to fetch users');
      console.error('Error fetching users:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateUser = () => {
    setModalMode('create');
    setCurrentUser(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEditUser = (user: UserWithRoles) => {
    setModalMode('edit');
    setCurrentUser(user);
    form.setFieldsValue({
      name: user.name,
      email: user.email,
      role: user.roles.length > 0 ? user.roles[0] : 'preparer'
    });
    setModalVisible(true);
  };

  const handleRoleChange = async (userId: string, role: string, operation: 'add' | 'remove') => {
    try {
      setIsLoading(true);
      const roleData: UpdateRoleRequest = { role, operation };
      await userService.updateUserRole(userId, roleData);
      message.success(`User role ${operation === 'add' ? 'added' : 'removed'} successfully`);
      fetchUsers(); // Refresh users list
    } catch (error) {
      message.error(`Failed to ${operation} role: ${error}`);
    } finally {
      setIsLoading(false);
    }
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      
      if (modalMode === 'create') {
        // Create new user
        const userData: CreateUserRequest = {
          name: values.name,
          email: values.email,
          password: values.password,
          role: values.role
        };
        
        await userService.createUser(userData);
        message.success('User created successfully');
      } else if (modalMode === 'edit' && currentUser) {
        // For editing, we only handle role changes as we don't have an update user endpoint
        // Check if role has changed
        const currentRole = currentUser.roles.length > 0 ? currentUser.roles[0] : '';
        if (currentRole !== values.role) {
          // Remove current role if it exists
          if (currentRole) {
            await userService.updateUserRole(currentUser.id, {
              role: currentRole,
              operation: 'remove'
            });
          }
          
          // Add new role
          await userService.updateUserRole(currentUser.id, {
            role: values.role,
            operation: 'add'
          });
          
          message.success('User role updated successfully');
        }
      }
      
      setModalVisible(false);
      fetchUsers(); // Refresh users list
    } catch (error) {
      message.error(`Operation failed: ${error}`);
    }
  };

  const renderRoleTags = (roles: string[]) => {
    if (!roles || roles.length === 0) return <Tag color="default">No Roles</Tag>;
    
    return roles.map(role => {
      let color = 'default';
      if (role === 'admin') color = 'red';
      else if (role === 'approver') color = 'green';
      else if (role === 'preparer') color = 'blue';
      
      return <Tag color={color} key={role}>{role}</Tag>;
    });
  };

  const columns = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'Email',
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: 'Roles',
      dataIndex: 'roles',
      key: 'roles',
      render: (roles: string[]) => renderRoleTags(roles),
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_: any, record: UserWithRoles) => (
        <Space>
          <Tooltip title="Edit User">
            <Button 
              icon={<EditOutlined />} 
              onClick={() => handleEditUser(record)} 
              size="small"
            />
          </Tooltip>
          
          {!record.roles.includes('admin') && (
            <Tooltip title="Make Admin">
              <Button
                type="primary"
                icon={<CheckCircleOutlined />}
                onClick={() => handleRoleChange(record.id, 'admin', 'add')}
                size="small"
                danger
              />
            </Tooltip>
          )}
          
          {record.roles.includes('admin') && (
            <Tooltip title="Remove Admin">
              <Button
                icon={<CloseCircleOutlined />}
                onClick={() => handleRoleChange(record.id, 'admin', 'remove')}
                size="small"
              />
            </Tooltip>
          )}
          
          {!record.roles.includes('preparer') && (
            <Tooltip title="Add Preparer Role">
              <Button
                type="primary"
                onClick={() => handleRoleChange(record.id, 'preparer', 'add')}
                size="small"
              >
                + Preparer
              </Button>
            </Tooltip>
          )}
          
          {!record.roles.includes('approver') && (
            <Tooltip title="Add Approver Role">
              <Button
                type="primary"
                onClick={() => handleRoleChange(record.id, 'approver', 'add')}
                size="small"
              >
                + Approver
              </Button>
            </Tooltip>
          )}
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Card 
        title={<Title level={4}>{t('userManagement.title', 'User Management')}</Title>}
        extra={
          <Button 
            type="primary" 
            icon={<UserAddOutlined />}
            onClick={handleCreateUser}
          >
            {t('userManagement.createUser', 'Create User')}
          </Button>
        }
      >
        {isLoading ? (
          <div style={{ textAlign: 'center', padding: '20px' }}>
            <Spin size="large" />
          </div>
        ) : (
          <Table 
            dataSource={users} 
            columns={columns} 
            rowKey="id"
            pagination={{ pageSize: 10 }}
          />
        )}
      </Card>

      <Modal
        title={modalMode === 'create' 
          ? t('userManagement.createUser', 'Create User') 
          : t('userManagement.editUser', 'Edit User')}
        visible={modalVisible}
        onOk={handleModalOk}
        onCancel={() => setModalVisible(false)}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
        >
          <Form.Item
            name="name"
            label={t('userManagement.name', 'Name')}
            rules={[{ required: true, message: 'Please enter a name' }]}
          >
            <Input />
          </Form.Item>
          
          <Form.Item
            name="email"
            label={t('userManagement.email', 'Email')}
            rules={[
              { required: true, message: 'Please enter an email' },
              { type: 'email', message: 'Please enter a valid email' }
            ]}
          >
            <Input disabled={modalMode === 'edit'} />
          </Form.Item>
          
          {modalMode === 'create' && (
            <Form.Item
              name="password"
              label={t('userManagement.password', 'Password')}
              rules={[{ required: true, message: 'Please enter a password' }]}
            >
              <Input.Password />
            </Form.Item>
          )}
          
          <Form.Item
            name="role"
            label={t('userManagement.role', 'Role')}
            rules={[{ required: true, message: 'Please select a role' }]}
          >
            <Select>
              <Option value="admin">{t('userManagement.roles.admin', 'Admin')}</Option>
              <Option value="preparer">{t('userManagement.roles.preparer', 'Preparer')}</Option>
              <Option value="approver">{t('userManagement.roles.approver', 'Approver')}</Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default UsersManagement; 