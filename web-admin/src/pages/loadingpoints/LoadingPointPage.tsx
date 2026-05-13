import { useState, useEffect } from 'react';
import {
  Table, Card, Button, Tag, Space, Typography, Modal, Form, Input, Select, message, Popconfirm,
} from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import apiClient from '../../api/client';

interface LoadingPointRecord {
  id: number;
  name: string;
  type: string;
  latitude: number;
  longitude: number;
  material: string;
  status: number;
}

const TYPE_OPTIONS = [
  { value: 'loading', label: '装载点' },
  { value: 'dumping', label: '卸载点' },
];

export default function LoadingPointPage() {
  const [points, setPoints] = useState<LoadingPointRecord[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<LoadingPointRecord | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();

  const fetch = async () => {
    setLoading(true);
    try {
      const res = await apiClient.get('/api/v1/loading-points');
      setPoints(res.data.data || []);
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
      const body = {
        name: values.name,
        type: values.type,
        material: values.material || '',
        latitude: values.latitude || 0,
        longitude: values.longitude || 0,
      };
      let res: any;
      if (editing) {
        res = await apiClient.put(`/api/v1/loading-points/${editing.id}`, body);
      } else {
        res = await apiClient.post('/api/v1/loading-points', body);
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
      const res = await apiClient.delete(`/api/v1/loading-points/${id}`);
      if (res.data.code === 0 || res.data.message === 'success') {
        message.success('删除成功');
        fetch();
      }
    } catch {
      message.error('删除失败');
    }
  };

  const openEdit = (p: LoadingPointRecord) => {
    setEditing(p);
    form.setFieldsValue(p);
    setModalOpen(true);
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 100 },
    { title: '名称', dataIndex: 'name', key: 'name' },
    {
      title: '类型', dataIndex: 'type', key: 'type',
      render: (t: string) => t === 'loading' ? <Tag color="blue">装载点</Tag> : <Tag color="orange">卸载点</Tag>,
    },
    { title: '物料', dataIndex: 'material', key: 'material' },
    { title: '经度', dataIndex: 'longitude', key: 'longitude', render: (v: number) => v.toFixed(4) },
    { title: '纬度', dataIndex: 'latitude', key: 'latitude', render: (v: number) => v.toFixed(4) },
    {
      title: '状态', dataIndex: 'status', key: 'status',
      render: (s: number) => <Tag color={s === 1 ? 'green' : 'default'}>{s === 1 ? '启用' : '停用'}</Tag>,
    },
    {
      title: '操作', key: 'action',
      render: (_: any, record: LoadingPointRecord) => (
        <Space>
          <a onClick={() => openEdit(record)}>编辑</a>
          <Popconfirm title="确认删除？" onConfirm={() => handleDelete(record.id)}>
            <a style={{ color: 'red' }}>删除</a>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Typography.Title level={4}>装载点管理</Typography.Title>
      <Card>
        <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'flex-end' }}>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditing(null); form.resetFields(); setModalOpen(true); }}>
            添加点
          </Button>
        </div>
        <Table columns={columns} dataSource={points} rowKey="id" loading={loading} />
      </Card>

      <Modal
        title={editing ? '编辑点' : '添加点'}
        open={modalOpen}
        onCancel={() => { setModalOpen(false); setEditing(null); }}
        onOk={handleSubmit}
        confirmLoading={submitting}
      >
        <Form form={form} layout="vertical" initialValues={{ type: 'loading' }}>
          <Form.Item label="名称" name="name" rules={[{ required: true, message: '请输入名称' }]}>
            <Input placeholder="例如：装载点A" />
          </Form.Item>
          <Form.Item label="类型" name="type" rules={[{ required: true }]}>
            <Select options={TYPE_OPTIONS} />
          </Form.Item>
          <Form.Item label="物料" name="material">
            <Input placeholder="例如：矿石" />
          </Form.Item>
          <Form.Item label="经度" name="longitude">
            <Input placeholder="116.4000" type="number" />
          </Form.Item>
          <Form.Item label="纬度" name="latitude">
            <Input placeholder="39.9000" type="number" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
