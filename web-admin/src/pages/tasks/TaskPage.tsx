import { useState, useEffect, useCallback } from 'react';
import {
  Table, Card, Tag, Button, Typography, Modal, Select, Form, message, Space, Popconfirm,
} from 'antd';
import { PlusOutlined, CheckCircleOutlined, CloseCircleOutlined, ReloadOutlined, BulbOutlined } from '@ant-design/icons';
import { Drawer, Descriptions } from 'antd';
import apiClient from '../../api/client';
import { fetchSuggestions } from '../../services/aiService';

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
  { value: 'genetic_algorithm', label: '遗传算法优化' },
];

export default function TaskPage() {
  const [tasks, setTasks] = useState<TaskRecord[]>([]);
  const [vehicles, setVehicles] = useState<{ value: number; label: string }[]>([]);
  const [loadPoints, setLoadPoints] = useState<{ value: number; label: string }[]>([]);
  const [dumpPoints, setDumpPoints] = useState<{ value: number; label: string }[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [suggestOpen, setSuggestOpen] = useState(false);
  const [suggestions, setSuggestions] = useState<any[]>([]);
  const [suggestLoading, setSuggestLoading] = useState(false);
  const [form] = Form.useForm();

  const fetchTasks = useCallback(async () => {
    setLoading(true);
    try {
      const res = await apiClient.get('/api/v1/dispatch/tasks');
      setTasks(res.data.data || []);
    } catch {
      message.error('获取调度任务失败');
    } finally {
      setLoading(false);
    }
  }, []);

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
  }, [fetchTasks]);

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

  const handleUpdateStatus = async (id: number, status: string) => {
    try {
      const endpoint = status === 'completed'
        ? `/api/v1/dispatch/tasks/${id}/complete`
        : `/api/v1/dispatch/tasks/${id}/cancel`;
      const res = await apiClient.post(endpoint);
      if (res.data.code === 0 || res.data.message === 'success') {
        message.success(status === 'completed' ? '任务已完成' : '任务已取消');
        fetchTasks();
      } else {
        message.error(res.data.message || '操作失败');
      }
    } catch {
      message.error('操作失败');
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
    {
      title: '操作', key: 'action', width: 160,
      render: (_: any, record: TaskRecord) => (
        <Space>
          {record.status === 'active' && (
            <Popconfirm title="确认完成任务？" onConfirm={() => handleUpdateStatus(record.id, 'completed')}>
              <a style={{ color: 'green' }}><CheckCircleOutlined /> 完成</a>
            </Popconfirm>
          )}
          {(record.status === 'pending' || record.status === 'active') && (
            <Popconfirm title="确认取消任务？" onConfirm={() => handleUpdateStatus(record.id, 'cancelled')}>
              <a style={{ color: 'red' }}><CloseCircleOutlined /> 取消</a>
            </Popconfirm>
          )}
          {record.status === 'completed' && (
            <Tag color="success" style={{ margin: 0 }}>已完成</Tag>
          )}
          {record.status === 'cancelled' && (
            <Tag color="error" style={{ margin: 0 }}>已取消</Tag>
          )}
        </Space>
      ),
    },
  ];

  const activeTasks = tasks.filter((t) => t.status === 'active').length;
  const pendingTasks = tasks.filter((t) => t.status === 'pending').length;

  const openSuggestions = async () => {
    setSuggestOpen(true);
    setSuggestLoading(true);
    try {
      const vRes = await apiClient.get('/api/v1/vehicles');
      const vehicles = (vRes.data.data || []).slice(0, 20).map((v: any) => ({
        vehicle_id: v.id,
        latitude: v.latitude || 39.9,
        longitude: v.longitude || 116.4,
        active_task_count: 0,
      }));
      const lpRes = await apiClient.get('/api/v1/loading-points');
      const points = (lpRes.data.data || []).filter((p: any) => p.type === 'loading');
      const candidates = points.map((p: any) => ({
        load_point_id: p.id,
        dump_point_id: points.find((x: any) => x.type === 'dumping')?.id || 3,
      }));
      if (candidates.length === 0) {
        candidates.push({ load_point_id: 1, dump_point_id: 3 });
      }
      const sug = await fetchSuggestions(vehicles, candidates, 1);
      setSuggestions(sug || []);
    } catch {
      message.error('获取 AI 建议失败');
    } finally {
      setSuggestLoading(false);
    }
  };

  return (
    <div>
      <Typography.Title level={4}>调度任务</Typography.Title>
      <Card>
        <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Space>
            <Tag color="processing">进行中 {activeTasks}</Tag>
            <Tag color="default">待分配 {pendingTasks}</Tag>
            <Tag color="success">已完成 {tasks.filter((t) => t.status === 'completed').length}</Tag>
          </Space>
          <Space>
            <Button icon={<ReloadOutlined />} onClick={fetchTasks}>刷新</Button>
            <Button icon={<BulbOutlined />} onClick={openSuggestions}>AI 建议</Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
              创建调度任务
            </Button>
          </Space>
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

      <Drawer
        title="AI 调度建议"
        placement="right"
        width={480}
        open={suggestOpen}
        onClose={() => setSuggestOpen(false)}
        loading={suggestLoading}
      >
        {suggestions.length === 0 && !suggestLoading ? (
          <Typography.Text type="secondary">暂无可用建议</Typography.Text>
        ) : (
          suggestions.map((s: any, i: number) => (
            <Card key={i} size="small" style={{ marginBottom: 8 }}>
              <Descriptions column={1} size="small">
                <Descriptions.Item label="车辆">#{s.vehicle_id}</Descriptions.Item>
                <Descriptions.Item label="装载点">#{s.load_point_id}</Descriptions.Item>
                <Descriptions.Item label="卸载点">#{s.dump_point_id}</Descriptions.Item>
                <Descriptions.Item label="评分">
                  <Tag color={s.score > 0.8 ? 'green' : s.score > 0.6 ? 'orange' : 'red'}>
                    {(s.score * 100).toFixed(0)}%
                  </Tag>
                </Descriptions.Item>
                <Descriptions.Item label="距离">
                  {(s.estimated_distance_m / 1000).toFixed(1)} km
                </Descriptions.Item>
                <Descriptions.Item label="时长">
                  {Math.round(s.estimated_duration_s / 60)} min
                </Descriptions.Item>
                <Descriptions.Item label="原因">{s.reason}</Descriptions.Item>
              </Descriptions>
            </Card>
          ))
        )}
      </Drawer>
    </div>
  );
}
