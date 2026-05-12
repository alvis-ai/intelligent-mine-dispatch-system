export interface User {
  id: number;
  username: string;
  realName: string;
  email: string;
  phone: string;
  role: number;
  status: number;
  mineId: number;
  createdAt: string;
}

export interface Vehicle {
  id: number;
  plate: string;
  type: number;
  status: number;
  latitude: number;
  longitude: number;
  fuelLevel: number;
  mineId: number;
  driverId: number;
}

export interface DispatchTask {
  id: number;
  vehicleId: number;
  loadPointId: number;
  dumpPointId: number;
  material: string;
  loadLat: number;
  loadLon: number;
  dumpLat: number;
  dumpLon: number;
  status: string;
  algorithm: string;
  createdAt: string;
}

export interface LocationData {
  vehicleId: number;
  latitude: number;
  longitude: number;
  speed: number;
  heading: number;
  timestamp: number;
}
