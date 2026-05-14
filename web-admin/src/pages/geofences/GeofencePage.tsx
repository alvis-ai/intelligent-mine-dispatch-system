import { useState, useEffect } from 'react';
import {
  Table, Card, Button, Space, Typography, Modal, Form, Input, InputNumber, Select, message, Popconfirm, Tag, Switch,
} from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import apiClient from '../../api/client';

interface GeofenceRecord {
  id: number;
  name: string;
  shape: string;
  center_lat: number;
  center_lon: number;
  radius_m: number;
  fence_type: string;
  max_speed_kmh: number;
  enabled: boolean;
}

const FENCE_TYPE_MAP: Record<string, string> = {
  restricted: '禁区',
  loading: '装载区',
  dumping: '卸载区',
  parking: '停车场',
};

export default function GeofencePage() {
  const [fences, setFences] = useState<GeofenceRecord[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<GeofenceRecord | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();

  const fetch = async () => {
    setLoading(true);
    try {
      const res = await apiClient.get('/api/v1/geofences');
      setFences(res.data.data || []);
    } catch {
      message.error('获取围栏失败');
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
        res = await apiClient.put(`/api/v1/geofences/${editing.id}`, values);
      } else {
        res = await apiClient.post('/api/v1/geofences', values);
      }
      if (res.data.code === 0) {
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
      const res = await apiClient.delete(`/api/v1/geofences/${id}`);
      if (res.data.code === 0) {
        message.success('删除成功');
        fetch();
      }
    } catch {
      message.error('删除失败');
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '名称', dataIndex: 'name', key: 'name' },
    {
      title: '形状', dataIndex: 'shape', key: 'shape', width: 80,
      render: (s: string) => s === 'circle' ? <Tag>圆形</Tag> : <Tag color="blue">多边形</Tag>,
    },
    {
      title: '类型', dataIndex: 'fence_type', key: 'fence_type', width: 100,
      render: (t: string) => FENCE_TYPE_MAP[t] || t,
    },
    { title: '中心经度', dataIndex: 'center_lon', key: 'center_lon', width: 120 },
    { title: '中心纬度', dataIndex: 'center_lat', key: 'center_lat', width: 120 },
    { title: '半径(m)', dataIndex: 'radius_m', key: 'radius_m', width: 100 },
    { title: '限速', dataIndex: 'max_speed_kmh', key: 'max_speed_kmh', width: 80, render: (v: number) => v ? `${v} km/h` : '-' },
    {
      title: '启用', dataIndex: 'enabled', key: 'enabled', width: 80,
      render: (v: boolean) => v ? <Tag color="green">是</Tag> : <Tag color="red">否</Tag>,
    },
    {
      title: '操作', key: 'action', width: 150,
      render: (_: any, record: GeofenceRecord) => (
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
      <Typography.Title level={4}>电子围栏管理</Typography.Title>
      <Card>
        <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'flex-end' }}>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditing(null); form.resetFields(); setModalOpen(true); }}>
            添加围栏
          </Button>
        </div>
        <Table columns={columns} dataSource={fences} rowKey="id" loading={loading} />
      </Card>

      <Modal
        title={editing ? '编辑电子围栏' : '添加电子围栏'}
        open={modalOpen}
        onCancel={() => { setModalOpen(false); setEditing(null); }}
        onOk={handleSubmit}
        confirmLoading={submitting}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item label="名称" name="name" rules={[{ required: true, message: '请输入名称' }]}>
            <Input placeholder="例如：东区禁区" />
          </Form.Item>
          <Form.Item label="形状" name="shape" initialValue="circle">
            <Select options={[
              { value: 'circle', label: '圆形' },
              { value: 'polygon', label: '多边形' },
            ]} />
          </Form.Item>
          <Form.Item label="围栏类型" name="fence_type" initialValue="restricted">
            <Select options={[
              { value: 'restricted', label: '禁区' },
              { value: 'loading', label: '装载区' },
              { value: 'dumping', label: '卸载区' },
              { value: 'parking', label: '停车场' },
            ]} />
          </Form.Item>
          <Space style={{ width: '100%' }} align="start">
            <Form.Item label="中心纬度" name="center_lat">
              <InputNumber style={{ width: 180 }} placeholder="39.9042" step={0.0001} />
            </Form.Item>
            <Form.Item label="中心经度" name="center_lon">
              <InputNumber style={{ width: 180 }} placeholder="116.4074" step={0.0001} />
            </Form.Item>
            <Form.Item label="半径(米)" name="radius_m">
              <InputNumber style={{ width: 120 }} min={0} placeholder="500" />
            </Form.Item>
          </Space>
          <Form.Item label="限速 (km/h)" name="max_speed_kmh">
            <InputNumber min={0} placeholder="0=不限速" />
          </Form.Item>
          <Form.Item label="启用" name="enabled" valuePropName="checked" initialValue={true}>
            <Switch />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
