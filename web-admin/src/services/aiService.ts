import apiClient from '../api/client';

export interface EdgeCongestion {
  edge_id: number;
  from_node_id: number;
  to_node_id: number;
  congestion_score: number;
  predicted_speed_kmh: number;
  predicted_vehicle_count: number;
  confidence: number;
}

export interface LoadingPointDemand {
  load_point_id: number;
  name: string;
  point_type: string;
  material: string;
  demand_score: number;
  pending_task_count: number;
  active_vehicle_count: number;
  confidence: number;
}

export interface AIVehicleInfo {
  vehicle_id: number;
  latitude: number;
  longitude: number;
  active_task_count: number;
}

export interface AITaskCandidate {
  load_point_id: number;
  dump_point_id: number;
}

export interface AISuggestion {
  vehicle_id: number;
  load_point_id: number;
  dump_point_id: number;
  score: number;
  estimated_distance_m: number;
  estimated_duration_s: number;
  reason: string;
}

export async function fetchCongestion(mineId = 1, lookbackMinutes = 60): Promise<EdgeCongestion[]> {
  const res = await apiClient.post('/api/v1/ai/congestion', { mine_id: mineId, lookback_minutes: lookbackMinutes });
  return res.data.data || [];
}

export async function fetchDemand(mineId = 1): Promise<LoadingPointDemand[]> {
  const res = await apiClient.post('/api/v1/ai/demand', { mine_id: mineId });
  return res.data.data || [];
}

export async function fetchSuggestions(
  vehicles: AIVehicleInfo[],
  tasks: AITaskCandidate[],
  mineId = 1,
): Promise<AISuggestion[]> {
  const res = await apiClient.post('/api/v1/ai/suggest-assign', { vehicles, tasks, mine_id: mineId });
  return res.data.suggestions || [];
}
