import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { Layout, Menu, Button, Typography } from 'antd';
import {
  DashboardOutlined,
  CarOutlined,
  EnvironmentOutlined,
  ScheduleOutlined,
  FlagOutlined,
  AppstoreOutlined,
  BellOutlined,
  SecurityScanOutlined,
  LogoutOutlined,
  BulbOutlined,
} from '@ant-design/icons';
import { useAuthStore } from '../stores/authStore';

const { Header, Sider, Content } = Layout;

const menuItems = [
  { key: '/dashboard', icon: <DashboardOutlined />, label: '调度看板' },
  { key: '/vehicles', icon: <CarOutlined />, label: '车辆管理' },
  { key: '/vehicle-types', icon: <AppstoreOutlined />, label: '车辆类型' },
  { key: '/loading-points', icon: <FlagOutlined />, label: '装载点管理' },
  { key: '/map', icon: <EnvironmentOutlined />, label: '实时地图' },
  { key: '/tasks', icon: <ScheduleOutlined />, label: '调度任务' },
  { key: '/alarms', icon: <BellOutlined />, label: '告警中心' },
  { key: '/geofences', icon: <SecurityScanOutlined />, label: '电子围栏' },
  { key: '/ai', icon: <BulbOutlined />, label: 'AI 智能分析' },
];

export default function AppLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout } = useAuthStore();

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider theme="dark" collapsible>
        <div style={{ height: 48, margin: 16, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
          <Typography.Title level={5} style={{ color: '#fff', margin: 0 }}>
            🚛 矿山调度
          </Typography.Title>
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header
          style={{
            background: '#fff',
            padding: '0 24px',
            display: 'flex',
            justifyContent: 'flex-end',
            alignItems: 'center',
            borderBottom: '1px solid #f0f0f0',
          }}
        >
          <Typography.Text style={{ marginRight: 16 }}>
            {user?.realName || user?.username}
          </Typography.Text>
          <Button type="text" icon={<LogoutOutlined />} onClick={logout}>
            退出
          </Button>
        </Header>
        <Content style={{ margin: 16, padding: 24, background: '#fff', borderRadius: 8 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
}
