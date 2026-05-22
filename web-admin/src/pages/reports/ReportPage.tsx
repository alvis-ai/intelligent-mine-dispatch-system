import { useState, useEffect } from 'react';
import { Row, Col, Card, Statistic, Table, Typography, Spin, Select, Space } from 'antd';
import { CarOutlined, ClockCircleOutlined, CheckCircleOutlined, WarningOutlined } from '@ant-design/icons';
import { Column, Line, Pie } from '@ant-design/charts';
import apiClient from '../../api/client';
import { getCurrentMineId } from '../../utils/mineContext';

export default function ReportPage() {
  const [loading, setLoading] = useState(true);
  const [summary, setSummary] = useState<any>({});
  const [dispatchData, setDispatchData] = useState<any[]>([]);
  const [utilData, setUtilData] = useState<any[]>([]);
  const [volumeData, setVolumeData] = useState<any[]>([]);
  const [alarmData, setAlarmData] = useState<any[]>([]);
  const [groupBy, setGroupBy] = useState('day');
  const [dateRange, _setDateRange] = useState<[string, string] | null>(null);

  useEffect(() => {
    const fetch = async () => {
      setLoading(true);
      try {
        const params: Record<string, string> = {};
        if (dateRange) {
          params.start_date = dateRange[0];
          params.end_date = dateRange[1];
        }
        const qs = new URLSearchParams(params).toString();

        const [sumRes, dispRes, utilRes, volRes, alarmRes] = await Promise.all([
          apiClient.get('/api/v1/reports/dashboard-summary'),
          apiClient.get(`/api/v1/reports/dispatch?group_by=${groupBy}&${qs}`),
          apiClient.get(`/api/v1/reports/vehicle-utilization?${qs}`),
          apiClient.get(`/api/v1/reports/transport-volume?${qs}`),
          apiClient.get(`/api/v1/reports/alarm-trend?${qs}`),
        ]);

        setSummary(sumRes.data.data || sumRes.data || {});
        setDispatchData(dispRes.data.rows || []);
        setUtilData(utilRes.data.rows || []);
        setVolumeData(volRes.data.rows || []);
        setAlarmData(alarmRes.data.rows || []);
      } catch {
        // ignore
      } finally {
        setLoading(false);
      }
    };
    fetch();
  }, [groupBy, dateRange]);

  const dispatchColumns = [
    { title: '维度', dataIndex: 'dimension', key: 'dimension' },
    { title: '总计', dataIndex: 'total', key: 'total' },
    { title: '已完成', dataIndex: 'completed', key: 'completed' },
    { title: '已取消', dataIndex: 'cancelled', key: 'cancelled' },
    { title: '进行中', dataIndex: 'active', key: 'active' },
    {
      title: '平均时长(min)', dataIndex: 'avgDurationMinutes', key: 'avg',
      render: (v: number) => v ? Math.round(v) : '-',
    },
  ];

  const utilColumns = [
    { title: '车辆', dataIndex: 'plate', key: 'plate' },
    { title: '总任务', dataIndex: 'totalTasks', key: 'totalTasks' },
    { title: '已完成', dataIndex: 'completedTasks', key: 'completedTasks' },
    {
      title: '利用率', dataIndex: 'utilizationRate', key: 'utilizationRate',
      render: (v: number) => (v * 100).toFixed(1) + '%',
    },
  ];

  const volumeColumns = [
    { title: '物料', dataIndex: 'material', key: 'material' },
    { title: '装载点', dataIndex: 'loadingPointName', key: 'loadingPointName' },
    { title: '任务数', dataIndex: 'taskCount', key: 'taskCount' },
  ];

  return (
    <div>
      <Typography.Title level={4}>BI 报表分析</Typography.Title>
      <Spin spinning={loading}>
        <Row gutter={[16, 16]}>
          <Col xs={24} sm={12} lg={4}>
            <Card size="small">
              <Statistic title="车辆总数" value={summary.totalVehicles} prefix={<CarOutlined />} />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={4}>
            <Card size="small">
              <Statistic title="进行中" value={summary.activeTasks} prefix={<ClockCircleOutlined />} />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={4}>
            <Card size="small">
              <Statistic title="今日调度" value={summary.todayDispatched} prefix={<CheckCircleOutlined />} />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={4}>
            <Card size="small">
              <Statistic
                title="未确认严重告警"
                value={summary.unacknowledgedCritical}
                prefix={<WarningOutlined />}
                valueStyle={{ color: summary.unacknowledgedCritical > 0 ? '#cf1322' : undefined }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={4}>
            <Card size="small">
              <Statistic title="待处理" value={summary.pendingTasks} />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={4}>
            <Card size="small">
              <Statistic title="已完成" value={summary.completedTasks} />
            </Card>
          </Col>
        </Row>

        <Card title="派车统计" size="small" style={{ marginTop: 16 }}
          extra={
            <Space>
              <Select value={groupBy} onChange={setGroupBy} size="small" style={{ width: 120 }}
                options={[
                  { value: 'day', label: '按天' },
                  { value: 'algorithm', label: '按算法' },
                  { value: 'status', label: '按状态' },
                ]}
              />
            </Space>
          }
        >
          <Row gutter={16}>
            <Col span={12}>
              <Column
                data={dispatchData}
                xField="dimension"
                yField="total"
                color="#1890ff"
                xAxis={{ label: { autoRotate: true } }}
                height={250}
              />
            </Col>
            <Col span={12}>
              <Table columns={dispatchColumns} dataSource={dispatchData} rowKey="dimension" pagination={false} size="small" />
            </Col>
          </Row>
        </Card>

        <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
          <Col span={12}>
            <Card title="车辆利用率" size="small">
              <Column
                data={utilData.slice(0, 10)}
                xField="plate"
                yField="utilizationRate"
                color="#52c41a"
                height={250}
                yAxis={{
                  label: { formatter: (v: string) => (parseFloat(v) * 100).toFixed(0) + '%' },
                }}
              />
              <Table columns={utilColumns} dataSource={utilData} rowKey="vehicleId" pagination={false} size="small" style={{ marginTop: 8 }} />
            </Card>
          </Col>
          <Col span={12}>
            <Card title="运输物料分布" size="small">
              <Pie
                data={volumeData}
                angleField="taskCount"
                colorField="material"
                radius={0.8}
                height={250}
                label={{ type: 'outer', content: '{name}: {value}' }}
              />
              <Table columns={volumeColumns} dataSource={volumeData} rowKey={(r: any) => r.material + r.loadingPointId} pagination={false} size="small" style={{ marginTop: 8 }} />
            </Card>
          </Col>
        </Row>

        <Card title="告警趋势" size="small" style={{ marginTop: 16 }}>
          <Line
            data={alarmData}
            xField="date"
            yField="count"
            seriesField="severity"
            height={250}
            color={['#cf1322', '#faad14', '#1890ff']}
            legend={{ position: 'top' }}
          />
        </Card>
      </Spin>
    </div>
  );
}
