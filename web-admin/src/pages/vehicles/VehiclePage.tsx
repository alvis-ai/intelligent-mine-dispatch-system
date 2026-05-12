import { useState, useEffect } from 'react';
import {
  Table, Button, Card, Tag, Space, Typography, Modal, Form, Input, Select, message,
} from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import apiClient from '../../api/client';

interface VehicleRecord {
  id: number;
  plate: string;
  type: number;
  typeName: string;
  status: number;
  statusName: string;
  fuelLevel: number;
}

const TYPE_MAP: Record<number, string> = { 1: '矿用卡车', 2: '挖掘机', 3: '装载机' };
const STATUS_MAP: Record<number, string> = { 1: '空闲', 2: '作业中', 3: '维修中', 4: '离线' };
const STATUS_COLORS: Record<number, string> = { 1: 'green', 2: 'blue', 3: 'orange', 4: 'default' };

export default function VehiclePage() {
  const [vehicles, setVehicles] = useState<VehicleRecord[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();

  const fetchVehicles = async () => {
    setLoading(true);
    try {
      const res = await apiClient.get('/api/v1/vehicles');
      const list: VehicleRecord[] = (res.data.data || []).map((v: any) => ({
        id: v.id,
        plate: v.plate,
        type: v.type,
        typeName: TYPE_MAP[v.type] || '未知',
        status: v.status,
        statusName: STATUS_MAP[v.status] || '未知',
        fuelLevel: v.fuel_level ?? 100,
      }));
      setVehicles(list);
    } catch {
      message.error('获取车辆列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchVehicles(); }, []);

  const handleAdd = async () => {
    try {
      const values = await form.validateFields();
      setSubmitting(true);
      const res = await apiClient.post('/api/v1/vehicles', {
        plate: values.plate,
        type: values.type,
        mineId: 1,
      });
      if (res.data.code === 0 || res.data.message === 'success') {
        message.success('添加成功');
        setModalOpen(false);
        form.resetFields();
        fetchVehicles();
      } else {
        message.error(res.data.message || '添加失败');
      }
    } catch (err: any) {
      if (err?.response?.data?.message) {
        message.error(err.response.data.message);
      } else if (err?.errorFields) {
        // form validation error, antd shows inline
      } else {
        message.error('添加失败');
      }
    } finally {
      setSubmitting(false);
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 100 },
    { title: '车牌/编号', dataIndex: 'plate', key: 'plate' },
    { title: '类型', dataIndex: 'typeName', key: 'type' },
    {
      title: '状态', dataIndex: 'statusName', key: 'status',
      render: (_: string, record: VehicleRecord) => (
        <Tag color={STATUS_COLORS[record.status]}>{record.statusName}</Tag>
      ),
    },
    {
      title: '油量', dataIndex: 'fuelLevel', key: 'fuelLevel',
      render: (v: number) => `${v}%`,
    },
    {
      title: '操作', key: 'action',
      render: () => (
        <Space>
          <a>编辑</a>
          <a>调度</a>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Typography.Title level={4}>车辆管理</Typography.Title>
      <Card>
        <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'flex-end' }}>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
            添加车辆
          </Button>
        </div>
        <Table columns={columns} dataSource={vehicles} rowKey="id" loading={loading} />
      </Card>

      <Modal
        title="添加车辆"
        open={modalOpen}
        onCancel={() => { setModalOpen(false); form.resetFields(); }}
        onOk={handleAdd}
        confirmLoading={submitting}
      >
        <Form form={form} layout="vertical">
          <Form.Item label="车牌/编号" name="plate" rules={[{ required: true, message: '请输入车牌' }]}>
            <Input placeholder="请输入车牌" />
          </Form.Item>
          <Form.Item label="车辆类型" name="type" rules={[{ required: true, message: '请选择车辆类型' }]}>
            <Select
              placeholder="请选择"
              options={[
                { value: 1, label: '矿用卡车' },
                { value: 2, label: '挖掘机' },
                { value: 3, label: '装载机' },
              ]}
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
