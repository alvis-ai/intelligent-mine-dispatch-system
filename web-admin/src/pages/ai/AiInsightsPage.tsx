import { useState, useEffect } from 'react';
import { Row, Col, Card, Table, Tag, Typography, Spin } from 'antd';
import {
  fetchCongestion, fetchDemand, fetchSuggestions,
} from '../../services/aiService';
import type {
  EdgeCongestion, LoadingPointDemand, AISuggestion,
} from '../../services/aiService';
import { getCurrentMineId } from '../../utils/mineContext';

const CONGESTION_COLOR = (score: number) => {
  if (score > 0.6) return 'red';
  if (score > 0.3) return 'orange';
  return 'green';
};

const CONGESTION_TEXT = (score: number) => {
  if (score > 0.6) return '拥堵';
  if (score > 0.3) return '缓行';
  return '畅通';
};

export default function AiInsightsPage() {
  const [congestion, setCongestion] = useState<EdgeCongestion[]>([]);
  const [demand, setDemand] = useState<LoadingPointDemand[]>([]);
  const [suggestions, setSuggestions] = useState<AISuggestion[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetch = async () => {
      setLoading(true);
      try {
        const mineId = getCurrentMineId();
        const [cong, dem] = await Promise.all([
          fetchCongestion(mineId, 60),
          fetchDemand(mineId),
        ]);
        setCongestion(cong);
        setDemand(dem);

        // Fetch AI suggestions for demo
        const sug = await fetchSuggestions(
          [{ vehicle_id: 1, latitude: 39.895, longitude: 116.405, active_task_count: 0 }],
          [{ load_point_id: 1, dump_point_id: 3 }],
        );
        setSuggestions(sug);
      } catch {
        // ignore
      } finally {
        setLoading(false);
      }
    };
    fetch();
    const interval = setInterval(fetch, 30000);
    return () => clearInterval(interval);
  }, []);

  const loadingPoints = demand.filter((d) => d.point_type === 'loading');
  const dumpingPoints = demand.filter((d) => d.point_type === 'dumping');

  const congestionColumns = [
    { title: '边ID', dataIndex: 'edge_id', key: 'edge_id', width: 80 },
    { title: '起点节点', dataIndex: 'from_node_id', key: 'from_node_id', width: 100 },
    { title: '终点节点', dataIndex: 'to_node_id', key: 'to_node_id', width: 100 },
    {
      title: '状态', key: 'status', width: 80,
      render: (_: any, r: EdgeCongestion) => (
        <Tag color={CONGESTION_COLOR(r.congestion_score)}>{CONGESTION_TEXT(r.congestion_score)}</Tag>
      ),
    },
    {
      title: '拥堵分', dataIndex: 'congestion_score', key: 'congestion_score', width: 90,
      render: (v: number) => (v * 100).toFixed(0) + '%',
    },
    {
      title: '预测速度', key: 'speed', width: 100,
      render: (_: any, r: EdgeCongestion) => `${r.predicted_speed_kmh.toFixed(1)} km/h`,
    },
    { title: '车辆数', dataIndex: 'predicted_vehicle_count', key: 'predicted_vehicle_count', width: 80 },
    {
      title: '置信度', dataIndex: 'confidence', key: 'confidence', width: 80,
      render: (v: number) => (v * 100).toFixed(0) + '%',
    },
  ];

  const demandColumns = [
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '类型', dataIndex: 'point_type', key: 'point_type' },
    { title: '物料', dataIndex: 'material', key: 'material' },
    {
      title: '需求分', dataIndex: 'demand_score', key: 'demand_score',
      render: (v: number) => (
        <Tag color={v > 0.5 ? 'volcano' : v > 0.2 ? 'gold' : 'green'}>
          {(v * 100).toFixed(0)}%
        </Tag>
      ),
    },
    { title: '待处理任务', dataIndex: 'pending_task_count', key: 'pending_task_count' },
    { title: '活跃车辆', dataIndex: 'active_vehicle_count', key: 'active_vehicle_count' },
    {
      title: '置信度', dataIndex: 'confidence', key: 'confidence',
      render: (v: number) => (v * 100).toFixed(0) + '%',
    },
  ];

  const suggestionColumns = [
    { title: '车辆ID', dataIndex: 'vehicle_id', key: 'vehicle_id' },
    { title: '装载点', dataIndex: 'load_point_id', key: 'load_point_id' },
    { title: '卸载点', dataIndex: 'dump_point_id', key: 'dump_point_id' },
    {
      title: '评分', dataIndex: 'score', key: 'score',
      render: (v: number) => <Tag color={v > 0.8 ? 'green' : v > 0.6 ? 'orange' : 'red'}>{(v * 100).toFixed(0)}%</Tag>,
    },
    {
      title: '预估距离', dataIndex: 'estimated_distance_m', key: 'distance',
      render: (v: number) => (v / 1000).toFixed(1) + ' km',
    },
    {
      title: '预估时长', dataIndex: 'estimated_duration_s', key: 'duration',
      render: (v: number) => Math.round(v / 60) + ' min',
    },
    { title: '原因', dataIndex: 'reason', key: 'reason' },
  ];

  return (
    <div>
      <Typography.Title level={4}>AI 智能分析</Typography.Title>
      <Spin spinning={loading}>
        <Row gutter={[16, 16]}>
          <Col span={12}>
            <Card title="装载点需求" size="small">
              <Table
                columns={demandColumns}
                dataSource={loadingPoints}
                rowKey="load_point_id"
                pagination={false}
                size="small"
              />
            </Card>
          </Col>
          <Col span={12}>
            <Card title="卸载点需求" size="small">
              <Table
                columns={demandColumns}
                dataSource={dumpingPoints}
                rowKey="load_point_id"
                pagination={false}
                size="small"
              />
            </Card>
          </Col>
        </Row>

        <Card title="路段拥堵预测" size="small" style={{ marginTop: 16 }}>
          <Table
            columns={congestionColumns}
            dataSource={congestion}
            rowKey="edge_id"
            pagination={{ pageSize: 10, size: 'small' }}
            size="small"
          />
        </Card>

        <Card title="AI 调度建议" size="small" style={{ marginTop: 16 }}>
          <Table
            columns={suggestionColumns}
            dataSource={suggestions}
            rowKey={(r) => `${r.vehicle_id}-${r.load_point_id}`}
            pagination={false}
            size="small"
          />
        </Card>
      </Spin>
    </div>
  );
}
