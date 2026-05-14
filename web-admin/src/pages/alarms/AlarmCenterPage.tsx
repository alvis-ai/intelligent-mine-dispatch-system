import { useState, useEffect } from 'react';
import {
  Table, Card, Tag, Button, Typography, Space, Select, message, Badge,
} from 'antd';
import { CheckCircleOutlined, ReloadOutlined } from '@ant-design/icons';
import apiClient from '../../api/client';

interface AlarmEvent {
  id: number;
  vehicle_id: number;
  vehicle_plate: string;
  alarm_type: string;
  severity: string;
  message: string;
  latitude: number;
  longitude: number;
  speed: number;
  acknowledged: boolean;
  acknowledged_by: string;
  acknowledged_at: string;
  created_at: string;
}

const SEVERITY_CONFIG: Record<string, { color: string; label: string }> = {
  critical: { color: 'red', label: '严重' },
  warning: { color: 'orange', label: '警告' },
  info: { color: 'blue', label: '信息' },
};

const ALARM_TYPE_MAP: Record<string, string> = {
  geofence_entry: '禁区闯入',
  speeding: '超速',
  offline: '离线',
  deviation: '偏离路线',
};

export default function AlarmCenterPage() {
  const [alarms, setAlarms] = useState<AlarmEvent[]>([]);
  const [loading, setLoading] = useState(false);
  const [severity, setSeverity] = useState<string>('');
  const [unackOnly, setUnackOnly] = useState(false);

  const fetchAlarms = async () => {
    setLoading(true);
    try {
      const params: any = { page: 1, page_size: 50 };
      if (severity) params.severity = severity;
      if (unackOnly) params.unacknowledged_only = 'true';
      const res = await apiClient.get('/api/v1/alarms', { params });
      setAlarms(res.data.data || []);
    } catch {
      message.error('获取告警失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchAlarms(); }, [severity, unackOnly]);

  const handleAcknowledge = async (id: number) => {
    try {
      const res = await apiClient.post(`/api/v1/alarms/${id}/acknowledge`, {
        acknowledged_by: 'admin',
      });
      if (res.data.code === 0 || res.data.message === 'success') {
        message.success('已确认');
        fetchAlarms();
      }
    } catch {
      message.error('确认失败');
    }
  };

  const columns = [
    {
      title: '级别', dataIndex: 'severity', key: 'severity', width: 80,
      render: (s: string) => {
        const cfg = SEVERITY_CONFIG[s] || { color: 'default', label: s };
        return <Tag color={cfg.color}>{cfg.label}</Tag>;
      },
    },
    { title: '车辆', dataIndex: 'vehicle_plate', key: 'vehicle_plate', width: 120 },
    {
      title: '类型', dataIndex: 'alarm_type', key: 'alarm_type', width: 100,
      render: (t: string) => ALARM_TYPE_MAP[t] || t,
    },
    { title: '告警信息', dataIndex: 'message', key: 'message', ellipsis: true },
    { title: '速度', dataIndex: 'speed', key: 'speed', width: 80, render: (v: number) => v ? `${v} km/h` : '-' },
    {
      title: '时间', dataIndex: 'created_at', key: 'created_at', width: 180,
      render: (t: string) => t ? new Date(t).toLocaleString('zh-CN') : '-',
    },
    {
      title: '状态', key: 'status', width: 80,
      render: (_: any, r: AlarmEvent) => r.acknowledged
        ? <Tag icon={<CheckCircleOutlined />} color="green">已确认</Tag>
        : <Badge status="processing" text="未确认" />,
    },
    {
      title: '操作', key: 'action', width: 80,
      render: (_: any, r: AlarmEvent) => !r.acknowledged && (
        <Button size="small" type="link" onClick={() => handleAcknowledge(r.id)}>
          确认
        </Button>
      ),
    },
  ];

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Typography.Title level={4} style={{ margin: 0 }}>告警中心</Typography.Title>
        <Space>
          <Select
            placeholder="告警级别"
            allowClear
            style={{ width: 120 }}
            value={severity || undefined}
            onChange={(v) => setSeverity(v || '')}
            options={[
              { value: 'critical', label: '严重' },
              { value: 'warning', label: '警告' },
              { value: 'info', label: '信息' },
            ]}
          />
          <Select
            placeholder="确认状态"
            allowClear
            style={{ width: 120 }}
            value={unackOnly ? 'unack' : undefined}
            onChange={(v) => setUnackOnly(v === 'unack')}
            options={[{ value: 'unack', label: '未确认' }]}
          />
          <Button icon={<ReloadOutlined />} onClick={fetchAlarms}>刷新</Button>
        </Space>
      </div>
      <Card>
        <Table
          columns={columns}
          dataSource={alarms}
          rowKey="id"
          loading={loading}
          pagination={{ pageSize: 20 }}
          size="small"
        />
      </Card>
    </div>
  );
}
