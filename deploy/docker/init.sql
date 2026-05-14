CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    username VARCHAR(64) UNIQUE NOT NULL,
    password VARCHAR(256) NOT NULL,
    real_name VARCHAR(64) DEFAULT '',
    email VARCHAR(128) DEFAULT '',
    phone VARCHAR(20) DEFAULT '',
    role INT DEFAULT 0,
    status INT DEFAULT 1,
    mine_id BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS vehicles (
    id BIGINT PRIMARY KEY,
    plate VARCHAR(32) UNIQUE NOT NULL,
    type INT DEFAULT 1,
    status INT DEFAULT 1,
    latitude DOUBLE PRECISION DEFAULT 0,
    longitude DOUBLE PRECISION DEFAULT 0,
    fuel_level DOUBLE PRECISION DEFAULT 100,
    mine_id BIGINT DEFAULT 0,
    driver_id BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS dispatch_tasks (
    id BIGINT PRIMARY KEY,
    vehicle_id BIGINT NOT NULL,
    load_point_id BIGINT NOT NULL,
    dump_point_id BIGINT NOT NULL,
    material VARCHAR(64) DEFAULT '',
    load_lat DOUBLE PRECISION DEFAULT 0,
    load_lon DOUBLE PRECISION DEFAULT 0,
    dump_lat DOUBLE PRECISION DEFAULT 0,
    dump_lon DOUBLE PRECISION DEFAULT 0,
    status VARCHAR(32) DEFAULT 'pending',
    algorithm VARCHAR(32) DEFAULT 'fifo',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_mine_id ON users(mine_id);
CREATE INDEX idx_vehicles_mine_id ON vehicles(mine_id);
CREATE INDEX idx_vehicles_status ON vehicles(status);
CREATE INDEX idx_dispatch_tasks_vehicle_id ON dispatch_tasks(vehicle_id);
CREATE INDEX idx_dispatch_tasks_status ON dispatch_tasks(status);

CREATE TABLE IF NOT EXISTS vehicle_types (
    id BIGINT PRIMARY KEY,
    name VARCHAR(64) UNIQUE NOT NULL,
    description VARCHAR(256) DEFAULT '',
    icon VARCHAR(64) DEFAULT '',
    capacity DOUBLE PRECISION DEFAULT 0,
    weight DOUBLE PRECISION DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS loading_points (
    id BIGINT PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    type VARCHAR(32) NOT NULL DEFAULT 'loading',
    latitude DOUBLE PRECISION DEFAULT 0,
    longitude DOUBLE PRECISION DEFAULT 0,
    material VARCHAR(64) DEFAULT '',
    status INT DEFAULT 1,
    mine_id BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_vehicle_types_name ON vehicle_types(name);
CREATE INDEX idx_loading_points_mine_id ON loading_points(mine_id);
CREATE INDEX idx_loading_points_type ON loading_points(type);

INSERT INTO vehicle_types (id, name, description, capacity, weight) VALUES
    (1, '矿用卡车', '大型矿山运输卡车', 60, 20),
    (2, '挖掘机', '矿山挖掘设备', 0, 30),
    (3, '装载机', '物料装载设备', 5, 15),
    (4, '推土机', '矿山推土设备', 0, 25)
ON CONFLICT (id) DO NOTHING;

INSERT INTO loading_points (id, name, type, material) VALUES
    (1, '装载点A', 'loading', '矿石'),
    (2, '装载点B', 'loading', '岩石'),
    (3, '卸载点C', 'dumping', '废石'),
    (4, '卸载点D', 'dumping', '矿石'),
    (5, '卸载点E', 'dumping', '岩石')
ON CONFLICT (id) DO NOTHING;

-- ── Alarm & Geofence ──

CREATE TABLE IF NOT EXISTS geofences (
    id BIGINT PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    shape VARCHAR(16) DEFAULT 'circle',
    center_lat DOUBLE PRECISION DEFAULT 0,
    center_lon DOUBLE PRECISION DEFAULT 0,
    radius_m DOUBLE PRECISION DEFAULT 0,
    points_json TEXT DEFAULT '',
    fence_type VARCHAR(32) DEFAULT 'restricted',
    min_speed_kmh INT DEFAULT 0,
    max_speed_kmh INT DEFAULT 0,
    time_range VARCHAR(32) DEFAULT '',
    enabled BOOLEAN DEFAULT TRUE,
    mine_id BIGINT DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS alarm_rules (
    id BIGINT PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    rule_type VARCHAR(32) NOT NULL,
    geofence_id BIGINT DEFAULT 0,
    severity VARCHAR(16) DEFAULT 'warning',
    description VARCHAR(256) DEFAULT '',
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS alarm_events (
    id BIGINT PRIMARY KEY,
    rule_id BIGINT DEFAULT 0,
    vehicle_id BIGINT DEFAULT 0,
    vehicle_plate VARCHAR(64) DEFAULT '',
    alarm_type VARCHAR(32) NOT NULL,
    severity VARCHAR(16) DEFAULT 'warning',
    message VARCHAR(512) DEFAULT '',
    latitude DOUBLE PRECISION DEFAULT 0,
    longitude DOUBLE PRECISION DEFAULT 0,
    speed DOUBLE PRECISION DEFAULT 0,
    acknowledged BOOLEAN DEFAULT FALSE,
    acknowledged_by VARCHAR(64) DEFAULT '',
    acknowledged_at TIMESTAMP WITH TIME ZONE,
    mine_id BIGINT DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_geofences_mine_id ON geofences(mine_id);
CREATE INDEX idx_geofences_fence_type ON geofences(fence_type);
CREATE INDEX idx_alarm_rules_rule_type ON alarm_rules(rule_type);
CREATE INDEX idx_alarm_events_vehicle_id ON alarm_events(vehicle_id);
CREATE INDEX idx_alarm_events_severity ON alarm_events(severity);
CREATE INDEX idx_alarm_events_acknowledged ON alarm_events(acknowledged);
CREATE INDEX idx_alarm_events_created_at ON alarm_events(created_at);

-- Seed: geofences (电子围栏)
INSERT INTO geofences (id, name, shape, center_lat, center_lon, radius_m, fence_type, max_speed_kmh, enabled) VALUES
    (1, '矿区-东区', 'circle', 39.9042, 116.4074, 500, 'restricted', 40, TRUE),
    (2, '矿区-西区', 'circle', 39.9142, 116.3974, 400, 'restricted', 30, TRUE),
    (3, '装载区安全围栏', 'circle', 39.9080, 116.4020, 200, 'loading', 20, TRUE),
    (4, '卸载区安全围栏', 'circle', 39.9000, 116.4120, 200, 'dumping', 15, TRUE)
ON CONFLICT (id) DO NOTHING;

-- Seed: alarm rules
INSERT INTO alarm_rules (id, name, rule_type, geofence_id, severity, description, enabled) VALUES
    (1, '东区禁区闯入', 'geofence', 1, 'critical', '车辆进入东区禁区触发告警', TRUE),
    (2, '西区禁区闯入', 'geofence', 2, 'critical', '车辆进入西区禁区触发告警', TRUE),
    (3, '装载区超速告警', 'geofence', 3, 'warning', '装载区内车速超过 20 km/h', TRUE),
    (4, '卸载区超速告警', 'geofence', 4, 'warning', '卸载区内车速超过 15 km/h', TRUE),
    (5, '严重超速告警', 'speeding', 0, 'critical', '车速超过 80 km/h', TRUE),
    (6, '超速警告', 'speeding', 0, 'warning', '车速超过 60 km/h', TRUE)
ON CONFLICT (id) DO NOTHING;

-- ── Road Network ──

CREATE TABLE IF NOT EXISTS road_nodes (
    id BIGINT PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    mine_id BIGINT DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS road_edges (
    id BIGINT PRIMARY KEY,
    from_node_id BIGINT NOT NULL,
    to_node_id BIGINT NOT NULL,
    distance_m DOUBLE PRECISION NOT NULL,
    max_speed_kmh INT DEFAULT 30,
    is_oneway BOOLEAN DEFAULT FALSE,
    mine_id BIGINT DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_road_nodes_mine_id ON road_nodes(mine_id);
CREATE INDEX idx_road_edges_from_node ON road_edges(from_node_id);
CREATE INDEX idx_road_edges_to_node ON road_edges(to_node_id);
CREATE INDEX idx_road_edges_mine_id ON road_edges(mine_id);

-- Seed: road network connecting loading/dumping points
INSERT INTO road_nodes (id, name, latitude, longitude) VALUES
    (100, '矿区入口', 39.8950, 116.4050),
    (101, '装载点A路口', 39.9080, 116.4020),
    (102, '装载点B路口', 39.9120, 116.3980),
    (103, '卸载点C路口', 39.9000, 116.4120),
    (104, '卸载点D路口', 39.8960, 116.4150),
    (105, '卸载点E路口', 39.9020, 116.4100),
    (106, '东区主干道-南', 39.9020, 116.4070),
    (107, '东区主干道-北', 39.9100, 116.4070),
    (108, '西区分岔口', 39.9080, 116.3950),
    (109, '停车场入口', 39.8980, 116.4000)
ON CONFLICT (id) DO NOTHING;

INSERT INTO road_edges (id, from_node_id, to_node_id, distance_m, max_speed_kmh, is_oneway) VALUES
    -- Main road south-north
    (1001, 100, 106, 600, 40, false),
    (1002, 106, 107, 950, 40, false),
    (1003, 106, 103, 500, 30, false),
    (1004, 106, 105, 350, 30, false),
    (1005, 107, 101, 450, 25, false),
    (1006, 107, 102, 700, 25, false),
    -- West branch
    (1007, 106, 108, 1100, 35, false),
    (1008, 108, 101, 700, 25, false),
    (1009, 108, 102, 550, 25, false),
    -- To dumping points
    (1010, 103, 104, 500, 30, false),
    (1011, 105, 104, 700, 30, false),
    (1012, 105, 103, 300, 20, false),
    -- Parking
    (1013, 100, 109, 400, 20, false),
    (1014, 106, 109, 500, 20, false)
ON CONFLICT (id) DO NOTHING;
