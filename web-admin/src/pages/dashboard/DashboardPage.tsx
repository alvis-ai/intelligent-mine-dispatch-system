import { useState, useEffect } from 'react';
import { Row, Col, Card, Statistic, Table, Typography, Tag } from 'antd';
import {
  CarOutlined, CheckCircleOutlined, ClockCircleOutlined, WarningOutlined,
  FireOutlined,
} from '@ant-design/icons';
import apiClient from '../../api/client';

interface TaskRecord {
  id: number;
  vehicle_id: number;
  status: string;
  algorithm: string;
}

const STATUS_TAG: Record<string, { color: string; text: string }> = {
  pending: { color: 'default', text: '待分配' },
  active: { color: 'processing', text: '进行中' },
  completed: { color: 'success', text: '已完成' },
  cancelled: { color: 'error', text: '已取消' },
};

export default function DashboardPage() {
  const [vehicles, setVehicles] = useState<any[]>([]);
  const [tasks, setTasks] = useState<TaskRecord[]>([]);
  const [alarmStats, setAlarmStats] = useState<any>({});
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetch = async () => {
      try {
        const [vRes, tRes, aRes] = await Promise.all([
          apiClient.get('/api/v1/vehicles'),
          apiClient.get('/api/v1/dispatch/tasks'),
          apiClient.get('/api/v1/alarms/stats'),
        ]);
        setVehicles(vRes.data.data || []);
        setTasks(tRes.data.data || []);
        setAlarmStats(aRes.data.data || {});
      } catch {
        // use defaults
      } finally {
        setLoading(false);
      }
    };
    fetch();
    const interval = setInterval(fetch, 10000);
    return () => clearInterval(interval);
  }, []);

  const onlineCount = vehicles.length;
  const activeCount = tasks.filter((t) => t.status === 'active' || t.status === 'pending').length;
  const completedCount = tasks.filter((t) => t.status === 'completed').length;
  const unackCritical = alarmStats.unacknowledged_critical || 0;
  const unackWarning = alarmStats.unacknowledged_warning || 0;

  const recentColumns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    {
      title: '状态', dataIndex: 'status', key: 'status', width: 100,
      render: (s: string) => {
        const m = STATUS_TAG[s] || { color: 'default', text: s };
        return <Tag color={m.color}>{m.text}</Tag>;
      },
    },
    { title: '算法', dataIndex: 'algorithm', key: 'algorithm', width: 120 },
  ];

  return (
    <div>
      <Typography.Title level={4}>调度看板</Typography.Title>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic title="在线车辆" value={onlineCount} prefix={<CarOutlined />} suffix="辆" />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic title="进行中任务" value={activeCount} prefix={<ClockCircleOutlined />} suffix="个" />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic title="已完成" value={completedCount} prefix={<CheckCircleOutlined />} suffix="个" />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="未确认告警"
              value={unackCritical + unackWarning}
              prefix={<WarningOutlined />}
              valueStyle={{ color: unackCritical > 0 ? '#cf1322' : '#faad14' }}
              suffix={
                <span style={{ fontSize: 14 }}>
                  {unackCritical > 0 && <Tag color="red" style={{ marginLeft: 4 }}>{unackCritical}严重</Tag>}
                  {unackWarning > 0 && <Tag color="orange">{unackWarning}警告</Tag>}
                </span>
              }
            />
          </Card>
        </Col>
      </Row>
      {unackCritical > 0 && (
        <Row gutter={[16, 16]} style={{ marginTop: 8 }}>
          <Col span={24}>
            <Card size="small" style={{ borderLeft: '4px solid #cf1322', background: '#fff2f0' }}>
              <Typography.Text type="danger" strong>
                <FireOutlined /> 存在 {unackCritical} 条严重告警未处理，请及时处理！
              </Typography.Text>
            </Card>
          </Col>
        </Row>
      )}
      <Card title="最近调度任务" style={{ marginTop: 16 }}>
        <Table
          columns={recentColumns}
          dataSource={tasks.slice(0, 10)}
          rowKey="id"
          loading={loading}
          pagination={false}
          size="small"
        />
      </Card>
    </div>
  );
}
