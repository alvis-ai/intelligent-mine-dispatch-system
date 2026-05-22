import { useEffect, useRef, useState, useCallback } from 'react';
import { Card, Typography, Tag, Space, Spin, Select, Drawer, Descriptions, Switch } from 'antd';
import {
  EnvironmentOutlined, CarOutlined, LoadingOutlined, FireOutlined,
} from '@ant-design/icons';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import 'leaflet.heat';
import apiClient from '../../api/client';
import { fetchCongestion } from '../../services/aiService';

// ── Types ──

interface VehicleMarker {
  id: number;
  plate: string;
  lat: number;
  lng: number;
  speed: number;
  status: number;
  fuel_level?: number;
  driver_id?: number;
  marker?: L.Marker;
}

interface RoadNode {
  id: number;
  name: string;
  latitude: number;
  longitude: number;
}

interface RoadEdge {
  id: number;
  from_node_id: number;
  to_node_id: number;
  distance_m: number;
  max_speed_kmh: number;
  is_oneway: boolean;
}

interface RouteData {
  total_distance_m: number;
  total_duration_s: number;
  points: { latitude: number; longitude: number }[];
  node_ids: number[];
  edge_ids: number[];
}

// ── Helpers ──

function vehicleIcon(active: boolean): L.DivIcon {
  const color = active ? '#52c41a' : '#bfbfbf';
  return L.divIcon({
    className: '',
    html: `<div style="
      width: 32px; height: 32px; border-radius: 50%;
      background: ${color}; border: 3px solid #fff;
      box-shadow: 0 2px 6px rgba(0,0,0,0.3);
      display: flex; align-items: center; justify-content: center;
      font-size: 16px; color: #fff; cursor: pointer;
    ">🚛</div>`,
    iconSize: [32, 32],
    iconAnchor: [16, 16],
  });
}

function formatDuration(sec: number): string {
  if (sec < 60) return `${Math.round(sec)}s`;
  if (sec < 3600) return `${Math.floor(sec / 60)}m${Math.round(sec % 60)}s`;
  return `${Math.floor(sec / 3600)}h${Math.floor((sec % 3600) / 60)}m`;
}

// ── Component ──

export default function MapPage() {
  const mapContainer = useRef<HTMLDivElement>(null);
  const mapInstance = useRef<L.Map | null>(null);
  const vehicleLayers = useRef<Map<number, L.Marker>>(new Map());
  const roadNodeLayer = useRef<L.LayerGroup>(L.layerGroup());
  const roadEdgeLayer = useRef<L.LayerGroup>(L.layerGroup());
  const routeLayer = useRef<L.LayerGroup>(L.layerGroup());
  const heatLayer = useRef<L.Layer | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const [heatEnabled, setHeatEnabled] = useState(false);

  const [loading, setLoading] = useState(true);
  const [vehicles, setVehicles] = useState<VehicleMarker[]>([]);
  const [wsConnected, setWsConnected] = useState(false);
  const [selectedVehicle, setSelectedVehicle] = useState<VehicleMarker | null>(null);
  const [routeData, setRouteData] = useState<RouteData | null>(null);
  const [routeLoading, setRouteLoading] = useState(false);
  const [drawerOpen, setDrawerOpen] = useState(false);

  // ── Initialize map ──

  useEffect(() => {
    if (!mapContainer.current || mapInstance.current) return;

    const map = L.map(mapContainer.current, {
      center: [39.906, 116.407],
      zoom: 14,
      zoomControl: true,
    });

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors',
      maxZoom: 19,
    }).addTo(map);

    roadNodeLayer.current.addTo(map);
    roadEdgeLayer.current.addTo(map);
    routeLayer.current.addTo(map);
    mapInstance.current = map;

    // Pre-add empty heat layer
    heatLayer.current = L.heatLayer([], { radius: 30, blur: 20, maxZoom: 17 }).addTo(map);

    return () => {
      map.remove();
      mapInstance.current = null;
    };
  }, []);

  // ── Fetch road network ──

  const fetchRoadNetwork = useCallback(async () => {
    try {
      const [nodeRes, edgeRes] = await Promise.all([
        apiClient.get('/api/v1/road-nodes'),
        apiClient.get('/api/v1/road-edges'),
      ]);
      const nodes: RoadNode[] = nodeRes.data.data || [];
      const edges: RoadEdge[] = edgeRes.data.data || [];

      const map = mapInstance.current;
      if (!map) return;

      roadNodeLayer.current.clearLayers();
      roadEdgeLayer.current.clearLayers();

      // Draw edges as polylines
      const nodeMap = new Map(nodes.map((n) => [n.id, n]));
      for (const e of edges) {
        const from = nodeMap.get(e.from_node_id);
        const to = nodeMap.get(e.to_node_id);
        if (!from || !to) continue;

        const color = e.is_oneway ? '#faad14' : '#1890ff';
        const dash = e.is_oneway ? '10, 6' : '';

        L.polyline(
          [[from.latitude, from.longitude], [to.latitude, to.longitude]],
          { color, weight: 2, opacity: 0.6, dashArray: dash },
        ).addTo(roadEdgeLayer.current);
      }

      // Draw nodes
      for (const n of nodes) {
        const label = n.name || `N${n.id}`;
        L.circleMarker([n.latitude, n.longitude], {
          radius: 4,
          color: '#1890ff',
          fillColor: '#fff',
          fillOpacity: 1,
          weight: 2,
        })
          .bindTooltip(label, { permanent: false, direction: 'top' })
          .addTo(roadNodeLayer.current);
      }
    } catch {
      // Road network unavailable
    }
  }, []);

  // ── Fetch vehicles and update markers ──

  const updateVehicleMarkers = useCallback((list: VehicleMarker[]) => {
    const map = mapInstance.current;
    if (!map) return;

    const newIds = new Set(list.map((v) => v.id));

    // Remove stale markers
    for (const [id, marker] of vehicleLayers.current) {
      if (!newIds.has(id)) {
        map.removeLayer(marker);
        vehicleLayers.current.delete(id);
      }
    }

    // Add/update markers
    for (const v of list) {
      const existing = vehicleLayers.current.get(v.id);
      if (existing) {
        existing.setLatLng([v.lat, v.lng]);
        existing.setIcon(vehicleIcon(v.speed > 0));
        existing.bindTooltip(`${v.plate} ${v.speed}km/h`, {
          permanent: false, direction: 'top',
        });
      } else {
        const marker = L.marker([v.lat, v.lng], {
          icon: vehicleIcon(v.speed > 0),
        })
          .bindTooltip(`${v.plate} ${v.speed}km/h`, {
            permanent: false, direction: 'top',
          })
          .on('click', () => handleVehicleClick(v))
          .addTo(map);

        vehicleLayers.current.set(v.id, marker);
      }
    }
  }, []);

  const fetchVehicles = useCallback(async () => {
    try {
      const res = await apiClient.get('/api/v1/vehicles');
      const list: VehicleMarker[] = (res.data.data || []).map((v: any) => ({
        id: v.id,
        plate: v.plate,
        lat: v.latitude || 39.9,
        lng: v.longitude || 116.4,
        speed: 0,
        status: v.status || 1,
        fuel_level: v.fuel_level,
        driver_id: v.driver_id,
      }));
      setVehicles(list);
      updateVehicleMarkers(list);
    } catch {
      // ignore
    } finally {
      setLoading(false);
    }
  }, [updateVehicleMarkers]);

  // ── WebSocket for real-time locations ──

  useEffect(() => {
    const envWsUrl = import.meta.env.VITE_WS_URL;
    const wsUrl = envWsUrl !== undefined
      ? envWsUrl
      : `ws://${window.location.hostname}:8080/ws/location`;
    const ws = new WebSocket(wsUrl);

    ws.onopen = () => setWsConnected(true);
    ws.onclose = () => setWsConnected(false);
    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (!data.vehicle_id) return;

        setVehicles((prev) => {
          const updated = prev.map((v) =>
            v.id === data.vehicle_id
              ? {
                ...v,
                lat: data.latitude ?? v.lat,
                lng: data.longitude ?? v.lng,
                speed: data.speed ?? v.speed,
              }
              : v,
          );
          updateVehicleMarkers(updated);
          return updated;
        });
      } catch {
        // ignore parse errors
      }
    };

    wsRef.current = ws;
    return () => {
      ws.close();
      wsRef.current = null;
    };
  }, [updateVehicleMarkers]);

  // ── Toggle heatmap ──

  const toggleHeatmap = useCallback(async (enabled: boolean) => {
    setHeatEnabled(enabled);
    if (!enabled) {
      // Clear heatmap
      const map = mapInstance.current;
      if (map && heatLayer.current) {
        map.removeLayer(heatLayer.current);
        heatLayer.current = null;
      }
      return;
    }
    // Show loading state, then fetch and render congestion
    const data = await fetchCongestion(1, 60);
    if (!data || data.length === 0) return;
    const points: [number, number, number][] = data.map((c: any) => [
      // Use edge midpoint approximate: from/to average
      c.latitude || 39.9,
      c.longitude || 116.4,
      Math.min(c.congestion_score * 2, 1),
    ]);
    const map = mapInstance.current;
    if (!map) return;
    // Remove old heat layer if exists
    if (heatLayer.current) {
      map.removeLayer(heatLayer.current);
    }
    heatLayer.current = L.heatLayer(points.length > 0 ? points : [[39.906, 116.407, 0]], {
      radius: 30,
      blur: 20,
      maxZoom: 17,
      gradient: { 0.3: 'green', 0.5: 'orange', 0.8: 'red' },
    }).addTo(map);
  }, []);

  // ── Initial load ──

  useEffect(() => {
    fetchVehicles();
    fetchRoadNetwork();
  }, [fetchVehicles, fetchRoadNetwork]);

  // ── Handle vehicle click ──

  const handleVehicleClick = async (v: VehicleMarker) => {
    setSelectedVehicle(v);
    setDrawerOpen(true);

    // Try to get vehicle's current task destination
    try {
      const taskRes = await apiClient.get('/api/v1/dispatch/tasks', {
        params: { vehicle_id: v.id, status: 'active', page: 1, page_size: 1 },
      });
      const tasks: any[] = taskRes.data.data || [];
      if (tasks.length === 0) {
        setRouteData(null);
        return;
      }

      const task = tasks[0];
      setRouteLoading(true);

      const routeRes = await apiClient.post('/api/v1/route/calculate', {
        from_lat: v.lat,
        from_lon: v.lng,
        to_lat: task.dump_lat || v.lat + 0.01,
        to_lon: task.dump_lon || v.lng + 0.01,
        algorithm: 'astar',
      });

      if (routeRes.data.data) {
        const rd: RouteData = routeRes.data.data;
        setRouteData(rd);

        // Draw route on map
        routeLayer.current.clearLayers();
        if (rd.points && rd.points.length > 0) {
          const coords: [number, number][] = rd.points.map((p) => [p.latitude, p.longitude]);
          L.polyline(coords, {
            color: '#eb2f96',
            weight: 4,
            opacity: 0.8,
          }).addTo(routeLayer.current);
        }
      }
    } catch {
      // route unavailable
    } finally {
      setRouteLoading(false);
    }
  };

  // ── Clear route ──

  const closeDrawer = () => {
    setDrawerOpen(false);
    setSelectedVehicle(null);
    routeLayer.current.clearLayers();
    setRouteData(null);
  };

  // ── Render ──

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <Typography.Title level={4} style={{ margin: 0 }}>
          <EnvironmentOutlined /> 实时地图
        </Typography.Title>
        <Space>
          <Space>
            <FireOutlined style={{ color: heatEnabled ? '#cf1322' : undefined }} />
            <Switch checked={heatEnabled} onChange={toggleHeatmap} size="small" />
            <Typography.Text type="secondary">拥堵热力</Typography.Text>
          </Space>
          <Select
            size="small"
            style={{ width: 140 }}
            placeholder="跳转车辆"
            showSearch
            filterOption={(input, option) =>
              (option?.label as string)?.includes(input) ?? false
            }
            options={vehicles.map((v) => ({
              value: v.id,
              label: `${v.plate} (${v.speed}km/h)`,
            }))}
            onChange={(id) => {
              const v = vehicles.find((x) => x.id === id);
              if (v && mapInstance.current) {
                mapInstance.current.setView([v.lat, v.lng], 16);
                handleVehicleClick(v);
              }
            }}
          />
          <Tag icon={<CarOutlined />}>{vehicles.length} 辆车</Tag>
          <Tag color={wsConnected ? 'green' : 'red'}>
            {wsConnected ? '实时' : '离线'}
          </Tag>
        </Space>
      </div>

      <Card bodyStyle={{ padding: 0 }}>
        {loading ? (
          <div style={{
            height: 600, display: 'flex', justifyContent: 'center',
            alignItems: 'center', flexDirection: 'column', gap: 12,
          }}>
            <Spin indicator={<LoadingOutlined style={{ fontSize: 32 }} spin />} />
            <Typography.Text type="secondary">加载地图...</Typography.Text>
          </div>
        ) : (
          <div
            ref={mapContainer}
            style={{ height: 620, width: '100%', borderRadius: 8 }}
          />
        )}
      </Card>

      <Drawer
        title={selectedVehicle ? `车辆 ${selectedVehicle.plate}` : ''}
        placement="right"
        width={360}
        open={drawerOpen}
        onClose={closeDrawer}
      >
        {selectedVehicle && (
          <>
            <Descriptions column={1} size="small" bordered>
              <Descriptions.Item label="车牌">{selectedVehicle.plate}</Descriptions.Item>
              <Descriptions.Item label="速度">
                <Tag color={selectedVehicle.speed > 0 ? 'green' : 'default'}>
                  {selectedVehicle.speed} km/h
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="纬度">{selectedVehicle.lat.toFixed(6)}</Descriptions.Item>
              <Descriptions.Item label="经度">{selectedVehicle.lng.toFixed(6)}</Descriptions.Item>
              <Descriptions.Item label="油量">
                {selectedVehicle.fuel_level != null ? `${selectedVehicle.fuel_level}%` : '-'}
              </Descriptions.Item>
            </Descriptions>

            <div style={{ marginTop: 16 }}>
              <Typography.Text strong>路径信息</Typography.Text>
              {routeLoading ? (
                <div style={{ padding: '16px 0', textAlign: 'center' }}>
                  <Spin /> <Typography.Text type="secondary" style={{ marginLeft: 8 }}>计算路径中...</Typography.Text>
                </div>
              ) : routeData ? (
                <Descriptions column={1} size="small" bordered style={{ marginTop: 8 }}>
                  <Descriptions.Item label="距离">
                    {(routeData.total_distance_m / 1000).toFixed(2)} km
                  </Descriptions.Item>
                  <Descriptions.Item label="预计耗时">
                    {formatDuration(routeData.total_duration_s)}
                  </Descriptions.Item>
                  <Descriptions.Item label="途经节点">
                    {routeData.node_ids.length} 个
                  </Descriptions.Item>
                </Descriptions>
              ) : (
                <Typography.Text type="secondary" style={{ display: 'block', marginTop: 8 }}>
                  暂无活跃任务或路径无法计算
                </Typography.Text>
              )}
            </div>
          </>
        )}
      </Drawer>
    </div>
  );
}
