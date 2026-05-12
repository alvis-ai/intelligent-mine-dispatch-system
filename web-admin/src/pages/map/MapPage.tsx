import { useEffect, useRef, useState } from 'react';
import { Card, Typography, Tag, Space } from 'antd';
import apiClient from '../../api/client';

interface VehicleMarker {
  id: number;
  plate: string;
  lat: number;
  lng: number;
  speed: number;
}

export default function MapPage() {
  const mapRef = useRef<HTMLDivElement>(null);
  const [vehicles, setVehicles] = useState<VehicleMarker[]>([]);
  const [wsConnected, setWsConnected] = useState(false);

  useEffect(() => {
    apiClient.get('/api/v1/vehicles').then((res) => {
      const list = (res.data.data || []).map((v: any) => ({
        id: v.id,
        plate: v.plate,
        lat: 39.9 + Math.random() * 0.02,
        lng: 116.4 + Math.random() * 0.03,
        speed: Math.floor(Math.random() * 50),
      }));
      setVehicles(list);
    }).catch(() => {});
  }, []);

  useEffect(() => {
    const ws = new WebSocket(`ws://${window.location.hostname}:8080/ws/location`);

    ws.onopen = () => setWsConnected(true);
    ws.onclose = () => setWsConnected(false);
    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.vehicle_id) {
          setVehicles((prev) =>
            prev.map((v) =>
              v.id === data.vehicle_id
                ? { ...v, lat: data.latitude ?? v.lat, lng: data.longitude ?? v.lng, speed: data.speed ?? v.speed }
                : v
            )
          );
        }
      } catch {
        // ignore parse errors
      }
    };

    return () => ws.close();
  }, []);

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <Typography.Title level={4} style={{ margin: 0 }}>
          实时地图
        </Typography.Title>
        <Space>
          <Tag color={wsConnected ? 'green' : 'red'}>
            {wsConnected ? 'WebSocket 已连接' : 'WebSocket 未连接'}
          </Tag>
          <Typography.Text type="secondary">{vehicles.length} 辆车在线</Typography.Text>
        </Space>
      </div>
      <Card>
        <div
          ref={mapRef}
          style={{
            height: 600,
            background: '#f5f5f5',
            borderRadius: 8,
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            flexDirection: 'column',
            color: '#999',
          }}
        >
          <Typography.Title level={3} type="secondary">
            地图容器
          </Typography.Title>
          <Typography.Text type="secondary">
            集成高德地图 / Mapbox 后显示实时车辆位置
          </Typography.Text>
          <div style={{ marginTop: 16 }}>
            {vehicles.map((v) => (
              <Tag key={v.id} style={{ marginBottom: 4 }}>
                {v.plate} — {v.speed > 0 ? `${v.speed}km/h` : '静止'}
              </Tag>
            ))}
          </div>
        </div>
      </Card>
    </div>
  );
}
