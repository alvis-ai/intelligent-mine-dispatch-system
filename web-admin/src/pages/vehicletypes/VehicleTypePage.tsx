import { useState, useEffect } from 'react';
import {
  Table, Card, Button, Space, Typography, Modal, Form, Input, InputNumber, message, Popconfirm,
} from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import apiClient from '../../api/client';

interface VehicleTypeRecord {
  id: number;
  name: string;
  description: string;
  icon: string;
  capacity: number;
  weight: number;
}

export default function VehicleTypePage() {
  const [types, setTypes] = useState<VehicleTypeRecord[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<VehicleTypeRecord | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();

  const fetch = async () => {
    setLoading(true);
    try {
      const res = await apiClient.get('/api/v1/vehicle-types');
      setTypes(res.data.data || []);
    } catch {
      message.error('获取失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetch(); }, []);

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setSubmitting(true);
      let res: any;
      if (editing) {
        res = await apiClient.put(`/api/v1/vehicle-types/${editing.id}`, values);
      } else {
        res = await apiClient.post('/api/v1/vehicle-types', values);
      }
      if (res.data.code === 0 || res.data.message === 'success') {
        message.success(editing ? '更新成功' : '添加成功');
        setModalOpen(false);
        form.resetFields();
        setEditing(null);
        fetch();
      } else {
        message.error(res.data.message || '操作失败');
      }
    } catch (err: any) {
      if (err?.errorFields) return;
      message.error('操作失败');
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async (id: number) => {
    try {
      const res = await apiClient.delete(`/api/v1/vehicle-types/${id}`);
      if (res.data.code === 0 || res.data.message === 'success') {
        message.success('删除成功');
        fetch();
      }
    } catch {
      message.error('删除失败');
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 100 },
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '描述', dataIndex: 'description', key: 'description' },
    { title: '载重(t)', dataIndex: 'capacity', key: 'capacity', render: (v: number) => v || '-' },
    { title: '自重(t)', dataIndex: 'weight', key: 'weight', render: (v: number) => v || '-' },
    {
      title: '操作', key: 'action',
      render: (_: any, record: VehicleTypeRecord) => (
        <Space>
          <a onClick={() => { setEditing(record); form.setFieldsValue(record); setModalOpen(true); }}>编辑</a>
          <Popconfirm title="确认删除？" onConfirm={() => handleDelete(record.id)}>
            <a style={{ color: 'red' }}>删除</a>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Typography.Title level={4}>车辆类型管理</Typography.Title>
      <Card>
        <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'flex-end' }}>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditing(null); form.resetFields(); setModalOpen(true); }}>
            添加类型
          </Button>
        </div>
        <Table columns={columns} dataSource={types} rowKey="id" loading={loading} />
      </Card>

      <Modal
        title={editing ? '编辑车辆类型' : '添加车辆类型'}
        open={modalOpen}
        onCancel={() => { setModalOpen(false); setEditing(null); }}
        onOk={handleSubmit}
        confirmLoading={submitting}
      >
        <Form form={form} layout="vertical">
          <Form.Item label="名称" name="name" rules={[{ required: true, message: '请输入名称' }]}>
            <Input placeholder="例如：矿用卡车" />
          </Form.Item>
          <Form.Item label="描述" name="description">
            <Input placeholder="简短描述" />
          </Form.Item>
          <Form.Item label="载重 (吨)" name="capacity">
            <InputNumber style={{ width: '100%' }} min={0} placeholder="0" />
          </Form.Item>
          <Form.Item label="自重 (吨)" name="weight">
            <InputNumber style={{ width: '100%' }} min={0} placeholder="0" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
