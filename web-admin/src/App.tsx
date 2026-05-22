import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import AppLayout from './components/AppLayout';
import LoginPage from './pages/login/LoginPage';
import DashboardPage from './pages/dashboard/DashboardPage';
import VehiclePage from './pages/vehicles/VehiclePage';
import MapPage from './pages/map/MapPage';
import TaskPage from './pages/tasks/TaskPage';
import LoadingPointPage from './pages/loadingpoints/LoadingPointPage';
import VehicleTypePage from './pages/vehicletypes/VehicleTypePage';
import AlarmCenterPage from './pages/alarms/AlarmCenterPage';
import GeofencePage from './pages/geofences/GeofencePage';
import AiInsightsPage from './pages/ai/AiInsightsPage';
import ReportPage from './pages/reports/ReportPage';

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const token = localStorage.getItem('token');
  return token ? <>{children}</> : <Navigate to="/login" replace />;
}

export default function App() {
  return (
    <ConfigProvider locale={zhCN}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route
            path="/"
            element={
              <PrivateRoute>
                <AppLayout />
              </PrivateRoute>
            }
          >
            <Route index element={<Navigate to="/dashboard" replace />} />
            <Route path="dashboard" element={<DashboardPage />} />
            <Route path="vehicles" element={<VehiclePage />} />
            <Route path="map" element={<MapPage />} />
            <Route path="tasks" element={<TaskPage />} />
            <Route path="loading-points" element={<LoadingPointPage />} />
            <Route path="vehicle-types" element={<VehicleTypePage />} />
            <Route path="alarms" element={<AlarmCenterPage />} />
            <Route path="geofences" element={<GeofencePage />} />
            <Route path="ai" element={<AiInsightsPage />} />
            <Route path="reports" element={<ReportPage />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </ConfigProvider>
  );
}
