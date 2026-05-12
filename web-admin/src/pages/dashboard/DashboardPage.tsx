import { useState, useEffect } from 'react';
import { Row, Col, Card, Statistic, Table, Typography } from 'antd';
import {
  CarOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
  WarningOutlined,
} from '@ant-design/icons';
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

export default function DashboardPage() {
  const [vehicles, setVehicles] = useState<any[]>([]);
  const [tasks, setTasks] = useState<TaskRecord[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetch = async () => {
      try {
        const [vRes, tRes] = await Promise.all([
          apiClient.get('/api/v1/vehicles'),
          apiClient.get('/api/v1/dispatch/tasks'),
        ]);
        setVehicles(vRes.data.data || []);
        setTasks(tRes.data.data || []);
      } catch {
        // use defaults
      } finally {
        setLoading(false);
      }
    };
    fetch();
  }, []);

  const onlineCount = vehicles.length;
  const activeCount = tasks.filter((t) => t.status === 'active' || t.status === 'pending').length;
  const completedCount = tasks.filter((t) => t.status === 'completed').length;

  const recentColumns = [
    { title: 'ID', dataIndex: 'id', key: 'id' },
    { title: '状态', dataIndex: 'status', key: 'status' },
    { title: '物料', dataIndex: 'material', key: 'material' },
    { title: '算法', dataIndex: 'algorithm', key: 'algorithm' },
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
            <Statistic title="告警" value={0} prefix={<WarningOutlined />} valueStyle={{ color: '#cf1322' }} />
          </Card>
        </Col>
      </Row>
      <Card title="最近调度任务" style={{ marginTop: 16 }}>
        <Table
          columns={recentColumns}
          dataSource={tasks.slice(0, 10)}
          rowKey="id"
          loading={loading}
          pagination={false}
        />
      </Card>
    </div>
  );
}
