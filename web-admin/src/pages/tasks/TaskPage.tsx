import { useState, useEffect } from 'react';
import {
  Table, Card, Tag, Button, Typography, Modal, Select, Form, message,
} from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import apiClient from '../../api/client';

interface TaskRecord {
  id: number;
  vehicle_id: number;
  load_point_id: number;
  dump_point_id: number;
  material: string;
  status: string;
  algorithm: string;
}

const STATUS_MAP: Record<string, { color: string; text: string }> = {
  pending: { color: 'default', text: '待分配' },
  active: { color: 'processing', text: '进行中' },
  completed: { color: 'success', text: '已完成' },
  cancelled: { color: 'error', text: '已取消' },
};

const ALGORITHMS = [
  { value: 'fifo', label: 'FIFO (先进先出)' },
  { value: 'nearest_first', label: '最近优先' },
  { value: 'weighted_round_robin', label: '加权轮询' },
];

export default function TaskPage() {
  const [tasks, setTasks] = useState<TaskRecord[]>([]);
  const [vehicles, setVehicles] = useState<{ value: number; label: string }[]>([]);
  const [loadPoints, setLoadPoints] = useState<{ value: number; label: string }[]>([]);
  const [dumpPoints, setDumpPoints] = useState<{ value: number; label: string }[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();

  const fetchTasks = async () => {
    setLoading(true);
    try {
      const res = await apiClient.get('/api/v1/dispatch/tasks');
      setTasks(res.data.data || []);
    } catch {
      message.error('获取调度任务失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchVehicles = async () => {
    try {
      const res = await apiClient.get('/api/v1/vehicles');
      const list = (res.data.data || []).map((v: any) => ({
        value: v.id,
        label: v.plate,
      }));
      setVehicles(list);
    } catch {
      // ignore
    }
  };

  const fetchPoints = async () => {
    try {
      const res = await apiClient.get('/api/v1/loading-points');
      const all: any[] = res.data.data || [];
      setLoadPoints(
        all.filter((p) => p.type === 'loading').map((p) => ({ value: p.id, label: p.name }))
      );
      setDumpPoints(
        all.filter((p) => p.type === 'dumping').map((p) => ({ value: p.id, label: p.name }))
      );
    } catch {
      // ignore
    }
  };

  useEffect(() => {
    fetchTasks();
    fetchVehicles();
    fetchPoints();
  }, []);

  const handleCreate = async () => {
    try {
      const values = await form.validateFields();
      setSubmitting(true);
      const res = await apiClient.post('/api/v1/dispatch/assign', {
        vehicle_id: values.vehicle_id,
        load_point_id: values.load_point_id,
        dump_point_id: values.dump_point_id,
        algorithm: values.algorithm || 'fifo',
      });
      if (res.data.code === 0 || res.data.message === 'success') {
        message.success('调度任务已创建');
        setModalOpen(false);
        form.resetFields();
        fetchTasks();
      } else {
        message.error(res.data.message || '创建失败');
      }
    } catch (err: any) {
      if (err?.response?.data?.message) {
        message.error(err.response.data.message);
      } else if (!err?.errorFields) {
        message.error('创建失败');
      }
    } finally {
      setSubmitting(false);
    }
  };

  const pointLabel = (id: number) =>
    [...loadPoints, ...dumpPoints].find((p) => p.value === id)?.label || `点${id}`;

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 100 },
    {
      title: '车辆', dataIndex: 'vehicle_id', key: 'vehicle_id',
      render: (id: number) => vehicles.find((v) => v.value === id)?.label || `车辆#${id}`,
    },
    {
      title: '装载点', dataIndex: 'load_point_id', key: 'load_point_id',
      render: (id: number) => pointLabel(id),
    },
    {
      title: '卸载点', dataIndex: 'dump_point_id', key: 'dump_point_id',
      render: (id: number) => pointLabel(id),
    },
    { title: '物料', dataIndex: 'material', key: 'material' },
    {
      title: '状态', dataIndex: 'status', key: 'status',
      render: (s: string) => {
        const m = STATUS_MAP[s] || { color: 'default', text: s };
        return <Tag color={m.color}>{m.text}</Tag>;
      },
    },
    { title: '算法', dataIndex: 'algorithm', key: 'algorithm' },
  ];

  return (
    <div>
      <Typography.Title level={4}>调度任务</Typography.Title>
      <Card>
        <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'flex-end' }}>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
            创建调度任务
          </Button>
        </div>
        <Table columns={columns} dataSource={tasks} rowKey="id" loading={loading} />
      </Card>

      <Modal
        title="创建调度任务"
        open={modalOpen}
        onCancel={() => { setModalOpen(false); form.resetFields(); }}
        onOk={handleCreate}
        confirmLoading={submitting}
      >
        <Form form={form} layout="vertical">
          <Form.Item label="车辆" name="vehicle_id" rules={[{ required: true, message: '请选择车辆' }]}>
            <Select placeholder="选择车辆" options={vehicles} />
          </Form.Item>
          <Form.Item label="装载点" name="load_point_id" rules={[{ required: true, message: '请选择装载点' }]}>
            <Select placeholder="选择装载点" options={loadPoints} />
          </Form.Item>
          <Form.Item label="卸载点" name="dump_point_id" rules={[{ required: true, message: '请选择卸载点' }]}>
            <Select placeholder="选择卸载点" options={dumpPoints} />
          </Form.Item>
          <Form.Item label="调度算法" name="algorithm" initialValue="fifo">
            <Select options={ALGORITHMS} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
